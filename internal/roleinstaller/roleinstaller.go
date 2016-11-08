package roleinstaller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gantsign/alt-galaxy/internal/application"
	"github.com/gantsign/alt-galaxy/internal/message"
	"github.com/gantsign/alt-galaxy/internal/restapi"
	"github.com/gantsign/alt-galaxy/internal/restclient"
	"github.com/gantsign/alt-galaxy/internal/rolesfile"
	"github.com/gantsign/alt-galaxy/internal/untar"
	"github.com/gantsign/alt-galaxy/internal/util"
)

const (
	maxConcurrentGitHubDownloads = 10
	maxConcurrentGalaxyRequests  = 10
	maxConcurrentUntar           = 2
)

type additionalRoleFields struct {
	Url         string
	Index       int
	ArchivePath string
}

type installerRole struct {
	rolesfile.Role
	additionalRoleFields
}

type roleInstaller struct {
	rolesPath               string
	roleLookupQueue         chan installerRole
	roleDownloadQueue       chan installerRole
	roleUntarQueue          chan installerRole
	roleOutputBuffers       []chan message.Message
	restClient              restclient.RestClient
	restApi                 restapi.RestApi
	galaxyRequestSemaphore  util.Semaphore
	gitHubDownloadSemaphore util.Semaphore
	untarSemaphore          util.Semaphore
	completionLatch         util.CompletionLatch
}

func (installer *roleInstaller) printOutput() {
	for _, output := range installer.roleOutputBuffers {
		for msg := range output {
			switch msg.MessageType {
			case message.OutMsg:
				fmt.Println("- ", msg.Body)
			case message.ErrorMsg:
				fmt.Fprintln(os.Stderr, "ERROR! ", msg.Body)
			default:
				fmt.Fprintln(os.Stderr, fmt.Sprintf("ERROR! Unsupported MessageType: %d", msg.MessageType))
			}
		}
	}
}

func (installer *roleInstaller) fail(role installerRole) {
	installer.roleLog(role).Progressf("%s install failed", role.Name)
	close(installer.roleOutputBuffers[role.Index])
	installer.completionLatch.Failure()
}

func (installer *roleInstaller) success(role installerRole) {
	installer.roleLog(role).Progressf("%s was installed successfully", role.Name)
	close(installer.roleOutputBuffers[role.Index])
	installer.completionLatch.Success()
}

func (installer *roleInstaller) lookupRole(role installerRole) {
	log := installer.roleLog(role)

	roleName, err := role.ParseRoleName()
	if err != nil {
		log.Errorf("Failed building query for role [%s].\nCaused by: %s", role.Name, err)
		installer.galaxyRequestSemaphore.Release()
		installer.fail(role)
		return
	}

	log.Progressf("downloading role '%s', owned by %s", roleName.RoleNamePart, roleName.UsernamePart)

	roleQueryResponse, err := installer.restApi.QueryRolesByName(roleName)
	if err != nil {
		log.Errorf("Failed querying details for role [%s].\nCaused by: %s", role.Name, err)
		installer.galaxyRequestSemaphore.Release()
		installer.fail(role)
		return
	}
	installer.galaxyRequestSemaphore.Release()

	roleDetails := roleQueryResponse.Results[0]
	if role.Version == "" {
		role.Version = roleDetails.LatestVersion()
	}
	role.Url = fmt.Sprintf("https://github.com/%s/%s/archive/%s.tar.gz", roleDetails.GitHubUser, roleDetails.GitHubRepo, role.Version)

	installer.roleDownloadQueue <- role
}

func (installer *roleInstaller) lookupRoles() {
	for role := range installer.roleLookupQueue {
		installer.galaxyRequestSemaphore.Acquire()

		go installer.lookupRole(role)
	}
}

func (installer *roleInstaller) downloadRole(role installerRole) {
	log := installer.roleLog(role)

	destFilePath := path.Join(installer.rolesPath, ".downloads", fmt.Sprint(role.Name, ".tar.gz"))

	log.Progressf("downloading role from %s", role.Url)
	destFilePath, err := installer.restClient.DownloadUrl(role.Url, destFilePath)
	if err != nil {
		log.Errorf("Failed to download URL [%s].\nCaused by: %s", role.Url, err)
		installer.gitHubDownloadSemaphore.Release()
		installer.fail(role)
		return
	}
	role.ArchivePath = destFilePath

	installer.gitHubDownloadSemaphore.Release()

	installer.roleUntarQueue <- role
}

func (installer *roleInstaller) downloadRoles() {
	for role := range installer.roleDownloadQueue {
		installer.gitHubDownloadSemaphore.Acquire()

		go installer.downloadRole(role)
	}
}

func (installer *roleInstaller) untarRole(role installerRole) {
	log := installer.roleLog(role)

	destDirPath := path.Join(installer.rolesPath, role.Name)

	log.Progressf("extracting %s to %s", role.Name, destDirPath)

	err := untar.Untar(log, role.ArchivePath, destDirPath)
	if err != nil {
		log.Errorf("Failed to untar archive [%s].\nCaused by: %s", role.ArchivePath, err)
		installer.untarSemaphore.Release()
		installer.fail(role)
		return
	}

	installer.untarSemaphore.Release()
	installer.success(role)
}

func (installer *roleInstaller) untarRoles() {
	for role := range installer.roleUntarQueue {
		installer.untarSemaphore.Acquire()

		go installer.untarRole(role)
	}
}

type RoleInstallerCmd struct {
	RoleFile  string
	RolesPath string
	ApiServer string
}

func (cmd RoleInstallerCmd) Execute() error {
	roles, err := rolesfile.ParseRolesFile(cmd.RoleFile)
	if err != nil {
		return fmt.Errorf("Failed to read role file [%s].\nCaused by: %s", cmd.RoleFile, err)
	}

	httpClient := &http.Client{}
	userAgent := fmt.Sprintf("%s/%s (+%s)", application.Name, application.Version, application.Repository)
	restClient, err := restclient.NewRestClient(httpClient, userAgent)
	if err != nil {
		return fmt.Errorf("Failed to create REST client.\nCaused by: %s", err)
	}

	baseUrl := fmt.Sprint(cmd.ApiServer, "/api/v1")
	restApi, err := restapi.NewRestApi(restClient, baseUrl)
	if err != nil {
		return fmt.Errorf("Failed to create REST API.\nCaused by: %s", err)
	}

	installer := &roleInstaller{
		rolesPath:               cmd.RolesPath,
		restClient:              restClient,
		restApi:                 restApi,
		roleLookupQueue:         make(chan installerRole, len(roles)),
		roleDownloadQueue:       make(chan installerRole, len(roles)),
		roleUntarQueue:          make(chan installerRole, len(roles)),
		completionLatch:         util.NewCompletionLatch(len(roles)),
		roleOutputBuffers:       make([]chan message.Message, len(roles)),
		galaxyRequestSemaphore:  util.NewSemaphore(maxConcurrentGalaxyRequests),
		gitHubDownloadSemaphore: util.NewSemaphore(maxConcurrentGitHubDownloads),
		untarSemaphore:          util.NewSemaphore(maxConcurrentUntar),
	}

	go installer.printOutput()
	go installer.lookupRoles()
	go installer.downloadRoles()
	go installer.untarRoles()

	for index, fileRole := range roles {
		role := installerRole{fileRole, additionalRoleFields{
			Index: index,
		}}
		installer.roleOutputBuffers[index] = make(chan message.Message, 20)

		if strings.HasPrefix(role.Src, "http://") || strings.HasPrefix(role.Src, "https://") {
			role.Url = role.Src
			installer.roleDownloadQueue <- role
		} else if strings.Contains(role.Src, "://") {
			installer.roleLog(role).Errorf("Unsupported protocol in URL [%s]; only 'http' and 'https' are supported.", role.Src)
			installer.fail(role)
		} else {
			role.Name = role.Src
			installer.roleLookupQueue <- role
		}
	}
	success := installer.completionLatch.Await()
	if !success {
		return errors.New("Failed to complete successfully. Any error output should be visible above. Please fix these errors and try again.")
	}

	return nil
}

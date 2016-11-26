package roleinstaller

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gantsign/alt-galaxy/internal/application"
	"github.com/gantsign/alt-galaxy/internal/message"
	"github.com/gantsign/alt-galaxy/internal/metadata"
	"github.com/gantsign/alt-galaxy/internal/restapi"
	"github.com/gantsign/alt-galaxy/internal/restclient"
	"github.com/gantsign/alt-galaxy/internal/rolesfile"
	"github.com/gantsign/alt-galaxy/internal/untar"
	"github.com/gantsign/alt-galaxy/internal/util"
)

const (
	maxConcurrentGitHubDownloads   = 10
	maxConcurrentGalaxyRequests    = 10
	maxConcurrentUntar             = 2
	maxConcurrentParseDependencies = 2
)

type additionalRoleFields struct {
	Url         string
	ArchivePath string
}

type installerRole struct {
	rolesfile.Role
	roleLog
	additionalRoleFields
}

type roleInstaller struct {
	rolesPath                  string
	roleLookupQueue            chan installerRole
	roleDownloadQueue          chan installerRole
	roleUntarQueue             chan installerRole
	roleParseDependenciesQueue chan installerRole
	roleOutputBuffers          chan chan message.Message
	restClient                 restclient.RestClient
	restApi                    restapi.RestApi
	galaxyRequestSemaphore     util.Semaphore
	gitHubDownloadSemaphore    util.Semaphore
	untarSemaphore             util.Semaphore
	parseDependenciesSemaphore util.Semaphore
	roleLatch                  util.CompletionLatch
	loggingLatch               util.CompletionLatch
	roleNames                  []string
}

func repoUrlToRoleName(repoUrl string) string {
	// gets the role name out of a repo like
	// http://git.example.com/repos/repo.git" => "repo"

	if !strings.Contains(repoUrl, "://") && !strings.Contains(repoUrl, "@") {
		return repoUrl
	}
	splitPath := strings.Split(repoUrl, "/")
	trailingPath := splitPath[len(splitPath)-1]
	if strings.HasSuffix(trailingPath, ".git") {
		trailingPath = trailingPath[:len(trailingPath)-4]
	}
	if strings.HasSuffix(trailingPath, ".tar.gz") {
		trailingPath = trailingPath[:len(trailingPath)-7]
	}
	if strings.Contains(trailingPath, ",") {
		trailingPath = strings.Split(trailingPath, ",")[0]
	}
	return trailingPath
}

func (installer *roleInstaller) printOutput() {
	stdOut := bufio.NewWriter(os.Stdout)
	stdErr := bufio.NewWriter(os.Stderr)
	for output := range installer.roleOutputBuffers {
		for msg := range output {
			switch msg.MessageType {
			case message.OutMsg:
				fmt.Fprintln(stdOut, "- ", msg.Body)
				stdOut.Flush()
			case message.ErrorMsg:
				// A short sleep helps the stdout and stderr render in the correct order
				time.Sleep(time.Second)

				fmt.Fprintln(stdErr, "ERROR! ", msg.Body)
				stdErr.Flush()

				// A short sleep helps the stdout and stderr render in the correct order
				time.Sleep(time.Second)
			default:
				// A short sleep helps the stdout and stderr render in the correct order
				time.Sleep(time.Second)

				fmt.Fprintln(stdErr, fmt.Sprintf("ERROR! Unsupported MessageType: %d", msg.MessageType))
				stdErr.Flush()

				// A short sleep helps the stdout and stderr render in the correct order
				time.Sleep(time.Second)
			}
		}
	}
	installer.loggingLatch.Success()
}

func (installer *roleInstaller) fail(role installerRole) {
	role.Progressf("%s install failed", role.Name)
	role.close()
	installer.roleLatch.Failure()
}

func (installer *roleInstaller) success(role installerRole) {
	role.Progressf("%s was installed successfully", role.Name)
	role.close()
	installer.roleLatch.Success()
}

func (installer *roleInstaller) lookupRole(role installerRole) {
	roleName, err := role.ParseRoleName()
	if err != nil {
		role.Errorf("Failed building query for role [%s].\nCaused by: %s", role.Name, err)
		installer.galaxyRequestSemaphore.Release()
		installer.fail(role)
		return
	}

	role.Progressf("downloading role '%s', owned by %s", roleName.RoleNamePart, roleName.UsernamePart)

	roleQueryResponse, err := installer.restApi.QueryRolesByName(roleName)
	if err != nil {
		role.Errorf("Failed querying details for role [%s].\nCaused by: %s", role.Name, err)
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
	destFilePath := path.Join(installer.rolesPath, ".downloads", fmt.Sprint(role.Name, ".tar.gz"))

	role.Progressf("downloading role from %s", role.Url)
	destFilePath, err := installer.restClient.DownloadUrl(role.Url, destFilePath)
	if err != nil {
		role.Errorf("Failed to download URL [%s].\nCaused by: %s", role.Url, err)
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
	destDirPath := path.Join(installer.rolesPath, role.Name)

	role.Progressf("extracting %s to %s", role.Name, destDirPath)

	err := untar.Untar(role, role.ArchivePath, destDirPath)
	if err != nil {
		role.Errorf("Failed to untar archive [%s].\nCaused by: %s", role.ArchivePath, err)
		installer.untarSemaphore.Release()
		installer.fail(role)
		return
	}

	installer.untarSemaphore.Release()

	installer.roleParseDependenciesQueue <- role
}

func (installer *roleInstaller) untarRoles() {
	for role := range installer.roleUntarQueue {
		installer.untarSemaphore.Acquire()

		go installer.untarRole(role)
	}
}

func (installer *roleInstaller) isDuplicateRole(roleName string) bool {
	for _, name := range installer.roleNames {
		if name == roleName {
			return true
		}
	}
	return false
}

func (installer *roleInstaller) addRole(fileRole rolesfile.Role) {
	if fileRole.Name == "" {
		fileRole.Name = repoUrlToRoleName(fileRole.Src)
	}

	outputBuffer := make(chan message.Message, 20)
	installer.roleOutputBuffers <- outputBuffer

	role := installerRole{fileRole, roleLog{outputBuffer}, additionalRoleFields{}}

	if strings.HasPrefix(role.Src, "http://") || strings.HasPrefix(role.Src, "https://") {
		role.Url = role.Src
		installer.roleDownloadQueue <- role
	} else if strings.Contains(role.Src, "://") {
		role.Errorf("Unsupported protocol in URL [%s]; only 'http' and 'https' are supported.", role.Src)
		installer.fail(role)
	} else {
		installer.roleLookupQueue <- role
	}
}

func (installer *roleInstaller) parseDependenciesForRole(role installerRole) {
	metadataPath := path.Join(installer.rolesPath, role.Name, "meta", "main.yml")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		// no metadata = no dependencies
		installer.parseDependenciesSemaphore.Release()
		installer.success(role)
		return
	}

	roleMetadata, err := metadata.ParseMetadataFile(metadataPath)
	if err != nil {
		role.Errorf("Failed to read role metadata [%s].\nCaused by: %s", metadataPath, err)
		installer.parseDependenciesSemaphore.Release()
		installer.fail(role)
		return
	}

	for _, dependency := range roleMetadata.Dependencies {
		if dependency.Name == "" {
			dependency.Name = repoUrlToRoleName(dependency.Src)
		}
		// We don't want role dependencies overwriting the versions of the roles
		// explicitly specified in the roles file.
		if installer.isDuplicateRole(dependency.Name) {
			continue
		}

		role.Progressf("adding dependency: %s", dependency.Name)

		installer.roleLatch.TaskAdded()

		installer.addRole(rolesfile.Role{
			Src:     dependency.Src,
			Name:    dependency.Name,
			Version: dependency.Version,
		})
	}

	installer.parseDependenciesSemaphore.Release()
	installer.success(role)
}

func (installer *roleInstaller) parseDependenciesForRoles() {
	for role := range installer.roleParseDependenciesQueue {
		installer.parseDependenciesSemaphore.Acquire()

		go installer.parseDependenciesForRole(role)
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

	roleNames := make([]string, len(roles))
	for index := range roles {
		if !strings.Contains(roles[index].Src, "://") {
			roles[index].Name = roles[index].Src
		}
		roleNames[index] = roles[index].Name
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

	queueSize := len(roles) + 100

	installer := &roleInstaller{
		rolesPath:                  cmd.RolesPath,
		restClient:                 restClient,
		restApi:                    restApi,
		roleLookupQueue:            make(chan installerRole, queueSize),
		roleDownloadQueue:          make(chan installerRole, queueSize),
		roleUntarQueue:             make(chan installerRole, queueSize),
		roleParseDependenciesQueue: make(chan installerRole, queueSize),
		roleLatch:                  util.NewCompletionLatch(len(roles)),
		loggingLatch:               util.NewCompletionLatch(1),
		roleOutputBuffers:          make(chan chan message.Message, queueSize),
		galaxyRequestSemaphore:     util.NewSemaphore(maxConcurrentGalaxyRequests),
		gitHubDownloadSemaphore:    util.NewSemaphore(maxConcurrentGitHubDownloads),
		untarSemaphore:             util.NewSemaphore(maxConcurrentUntar),
		parseDependenciesSemaphore: util.NewSemaphore(maxConcurrentParseDependencies),
		roleNames:                  roleNames,
	}

	go installer.printOutput()
	go installer.lookupRoles()
	go installer.downloadRoles()
	go installer.untarRoles()
	go installer.parseDependenciesForRoles()

	for _, fileRole := range roles {
		installer.addRole(fileRole)
	}

	success := installer.roleLatch.Await()

	close(installer.roleOutputBuffers)
	installer.loggingLatch.Await()

	if !success {
		return errors.New("Failed to complete successfully. Any error output should be visible above. Please fix these errors and try again.")
	}

	return nil
}

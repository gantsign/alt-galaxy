package roleinstaller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gantsign/alt-galaxy/internal/application"
	"github.com/gantsign/alt-galaxy/internal/logging"
	"github.com/gantsign/alt-galaxy/internal/metadata"
	"github.com/gantsign/alt-galaxy/internal/restapi"
	"github.com/gantsign/alt-galaxy/internal/restclient"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
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

type context struct {
	rolesPath                  string
	roleLookupQueue            chan model.Role
	roleDownloadQueue          chan model.Role
	roleUntarQueue             chan model.Role
	roleParseDependenciesQueue chan model.Role
	loggerFactory              logging.SerialLoggerFactory
	restClient                 restclient.RestClient
	restApi                    restapi.RestApi
	galaxyRequestSemaphore     util.Semaphore
	gitHubDownloadSemaphore    util.Semaphore
	untarSemaphore             util.Semaphore
	parseDependenciesSemaphore util.Semaphore
	roleLatch                  util.CompletionLatch
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

func (ctx *context) fail(role model.Role) {
	role.Progressf("%s install failed", role.Name)
	role.Close()
	ctx.roleLatch.Failure()
}

func (ctx *context) success(role model.Role) {
	role.Progressf("%s was installed successfully", role.Name)
	role.Close()
	ctx.roleLatch.Success()
}

func (ctx *context) lookupRole(role model.Role) {
	roleName, err := role.ParseRoleName()
	if err != nil {
		role.Errorf("Failed building query for role [%s].\nCaused by: %s", role.Name, err)
		ctx.galaxyRequestSemaphore.Release()
		ctx.fail(role)
		return
	}

	role.Progressf("downloading role '%s', owned by %s", roleName.RoleNamePart, roleName.UsernamePart)

	roleQueryResponse, err := ctx.restApi.QueryRolesByName(roleName)
	if err != nil {
		role.Errorf("Failed querying details for role [%s].\nCaused by: %s", role.Name, err)
		ctx.galaxyRequestSemaphore.Release()
		ctx.fail(role)
		return
	}
	ctx.galaxyRequestSemaphore.Release()

	roleDetails := roleQueryResponse.Results[0]
	if role.Version == "" {
		role.Version = roleDetails.LatestVersion()
	}
	role.Url = fmt.Sprintf("https://github.com/%s/%s/archive/%s.tar.gz", roleDetails.GitHubUser, roleDetails.GitHubRepo, role.Version)

	ctx.roleDownloadQueue <- role
}

func (ctx *context) lookupRoles() {
	for role := range ctx.roleLookupQueue {
		ctx.galaxyRequestSemaphore.Acquire()

		go ctx.lookupRole(role)
	}
}

func (ctx *context) downloadRole(role model.Role) {
	destFilePath := path.Join(ctx.rolesPath, ".downloads", fmt.Sprint(role.Name, ".tar.gz"))

	role.Progressf("downloading role from %s", role.Url)
	destFilePath, err := ctx.restClient.DownloadUrl(role.Url, destFilePath)
	if err != nil {
		role.Errorf("Failed to download URL [%s].\nCaused by: %s", role.Url, err)
		ctx.gitHubDownloadSemaphore.Release()
		ctx.fail(role)
		return
	}
	role.ArchivePath = destFilePath

	ctx.gitHubDownloadSemaphore.Release()

	ctx.roleUntarQueue <- role
}

func (ctx *context) downloadRoles() {
	for role := range ctx.roleDownloadQueue {
		ctx.gitHubDownloadSemaphore.Acquire()

		go ctx.downloadRole(role)
	}
}

func (ctx *context) untarRole(role model.Role) {
	destDirPath := path.Join(ctx.rolesPath, role.Name)

	role.Progressf("extracting %s to %s", role.Name, destDirPath)

	err := untar.Untar(role, role.ArchivePath, destDirPath)
	if err != nil {
		role.Errorf("Failed to untar archive [%s].\nCaused by: %s", role.ArchivePath, err)
		ctx.untarSemaphore.Release()
		ctx.fail(role)
		return
	}

	ctx.untarSemaphore.Release()

	ctx.roleParseDependenciesQueue <- role
}

func (ctx *context) untarRoles() {
	for role := range ctx.roleUntarQueue {
		ctx.untarSemaphore.Acquire()

		go ctx.untarRole(role)
	}
}

func (ctx *context) isDuplicateRole(roleName string) bool {
	for _, name := range ctx.roleNames {
		if name == roleName {
			return true
		}
	}
	return false
}

func (ctx *context) addRole(fileRole rolesfile.Role) {
	if fileRole.Name == "" {
		fileRole.Name = repoUrlToRoleName(fileRole.Src)
	}

	logger := ctx.loggerFactory.NewLogger()

	role := model.NewRole(fileRole, logger)

	if strings.HasPrefix(role.Src, "http://") || strings.HasPrefix(role.Src, "https://") {
		role.Url = role.Src
		ctx.roleDownloadQueue <- role
	} else if strings.Contains(role.Src, "://") {
		role.Errorf("Unsupported protocol in URL [%s]; only 'http' and 'https' are supported.", role.Src)
		ctx.fail(role)
	} else {
		ctx.roleLookupQueue <- role
	}
}

func (ctx *context) parseDependenciesForRole(role model.Role) {
	metadataPath := path.Join(ctx.rolesPath, role.Name, "meta", "main.yml")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		// no metadata = no dependencies
		ctx.parseDependenciesSemaphore.Release()
		ctx.success(role)
		return
	}

	roleMetadata, err := metadata.ParseMetadataFile(metadataPath)
	if err != nil {
		role.Errorf("Failed to read role metadata [%s].\nCaused by: %s", metadataPath, err)
		ctx.parseDependenciesSemaphore.Release()
		ctx.fail(role)
		return
	}

	for _, dependency := range roleMetadata.Dependencies {
		if dependency.Name == "" {
			dependency.Name = repoUrlToRoleName(dependency.Src)
		}
		// We don't want role dependencies overwriting the versions of the roles
		// explicitly specified in the roles file.
		if ctx.isDuplicateRole(dependency.Name) {
			continue
		}

		role.Progressf("adding dependency: %s", dependency.Name)

		ctx.roleLatch.TaskAdded()

		ctx.addRole(rolesfile.Role{
			Src:     dependency.Src,
			Name:    dependency.Name,
			Version: dependency.Version,
		})
	}

	ctx.parseDependenciesSemaphore.Release()
	ctx.success(role)
}

func (ctx *context) parseDependenciesForRoles() {
	for role := range ctx.roleParseDependenciesQueue {
		ctx.parseDependenciesSemaphore.Acquire()

		go ctx.parseDependenciesForRole(role)
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
		if roles[index].Name == "" {
			roles[index].Name = repoUrlToRoleName(roles[index].Src)
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

	ctx := &context{
		rolesPath:                  cmd.RolesPath,
		restClient:                 restClient,
		restApi:                    restApi,
		roleLookupQueue:            make(chan model.Role, queueSize),
		roleDownloadQueue:          make(chan model.Role, queueSize),
		roleUntarQueue:             make(chan model.Role, queueSize),
		roleParseDependenciesQueue: make(chan model.Role, queueSize),
		roleLatch:                  util.NewCompletionLatch(len(roles)),
		loggerFactory:              logging.NewSerialLoggerFactory(queueSize),
		galaxyRequestSemaphore:     util.NewSemaphore(maxConcurrentGalaxyRequests),
		gitHubDownloadSemaphore:    util.NewSemaphore(maxConcurrentGitHubDownloads),
		untarSemaphore:             util.NewSemaphore(maxConcurrentUntar),
		parseDependenciesSemaphore: util.NewSemaphore(maxConcurrentParseDependencies),
		roleNames:                  roleNames,
	}

	ctx.loggerFactory.StartOutput()

	go ctx.lookupRoles()
	go ctx.downloadRoles()
	go ctx.untarRoles()
	go ctx.parseDependenciesForRoles()

	for _, fileRole := range roles {
		ctx.addRole(fileRole)
	}

	success := ctx.roleLatch.Await()

	ctx.loggerFactory.Close()
	ctx.loggerFactory.AwaitOutputComplete()

	if !success {
		return errors.New("Failed to complete successfully. Any error output should be visible above. Please fix these errors and try again.")
	}

	return nil
}

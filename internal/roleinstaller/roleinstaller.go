package roleinstaller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gantsign/alt-galaxy/internal/application"
	"github.com/gantsign/alt-galaxy/internal/restapi"
	"github.com/gantsign/alt-galaxy/internal/restclient"
	installpipeline "github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/pipeline"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/step"
	"github.com/gantsign/alt-galaxy/internal/rolesfile"
)

const (
	maxConcurrentGitHubDownloads   = 10
	maxConcurrentGalaxyRequests    = 10
	maxConcurrentUntar             = 2
	maxConcurrentParseDependencies = 2
)

type context struct {
	rolesPath  string
	restClient restclient.RestClient
	restApi    restapi.RestApi
	roleNames  []string
	pipeline   installpipeline.Pipeline
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

func (ctx *context) RepoUrlToRoleName(repoUrl string) string {
	return repoUrlToRoleName(repoUrl)
}

func (ctx *context) RolesPath() string {
	return ctx.rolesPath
}

func (ctx *context) RestClient() restclient.RestClient {
	return ctx.restClient
}

func (ctx *context) RestApi() restapi.RestApi {
	return ctx.restApi
}

func (ctx *context) IsDuplicateRole(roleName string) bool {
	for _, name := range ctx.roleNames {
		if name == roleName {
			return true
		}
	}
	return false
}

func (ctx *context) InstallRole(role rolesfile.Role) {
	ctx.pipeline.InstallRole(role)
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
		rolesPath:  cmd.RolesPath,
		restClient: restClient,
		restApi:    restApi,
		roleNames:  roleNames,
	}

	pipeline := installpipeline.NewInstallPipeline(ctx, queueSize,
		[]installpipeline.Step{
			step.NewLookupRole(maxConcurrentGalaxyRequests),
			step.NewDownloadRole(maxConcurrentGitHubDownloads),
			step.NewUntarRole(maxConcurrentUntar),
			step.NewInstallDependencies(maxConcurrentParseDependencies),
		})
	ctx.pipeline = pipeline

	pipeline.Start()

	for _, role := range roles {
		pipeline.InstallRole(role)
	}

	success := pipeline.Await()

	if !success {
		return errors.New("Failed to complete successfully. Any error output should be visible above. Please fix these errors and try again.")
	}

	return nil
}

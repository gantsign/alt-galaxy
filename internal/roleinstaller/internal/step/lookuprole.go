package step

import (
	"fmt"
	"strings"

	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/pipeline"
	"github.com/gantsign/alt-galaxy/internal/util"
)

type lookupRole struct {
	pipeline.StepBase
}

func lookupRoleDetails(ctx model.Context, role model.Role) (model.Role, error) {
	roleName, err := role.ParseRoleName()
	if err != nil {
		return role, fmt.Errorf("Failed building query for role [%s].\nCaused by: %s", role.Name, err)
	}

	role.Progressf("downloading role '%s', owned by %s", roleName.RoleNamePart, roleName.UsernamePart)

	roleQueryResponse, err := ctx.RestApi().QueryRolesByName(roleName)
	if err != nil {
		return role, fmt.Errorf("Failed querying details for role [%s].\nCaused by: %s", role.Name, err)
	}

	roleDetails := roleQueryResponse.Results[0]
	if role.Version == "" {
		role.Version, err = roleDetails.LatestVersion()
		if err != nil {
			return role, fmt.Errorf("Failed to determine latest version for role [%s].\nCaused by: %s", role.Name, err)
		}
	}
	role.Url = fmt.Sprintf("https://github.com/%s/%s/archive/%s.tar.gz", roleDetails.GitHubUser, roleDetails.GitHubRepo, role.Version)

	return role, nil
}

func (step *lookupRole) processRoles() {
	ctx := step.Context()

	roleLookupQueue := make(chan model.Role, cap(step.RoleQueue))

	go func() {
		for role := range step.RoleQueue {
			if role.Name == "" {
				role.Name = ctx.RepoUrlToRoleName(role.Src)
			}
			if strings.HasPrefix(role.Src, "http://") || strings.HasPrefix(role.Src, "https://") {
				role.Url = role.Src
				step.Success(role)
			} else if strings.Contains(role.Src, "://") {
				role.Errorf("Unsupported protocol in URL [%s]; only 'http' and 'https' are supported.", role.Src)
				step.Fail(role)
			} else {
				roleLookupQueue <- role
			}
		}
	}()

	for role := range roleLookupQueue {
		step.ConcurrentlyProcessRole(role, lookupRoleDetails)
	}
}

func (step *lookupRole) Start() {
	go step.processRoles()
}

func NewLookupRole(maxConcurrent int) pipeline.Step {
	return &lookupRole{
		StepBase: pipeline.StepBase{
			Semaphore: util.NewSemaphore(maxConcurrent),
		},
	}
}

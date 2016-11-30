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
	semaphore util.Semaphore
}

func (step *lookupRole) lookupRole(ctx model.Context, role model.Role) {
	roleName, err := role.ParseRoleName()
	if err != nil {
		role.Errorf("Failed building query for role [%s].\nCaused by: %s", role.Name, err)
		step.semaphore.Release()
		step.Fail(role)
		return
	}

	role.Progressf("downloading role '%s', owned by %s", roleName.RoleNamePart, roleName.UsernamePart)

	roleQueryResponse, err := ctx.RestApi().QueryRolesByName(roleName)
	if err != nil {
		role.Errorf("Failed querying details for role [%s].\nCaused by: %s", role.Name, err)
		step.semaphore.Release()
		step.Fail(role)
		return
	}
	step.semaphore.Release()

	roleDetails := roleQueryResponse.Results[0]
	if role.Version == "" {
		role.Version = roleDetails.LatestVersion()
	}
	role.Url = fmt.Sprintf("https://github.com/%s/%s/archive/%s.tar.gz", roleDetails.GitHubUser, roleDetails.GitHubRepo, role.Version)

	step.Success(role)
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
		step.semaphore.Acquire()

		go step.lookupRole(ctx, role)
	}
}

func (step *lookupRole) Start() {
	go step.processRoles()
}

func NewLookupRole(maxConcurrent int) pipeline.Step {
	return &lookupRole{
		semaphore: util.NewSemaphore(maxConcurrent),
	}
}

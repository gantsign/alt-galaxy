package step

import (
	"path"

	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/pipeline"
	"github.com/gantsign/alt-galaxy/internal/untar"
)

func untarRole(ctx model.Context, step pipeline.Step, role model.Role) {
	destDirPath := path.Join(ctx.RolesPath(), role.Name)

	role.Progressf("extracting %s to %s", role.Name, destDirPath)

	err := untar.Untar(role, role.ArchivePath, destDirPath)
	if err != nil {
		role.Errorf("Failed to untar archive [%s].\nCaused by: %s", role.ArchivePath, err)
		step.Fail(role)
		return
	}

	step.Success(role)
}

func NewUntarRole(maxConcurrent int) pipeline.Step {
	return pipeline.NewStep(untarRole, maxConcurrent)
}

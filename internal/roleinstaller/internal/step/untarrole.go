package step

import (
	"fmt"
	"path"

	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/pipeline"
	"github.com/gantsign/alt-galaxy/internal/untar"
)

func untarRole(ctx model.Context, role model.Role) (model.Role, error) {
	destDirPath := path.Join(ctx.RolesPath(), role.Name)

	role.Progressf("extracting %s to %s", role.Name, destDirPath)

	err := untar.Untar(role, role.ArchivePath, destDirPath)
	if err != nil {
		return role, fmt.Errorf("Failed to untar archive [%s].\nCaused by: %s", role.ArchivePath, err)
	}

	return role, nil
}

func NewUntarRole(maxConcurrent int) pipeline.Step {
	return pipeline.NewStep(untarRole, maxConcurrent)
}

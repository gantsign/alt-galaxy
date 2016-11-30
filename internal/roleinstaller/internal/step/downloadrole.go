package step

import (
	"fmt"
	"path"

	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/pipeline"
)

func downloadRole(ctx model.Context, step pipeline.Step, role model.Role) {
	destFilePath := path.Join(ctx.RolesPath(), ".downloads", fmt.Sprint(role.Name, ".tar.gz"))

	role.Progressf("downloading role from %s", role.Url)
	destFilePath, err := ctx.RestClient().DownloadUrl(role.Url, destFilePath)
	if err != nil {
		role.Errorf("Failed to download URL [%s].\nCaused by: %s", role.Url, err)
		step.Fail(role)
		return
	}
	role.ArchivePath = destFilePath

	step.Success(role)
}

func NewDownloadRole(maxConcurrent int) pipeline.Step {
	return pipeline.NewStep(downloadRole, maxConcurrent)
}

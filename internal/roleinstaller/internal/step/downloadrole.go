package step

import (
	"fmt"
	"path"

	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/pipeline"
)

func downloadRole(ctx model.Context, role model.Role) (model.Role, error) {
	destFilePath := path.Join(ctx.RolesPath(), ".downloads", fmt.Sprint(role.Name, ".tar.gz"))

	role.Progressf("downloading role from %s", role.Url)
	destFilePath, err := ctx.RestClient().DownloadUrl(role.Url, destFilePath)
	if err != nil {
		return role, fmt.Errorf("Failed to download URL [%s].\nCaused by: %s", role.Url, err)
	}
	role.ArchivePath = destFilePath

	return role, nil
}

func NewDownloadRole(maxConcurrent int) pipeline.Step {
	return pipeline.NewStep(downloadRole, maxConcurrent)
}

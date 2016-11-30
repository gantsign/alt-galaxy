package step

import (
	"fmt"
	"os"
	"path"

	"github.com/gantsign/alt-galaxy/internal/metadata"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/pipeline"
	"github.com/gantsign/alt-galaxy/internal/rolesfile"
)

func parseDependenciesForRole(ctx model.Context, role model.Role) (model.Role, error) {
	metadataPath := path.Join(ctx.RolesPath(), role.Name, "meta", "main.yml")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		// no metadata = no dependencies
		return role, nil
	}

	roleMetadata, err := metadata.ParseMetadataFile(metadataPath)
	if err != nil {
		return role, fmt.Errorf("Failed to read role metadata [%s].\nCaused by: %s", metadataPath, err)
	}

	for _, dependency := range roleMetadata.Dependencies {
		if dependency.Name == "" {
			dependency.Name = ctx.RepoUrlToRoleName(dependency.Src)
		}
		// We don't want role dependencies overwriting the versions of the roles
		// explicitly specified in the roles file.
		if ctx.IsDuplicateRole(dependency.Name) {
			continue
		}

		role.Progressf("adding dependency: %s", dependency.Name)

		ctx.InstallRole(rolesfile.Role{
			Src:     dependency.Src,
			Name:    dependency.Name,
			Version: dependency.Version,
		})
	}

	return role, nil
}

func NewInstallDependencies(maxConcurrent int) pipeline.Step {
	return pipeline.NewStep(parseDependenciesForRole, maxConcurrent)
}

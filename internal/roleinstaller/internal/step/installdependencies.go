package step

import (
	"os"
	"path"

	"github.com/gantsign/alt-galaxy/internal/metadata"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/pipeline"
	"github.com/gantsign/alt-galaxy/internal/rolesfile"
)

func parseDependenciesForRole(ctx model.Context, step pipeline.Step, role model.Role) {
	metadataPath := path.Join(ctx.RolesPath(), role.Name, "meta", "main.yml")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		// no metadata = no dependencies
		step.Success(role)
		return
	}

	roleMetadata, err := metadata.ParseMetadataFile(metadataPath)
	if err != nil {
		role.Errorf("Failed to read role metadata [%s].\nCaused by: %s", metadataPath, err)
		step.Fail(role)
		return
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

	step.Success(role)
}

func NewInstallDependencies(maxConcurrent int) pipeline.Step {
	return pipeline.NewStep(parseDependenciesForRole, maxConcurrent)
}

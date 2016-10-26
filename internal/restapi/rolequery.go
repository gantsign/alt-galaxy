package restapi

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/gantsign/alt-galaxy/internal/rolesfile"
	"github.com/hashicorp/go-version"
)

type RoleQueryVersion struct {
	ReleaseDate string `json:"release_date"`
	Name        string
	Id          int
}

type RoleQuerySummaryFields struct {
	Versions []RoleQueryVersion `json:"versions"`
}

type RoleQueryResult struct {
	SummaryFields RoleQuerySummaryFields `json:"summary_fields"`
	GitHubUser    string                 `json:"github_user"`
	GitHubRepo    string                 `json:"github_repo"`
}

type RoleQueryResponse struct {
	Results []RoleQueryResult `json:"results"`
}

func ParseRoleQueryResponse(bytes []byte) (RoleQueryResponse, error) {
	var response RoleQueryResponse
	err := json.Unmarshal(bytes, &response)
	return response, err
}

func (result RoleQueryResult) LatestVersion() string {
	versions := result.SummaryFields.Versions
	if len(versions) == 0 {
		return "master"
	}
	libVersions := make([]*version.Version, len(versions))
	for i, rawVersion := range versions {
		ver, _ := version.NewVersion(rawVersion.Name)
		libVersions[i] = ver
	}

	sort.Sort(version.Collection(libVersions))
	return libVersions[len(libVersions)-1].String()
}

func (restApi restApiImpl) QueryRolesByName(roleName rolesfile.RoleName) (RoleQueryResponse, error) {
	// https://galaxy.ansible.com/api/v1/roles/?owner__username=gantsign&name=apt
	url := fmt.Sprintf("%s/roles/?owner__username=%s&name=%s", restApi.baseUrl, roleName.UsernamePart, roleName.RoleNamePart)

	_, respBytes, err := restApi.restClient.JsonHttpGet(url)
	if err != nil {
		return RoleQueryResponse{}, fmt.Errorf("GET request to [%s] failed.\nCaused by: %s", url, err)
	}

	return ParseRoleQueryResponse(respBytes)
}

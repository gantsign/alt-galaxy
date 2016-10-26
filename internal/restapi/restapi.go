package restapi

import (
	"fmt"

	"github.com/gantsign/alt-galaxy/internal/restclient"
	"github.com/gantsign/alt-galaxy/internal/rolesfile"
)

type RestApi interface {
	QueryRolesByName(roleName rolesfile.RoleName) (RoleQueryResponse, error)
}

type restApiImpl struct {
	restClient restclient.RestClient
	baseUrl    string
}

func NewRestApi(restClient restclient.RestClient, baseUrl string) (RestApi, error) {
	if restClient == nil {
		return nil, fmt.Errorf("Parameter [%s] is required.", "restClient")
	}
	if baseUrl == "" {
		return nil, fmt.Errorf("Parameter [%s] is required.", "baseUrl")
	}
	restApiImpl := &restApiImpl{
		restClient: restClient,
		baseUrl:    baseUrl,
	}
	return restApiImpl, nil
}

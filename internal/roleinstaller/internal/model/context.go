package model

import (
	"github.com/gantsign/alt-galaxy/internal/restapi"
	"github.com/gantsign/alt-galaxy/internal/restclient"
	"github.com/gantsign/alt-galaxy/internal/rolesfile"
)

type Context interface {
	RolesPath() string

	RestClient() restclient.RestClient

	RestApi() restapi.RestApi

	RepoUrlToRoleName(repoUrl string) string

	IsDuplicateRole(name string) bool

	InstallRole(role rolesfile.Role)
}

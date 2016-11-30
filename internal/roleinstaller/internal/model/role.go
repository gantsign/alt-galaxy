package model

import (
	"github.com/gantsign/alt-galaxy/internal/logging"
	"github.com/gantsign/alt-galaxy/internal/rolesfile"
)

type additionalRoleFields struct {
	Url         string
	ArchivePath string
}

type Role struct {
	rolesfile.Role
	logging.SerialLogger
	additionalRoleFields
}

func NewRole(fileRole rolesfile.Role, logger logging.SerialLogger) Role {
	return Role{fileRole, logger, additionalRoleFields{}}
}

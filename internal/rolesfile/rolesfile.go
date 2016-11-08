package rolesfile

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Role struct {
	Src     string
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
}

func ParseRolesFile(filename string) ([]Role, error) {
	yamlBytes, readErr := ioutil.ReadFile(filename)
	if readErr != nil {
		return nil, readErr
	}
	var roles []Role
	parseErr := yaml.Unmarshal(yamlBytes, &roles)

	return roles, parseErr
}

type RoleName struct {
	UsernamePart string
	RoleNamePart string
}

func (role Role) ParseRoleName() (RoleName, error) {
	parts := strings.Split(role.Name, ".")
	if len(parts) != 2 {
		return RoleName{}, fmt.Errorf("Unable to parse role name [%s].", role.Name)
	}
	roleName := RoleName{
		UsernamePart: parts[0],
		RoleNamePart: parts[1],
	}
	return roleName, nil
}

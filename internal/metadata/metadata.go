package metadata

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ComplexRole struct {
	Name string `yaml:"role,omitempty"`
}

type Role struct {
	Name string
}

func (role *Role) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Maybe a simple string
	var roleName string
	err := unmarshal(&roleName)
	if err == nil {
		role.Name = roleName
		return nil
	}

	// Maybe an object literal
	var complexRole ComplexRole
	err = unmarshal(&complexRole)
	if err != nil {
		return err
	}
	role.Name = complexRole.Name
	return nil
}

type Metadata struct {
	Dependencies []Role `yaml:"dependencies,omitempty"`
}

func ParseMetadataFile(filename string) (Metadata, error) {
	yamlBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Metadata{}, err
	}
	var metadata Metadata
	err = yaml.Unmarshal(yamlBytes, &metadata)

	return metadata, err
}

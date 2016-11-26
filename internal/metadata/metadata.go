package metadata

import (
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type complexRole struct {
	Role    string `yaml:"role,omitempty"`
	Src     string `yaml:"src,omitempty"`
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type Role struct {
	Src     string
	Name    string
	Version string
}

func (role *Role) setPropertiesFromString(specification string) {
	tokens := strings.Split(specification, ",")
	role.Src = tokens[0]

	if len(tokens) >= 2 {
		role.Version = tokens[1]
	}

	if len(tokens) >= 3 {
		role.Name = tokens[2]
	}
}

func (role *Role) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Maybe a simple string
	var roleSpecification string
	err := unmarshal(&roleSpecification)
	if err == nil {
		role.setPropertiesFromString(roleSpecification)
		return nil
	}

	// Maybe an object literal
	var complexRole complexRole
	err = unmarshal(&complexRole)
	if err != nil {
		return err
	}

	if complexRole.Role != "" {
		role.setPropertiesFromString(complexRole.Role)
		return nil
	}

	if complexRole.Src == "" {
		return errors.New("Unable to parse dependencies: one of property 'role' or 'src' must be specified.")
	}

	role.Src = complexRole.Src
	role.Version = complexRole.Version
	role.Name = complexRole.Name
	return nil
}

type Metadata struct {
	Dependencies []Role `yaml:"dependencies,omitempty"`
}

func parseMetadataBytes(yamlBytes []byte) (Metadata, error) {
	var metadata Metadata
	err := yaml.Unmarshal(yamlBytes, &metadata)

	return metadata, err
}

func ParseMetadataFile(filename string) (Metadata, error) {
	yamlBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Metadata{}, err
	}
	return parseMetadataBytes(yamlBytes)
}

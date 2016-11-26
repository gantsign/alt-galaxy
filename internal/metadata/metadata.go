package metadata

import (
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type complexRole struct {
	Role string `yaml:"role"`
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

	if complexRole.Role == "" {
		return errors.New("Unable to parse dependencies: expected property 'role' is missing or empty.")
	}
	role.setPropertiesFromString(complexRole.Role)
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

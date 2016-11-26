package metadata

import (
	"reflect"
	"testing"
)

func parseMetadataString(yaml string) (Metadata, error) {
	return parseMetadataBytes([]byte(yaml))
}

func (role Role) toMetadata() Metadata {
	dependencies := make([]Role, 1)
	dependencies[0] = role

	var metadata Metadata
	metadata.Dependencies = dependencies
	return metadata
}

func TestStringDependency(t *testing.T) {
	actual, err := parseMetadataString(`{
		"dependencies": [
			"gantsign.test"
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	expected := Role{
		Src: "gantsign.test",
	}.toMetadata()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%+v], actual [%+v].", expected, actual)
	}
}

func TestStringDependencyVersion(t *testing.T) {
	actual, err := parseMetadataString(`{
		"dependencies": [
			"gantsign.test,1.0"
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	expected := Role{
		Src:     "gantsign.test",
		Version: "1.0",
	}.toMetadata()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%+v], actual [%+v].", expected, actual)
	}
}

func TestStringDependencyVersionName(t *testing.T) {
	actual, err := parseMetadataString(`{
		"dependencies": [
			"gantsign.test,1.0,test"
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	expected := Role{
		Src:     "gantsign.test",
		Version: "1.0",
		Name:    "test",
	}.toMetadata()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%+v], actual [%+v].", expected, actual)
	}
}

func TestRoleStringDependency(t *testing.T) {
	actual, err := parseMetadataString(`{
		"dependencies": [
			{
				"role": "gantsign.test"
			}
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	expected := Role{
		Src: "gantsign.test",
	}.toMetadata()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%+v], actual [%+v].", expected, actual)
	}
}

func TestRoleStringDependencyVersion(t *testing.T) {
	actual, err := parseMetadataString(`{
		"dependencies": [
			{
				"role": "gantsign.test,1.0"
			}
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	expected := Role{
		Src:     "gantsign.test",
		Version: "1.0",
	}.toMetadata()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%+v], actual [%+v].", expected, actual)
	}
}

func TestRoleStringDependencyVersionName(t *testing.T) {
	actual, err := parseMetadataString(`{
		"dependencies": [
			{
				"role": "gantsign.test,1.0,test"
			}
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	expected := Role{
		Src:     "gantsign.test",
		Version: "1.0",
		Name:    "test",
	}.toMetadata()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%+v], actual [%+v].", expected, actual)
	}
}

func TestObjectDependency(t *testing.T) {
	actual, err := parseMetadataString(`{
		"dependencies": [
			{
				"src": "gantsign.test"
			}
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	expected := Role{
		Src: "gantsign.test",
	}.toMetadata()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%+v], actual [%+v].", expected, actual)
	}
}

func TestObjectDependencyVersion(t *testing.T) {
	actual, err := parseMetadataString(`{
		"dependencies": [
			{
				"src": "gantsign.test",
				"version": "1.0"
			}
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	expected := Role{
		Src:     "gantsign.test",
		Version: "1.0",
	}.toMetadata()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%+v], actual [%+v].", expected, actual)
	}
}

func TestObjectDependencyVersionName(t *testing.T) {
	actual, err := parseMetadataString(`{
		"dependencies": [
			{
				"src": "gantsign.test",
				"version": "1.0",
				"name": "test"
			}
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	expected := Role{
		Src:     "gantsign.test",
		Version: "1.0",
		Name:    "test",
	}.toMetadata()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%+v], actual [%+v].", expected, actual)
	}
}

func TestMixedDependencies(t *testing.T) {
	actual, err := parseMetadataString(`{
		"dependencies": [
			"gantsign.test,1.0,test",
			{
				"role": "gantsign.test2,2.0,test2"
			},
			{
				"src": "gantsign.test3",
				"version": "3.0",
				"name": "test3"
			}
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	dependencies := make([]Role, 3)
	dependencies[0] = Role{
		Src:     "gantsign.test",
		Version: "1.0",
		Name:    "test",
	}
	dependencies[1] = Role{
		Src:     "gantsign.test2",
		Version: "2.0",
		Name:    "test2",
	}
	dependencies[2] = Role{
		Src:     "gantsign.test3",
		Version: "3.0",
		Name:    "test3",
	}
	var expected Metadata
	expected.Dependencies = dependencies

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%+v], actual [%+v].", expected, actual)
	}
}

func TestBadYaml(t *testing.T) {
	_, err := parseMetadataString("bad YAML!")
	if err == nil {
		t.Error("Expected error parsing bad YAML")
	}
}

package restapi

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gantsign/alt-galaxy/internal/application"
	"github.com/gantsign/alt-galaxy/internal/restclient"
)

func TestNewRestApi(t *testing.T) {
	httpClient := &http.Client{}
	userAgent := fmt.Sprintf("%s/%s (+%s)", application.Name, application.Version, application.Repository)
	restClient, err := restclient.NewRestClient(httpClient, userAgent)
	if err != nil {
		t.Errorf("Failed to create REST client.\nCaused by: %s", err)
		return
	}

	baseUrl := "https://galaxy.ansible.com/api/v1"
	_, err = NewRestApi(restClient, baseUrl)
	if err != nil {
		t.Errorf("Failed to create REST API.\nCaused by: %s", err)
	}
}

func TestNewRestApiNilClient(t *testing.T) {
	baseUrl := "https://galaxy.ansible.com/api/v1"
	_, err := NewRestApi(nil, baseUrl)
	if err == nil {
		t.Error("Error was expected")
		return
	}
	expected := "Parameter [restClient] is required."
	actual := err.Error()
	if actual != expected {
		t.Errorf("Expected [%s] != actual [%s]", expected, actual)
	}
}

func TestNewRestApiEmptyBaseUrl(t *testing.T) {
	httpClient := &http.Client{}
	userAgent := fmt.Sprintf("%s/%s (+%s)", application.Name, application.Version, application.Repository)
	restClient, err := restclient.NewRestClient(httpClient, userAgent)
	if err != nil {
		t.Errorf("Failed to create REST client.\nCaused by: %s", err)
		return
	}

	baseUrl := ""
	_, err = NewRestApi(restClient, baseUrl)
	if err == nil {
		t.Error("Error was expected")
		return
	}
	expected := "Parameter [baseUrl] is required."
	actual := err.Error()
	if actual != expected {
		t.Errorf("Expected [%s] != actual [%s]", expected, actual)
	}
}

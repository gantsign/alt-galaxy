package restclient

import (
	"fmt"
	"net/http"
)

type RestClient interface {
	JsonHttpGet(url string) (*http.Response, []byte, error)

	DownloadUrl(url string, destFilePath string) (string, error)
}

type restClientImpl struct {
	httpClient *http.Client
}

func NewRestClient(httpClient *http.Client) (RestClient, error) {
	if httpClient == nil {
		return nil, fmt.Errorf("Parameter [%s] is required.", "httpClient")
	}
	return &restClientImpl{httpClient: httpClient}, nil
}

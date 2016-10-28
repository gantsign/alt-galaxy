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
	userAgent  string
}

func NewRestClient(httpClient *http.Client, userAgent string) (RestClient, error) {
	if httpClient == nil {
		return nil, fmt.Errorf("Parameter [%s] is required.", "httpClient")
	}
	restClient := &restClientImpl{
		httpClient: httpClient,
		userAgent:  userAgent,
	}
	return restClient, nil
}

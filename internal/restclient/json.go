package restclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (restClient *restClientImpl) JsonHttpGet(url string) (*http.Response, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Invalid URL [%s].\nCaused by: %s", url, err)
	}

	req.Header.Add("User-Agent", restClient.userAgent)
	req.Header.Add("Accept", "application/json")

	resp, err := restClient.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("GET request to URL [%s] failed.\nCaused by: %s", url, err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed reading response body for URL [%s].\nCaused by: %s", url, err)
	}
	return resp, respBytes, nil
}

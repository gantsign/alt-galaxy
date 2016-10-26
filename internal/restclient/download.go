package restclient

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

func (restClient restClientImpl) DownloadUrl(url string, destFilePath string) (string, error) {
	absoluteDestFilePath, err := filepath.Abs(destFilePath)
	if err != nil {
		return "", fmt.Errorf("Failed to determine absolute file path for [%s].\nCaused by: %s", destFilePath, err)
	}

	destDirPath, err := filepath.Abs(path.Join(destFilePath, ".."))
	if err != nil {
		return "", fmt.Errorf("Failed to determine parent path for [%s].\nCaused by: %s", destFilePath, err)
	}

	err = os.MkdirAll(destDirPath, 0755)
	if err != nil {
		return "", fmt.Errorf("Failed to create dir [%s].\nCaused by: %s", destDirPath, err)
	}

	destFile, err := os.OpenFile(absoluteDestFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to creating file [%s].\nCaused by: %s", absoluteDestFilePath, err)
	}
	defer destFile.Close()

	resp, err := restClient.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("Failed sending GET request to URL [%s].\nCaused by: %s", url, err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(destFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to download URL [%s].\nCaused by: %s", url, err)
	}
	return absoluteDestFilePath, nil
}

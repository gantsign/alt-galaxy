package untar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/gantsign/alt-galaxy/internal/logging"
)

func extractEntry(log logging.Log, destRootDirPath string, header *tar.Header, tarReader io.Reader) error {

	entryPath := strings.TrimRight(header.Name, "/")
	parts := strings.Split(entryPath, "/")[1:] // strip off root directory

	// Skip root directory
	if len(parts) == 0 {
		return nil
	}

	for _, part := range parts {
		if part == ".." || strings.Contains(part, "~") || strings.Contains(part, "$") {
			return fmt.Errorf("Illegal path [%s] in tar element path [%s].", part, entryPath)
		}
	}

	entryPath = path.Join(parts...)

	destPath := path.Join(destRootDirPath, entryPath)

	switch header.Typeflag {
	case tar.TypeDir:
		err := os.MkdirAll(destPath, os.FileMode(header.Mode))
		if err != nil {
			return fmt.Errorf("Failed to create directory [%s].\nCaused by: %s", destPath, err)
		}
	case tar.TypeReg:
		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return fmt.Errorf("Failed to create file [%s].\nCaused by: %s", destPath, err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, tarReader)
		if err != nil {
			return fmt.Errorf("Failed writing file [%s].\nCaused by: %s", destPath, err)
		}
	case tar.TypeXGlobalHeader:
	// ignore
	default:
		log.Errorf("Unsupported tar entry type [%d] for entry [%s].", header.Typeflag, header.Name)
	}
	return nil
}

func Untar(log logging.Log, archivePath string, destDirPath string) error {
	if _, err := os.Stat(destDirPath); !os.IsNotExist(err) {
		err = os.RemoveAll(destDirPath)
		if err != nil {
			return fmt.Errorf("Failed to delete directory [%s].\nCaused by: %s", destDirPath, err)
		}
	}

	err := os.MkdirAll(destDirPath, 0755)
	if err != nil {
		return fmt.Errorf("Error creating directory [%s].\nCaused by: %s", destDirPath, err)
	}

	tarFile, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("Error opening file [%s].\nCaused by: %s", archivePath, err)
	}
	defer tarFile.Close()

	gzipReader, err := gzip.NewReader(tarFile)
	if err != nil {
		return fmt.Errorf("Error opening gzip reader for file [%s].\nCaused by: %s", archivePath, err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// Iterate through the files in the archive.
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		if err != nil {
			return fmt.Errorf("Error reading tar archive [%s].\nCaused by: %s", archivePath, err)
		}

		err = extractEntry(log, destDirPath, header, tarReader)
		if err != nil {
			return fmt.Errorf("Error processing tar entry [%s].\nCaused by: %s", header.Name, err)
		}
	}
	return nil
}

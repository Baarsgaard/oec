package test_util

import (
	"os"
)

func CreateTempTestFile(content []byte, fileExtension string) (string, error) {

	tempFile, err := os.CreateTemp("", "*"+fileExtension)
	if err != nil {
		return "", err
	}

	defer tempFile.Close()

	if _, err := tempFile.Write(content); err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

package utils

import "path/filepath"

func getFileNameAndAbsPath(path string) (string, string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", "", err
	}

	fileName := filepath.Base(path)

	return absPath, fileName, nil
}

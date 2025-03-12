// Package os provides a set of functions to interact with the operating system.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"os"
	"path/filepath"
)

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return err
		}
	}
	return nil
}

func CopyFileToDestination(srcFile, destFile string) error {
	input, err := os.ReadFile(srcFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(destFile, input, 0600)
	if err != nil {
		return err
	}
	return nil
}

// CreateFileWithPath creates a file with the specified path.
// It will create any required directories in the path if they don't exist.
func CreateFileWithPath(path string) (*os.File, error) {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, err
	}

	// Create the file
	return os.Create(path)
}

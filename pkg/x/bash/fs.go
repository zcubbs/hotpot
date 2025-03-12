// Package bash
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package bash

import (
	"fmt"
	"os"
)

func Chmod(path string, perm os.FileMode, debug bool) error {
	if perm == 0 {
		perm = 0755
	}
	stats, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting file stats: %s", err)
	}

	if debug {
		fmt.Printf("File permissions before: %s\n", stats.Mode())
	}

	err = os.Chmod(path, perm)
	if err != nil {
		return fmt.Errorf("error changing file permissions: %s", err)
	}

	stats, err = os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting file stats: %s", err)
	}

	if debug {
		fmt.Printf("File permissions after: %s\n", stats.Mode())
	}

	return nil
}

func RmDir(path string) error {
	return os.RemoveAll(path)
}

func Mkdir(path string, perm os.FileMode) error {
	if perm == 0 {
		perm = 0755
	}
	return os.MkdirAll(path, perm)
}

func FileExists(path string) bool {
	// check if file exists
	// if it does, return true
	// if it doesn't, return false
	// TODO: implement check if file exists

	return false
}

func DirExists(path string) bool {
	// check if directory exists
	// if it does, return true
	// if it doesn't, return false
	// TODO: implement check if directory exists

	return false
}

func ExtractTarGz(tarPath string, destPath string) error {
	// tar -zxvf <tarPath> -C <destPath>
	err := ExecuteCmd("tar", true, "-zxvf", tarPath, "-C", destPath)
	if err != nil {
		fmt.Println("Error extracting tar.gz file: ", tarPath, err)
		return err
	}

	return nil
}

func ExtractTarGzWithFile(tarPath string, file string, destPath string, debug bool) error {
	// tar -zxvf <tarPath> -C <destPath>
	err := ExecuteCmd("tar", debug, "-zxvf", tarPath, "-C", destPath, file)
	if err != nil {
		return fmt.Errorf("error extracting tar.gz file: %s\n%w", tarPath, err)
	}

	return nil
}

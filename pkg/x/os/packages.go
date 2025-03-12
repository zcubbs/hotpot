// Package os provides a set of functions to interact with the operating system.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"fmt"
	"os/exec"
)

const (
	BinSH = "/bin/sh"
)

func isValidPackageName(name string) bool {
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			r == '.' || r == '+' || r == '-' || r == '~') {
			return false
		}
	}
	return true
}

func Install(packages ...string) error {
	for _, p := range packages {
		// Validate package name (alphanumeric, dots, plus, minus, and tilde are allowed in package names)
		if !isValidPackageName(p) {
			return fmt.Errorf("invalid package name: %s", p)
		}

		// #nosec G204 -- Package name has been validated with isValidPackageName
		stdout, err := exec.Command("apt", "install", "-y", p).Output()
		if err != nil {
			return err
		}
		fmt.Println(string(stdout))
	}
	return nil
}

func Update() error {
	stdout, err := exec.Command(BinSH, "-c", "apt update -y").Output()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))
	return nil
}

func Upgrade() error {
	stdout, err := exec.Command(BinSH, "-c", "apt upgrade -y").Output()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))
	return nil
}

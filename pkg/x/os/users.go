// Package os provides a set of functions to interact with the operating system.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
)

const (
	userFile string = "/etc/passwd"
)

// ReadEtcPasswd file /etc/passwd and return slice of users
func ReadEtcPasswd() (list []string) {

	file, err := os.Open(userFile)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	r := bufio.NewScanner(file)

	for r.Scan() {
		lines := r.Text()
		parts := strings.Split(lines, ":")
		list = append(list, parts[0])
	}
	return list
}

// check if user on the host
func check(s []string, u string) bool {
	for _, w := range s {
		if u == w {
			return true
		}
	}
	return false
}

// Return securely generated random bytes

func CreateRandom(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	return string(b)
}

// AddUserIfNotExist check if user exist on the host
func AddUserIfNotExist(name string) (string, error) {

	users := ReadEtcPasswd()

	if check(users, name) {
		return "", errors.New("User already exists")
	} else {
		return AddNewUser(name)
	}
}

// AddNewUser is created by executing shell command useradd
func validateUsername(name string) error {
	// Username must be alphanumeric and can contain underscores and dashes
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
			return fmt.Errorf("invalid username: must only contain alphanumeric characters, underscores, or dashes")
		}
	}
	return nil
}

func AddNewUser(name string) (string, error) {
	// Validate username
	if err := validateUsername(name); err != nil {
		return "", err
	}

	// Generate a secure random password
	encrypt := base64.StdEncoding.EncodeToString([]byte(CreateRandom(9)))

	// Use adduser with validated arguments
	userCmd := exec.Command("adduser", "--disabled-password", "--gecos", "", name)

	// Run adduser command
	if out, err := userCmd.Output(); err != nil {
		return "", fmt.Errorf("failed to add user: %w", err)
	} else {
		fmt.Printf("Output: %s\n", out)

		// Set password using chpasswd
		passCmd := exec.Command("chpasswd")
		passCmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s\n", name, encrypt))

		if err := passCmd.Run(); err != nil {
			return "", fmt.Errorf("failed to set password: %w", err)
		}

		return encrypt, nil
	}
}

// DeleteUserIfExist check if user exist on the host
func DeleteUserIfExist(name string) error {
	users := ReadEtcPasswd()

	if check(users, name) {
		return DeleteUser(name)
	} else {
		return errors.New("user doesn't exists")
	}
}

// DeleteUser is created by executing shell command userdel
func DeleteUser(name string) error {
	// Validate username
	if err := validateUsername(name); err != nil {
		return err
	}

	// Use deluser with validated argument
	cmd := exec.Command("deluser", name)

	if out, err := cmd.Output(); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	} else {
		fmt.Printf("Output: %s\n", out)
		return nil
	}
}

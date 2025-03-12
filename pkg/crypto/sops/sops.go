package sops

import (
	"fmt"
	"os"
	"os/exec"
)

// Encrypt takes a string of plaintext data and returns the encrypted data
func Encrypt(plaintext, pubKeyPath string) (string, error) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "sops-encrypt-")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	// Write plaintext to the temporary file
	if _, err := tempFile.WriteString(plaintext); err != nil {
		return "", fmt.Errorf("error writing to temp file: %v", err)
	}
	err = tempFile.Close()
	if err != nil {
		return "", fmt.Errorf("error closing temp file: %v", err)
	}

	// Run sops to encrypt the file
	encrypted, err := execSopsCommand("--pgp", pubKeyPath, "--encrypt", tempFile.Name())
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

// Decrypt takes a string of encrypted data and returns the plaintext
func Decrypt(encrypted, privKeyPath string) (string, error) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "sops-decrypt-")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	// Write encrypted data to the temporary file
	if _, err := tempFile.WriteString(encrypted); err != nil {
		return "", fmt.Errorf("error writing to temp file: %v", err)
	}
	err = tempFile.Close()
	if err != nil {
		return "", fmt.Errorf("error closing temp file: %v", err)
	}

	// Run sops to decrypt the file
	decrypted, err := execSopsCommand("--pgp", privKeyPath, "--decrypt", tempFile.Name())
	if err != nil {
		return "", err
	}

	return decrypted, nil
}

// helper function to execute sops command
func execSopsCommand(args ...string) (string, error) {
	cmd := exec.Command("sops", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("sops command error: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

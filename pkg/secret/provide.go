package secret

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/crypto/sops"
	"os"
	"strings"
)

// Provide returns a secret value for a given key.
// if the key starts with "file://" then the value is read from the file.
// if the key starts with "env." then the value is read from the environment variable.
// if the key starts with "sops." then the value is read from the sops secret using the private key.
// if the key starts with "zkv." then the value is read from github.com/zcubbs/zkv.
// if the key starts with "hcv." the value is read from the hashicorp vault.
// if the key starts with "gcp." then the value is read from the gcp.
// if the key starts with "aws." then the value is read from the aws.
// if the key starts with "azure." then the value is read from the azure.
// if the key starts with "k8s." then the value is read from the k8s.
func Provide(key string, args ...interface{}) (string, error) {
	prefix := strings.Split(key, ".")[0]
	switch prefix {
	case "file":
		return provideFromFile(key)
	case "env":
		return provideFromEnv(key)
	case "sops":
		return provideFromSops(key, args...)
	case "zkv":
		return provideFromZkv(key)
	case "hcv":
		return provideFromHcv(key)
	case "gcp":
		return provideFromGcp(key)
	case "aws":
		return provideFromAws(key)
	case "azure":
		return provideFromAzure(key)
	case "k8s":
		return provideFromK8s(key)
	default:
		return key, nil
	}
}

// ProvideFromFile returns a secret value for a given key.
// if the key starts with "file." then the value is read from the file.
func provideFromFile(_ string) (string, error) {
	return "", fmt.Errorf("get secret from file not implemented")
}

// ProvideFromEnv returns a secret value for a given key.
// if the key starts with "env." then the value is read from the environment variable.
func provideFromEnv(key string) (string, error) {
	v := os.Getenv(strings.ReplaceAll(key, "env.", ""))
	if v == "" {
		return "", fmt.Errorf("failed to get secret from env")
	}

	return v, nil
}

// ProvideFromSops returns a secret value for a given key.
// if the key starts with "sops." then the value is read from the sops.
func provideFromSops(key string, args ...interface{}) (string, error) {
	v, err := sops.Decrypt(key, args[0].(string)) // args[0] is the path to the private key
	if err != nil {
		return "", fmt.Errorf("failed to get secret from sops")
	}

	return v, nil
}

// ProvideFromZkv returns a secret value for a given key.
// if the key starts with "zkv." then the value is read from github.com/zcubbs/zkv.
func provideFromZkv(_ string) (string, error) {
	return "", fmt.Errorf("get secret from zkv not implemented")
}

// ProvideFromHcv returns a secret value for a given key.
// if the key starts with "hcv." the value is read from the hashicorp vault.
func provideFromHcv(_ string) (string, error) {
	return "", fmt.Errorf("get secret from hcv not implemented")
}

// ProvideFromGcp returns a secret value for a given key.
// if the key starts with "gcp." then the value is read from the gcp.
func provideFromGcp(_ string) (string, error) {
	return "", fmt.Errorf("get secret from gcp not implemented")
}

// ProvideFromAws returns a secret value for a given key.
// if the key starts with "aws." then the value is read from the aws.
func provideFromAws(_ string) (string, error) {
	return "", fmt.Errorf("get secret from aws not implemented")
}

// ProvideFromAzure returns a secret value for a given key.
// if the key starts with "azure." then the value is read from the azure.
func provideFromAzure(_ string) (string, error) {
	return "", fmt.Errorf("get secret from azure not implemented")
}

// ProvideFromK8s returns a secret value for a given key.
// if the key starts with "k8s." then the value is read from the k8s.
func provideFromK8s(_ string) (string, error) {
	return "", fmt.Errorf("get secret from k8s not implemented")
}

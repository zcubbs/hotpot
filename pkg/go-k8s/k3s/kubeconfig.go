package k3s

import (
	"fmt"
	"os"
	"strings"
)

func PrintKubeconfig(kubeconfig, serverUrl string) error {
	const defaultUrl = "https://127.0.0.1:6443"
	// read kubeconfig
	kubeconfigContent, err := readKubeconfig(kubeconfig)
	if err != nil {
		return err
	}

	if serverUrl == "" {
		serverUrl = defaultUrl
	}

	// replace server url
	kubeconfigContent = replaceValueInString(kubeconfigContent,
		defaultUrl, serverUrl)

	// print kubeconfig
	fmt.Println(kubeconfigContent)

	return nil
}

func readKubeconfig(kubeconfig string) (string, error) {
	// read kubeconfig
	kubeconfigContent, err := os.ReadFile(kubeconfig)
	if err != nil {
		return "", err
	}
	return string(kubeconfigContent), nil
}

func replaceValueInString(str string, oldValue string, newValue string) string {
	return strings.ReplaceAll(str, oldValue, newValue)
}

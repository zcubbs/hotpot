package kc

import (
	"fmt"
	"os"
)

func getKubeConfig(path string, debug bool) (string, error) {
	const (
		found             = "kubeconfig found in default location %s\n"
		notFound          = "kubeconfig not found in default location %s\n"
		rancherKubeconfig = "/etc/rancher/k3s/k3s.yaml"
	)
	if path != "" {
		return path, nil
	}

	hd, err := getUserHomeDir()
	if err != nil {
		return "", err
	}

	kc := fmt.Sprintf("%s/.kube/config", hd)
	fi, err := os.Stat(kc)
	if err != nil && os.IsNotExist(err) && debug {
		fmt.Printf(notFound, kc)
	}

	if fi != nil {
		fmt.Printf(found, kc)
		return kc, nil
	}

	fi, err = os.Stat(rancherKubeconfig)
	if err != nil && os.IsNotExist(err) && debug {
		fmt.Printf(notFound, rancherKubeconfig)
	}

	if fi != nil {
		fmt.Printf(found, rancherKubeconfig)
		return rancherKubeconfig, nil
	}

	kc = os.Getenv("KUBECONFIG")
	if kc == "" {
		return "", fmt.Errorf("KUBECONFIG variable not set, and no kubeconfig found in default locations")
	}

	fmt.Printf(found, kc)
	return kc, nil
}

func getUserHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir")
	}
	return home, nil
}

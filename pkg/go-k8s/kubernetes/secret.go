package kubernetes

import (
	"context"
	"encoding/base64"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type ContainerRegistrySecret struct {
	Name     string
	Server   string
	Username string
	Password string
	Email    string
}

func CreateContainerRegistrySecret(
	ctx context.Context,
	kubeconfig string,
	secretConfig ContainerRegistrySecret,
	namespaces []string,
	replace bool,
	debug bool) error {
	auth := fmt.Sprintf("%s:%s", secretConfig.Username, secretConfig.Password)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	data := map[string][]byte{
		".dockerconfigjson": []byte(fmt.Sprintf(`{
			"auths": {
				"%s": {
					"username": "%s",
					"password": "%s",
					"email": "%s",
					"auth": "%s"
				}
			}
		}`, secretConfig.Server, secretConfig.Username,
			secretConfig.Password, secretConfig.Email, encodedAuth)),
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretConfig.Name,
		},
		Data: data,
		Type: v1.SecretTypeDockerConfigJson,
	}

	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return err
	}

	for _, namespace := range namespaces {
		_, err := cs.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "already exists") && replace {
				_, err := cs.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
				if err != nil {
					return fmt.Errorf("failed to update secret: %v, namespace: %s", err, namespace)
				}
				return nil
			}
			return fmt.Errorf("failed to create secret: %v", err)
		}
	}

	return nil
}

func CreateGenericSecret(
	ctx context.Context,
	kubeconfig string,
	secretConfig v1.Secret,
	namespaces []string,
	replace bool,
	debug bool) error {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return err
	}

	for _, namespace := range namespaces {
		created, err := cs.CoreV1().Secrets(namespace).Create(ctx, &secretConfig, metav1.CreateOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "already exists") && replace {
				_, err := cs.CoreV1().Secrets(namespace).Update(ctx, &secretConfig, metav1.UpdateOptions{})
				return err
			}
			return fmt.Errorf("failed to create secret: %v", err)
		}

		if debug {
			fmt.Printf("Created secret %s\n", created.String())
		}
	}

	return nil
}

func GetSecret(kubeconfig, namespace, secretName string) (*v1.Secret, error) {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return nil, err
	}
	secret, err := cs.CoreV1().Secrets(namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// GetSecretByName retrieves a Kubernetes Secret by its name.
func GetSecretByName(kubeconfig, namespace, secretName string) (*v1.Secret, error) {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return nil, err
	}
	return cs.CoreV1().
		Secrets(namespace).
		Get(context.TODO(), secretName, metav1.GetOptions{})
}

package traefik

import (
	"context"
	"fmt"
	"github.com/zcubbs/x/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

type DnsProvider string

/* #nosec */
const traefikProviderCredentialsSecretName = "traefik-dns-provider-credentials"

const (
	Cloudflare DnsProvider = "cloudflare"
	OVH        DnsProvider = "ovh"
	Azure      DnsProvider = "azure"
)

const (
	ovhEndpointEnvKey = "OVH_ENDPOINT"
	ovhAppKeyEnvKey   = "OVH_APPLICATION_KEY"
	/* #nosec */
	ovhAppSecretEnvKey   = "OVH_APPLICATION_SECRET"
	ovhConsumerKeyEnvKey = "OVH_CONSUMER_KEY"

	//azureClientIDEnvKey     = "AZURE_CLIENT_ID"
	///* #nosec */
	//azureClientSecretEnvKey = "AZURE_CLIENT_SECRET"
)

func configureDNSChallengeVars(values Values, kubeconfig string, debug bool) error {
	if values.DnsProvider == "" {
		return fmt.Errorf("dns provider is required")
	}

	if values.DnsProvider == string(Cloudflare) {
		return configureCloudflare(values, kubeconfig, debug)
	}

	if values.DnsProvider == string(OVH) {
		return configureOVH(values, kubeconfig, debug)
	}

	if values.DnsProvider == string(Azure) {
		return configureAzure(values, kubeconfig, debug)
	}

	return fmt.Errorf("unknown dns provider: %s", values.DnsProvider)
}

func configureCloudflare(_ Values, _ string, _ bool) error {
	return fmt.Errorf("cloudflare provider not implemented")
}

func configureOVH(values Values, kubeconfig string, debug bool) error {
	// load env vars
	ovhEndpoint := os.Getenv(ovhEndpointEnvKey)
	ovhAppKey := os.Getenv(ovhAppKeyEnvKey)
	ovhAppSecret := os.Getenv(ovhAppSecretEnvKey)
	ovhConsumerKey := os.Getenv(ovhConsumerKeyEnvKey)

	// validate env vars
	if ovhEndpoint == "" {
		return fmt.Errorf("ovh endpoint is required")
	}

	if ovhAppKey == "" {
		return fmt.Errorf("ovh app key is required")
	}

	if ovhAppSecret == "" {
		return fmt.Errorf("ovh app secret is required")
	}

	if ovhConsumerKey == "" {
		return fmt.Errorf("ovh consumer key is required")
	}

	// create namespace
	err := kubernetes.CreateNamespace(
		kubeconfig,
		[]string{traefikNamespace},
	)
	if err != nil {
		return fmt.Errorf("failed to create namespace %s \n %w", traefikNamespace, err)
	}

	// create secret
	err = kubernetes.CreateGenericSecret(
		context.Background(),
		kubeconfig,
		v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: traefikProviderCredentialsSecretName,
			},
			Data: map[string][]byte{
				"OVH_ENDPOINT":           []byte(ovhEndpoint),
				"OVH_APPLICATION_KEY":    []byte(ovhAppKey),
				"OVH_APPLICATION_SECRET": []byte(ovhAppSecret),
				"OVH_CONSUMER_KEY":       []byte(ovhConsumerKey),
				"TZ":                     []byte(values.DnsTZ),
			},
		},
		[]string{traefikNamespace},
		true,
		debug,
	)
	if err != nil {
		return fmt.Errorf("failed to create secret %s \n %w", traefikProviderCredentialsSecretName, err)
	}

	return nil
}

func configureAzure(_ Values, _ string, _ bool) error {
	return fmt.Errorf("azure provider not implemented")
}

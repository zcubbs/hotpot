package traefik

import (
	"context"
	"fmt"
	"github.com/zcubbs/hotpot/pkg/go-k8s/kubernetes"
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

	azureClientIDEnvKey = "AZURE_CLIENT_ID"
	/* #nosec */
	azureClientSecretEnvKey = "AZURE_CLIENT_SECRET"
	azureResourceGroupKey   = "AZURE_RESOURCE_GROUP"
	azureSubscriptionIDKey  = "AZURE_SUBSCRIPTION_ID"
	azureTenantIDKey        = "AZURE_TENANT_ID"
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

func configureAzure(values Values, kubeconfig string, debug bool) error {
	// load env vars
	azureClientId := os.Getenv(azureClientIDEnvKey)
	azureClientSecret := os.Getenv(azureClientSecretEnvKey)
	azureResourceGroup := os.Getenv(azureResourceGroupKey)
	azureSubscriptionID := os.Getenv(azureSubscriptionIDKey)
	azureTenantID := os.Getenv(azureTenantIDKey)

	// validate env vars
	if azureClientId == "" {
		return fmt.Errorf("azure client id is required")
	}

	if azureClientSecret == "" {
		return fmt.Errorf("azure client secret is required")
	}

	if azureResourceGroup == "" {
		return fmt.Errorf("azure resource group is required")
	}

	if azureSubscriptionID == "" {
		return fmt.Errorf("azure subscription id is required")
	}

	if azureTenantID == "" {
		return fmt.Errorf("azure tenant id is required")
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
	return createSecret(
		map[string][]byte{
			"AZURE_CLIENT_ID":       []byte(azureClientId),
			"AZURE_CLIENT_SECRET":   []byte(azureClientSecret),
			"AZURE_RESOURCE_GROUP":  []byte(azureResourceGroup),
			"AZURE_SUBSCRIPTION_ID": []byte(azureSubscriptionID),
			"AZURE_TENANT_ID":       []byte(azureTenantID),
			"TZ":                    []byte(values.DnsTZ),
		},
		kubeconfig,
		debug,
	)
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
	return createSecret(
		map[string][]byte{
			"OVH_ENDPOINT":           []byte(ovhEndpoint),
			"OVH_APPLICATION_KEY":    []byte(ovhAppKey),
			"OVH_APPLICATION_SECRET": []byte(ovhAppSecret),
			"OVH_CONSUMER_KEY":       []byte(ovhConsumerKey),
			"TZ":                     []byte(values.DnsTZ),
		},
		kubeconfig,
		debug,
	)
}

func createSecret(data map[string][]byte, kubeconfig string, debug bool) error {
	// create secret
	err := kubernetes.CreateGenericSecret(
		context.Background(),
		kubeconfig,
		v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: traefikProviderCredentialsSecretName,
			},
			Data: data,
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

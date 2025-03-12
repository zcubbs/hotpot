package rancher

import (
	"context"
	"fmt"
	"github.com/zcubbs/hotpot/pkg/go-k8s/helm"
	"github.com/zcubbs/hotpot/pkg/go-k8s/kubernetes"
	"github.com/zcubbs/hotpot/pkg/x/yaml"
	"os"
	"time"
)

const (
	defaultVersion    = ""
	helmRepoURL       = "https://releases.rancher.com/server-charts/stable"
	helmRepoName      = "rancher-stable"
	defaultNamespace  = "cattle-system"
	defaultChartName  = "rancher"
	defaultValuesFile = "values.yaml"
)

type Values struct {
	Version  string
	Hostname string
}

func Install(values *Values, kubeconfig string, debug bool) error {
	err := validateValues(values)
	if err != nil {
		return err
	}

	// create values file
	valuesFileData, err := yaml.ApplyTmpl(
		valuesTmpl,
		values,
		debug,
	)
	if err != nil {
		return fmt.Errorf("failed to parse values template file: %w", err)
	}

	valuesFilePath := fmt.Sprintf("%s/%s", os.TempDir(), defaultValuesFile)
	// write values file
	err = os.WriteFile(valuesFilePath, valuesFileData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write values file: %w", err)
	}

	helmClient := helm.NewClient()
	helmClient.Settings.KubeConfig = kubeconfig
	helmClient.Settings.SetNamespace(defaultNamespace)
	helmClient.Settings.Debug = debug

	err = helmClient.RepoAddAndUpdate(helmRepoName, helmRepoURL)
	if err != nil {
		return fmt.Errorf("failed to add helm repo: %w", err)
	}

	err = helmClient.InstallChart(helm.Chart{
		ChartName:       defaultChartName,
		ReleaseName:     defaultChartName,
		RepoName:        helmRepoName,
		Values:          nil,
		ValuesFiles:     []string{valuesFilePath},
		Debug:           debug,
		CreateNamespace: true,
		Upgrade:         true,
	})
	if err != nil {
		return fmt.Errorf("failed to install helm chart: %w", err)
	}

	// wait for rancher server to be ready
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	err = kubernetes.IsDeploymentReady(
		ctxWithTimeout,
		kubeconfig,
		defaultNamespace,
		[]string{
			"rancher",
		},
		debug,
	)
	if err != nil {
		return fmt.Errorf("failed to wait for rancher server to be ready \n %w", err)
	}

	return nil
}

func validateValues(values *Values) error {
	if values.Version == "" {
		values.Version = defaultVersion
	}
	if values.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	return nil
}

const valuesTmpl = `
---

replicas: 1
hostname: {{ .Hostname }}
ingress:
  enabled: false
tls: external
`

// Uninstall uninstalls rancher
func Uninstall(kubeconfig string, debug bool) error {
	helmClient := helm.NewClient()
	helmClient.Settings.KubeConfig = kubeconfig
	helmClient.Settings.SetNamespace(defaultNamespace)
	helmClient.Settings.Debug = debug

	err := helmClient.UninstallChart(defaultChartName)
	if err != nil {
		return fmt.Errorf("failed to uninstall helm chart: %w", err)
	}

	return nil
}

package helm

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zcubbs/hotpot/pkg/x/pretty"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/strvals"
	"os"
)

type Chart struct {
	ChartName       string
	ReleaseName     string
	RepoName        string
	Values          map[string]string
	ValuesFiles     []string
	Debug           bool
	CreateNamespace bool
	Upgrade         bool
}

func (c *Client) InstallChart(chartInput Chart) error {
	actionConfig, err := c.initActionConfig()
	if err != nil {
		return err
	}

	cp, err := c.locateChart(chartInput.RepoName, chartInput.ChartName, actionConfig)
	if err != nil {
		return err
	}

	vals, err := c.loadValues(chartInput.Values, chartInput.ValuesFiles)
	if err != nil {
		return err
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return fmt.Errorf("chart is not installable: %w", err)
	}

	if err := c.handleDependencies(chartRequested, cp, actionConfig); err != nil {
		return err
	}

	if chartInput.Upgrade && c.releaseExists(chartInput.ReleaseName, actionConfig) {
		return c.upgradeChart(chartRequested, vals, chartInput, actionConfig)
	} else if !c.releaseExists(chartInput.ReleaseName, actionConfig) {
		return c.installNewChart(chartRequested, vals, chartInput, actionConfig)
	} else {
		return fmt.Errorf("release %s already exists, and upgrade was not specified", chartInput.ReleaseName)
	}
}

func (c *Client) initActionConfig() (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(
		c.Settings.RESTClientGetter(),
		c.Settings.Namespace(),
		os.Getenv("HELM_DRIVER"),
		getDebug(c.Settings.Debug)); err != nil {
		return nil, fmt.Errorf("failed to initialize helm action configuration: %w", err)
	}
	return actionConfig, nil
}

func (c *Client) locateChart(repo, chart string, actionConfig *action.Configuration) (string, error) {
	client := action.NewInstall(actionConfig)
	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", repo, chart), c.Settings)
	if err != nil {
		return "", fmt.Errorf("failed to locate chart: %w", err)
	}
	return cp, nil
}

func (c *Client) loadValues(valuesMap map[string]string, valuesFiles []string) (map[string]interface{}, error) {
	vals := make(map[string]interface{})
	for k, v := range valuesMap {
		setString := fmt.Sprintf("%s=%s", k, v)
		if err := strvals.ParseInto(setString, vals); err != nil {
			return nil, errors.Wrapf(err, "failed setting value for %s", k)
		}
	}

	valueOpts := &values.Options{
		ValueFiles: valuesFiles,
	}
	finalVals, err := valueOpts.MergeValues(getter.All(c.Settings))
	if err != nil {
		return nil, fmt.Errorf("failed to merge values: %w", err)
	}

	return finalVals, nil
}

func (c *Client) handleDependencies(ch *chart.Chart, cp string, actionConfig *action.Configuration) error {
	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			installClient := action.NewInstall(actionConfig)

			man := &downloader.Manager{
				Out:              os.Stdout,
				ChartPath:        cp,
				Keyring:          installClient.ChartPathOptions.Keyring,
				SkipUpdate:       false,
				Getters:          getter.All(c.Settings),
				RepositoryConfig: c.Settings.RepositoryConfig,
				RepositoryCache:  c.Settings.RepositoryCache,
			}
			if err := man.Update(); err != nil {
				return fmt.Errorf("failed to update dependencies: %w", err)
			}
		}
	}
	return nil
}

func (c *Client) installNewChart(ch *chart.Chart, vals map[string]interface{}, chartInput Chart, actionConfig *action.Configuration) error {
	client := action.NewInstall(actionConfig)
	client.ReleaseName = chartInput.ReleaseName
	client.Namespace = c.Settings.Namespace()
	client.CreateNamespace = chartInput.CreateNamespace

	release, err := client.Run(ch, vals)
	if err != nil {
		return fmt.Errorf("failed to install chart: %w", err)
	}

	if chartInput.Debug {
		pretty.PrintJson(release.Manifest)
	}

	return nil
}

func (c *Client) upgradeChart(ch *chart.Chart, vals map[string]interface{}, chartInput Chart, actionConfig *action.Configuration) error {
	upgradeClient := action.NewUpgrade(actionConfig)
	upgradeClient.Namespace = c.Settings.Namespace()

	release, err := upgradeClient.Run(chartInput.ReleaseName, ch, vals)
	if err != nil {
		return fmt.Errorf("failed to upgrade chart: %w", err)
	}

	if chartInput.Debug {
		pretty.PrintJson(release.Manifest)
	}

	return nil
}

// UninstallChart uninstalls a helm chart
func (c *Client) UninstallChart(name string) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(
		c.Settings.RESTClientGetter(),
		c.Settings.Namespace(),
		os.Getenv("HELM_DRIVER"),
		getDebug(c.Settings.Debug)); err != nil {
		return fmt.Errorf("failed to initialize helm action configuration: %w", err)
	}
	client := action.NewUninstall(actionConfig)

	_, err := client.Run(name)
	if err != nil {
		return fmt.Errorf("failed to uninstall chart: %w", err)
	}

	return nil
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

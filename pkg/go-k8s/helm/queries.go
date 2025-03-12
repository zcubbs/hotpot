// Package helm.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	"log"
)

func GetAllReleases(kubeconfig string) ([]*release.Release, error) {
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(kube.GetConfig(kubeconfig, "", ""), "", "", log.Printf)
	if err != nil {
		return nil, err
	}

	_releases, err := action.NewList(actionConfig).Run()
	if err != nil {
		return nil, err
	}

	return _releases, nil
}

func (c *Client) releaseExists(releaseName string, actionConfig *action.Configuration) bool {
	list := action.NewList(actionConfig)
	list.All = true           // List all releases
	list.AllNamespaces = true // Look in all namespaces
	releases, err := list.Run()
	if err != nil {
		return false
	}

	for _, rel := range releases {
		if rel.Name == releaseName {
			return true
		}
	}

	return false
}

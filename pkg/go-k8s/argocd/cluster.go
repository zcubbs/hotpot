package argocd

import (
	"fmt"

	"github.com/zcubbs/hotpot/pkg/go-k8s/kubernetes"
	"github.com/zcubbs/hotpot/pkg/secret"
)

type Cluster struct {
	Name       string `mapstructure:"name" json:"name" yaml:"name"`
	Namespace  string `mapstructure:"namespace" json:"namespace" yaml:"namespace"`
	ServerName string `mapstructure:"serverName" json:"serverName" yaml:"serverName"` // Base64
	ServerUrl  string `mapstructure:"serverUrl" json:"serverUrl" yaml:"serverUrl"`    // Base64
	Config     string `mapstructure:"config" json:"config" yaml:"config"`             // Base64
}

func CreateCluster(cluster Cluster, _ string, debug bool) error {
	if cluster.Namespace == "" {
		cluster.Namespace = argocdNamespace
	}

	config, err := secret.Provide(cluster.Config)
	if err != nil {
		return fmt.Errorf("failed to provide argocd repository config \n %w", err)
	}

	tmpValues := Cluster{
		Name:       cluster.Name,
		Namespace:  cluster.Namespace,
		ServerName: cluster.ServerName,
		ServerUrl:  cluster.ServerUrl,
		Config:     config,
	}

	// Apply template
	err = kubernetes.ApplyManifest(clusterTmpl, tmpValues, debug)
	if err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}
	return nil
}

var clusterTmpl = `---

apiVersion: v1
kind: Secret
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    argocd.argoproj.io/secret-type: cluster
data:
  config: {{ .Config }}
  name: {{ .ServerName }}
  server: {{ .ServerUrl }}
type: Opaque
`

package argocd

import (
	"fmt"

	"github.com/zcubbs/hotpot/pkg/go-k8s/kubernetes"
)

type Project struct {
	Name        string   `mapstructure:"name" json:"name" yaml:"name"`
	Namespace   string   `mapstructure:"namespace" json:"namespace" yaml:"namespace"`
	ClustersUrl []string `mapstructure:"clustersUrl" json:"clustersUrl" yaml:"clustersUrl"`
}

func CreateProject(project Project, kubeconfig string, debug bool) error {
	if project.Namespace == "" {
		project.Namespace = argocdNamespace
	}
	// Apply template
	err := kubernetes.ApplyManifestWithKc(projectTmpl, project, kubeconfig, debug)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

var projectTmpl = `---

apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  description: {{ .Name }}
  sourceRepos:
    - '*'
  clusterResourceWhitelist:
    - group: '*'
      kind: '*'
  destinations:
    - namespace: '*'
      server: https://kubernetes.default.svc
    {{ if .ClustersUrl }}
    {{- range .ClustersUrl }}
    - namespace: '*'
      server: {{ . }}
    {{- end }}
    {{- end }}
`

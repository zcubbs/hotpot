package argocd

import (
	"fmt"
	"github.com/zcubbs/x/kubernetes"
)

type Project struct {
	Name      string `mapstructure:"name" json:"name" yaml:"name"`
	Namespace string `mapstructure:"namespace" json:"namespace" yaml:"namespace"`
}

func CreateProject(project Project, _ string, debug bool) error {
	project.Namespace = argocdNamespace
	// Apply template
	err := kubernetes.ApplyManifest(projectTmpl, project, debug)
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

`

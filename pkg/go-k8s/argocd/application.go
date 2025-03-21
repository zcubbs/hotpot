package argocd

import (
	"fmt"

	"github.com/zcubbs/hotpot/pkg/go-k8s/kubernetes"
	"github.com/zcubbs/hotpot/pkg/x/pretty"
)

type Application struct {
	Name             string   `json:"name"`
	Namespace        string   `json:"namespace"`
	IsOCI            bool     `json:"isOCI"`
	OCIChartName     string   `json:"ociChartName"`
	OCIChartRevision string   `json:"ociChartRevision"`
	OCIRepoURL       string   `json:"ociRepoURL"`
	IsHelm           bool     `json:"isHelm"`
	HelmValueFiles   []string `json:"helmValueFiles"`
	Project          string   `json:"project"`
	Cluster          string   `json:"cluster"`
	RepoURL          string   `json:"repoURL"`
	TargetRevision   string   `json:"targetRevision"`
	Path             string   `json:"path"`
	Recurse          bool     `json:"recurse"`
	CreateNamespace  bool     `json:"createNamespace"`
	Prune            bool     `json:"prune"`
	SelfHeal         bool     `json:"selfHeal"`
	AllowEmpty       bool     `json:"allowEmpty"`

	ArgoNamespace string `json:"argoNamespace"`
}

func CreateApplication(app Application, kubeconfig string, debug bool) error {
	if err := validateApp(&app); err != nil {
		return err
	}

	if debug {
		pretty.PrintJson(app)
	}

	if app.Cluster == "" {
		app.Cluster = "https://kubernetes.default.svc"
	}

	// create app
	if app.IsOCI {
		// Apply template
		err := kubernetes.ApplyManifestWithKc(argoAppOciTmpl, app, kubeconfig, debug)
		if err != nil {
			return fmt.Errorf("failed to create application: %s, %w", app.Name, err)
		}
		return nil
	}

	return kubernetes.ApplyManifestWithKc(argoAppTmpl, app, kubeconfig, debug)
}

func validateApp(app *Application) error {
	if !app.IsHelm && app.IsOCI {
		return fmt.Errorf("oci flag can only be used with helm charts. helm is false")
	}

	if app.IsOCI && app.OCIChartName == "" {
		return fmt.Errorf("oci chart name cannot be empty, when oci is true")
	}

	if (!app.IsOCI && !app.IsHelm) && app.Path == "" {
		return fmt.Errorf("path cannot be empty, when helm is false")
	}

	if app.ArgoNamespace == "" {
		app.ArgoNamespace = argocdNamespace
	}

	if app.Namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}

	return nil
}

var argoAppTmpl = `---

apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{ .Name }}
  namespace: {{ .ArgoNamespace }}
spec:
  project: {{ .Project }}
  source:
    repoURL: {{ .RepoURL }}
    targetRevision: {{ .TargetRevision }}
    path: {{ .Path }}
    {{ if .IsHelm }}
    helm:
      passCredentials: true
      valueFiles:
      {{- range .HelmValueFiles }}
        - {{ . }}
      {{- end }}
    {{ else }}
    directory:
      recurse: {{ .Recurse }}
    {{ end }}
  destination:
    server: {{ .Cluster }}
    namespace: {{ .Namespace }}
  syncPolicy:
    syncOptions:
      - CreateNamespace={{ .CreateNamespace }}
    automated:
      prune: {{ .Prune }}
      selfHeal: {{ .SelfHeal }}
      allowEmpty: {{ .AllowEmpty }}
`

var argoAppOciTmpl = `---

apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{ .Name }}
  namespace: {{ .ArgoNamespace }}
spec:
  project: {{ .Project }}
  sources:
    - repoURL: {{ .OciRepoURL }}
      targetRevision: {{ .OciChartRevision }}
      chart: {{ .OciChartName }}
      helm:
        passCredentials: true
        valueFiles:
        {{- range .HelmValueFiles }}
          - {{ . }}
        {{- end }}
    - repoURL: {{ .RepoURL }}
      targetRevision: {{ .TargetRevision }}
      ref: values
  destination:
    server: {{ .Cluster }}
    namespace: {{ .AppNamespace }}
  syncPolicy:
    syncOptions:
      - CreateNamespace={{ .CreateNamespace }}
    automated:
      prune: {{ .Prune }}
      selfHeal: {{ .SelfHeal }}
      allowEmpty: {{ .AllowEmpty }}
`

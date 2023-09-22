package argocd

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
	RepoURL          string   `json:"repoURL"`
	TargetRevision   string   `json:"targetRevision"`
	Path             string   `json:"path"`
	Recurse          bool     `json:"recurse"`
	CreateNamespace  bool     `json:"createNamespace"`
	Prune            bool     `json:"prune"`
	SelfHeal         bool     `json:"selfHeal"`
	AllowEmpty       bool     `json:"allowEmpty"`
}

func CreateApplication(app Application, kubeconfig string, debug bool) error {

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
    server: https://kubernetes.default.svc
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
  name: {{ .AppName }}
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
    server: https://kubernetes.default.svc
    namespace: {{ .AppNamespace }}
  syncPolicy:
    syncOptions:
      - CreateNamespace={{ .CreateNamespace }}
    automated:
      prune: {{ .Prune }}
      selfHeal: {{ .SelfHeal }}
      allowEmpty: {{ .AllowEmpty }}
`

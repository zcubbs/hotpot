package argocd

import (
	"fmt"
	"strings"

	"github.com/zcubbs/hotpot/pkg/go-k8s/kubernetes"
	"github.com/zcubbs/hotpot/pkg/secret"
)

const Git = "git"
const Helm = "helm"

type Repository struct {
	Name string `json:"name"`
	Url  string `json:"url"`

	Username string `json:"username"`
	Password string `json:"password"`

	Type  string `json:"type"`
	IsOci bool   `json:"isOci"`

	Namespace string `json:"namespace"`
}

func CreateRepository(repo Repository, kubeconfig string, debug bool) error {
	if repo.Namespace == "" {
		repo.Namespace = argocdNamespace
	}
	if repo.Type != Git && repo.Type != Helm {
		return fmt.Errorf("invalid repository type: %s, must be git of helm", repo.Type)
	}

	urlValid := strings.HasPrefix(repo.Url, "http://") || strings.HasPrefix(repo.Url, "https://")
	if !urlValid && repo.Type == Git {
		return fmt.Errorf("error: repository url must be valid url: %s, (http://... or https://...)", repo.Url)
	}

	if repo.Type == Git {
		urlValid = strings.HasSuffix(repo.Url, ".git")
		if !urlValid {
			repo.Url = repo.Url + ".git"
		}
	}

	username, err := secret.Provide(repo.Username)
	if err != nil {
		return fmt.Errorf("failed to provide argocd repository username \n %w", err)
	}

	password, err := secret.Provide(repo.Password)
	if err != nil {
		return fmt.Errorf("failed to provide argocd repository password \n %w", err)
	}

	tmpValues := repoTmplValues{
		Name:      repo.Name,
		Namespace: repo.Namespace,
		Type:      repo.Type,
		IsOci:     repo.IsOci,
		Url:       repo.Url,
		Username:  username,
		Password:  password,
	}

	err = kubernetes.ApplyManifestWithKc(repoTmpl, tmpValues, kubeconfig, debug)
	if err != nil {
		return err
	}

	return nil
}

type repoTmplValues struct {
	Name      string
	Namespace string
	Type      string
	IsOci     bool
	Url       string
	Username  string
	Password  string
}

var repoTmpl = `---

apiVersion: v1
kind: Secret
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    argocd.argoproj.io/secret-type: repository
stringData:
  type: {{ .Type }}
  name: {{ .Name }}
{{- if eq .Type "helm" }}
  enableOCI: "{{ .IsOci }}"
{{- end }}
  url: {{ .Url }}
  username: {{ .Username }}
  password: {{ .Password }}
type: Opaque
`

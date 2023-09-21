package argocd

import (
	"fmt"
	"github.com/zcubbs/x/kubernetes"
	"strings"
)

const Git = "git"
const Helm = "helm"

type Repository struct {
	Name string `json:"name"`
	Url  string `json:"url"`

	Username string `json:"username"`
	Password string `json:"password"`

	Type string `json:"type"`
}

func CreateRepository(repo Repository, _ string, debug bool) error {
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
			return fmt.Errorf("error: url must be valid git url: %s. %s",
				repo.Url,
				"example: https://example.com/example.git",
			)
		}
	}

	tmpValues := repoTmplValues{
		Name:      repo.Name,
		Namespace: argocdNamespace,
		Type:      repo.Type,
		Url:       repo.Url,
		Username:  repo.Username,
		Password:  repo.Password,
	}

	err := kubernetes.ApplyManifest(repoTmpl, tmpValues, debug)
	if err != nil {
		return err
	}

	return nil
}

type repoTmplValues struct {
	Name      string
	Namespace string
	Type      string
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
{{- if eq .Type "helm" }}
  name: {{ .Name }}
  enableOCI: "true"
{{- end }}
  url: {{ .Url }}
  username: {{ .Username }}
  password: {{ .Password }}
type: Opaque
`

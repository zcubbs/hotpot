package recipe

import (
	"github.com/zcubbs/hotpot/pkg/go-k8s/argocd"
	"github.com/zcubbs/hotpot/pkg/go-k8s/certmanager"
	"github.com/zcubbs/hotpot/pkg/go-k8s/k3s"
	"github.com/zcubbs/hotpot/pkg/go-k8s/rancher"
	"github.com/zcubbs/hotpot/pkg/go-k8s/traefik"
)

// SystemInfo provides system-related operations
type SystemInfo interface {
	IsOS(os string) error
	IsArchIn(archs []string) error
	IsRAMEnough(minRAM string) error
	IsCPUEnough(minCPU int) error
	IsDiskSpaceEnough(path, size string) error
	IsCurlOK(urls []string) error
}

// K3sManager handles K3s operations
type K3sManager interface {
	Install(cfg k3s.Config, debug bool) error
	Uninstall(debug bool) error
}

// HelmManager handles Helm operations
type HelmManager interface {
	IsHelmInstalled() (bool, error)
	InstallCli(debug bool) error
}

// CertManager handles cert-manager operations
type CertManager interface {
	Install(values certmanager.Values, kubeconfig string, debug bool) error
	Uninstall(kubeconfig string, debug bool) error
}

// TraefikManager handles Traefik operations
type TraefikManager interface {
	Install(values traefik.Values, kubeconfig string, debug bool) error
	Uninstall(kubeconfig string, debug bool) error
}

// ArgoCDManager handles ArgoCD operations
type ArgoCDManager interface {
	Install(values argocd.Values, kubeconfig string, debug bool) error
	Uninstall(kubeconfig string, debug bool) error
	CreateProject(project argocd.Project, kubeconfig string, debug bool) error
	CreateApplication(app argocd.Application, kubeconfig string, debug bool) error
	CreateRepository(repo argocd.Repository, kubeconfig string, debug bool) error
}

// RancherManager handles Rancher operations
type RancherManager interface {
	Install(values rancher.Values, kubeconfig string, debug bool) error
	Uninstall(kubeconfig string, debug bool) error
}

// K9sManager handles K9s operations
type K9sManager interface {
	Install(debug bool) error
}

// FileSystem handles file system operations
type FileSystem interface {
	RemoveAll(path string) error
}

// Dependencies holds all external dependencies
type Dependencies struct {
	SystemInfo  SystemInfo
	K3s         K3sManager
	Helm        HelmManager
	CertManager CertManager
	Traefik     TraefikManager
	ArgoCD      ArgoCDManager
	Rancher     RancherManager
	K9s         K9sManager
	FileSystem  FileSystem
}

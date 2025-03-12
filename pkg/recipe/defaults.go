package recipe

import (
	"github.com/zcubbs/hotpot/pkg/go-k8s/argocd"
	"github.com/zcubbs/hotpot/pkg/go-k8s/certmanager"
	"github.com/zcubbs/hotpot/pkg/go-k8s/helm"
	"github.com/zcubbs/hotpot/pkg/go-k8s/k3s"
	"github.com/zcubbs/hotpot/pkg/go-k8s/k9s"
	"github.com/zcubbs/hotpot/pkg/go-k8s/rancher"
	"github.com/zcubbs/hotpot/pkg/go-k8s/traefik"
	"github.com/zcubbs/hotpot/pkg/x/host"
	"os"
)

// DefaultDependencies returns a Dependencies struct with default implementations
func DefaultDependencies() Dependencies {
	return Dependencies{
		SystemInfo:  host.DefaultSystemInfo{},
		K3s:         k3s.DefaultManager{},
		Helm:        helm.DefaultManager{},
		CertManager: certmanager.DefaultManager{},
		Traefik:     traefik.DefaultManager{},
		ArgoCD:      argocd.DefaultManager{},
		Rancher:     rancher.DefaultManager{},
		K9s:         k9s.DefaultManager{},
		FileSystem:  defaultFileSystem{},
	}
}

type defaultFileSystem struct{}

func (d defaultFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

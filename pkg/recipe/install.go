package recipe

import (
	"github.com/zcubbs/hotpot/pkg/go-k8s/argocd"
	"github.com/zcubbs/hotpot/pkg/go-k8s/certmanager"
	"github.com/zcubbs/hotpot/pkg/go-k8s/k3s"
	"github.com/zcubbs/hotpot/pkg/go-k8s/rancher"
	"github.com/zcubbs/hotpot/pkg/go-k8s/traefik"
)

func installK3s(r *Recipe, k3sMgr K3sManager, helmMgr HelmManager, fs FileSystem) error {
	if r.K3s.PurgeExisting {
		if err := k3sMgr.Uninstall(r.Debug); err != nil {
			return err
		}
	}

	err := k3sMgr.Install(k3s.Config{
		Disable:                 r.K3s.Disable,
		Version:                 r.K3s.Version,
		TlsSan:                  r.K3s.TlsSan,
		DataDir:                 r.K3s.DataDir,
		DefaultLocalStoragePath: r.K3s.DefaultLocalStoragePath,
		WriteKubeconfigMode:     r.K3s.WriteKubeconfigMode,
		ResolvConfPath:          r.K3s.ResolvConfPath,
		HttpsListenPort:         r.K3s.HttpsListenPort,
	}, r.Debug)
	if err != nil {
		return err
	}

	// Check if helm is installed
	installed, err := helmMgr.IsHelmInstalled()
	if err != nil {
		return err
	}

	// Install helm if not already installed
	if !installed {
		return helmMgr.InstallCli(r.Debug)
	}

	return nil
}

func installK9s(r *Recipe, k9sMgr K9sManager) error {
	return k9sMgr.Install(r.Debug)
}

func installCertManager(r *Recipe, certMgr CertManager) error {
	if r.CertManager.PurgeExisting {
		if err := certMgr.Uninstall(r.Kubeconfig, r.Debug); err != nil {
			return err
		}
	}

	return certMgr.Install(certmanager.Values{
		Version:                     r.CertManager.Version,
		LetsencryptIssuerEnabled:    r.CertManager.LetsencryptIssuerEnabled,
		LetsencryptIssuerEmail:      r.CertManager.LetsencryptIssuerEmail,
		HttpChallengeEnabled:        r.CertManager.HttpChallengeEnabled,
		DnsChallengeEnabled:         r.CertManager.DnsChallengeEnabled,
		DnsProvider:                 r.CertManager.DnsProvider,
		DnsRecursiveNameservers:     r.CertManager.DnsRecursiveNameservers,
		DnsRecursiveNameserversOnly: r.CertManager.DnsRecursiveNameserversOnly,
		DnsAzureClientID:            r.CertManager.DnsAzureClientID,
		DnsAzureClientSecret:        r.CertManager.DnsAzureClientSecret,
		DnsAzureHostedZoneName:      r.CertManager.DnsAzureHostedZoneName,
		DnsAzureResourceGroupName:   r.CertManager.DnsAzureResourceGroupName,
		DnsAzureSubscriptionID:      r.CertManager.DnsAzureSubscriptionID,
		DnsAzureTenantID:            r.CertManager.DnsAzureTenantID,
		DnsOvhEndpoint:              r.CertManager.DnsOvhEndpoint,
		DnsOvhApplicationKey:        r.CertManager.DnsOvhApplicationKey,
		DnsOvhApplicationSecret:     r.CertManager.DnsOvhApplicationSecret,
		DnsOvhConsumerKey:           r.CertManager.DnsOvhConsumerKey,
		DnsOvhZone:                  r.CertManager.DnsOvhZone,
	}, r.Kubeconfig, r.Debug)
}

func installTraefik(r *Recipe, traefikMgr TraefikManager) error {
	return traefikMgr.Install(traefik.Values{
		AdditionalArguments: []string{},
		IngressProvider:     r.Traefik.IngressProvider,
		TlsStrictSNI:        false,
	}, r.Kubeconfig, r.Debug)
}

func installRancher(r *Recipe, rancherMgr RancherManager) error {
	return rancherMgr.Install(rancher.Values{
		Version:  r.Rancher.Version,
		Hostname: r.Rancher.Hostname,
	}, r.Kubeconfig, r.Debug)
}

func installArgocd(r *Recipe, argocdMgr ArgoCDManager) error {
	if r.ArgoCD.Enabled {
		return argocdMgr.Install(argocd.Values{
			Insecure:      r.ArgoCD.Insecure,
			ChartVersion:  r.ArgoCD.ChartVersion,
			AdminPassword: r.ArgoCD.AdminPassword,
		}, r.Kubeconfig, r.Debug)
	}
	return nil
}

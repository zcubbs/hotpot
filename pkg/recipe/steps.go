package recipe

import (
	"github.com/zcubbs/hotpot/pkg/argocd"
	"github.com/zcubbs/hotpot/pkg/traefik"
	"github.com/zcubbs/x/helm"
	"github.com/zcubbs/x/k3s"
)

type step struct {
	f func(*Recipe) error // function
	c bool                // condition
}

func checkPrerequisites(_ *Recipe) error {
	return nil
}

func installK3s(r *Recipe) error {
	k3sCfg := r.Ingredients.K3s
	return k3s.Install(k3s.Config{
		Disable:                 k3sCfg.Disable,
		TlsSan:                  k3sCfg.TlsSan,
		DataDir:                 k3sCfg.DataDir,
		DefaultLocalStoragePath: k3sCfg.DefaultLocalStoragePath,
		WriteKubeconfigMode:     k3sCfg.WriteKubeconfigMode,
	}, r.Debug)
}

func installHelm(cfg *Recipe) error {
	ok, err := helm.IsHelmInstalled()
	if err != nil {
		return err
	}

	if !ok {
		err = helm.InstallCli(cfg.Debug)
		if err != nil {
			return err
		}
	}

	return nil
}

func installCertManager(_ *Recipe) error {
	return nil
}

func installTraefik(r *Recipe) error {
	traefikCfg := r.Ingredients.Traefik
	if traefikCfg.PurgeExisting {
		err := traefik.Uninstall(r.Kubeconfig, r.Debug)
		if err != nil {
			return err
		}
	}

	var ingressProvider string
	if r.Ingredients.CertManager.Enabled {
		ingressProvider = CertResolver
	}

	return traefik.Install(
		traefik.Values{
			AdditionalArguments:                nil,
			IngressProvider:                    ingressProvider,
			DnsProvider:                        traefikCfg.DnsChallengeProvider,
			DnsResolverEmail:                   traefikCfg.DnsChallengeResolverEmail,
			EnableDashboard:                    traefikCfg.EnableDashboard,
			EnableAccessLog:                    traefikCfg.EnableAccessLog,
			DebugLog:                           traefikCfg.Debug,
			EndpointsWeb:                       traefikCfg.EndpointsWeb,
			EndpointsWebsecure:                 traefikCfg.EndpointsWebsecure,
			ServersTransportInsecureSkipVerify: traefikCfg.TransportInsecure,
			ForwardedHeaders:                   traefikCfg.ForwardHeaders,
			ForwardedHeadersInsecure:           traefikCfg.ForwardHeadersInsecure,
			ForwardedHeadersTrustedIPs:         traefikCfg.ForwardHeadersTrustedIPs,
			ProxyProtocol:                      traefikCfg.ProxyProtocol,
			ProxyProtocolInsecure:              traefikCfg.ProxyProtocolInsecure,
			ProxyProtocolTrustedIPs:            traefikCfg.ProxyProtocolTrustedIPs,
			DnsTZ:                              traefikCfg.DnsChallengeTZ,
		},
		r.Kubeconfig,
		r.Debug,
	)
}

func installArgocd(r *Recipe) error {
	if r.Ingredients.ArgoCD.PurgeExisting {
		err := argocd.Uninstall(r.Kubeconfig, r.Debug)
		if err != nil {
			return err
		}
	}

	v := argocd.Values{Insecure: true}
	return argocd.Install(v, r.Kubeconfig, r.Debug)
}

func configureArgocdRepos(r *Recipe) error {
	return nil
}

func configureArgocdProjects(r *Recipe) error {
	return nil
}

func configureArgocdApps(r *Recipe) error {
	return nil
}

func printKubeconfig(r *Recipe) error {
	return nil
}

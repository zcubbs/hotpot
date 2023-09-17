package recipe

import (
	"github.com/zcubbs/hotpot/pkg/traefik"
	"github.com/zcubbs/x/helm"
	"github.com/zcubbs/x/k3s"
)

const (
	// CertResolver is the name of the cert-manager resolver
	CertResolver = "certResolver"
)

type PreHook func() error

// Cook runs recipe
func Cook(cfgPath string, hooks ...PreHook) error {
	// load config
	cfg, err := Load(cfgPath)
	if err != nil {
		return err
	}

	// validate config
	if err := validate(cfg); err != nil {
		return err
	}

	// preheat hooks
	for _, hook := range hooks {
		if err := hook(); err != nil {
			return err
		}
	}

	if err := checkPrerequisites(cfg); err != nil {
		return err
	}

	if err := installK3s(cfg); err != nil {
		return err
	}

	if err := installHelm(cfg); err != nil {
		return err
	}

	if err := installCertManager(cfg); err != nil {
		return err
	}

	if err := installTraefik(cfg); err != nil {
		return err
	}

	if err := installArgocd(cfg); err != nil {
		return err
	}

	if err := configureArgocdRepos(cfg); err != nil {
		return err
	}

	if err := configureArgocdProjects(cfg); err != nil {
		return err
	}

	if err := configureArgocdApps(cfg); err != nil {
		return err
	}

	if err := printKubeconfig(cfg); err != nil {
		return err
	}

	return nil
}

func checkPrerequisites(_ *Recipe) error {
	return nil
}

func installK3s(cfg *Recipe) error {
	k3sCfg := cfg.Ingredients.K3s
	return k3s.Install(k3s.Config{
		Disable:                 k3sCfg.Disable,
		TlsSan:                  k3sCfg.TlsSan,
		DataDir:                 k3sCfg.DataDir,
		DefaultLocalStoragePath: k3sCfg.DefaultLocalStoragePath,
		WriteKubeconfigMode:     k3sCfg.WriteKubeconfigMode,
	}, cfg.Debug)
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

func installTraefik(cfg *Recipe) error {
	var ingressProvider string
	if cfg.Ingredients.CertManager.Enabled {
		ingressProvider = CertResolver
	}
	traefikCfg := cfg.Ingredients.Traefik
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
		cfg.Kubeconfig,
		cfg.Debug,
	)
}

func installArgocd(cfg *Recipe) error {
	return nil
}

func configureArgocdRepos(cfg *Recipe) error {
	return nil
}

func configureArgocdProjects(cfg *Recipe) error {
	return nil
}

func configureArgocdApps(cfg *Recipe) error {
	return nil
}

func printKubeconfig(cfg *Recipe) error {
	return nil
}

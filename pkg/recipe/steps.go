package recipe

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/argocd"
	"github.com/zcubbs/hotpot/pkg/host"
	"github.com/zcubbs/hotpot/pkg/traefik"
	"github.com/zcubbs/x/helm"
	"github.com/zcubbs/x/k3s"
	"strings"
)

type step struct {
	f func(*Recipe) error // function
	c bool                // condition
}

func checkPrerequisites(r *Recipe) error {
	fmt.Printf("🍳 Checking prerequisites... \n")
	// check if os is linux
	for _, v := range r.Ingredients.Node.SupportedOs {
		if err := host.IsOS(v); err != nil {
			return err
		}
	}
	fmt.Printf(" - os: ok\n")

	// check if arch is amd64
	if err := host.IsArchIn(r.Ingredients.Node.SupportedArch); err != nil {
		return err
	}
	fmt.Printf(" - arch: ok\n")

	// check if ram is enough
	if err := host.IsRAMEnough(r.Ingredients.Node.MinMemory); err != nil {
		return err
	}
	fmt.Printf(" - ram: ok\n")

	// check if cpu is enough
	if err := host.IsCPUEnough(r.Ingredients.Node.MinCpu); err != nil {
		return err
	}
	fmt.Printf(" - cpu: ok\n")

	// check if disk is enough, check all disks
	for _, v := range r.Ingredients.Node.MinDiskSize {
		if err := host.IsDiskSpaceEnough(v.Path, v.Size); err != nil {
			return err
		}
	}
	fmt.Printf(" - disk: ok\n")

	// check if curl ok for list of url (curl <url>)
	if err := host.IsCurlOK(r.Ingredients.Node.Curl); err != nil {
		return err
	}
	fmt.Printf(" - curl: ok\n")

	return nil
}

func installK3s(r *Recipe) error {
	fmt.Printf("🍕 Adding k3s... \n")
	k3sCfg := r.Ingredients.K3s
	if k3sCfg.PurgeExisting {
		fmt.Printf("purging existing k3s cluster... \n")
		err := k3s.Uninstall(r.Debug)
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") { // ignore if k3s is not installed
			return err
		}
	}
	disableOpts := ensureTraefikIsDisabled(k3sCfg.Disable)
	if r.Debug {
		fmt.Printf("k3s disable options: %+v\n", disableOpts)
	}
	return k3s.Install(k3s.Config{
		Disable:                 disableOpts,
		TlsSan:                  k3sCfg.TlsSan,
		DataDir:                 k3sCfg.DataDir,
		DefaultLocalStoragePath: k3sCfg.DefaultLocalStoragePath,
		WriteKubeconfigMode:     k3sCfg.WriteKubeconfigMode,
		HttpsListenPort:         k3sCfg.HttpsListenPort,
	}, r.Debug)
}

func ensureTraefikIsDisabled(options []string) []string {
	var found bool
	var updatedOptions []string
	for _, v := range options {
		if v == "traefik" {
			found = true
			break
		}
	}
	if !found {
		updatedOptions = append(options, "traefik")
	}
	return updatedOptions
}

func installHelm(cfg *Recipe) error {
	fmt.Printf("🍉 Adding helm cli... \n")
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
	fmt.Printf("🍙 Adding cert-manager... \n")
	return nil
}

func installTraefik(r *Recipe) error {
	fmt.Printf("🌶️  Adding traefik... \n")
	traefikCfg := r.Ingredients.Traefik
	if traefikCfg.PurgeExisting {
		err := traefik.Uninstall(r.Kubeconfig, r.Debug)
		if err != nil && !strings.Contains(err.Error(), "not found") {
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
	fmt.Printf("🥪 Adding argocd... \n")
	if r.Ingredients.ArgoCD.PurgeExisting {
		err := argocd.Uninstall(r.Kubeconfig, r.Debug)
		if err != nil && !strings.Contains(err.Error(), "not found") {
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

package recipe

import (
	"fmt"
	"github.com/zcubbs/hotpot/pkg/argocd"
	"github.com/zcubbs/hotpot/pkg/host"
	"github.com/zcubbs/hotpot/pkg/traefik"
	"github.com/zcubbs/x/helm"
	"github.com/zcubbs/x/k3s"
	"github.com/zcubbs/x/secret"
	"strings"
)

type step struct {
	f func(*Recipe) error // function
	c bool                // condition
}

func checkPrerequisites(r *Recipe) error {
	fmt.Printf("🍳 Checking prerequisites... \n")
	// check if os is linux
	for _, v := range r.Node.SupportedOs {
		if err := host.IsOS(v); err != nil {
			return err
		}
	}
	fmt.Printf(" - os: ok\n")

	// check if arch is amd64
	if err := host.IsArchIn(r.Node.SupportedArch); err != nil {
		return err
	}
	fmt.Printf(" - arch: ok\n")

	// check if ram is enough
	if err := host.IsRAMEnough(r.Node.MinMemory); err != nil {
		return err
	}
	fmt.Printf(" - ram: ok\n")

	// check if cpu is enough
	if err := host.IsCPUEnough(r.Node.MinCpu); err != nil {
		return err
	}
	fmt.Printf(" - cpu: ok\n")

	// check if disk is enough, check all disks
	for _, v := range r.Node.MinDiskSize {
		if err := host.IsDiskSpaceEnough(v.Path, v.Size); err != nil {
			return err
		}
	}
	fmt.Printf(" - disk: ok\n")

	// check if curl ok for list of url (curl <url>)
	if err := host.IsCurlOK(r.Node.Curl); err != nil {
		return err
	}
	fmt.Printf(" - curl: ok\n")

	return nil
}

func installK3s(r *Recipe) error {
	fmt.Printf("🍕 Adding k3s... \n")
	k3sCfg := r.K3s
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
	err := k3s.Install(k3s.Config{
		Disable:                 disableOpts,
		TlsSan:                  k3sCfg.TlsSan,
		DataDir:                 k3sCfg.DataDir,
		DefaultLocalStoragePath: k3sCfg.DefaultLocalStoragePath,
		WriteKubeconfigMode:     k3sCfg.WriteKubeconfigMode,
		HttpsListenPort:         k3sCfg.HttpsListenPort,
	}, r.Debug)
	if err != nil {
		return err
	}

	return installHelm(r)
}

func ensureTraefikIsDisabled(options []string) []string {
	var found bool
	var updatedOptions []string
	updatedOptions = options
	for _, v := range options {
		if v == "traefik" {
			found = true
			break
		}
	}
	if !found {
		updatedOptions = append(updatedOptions, "traefik")
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
	fmt.Printf("🍔 Adding traefik... \n")
	traefikCfg := r.Traefik
	if traefikCfg.PurgeExisting {
		err := traefik.Uninstall(r.Kubeconfig, r.Debug)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	var ingressProvider string
	if r.CertManager.Enabled {
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
	if r.ArgoCD.PurgeExisting {
		err := argocd.Uninstall(r.Kubeconfig, r.Debug)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	v := argocd.Values{
		Insecure: true,
	}
	err := argocd.Install(v, r.Kubeconfig, r.Debug)
	if err != nil {
		return err
	}

	if r.ArgoCD.AdminPassword != "" {
		password, err := secret.Provide(r.ArgoCD.AdminPassword)
		if err != nil {
			return fmt.Errorf("failed to provide argocd admin password \n %w", err)
		}

		v.AdminPassword = password
		err = argocd.PatchPassword(v, r.Kubeconfig, r.Debug)
		if err != nil {
			return fmt.Errorf("failed to patch argocd admin password \n %w", err)
		}

		fmt.Printf(" - argocd admin password: ok\n")
	}

	return nil
}

func configureGitopsRepos(r *Recipe, repos []ArgocdRepository) error {
	for _, ar := range repos {
		repo := argocd.Repository{
			Name:     ar.Name,
			Url:      ar.Url,
			Username: ar.Credentials.Username,
			Password: ar.Credentials.Password,
			Type:     string(ar.Type),
		}
		err := argocd.CreateRepository(repo, r.Kubeconfig, r.Debug)
		if err != nil {
			return err
		}
		fmt.Printf("    ├─ repository: %s ok\n", ar.Name)
	}
	return nil
}

func configureGitopsProjects(r *Recipe) error {
	fmt.Printf("🌭 Adding gitops... \n")
	for _, project := range r.Gitops.Projects {
		p := argocd.Project{
			Name: project.Name,
		}
		err := argocd.CreateProject(p, r.Kubeconfig, r.Debug)
		if err != nil {
			return err
		}
		fmt.Printf(" - project: %s ok\n", project.Name)

		if err := configureGitopsRepos(r, project.Repositories); err != nil {
			return err
		}

		if err := configureGitopsApps(r, project.Name, project.Apps); err != nil {
			return err
		}
	}
	return nil
}

func configureGitopsApps(r *Recipe, project string, apps []App) error {
	for _, app := range apps {
		a := argocd.Application{
			Name:             app.Name,
			Namespace:        app.Namespace,
			IsOCI:            app.IsOci,
			OCIChartName:     app.OciChartName,
			OCIChartRevision: app.OCIChartRevision,
			OCIRepoURL:       app.Repo,
			IsHelm:           app.IsHelm,
			HelmValueFiles:   app.ValuesFiles,
			Project:          project,
			RepoURL:          app.Repo,
			TargetRevision:   app.Revision,
			Path:             app.Path,
			Recurse:          app.Recurse,
			CreateNamespace:  app.CreateNamespace,
			Prune:            app.Prune,
			SelfHeal:         app.SelfHeal,
			AllowEmpty:       app.AllowEmpty,
		}
		err := argocd.CreateApplication(a, r.Kubeconfig, r.Debug)
		if err != nil {
			return err
		}
		fmt.Printf("    ├─ application: %s ok\n", app.Name)
	}
	return nil
}

func createSecrets(r *Recipe) error {
	fmt.Printf("🌶️  Adding secrets... \n")
	return nil
}

func printKubeconfig(r *Recipe) error {
	return nil
}

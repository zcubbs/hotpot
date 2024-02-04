package recipe

import (
	"context"
	"fmt"
	"github.com/zcubbs/go-k8s/argocd"
	"github.com/zcubbs/go-k8s/certmanager"
	"github.com/zcubbs/go-k8s/helm"
	"github.com/zcubbs/go-k8s/k3s"
	"github.com/zcubbs/go-k8s/kubernetes"
	"github.com/zcubbs/go-k8s/rancher"
	"github.com/zcubbs/go-k8s/traefik"
	"github.com/zcubbs/secret"
	"github.com/zcubbs/x/host"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
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
	fmt.Printf("    ├─ os: ok\n")

	// check if arch is amd64
	if err := host.IsArchIn(r.Node.SupportedArch); err != nil {
		return err
	}
	fmt.Printf("    ├─ arch: ok\n")

	// check if ram is enough
	if err := host.IsRAMEnough(r.Node.MinMemory); err != nil {
		return err
	}
	fmt.Printf("    ├─ ram: ok\n")

	// check if cpu is enough
	if err := host.IsCPUEnough(r.Node.MinCpu); err != nil {
		return err
	}
	fmt.Printf("    ├─ cpu: ok\n")

	// check if disk is enough, check all disks
	for _, v := range r.Node.MinDiskSize {
		if err := host.IsDiskSpaceEnough(v.Path, v.Size); err != nil {
			return err
		}
	}
	fmt.Printf("    ├─ disk: ok\n")

	// check if curl ok for list of url (curl <url>)
	if err := host.IsCurlOK(r.Node.Curl); err != nil {
		return err
	}
	fmt.Printf("    ├─ curl: ok\n")

	fmt.Printf("    └─ prerequisites ok\n")

	return nil
}

func installK3s(r *Recipe) error {
	fmt.Printf("🍕 Adding k3s... \n")
	k3sCfg := r.K3s
	if k3sCfg.PurgeExisting {
		fmt.Printf("    ├─ uninstalling k3s... \n")
		err := k3s.Uninstall(r.Debug)
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") { // ignore if k3s is not installed
			return err
		}

		// purge extra dirs
		for _, v := range k3sCfg.PurgeExtraDirs {
			fmt.Printf("    ├─ purging extra dir %s... \n", v)
			err := purgeDir(v)
			if err != nil {
				return err
			}
		}

		fmt.Printf("    └─ uninstall ok\n")
	}
	disableOpts := ensureTraefikIsDisabled(k3sCfg.Disable)
	if r.Debug {
		fmt.Printf("k3s disable options: %+v\n", disableOpts)
	}
	err := k3s.Install(k3s.Config{
		Disable:                 disableOpts,
		Version:                 k3sCfg.Version,
		TlsSan:                  k3sCfg.TlsSan,
		ResolvConfPath:          k3sCfg.ResolvConfPath,
		DataDir:                 k3sCfg.DataDir,
		DefaultLocalStoragePath: k3sCfg.DefaultLocalStoragePath,
		WriteKubeconfigMode:     k3sCfg.WriteKubeconfigMode,
		HttpsListenPort:         k3sCfg.HttpsListenPort,
	}, r.Debug)
	if err != nil {
		return err
	}

	fmt.Printf("    └─ install ok \n")

	return installHelm(r)
}

func purgeDir(dir string) error {
	// delete dir
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Printf("failed to delete dir %s: %s\n", dir, err)
	}

	return nil
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

func installCertManager(r *Recipe) error {
	fmt.Printf("🍙 Adding cert-manager... \n")
	certmanagerCfg := r.CertManager
	if certmanagerCfg.PurgeExisting {
		err := certmanager.Uninstall(r.Kubeconfig, r.Debug)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	err := certmanager.Install(
		certmanager.Values{
			Version:                         certmanagerCfg.Version,
			LetsencryptIssuerEnabled:        certmanagerCfg.LetsencryptIssuerEnabled,
			LetsencryptIssuerEmail:          certmanagerCfg.LetsencryptIssuerEmail,
			LetsEncryptIngressClassResolver: "cert-manager",
			HttpChallengeEnabled:            certmanagerCfg.HttpChallengeEnabled,
			DnsChallengeEnabled:             certmanagerCfg.DnsChallengeEnabled,
			DnsProvider:                     certmanagerCfg.DnsProvider,
			DnsRecursiveNameservers:         certmanagerCfg.DnsRecursiveNameservers,
			DnsRecursiveNameserversOnly:     certmanagerCfg.DnsRecursiveNameserversOnly,
			DnsAzureClientID:                certmanagerCfg.DnsAzureClientID,
			DnsAzureClientSecret:            certmanagerCfg.DnsAzureClientSecret,
			DnsAzureHostedZoneName:          certmanagerCfg.DnsAzureHostedZoneName,
			DnsAzureResourceGroupName:       certmanagerCfg.DnsAzureResourceGroupName,
			DnsAzureSubscriptionID:          certmanagerCfg.DnsAzureSubscriptionID,
			DnsAzureTenantID:                certmanagerCfg.DnsAzureTenantID,
			DnsOvhEndpoint:                  certmanagerCfg.DnsOvhEndpoint,
			DnsOvhApplicationKey:            certmanagerCfg.DnsOvhApplicationKey,
			DnsOvhApplicationSecret:         certmanagerCfg.DnsOvhApplicationSecret,
			DnsOvhConsumerKey:               certmanagerCfg.DnsOvhConsumerKey,
			DnsOvhZone:                      certmanagerCfg.DnsOvhZone,
		},
		r.Kubeconfig,
		r.Debug,
	)
	if err != nil {
		return fmt.Errorf("failed to install cert-manager \n %w", err)
	}

	fmt.Printf("    └─ install ok\n")
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

	if r.CertManager.Enabled && r.Traefik.IngressProvider == "" {
		fmt.Println("warn: cert-manager is enabled but traefik ingress provider is not set")
	}

	err := traefik.Install(
		traefik.Values{
			AdditionalArguments:                nil,
			IngressProvider:                    traefikCfg.IngressProvider,
			TlsChallengeEnabled:                traefikCfg.TlsChallenge,
			TlsResolver:                        traefikCfg.TlsChallengeResolver,
			TlsResolverEmail:                   traefikCfg.TlsChallengeResolverEmail,
			DnsChallengeEnabled:                traefikCfg.DnsChallenge,
			DnsProvider:                        traefikCfg.DnsChallengeProvider,
			DnsResolverEmail:                   traefikCfg.DnsChallengeResolverEmail,
			DnsResolverIPs:                     traefikCfg.DnsChallengeResolverIPs,
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
			DefaultCertificateEnabled:          traefikCfg.DefaultCertificateEnabled,
			DefaultCertificateCert:             traefikCfg.DefaultCertificateCert,
			DefaultCertificateKey:              traefikCfg.DefaultCertificateKey,
		},
		r.Kubeconfig,
		r.Debug,
	)
	if err != nil {
		return fmt.Errorf("failed to install traefik \n %w", err)
	}

	fmt.Printf("    └─ install ok\n")

	return nil
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

		fmt.Printf("    ├─ argocd admin password: ok\n")
	}

	fmt.Printf("    └─ install ok\n")

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
			IsOci:    ar.IsOci,
		}
		err := argocd.CreateRepository(repo, r.Kubeconfig, r.Debug)
		if err != nil {
			return fmt.Errorf("failed to create argocd repository %s \n %w", ar.Name, err)
		}
		fmt.Printf("    │  ├─ repository: %s ok\n", ar.Name)
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
		fmt.Printf("    ├─ project: %s ok\n", project.Name)

		if err := configureGitopsRepos(r, project.Repositories); err != nil {
			return err
		}

		if err := configureGitopsApps(r, project.Name, project.Apps); err != nil {
			return err
		}
	}

	fmt.Printf("    └─ gitops ok\n")

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
		fmt.Printf("    │  ├─ application: %s ok\n", app.Name)
	}
	return nil
}

func createSecrets(r *Recipe) error {
	fmt.Printf("🌶️  Adding secrets... \n")

	if err := createContainerRegistrySecrets(r.Secrets.ContainerRegistries, r.Kubeconfig, r.Debug); err != nil {
		return err
	}

	if err := createGenericSecrets(r.Secrets.GenericSecrets, r.Kubeconfig, r.Debug); err != nil {
		return err
	}

	if err := createGenericKeyValueSecrets(r.Secrets.GenericKeyValueSecrets, r.Kubeconfig, r.Debug); err != nil {
		return err
	}

	fmt.Printf("    └─ secrets ok\n")

	return nil
}

func createContainerRegistrySecrets(secrets []ContainerRegistryCredentials, kubeconfig string, debug bool) error {
	for _, crs := range secrets {
		fmt.Printf("    ├─ container registry credentials: %s \n", crs.Name)

		err := kubernetes.CreateNamespace(kubeconfig, crs.Namespaces)
		if err != nil {
			return fmt.Errorf("failed to create namespace: %s %w", crs.Namespaces, err)
		}

		fmt.Printf("    │  ├─ namespaces: %s ok\n", crs.Namespaces)

		username, err := secret.Provide(crs.Username)
		if err != nil {
			return fmt.Errorf("failed to provide container registry username: %w", err)
		}

		password, err := secret.Provide(crs.Password)
		if err != nil {
			return fmt.Errorf("failed to provide container registry password: %w", err)
		}

		err = kubernetes.CreateContainerRegistrySecret(
			context.Background(),
			kubeconfig,
			kubernetes.ContainerRegistrySecret{
				Name:     crs.Name,
				Server:   crs.Url,
				Username: username,
				Password: password,
			},
			crs.Namespaces,
			true,
			debug,
		)
		if err != nil {
			return fmt.Errorf("failed to create container registry secret: %s %w", crs.Name, err)
		}

		fmt.Printf("    │  └─ secret ok\n")
	}

	return nil
}

func createGenericSecrets(secrets []GenericSecret, kubeconfig string, debug bool) error {
	for _, s := range secrets {
		var data = make(map[string][]byte)
		for k, v := range s.Data {
			value, err := secret.Provide(v)
			if err != nil {
				return fmt.Errorf("failed to provide secret %s: %w", k, err)
			}
			data[k] = []byte(value)
		}

		err := createSecret(
			"generic secret",
			s.Name,
			s.Namespace,
			kubeconfig,
			debug,
			data,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func createGenericKeyValueSecrets(secrets []GenericKeyValueSecret, kubeconfig string, debug bool) error {
	for _, s := range secrets {
		var data = make(map[string][]byte)
		for _, d := range s.Data {
			value, err := secret.Provide(d.Value)
			if err != nil {
				return fmt.Errorf("failed to provide secret %s: %w", d.Value, err)
			}
			data[d.Key] = []byte(value)
		}

		err := createSecret(
			"generic (key-value) secret",
			s.Name,
			s.Namespace,
			kubeconfig,
			debug,
			data,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func createSecret(
	secretType, secretName, namespace, kubeconfig string,
	debug bool,
	data map[string][]byte) error {
	fmt.Printf("    ├─ %s: %s \n", secretType, secretName)

	err := kubernetes.CreateNamespace(kubeconfig, []string{namespace})
	if err != nil {
		return fmt.Errorf("failed to create namespace: %s %w", namespace, err)
	}

	fmt.Printf("    │  ├─ namespaces: %s ok\n", namespace)

	err = kubernetes.CreateGenericSecret(
		context.Background(),
		kubeconfig,
		v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespace,
				Annotations: map[string]string{
					"createdBy": "hotpot",
				},
			},
			Data: data,
		},
		[]string{namespace},
		true,
		debug,
	)
	if err != nil {
		return fmt.Errorf("failed to create %s: %s %w", secretType, secretName, err)
	}

	fmt.Printf("    │  └─ secret ok\n")

	return nil
}

func installK9s(r *Recipe) error {
	fmt.Printf("🍣 Adding k9s... \n")
	err := k3s.InstallK9s(r.Debug)
	if err != nil {
		return err
	}
	fmt.Printf(" └─ install ok\n")
	return nil
}

func installRancher(r *Recipe) error {
	fmt.Printf("🍍 Adding rancher... \n")
	values := &rancher.Values{
		Version:  r.Rancher.Version,
		Hostname: r.Rancher.Hostname,
	}
	err := rancher.Install(values, r.Kubeconfig, r.Debug)
	if err != nil {
		return err
	}
	fmt.Printf(" └─ install ok\n")
	return nil
}

func printKubeconfig(r *Recipe) error {
	fmt.Printf("🍹 Printing kubeconfig... \n")
	err := k3s.PrintKubeconfig(r.Kubeconfig, r.K3s.KubeApiAddress)
	if err != nil {
		return fmt.Errorf("failed to print kubeconfig \n %w", err)
	}
	return nil
}

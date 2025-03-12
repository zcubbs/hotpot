package recipe

import (
	"context"
	"fmt"

	"github.com/zcubbs/hotpot/pkg/go-k8s/argocd"
	"github.com/zcubbs/hotpot/pkg/go-k8s/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type step struct {
	f func(*Recipe) error // function
	c bool                // condition
}

func checkPrerequisites(r *Recipe, sysInfo SystemInfo) error {
	fmt.Printf("🍳 Checking prerequisites... \n")
	// check if os is linux
	for _, v := range r.Node.SupportedOs {
		if err := sysInfo.IsOS(v); err != nil {
			return err
		}
	}
	fmt.Printf("    ├─ os: ok\n")

	// check if arch is amd64
	if err := sysInfo.IsArchIn(r.Node.SupportedArch); err != nil {
		return err
	}
	fmt.Printf("    ├─ arch: ok\n")

	// check if ram is enough
	if err := sysInfo.IsRAMEnough(r.Node.MinMemory); err != nil {
		return err
	}
	fmt.Printf("    ├─ ram: ok\n")

	// check if cpu is enough
	if err := sysInfo.IsCPUEnough(r.Node.MinCpu); err != nil {
		return err
	}
	fmt.Printf("    ├─ cpu: ok\n")

	// check if disk is enough, check all disks
	for _, v := range r.Node.MinDiskSize {
		if err := sysInfo.IsDiskSpaceEnough(v.Path, v.Size); err != nil {
			return err
		}
	}
	fmt.Printf("    ├─ disk: ok\n")

	// check if curl ok for list of url (curl <url>)
	if err := sysInfo.IsCurlOK(r.Node.Curl); err != nil {
		return err
	}
	fmt.Printf("    ├─ curl: ok\n")

	fmt.Printf("    └─ prerequisites ok\n")

	return nil
}

func configureGitopsRepos(r *Recipe, namespace string, repos []ArgocdRepository) error {
	fmt.Printf("🍲 Configuring gitops repos... \n")
	for _, repo := range repos {
		err := r.Dependencies.ArgoCD.CreateRepository(argocd.Repository{
			Name:      repo.Name,
			Url:       repo.Url,
			Type:      string(repo.Type),
			Username:  repo.Credentials.Username,
			Password:  repo.Credentials.Password,
			Namespace: namespace,
			IsOci:     repo.IsOci,
		}, r.Kubeconfig, r.Debug)
		if err != nil {
			return err
		}
	}
	return nil
}

func configureGitopsProjects(r *Recipe) error {
	fmt.Printf("🍱 Configuring gitops projects... \n")
	for _, project := range r.Gitops.Projects {
		err := r.Dependencies.ArgoCD.CreateProject(argocd.Project{
			Name:        project.Name,
			Namespace:   project.Namespace,
			ClustersUrl: project.ClustersUrl,
		}, r.Kubeconfig, r.Debug)
		if err != nil {
			return err
		}

		if err = configureGitopsRepos(r, project.Namespace, project.Repositories); err != nil {
			return err
		}

		if err = configureGitopsApps(r, project.Name, project.Namespace, project.Apps); err != nil {
			return err
		}
	}
	return nil
}

func configureGitopsApps(r *Recipe, project string, namespace string, apps []App) error {
	fmt.Printf("🍛 Configuring gitops apps... \n")
	for _, app := range apps {
		err := r.Dependencies.ArgoCD.CreateApplication(argocd.Application{
			Name:            app.Name,
			Namespace:       app.Namespace,
			Project:         project,
			Path:            app.Path,
			IsHelm:          app.IsHelm,
			IsOCI:           app.IsOci,
			OCIChartName:    app.OciChartName,
			Cluster:         app.Cluster,
			Recurse:         app.Recurse,
			CreateNamespace: app.CreateNamespace,
			Prune:           app.Prune,
			SelfHeal:        app.SelfHeal,
			AllowEmpty:      app.AllowEmpty,
			ArgoNamespace:   namespace,
		}, r.Kubeconfig, r.Debug)
		if err != nil {
			return err
		}
	}
	return nil
}

func createSecrets(r *Recipe) error {
	fmt.Printf("🍝 Creating secrets... \n")
	if r.Secrets.Enabled {
		if err := createContainerRegistrySecrets(r.Secrets.ContainerRegistries, r.Kubeconfig, r.Debug); err != nil {
			return err
		}
		if err := createGenericSecrets(r.Secrets.GenericSecrets, r.Kubeconfig, r.Debug); err != nil {
			return err
		}
		if err := createGenericKeyValueSecrets(r.Secrets.GenericKeyValueSecrets, r.Kubeconfig, r.Debug); err != nil {
			return err
		}
	}
	return nil
}

func createContainerRegistrySecrets(secrets []ContainerRegistryCredentials, kubeconfig string, debug bool) error {
	fmt.Printf("🍜 Creating container registry secrets... \n")
	for _, secret := range secrets {
		// create secret
		for _, namespace := range secret.Namespaces {
			err := kubernetes.CreateContainerRegistrySecret(
				context.Background(),
				kubeconfig,
				kubernetes.ContainerRegistrySecret{
					Name:     secret.Name,
					Server:   secret.Url,
					Username: secret.Username,
					Password: secret.Password,
					Email:    "",
				},
				[]string{namespace},
				true,
				debug,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createGenericSecrets(secrets []GenericSecret, kubeconfig string, debug bool) error {
	fmt.Printf("🍡 Creating generic secrets... \n")
	for _, secret := range secrets {
		data := make(map[string][]byte)
		for k, v := range secret.Data {
			data[k] = []byte(v)
		}
		err := kubernetes.CreateGenericSecret(
			context.Background(),
			kubeconfig,
			v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secret.Name,
					Namespace: secret.Namespace,
				},
				Type: v1.SecretType(secret.Type),
				Data: data,
			},
			[]string{secret.Namespace},
			true,
			debug,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func createGenericKeyValueSecrets(secrets []GenericKeyValueSecret, kubeconfig string, debug bool) error {
	fmt.Printf("🍢 Creating generic key value secrets... \n")
	for _, secret := range secrets {
		data := make(map[string][]byte)
		for _, v := range secret.Data {
			data[v.Key] = []byte(v.Value)
		}
		err := kubernetes.CreateGenericSecret(
			context.Background(),
			kubeconfig,
			v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secret.Name,
					Namespace: secret.Namespace,
				},
				Type: v1.SecretType(secret.Type),
				Data: data,
			},
			[]string{secret.Namespace},
			true,
			debug,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func printKubeconfig(r *Recipe) error {
	fmt.Printf("🍳 Kubeconfig: %s\n", r.Kubeconfig)
	return nil
}

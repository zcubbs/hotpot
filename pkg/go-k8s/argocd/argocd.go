package argocd

import (
	"context"
	"fmt"
	"github.com/zcubbs/hotpot/pkg/go-k8s/helm"
	"github.com/zcubbs/hotpot/pkg/go-k8s/kubernetes"
	"golang.org/x/crypto/bcrypt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	argocdChartName    = "argo-cd"
	argocdHelmRepoName = "argocd"
	argocdHelmRepoURL  = "https://argoproj.github.io/argo-helm"
	argocdChartVersion = "" // latest
	argocdNamespace    = "argocd"
)

const (
	argocdServerDeploymentName                   = "argo-cd-argocd-server"
	argocdRepoServerDeploymentName               = "argo-cd-argocd-repo-server"
	argocdRedisDeploymentName                    = "argo-cd-argocd-redis"
	argocdDexServerDeploymentName                = "argo-cd-argocd-dex-server"
	argocdApplicationsetControllerDeploymentName = "argo-cd-argocd-applicationset-controller"
	argocdNotificationsControllerDeploymentName  = "argo-cd-argocd-notifications-controller"
)

func Install(values Values, kubeconfig string, debug bool) error {
	if err := validateValues(values); err != nil {
		return err
	}

	// install argocd
	helmClient := helm.NewClient()
	helmClient.Settings.KubeConfig = kubeconfig
	helmClient.Settings.SetNamespace(argocdNamespace)
	helmClient.Settings.Debug = debug

	// add argocd helm repo
	err := helmClient.RepoAddAndUpdate(argocdHelmRepoName, argocdHelmRepoURL)
	if err != nil {
		return fmt.Errorf("failed to add helm repo: %w", err)
	}

	// install argocd
	err = helmClient.InstallChart(helm.Chart{
		ChartName:       argocdChartName,
		ReleaseName:     argocdChartName,
		RepoName:        argocdHelmRepoName,
		Values:          nil,
		ValuesFiles:     nil,
		Debug:           debug,
		CreateNamespace: true,
		Upgrade:         true,
	})
	if err != nil {
		return fmt.Errorf("failed to install argocd \n %w", err)
	}

	// wait for argocd server to be ready
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	err = kubernetes.IsDeploymentReady(
		ctxWithTimeout,
		kubeconfig,
		argocdNamespace,
		[]string{
			argocdServerDeploymentName,
			argocdRepoServerDeploymentName,
			argocdRedisDeploymentName,
			argocdDexServerDeploymentName,
			argocdApplicationsetControllerDeploymentName,
			argocdNotificationsControllerDeploymentName,
		},
		debug,
	)
	if err != nil {
		return fmt.Errorf("failed to wait for argocd server to be ready \n %w", err)
	}

	// patch argocd-cmd-params-cm configmap to set insecure flag
	err = patchConfigMap(kubeconfig, argocdNamespace, "argocd-cmd-params-cm", map[string]string{
		"server.insecure": fmt.Sprintf("%t", values.Insecure),
	}, debug)
	if err != nil {
		return fmt.Errorf("failed to patch argocd-cmd-params-cm: %w", err)
	}

	// restart argocd server pod
	err = kubernetes.RestartPods(kubeconfig, argocdNamespace, []string{
		argocdServerDeploymentName,
	}, debug)
	if err != nil {
		return fmt.Errorf("failed to restart argocd server pod: %w", err)
	}

	return nil
}

func Uninstall(kubeconfig string, debug bool) error {
	helmClient := helm.NewClient()
	helmClient.Settings.KubeConfig = kubeconfig
	helmClient.Settings.SetNamespace(argocdNamespace)
	helmClient.Settings.Debug = debug

	// uninstall argocd
	return helmClient.UninstallChart(argocdChartName)
}

type Values struct {
	Insecure      bool
	ChartVersion  string
	AdminPassword string
}

const patchPasswordAnnotation = "patched-password"

func PatchPassword(values Values, kubeconfig string, debug bool) error {
	secret, err := kubernetes.GetSecret(kubeconfig, argocdNamespace, "argocd-secret")
	if err != nil {
		return fmt.Errorf("failed to get argocd-secret: %w", err)
	}

	if _, ok := secret.Annotations[patchPasswordAnnotation]; ok {
		if debug {
			fmt.Println("argocd-secret already patched")
		}

		return nil
	}

	hashedPassword, err := hashPassword(values.AdminPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	err = kubernetes.CreateGenericSecret(
		context.Background(),
		kubeconfig,
		v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "argocd-secret",
				Namespace: argocdNamespace,
				Annotations: map[string]string{
					patchPasswordAnnotation: "true",
				},
			},
			StringData: map[string]string{
				"admin.password":      hashedPassword,
				"admin.passwordMtime": "'$(date +%FT%T%Z)'",
			},
		},
		[]string{argocdNamespace},
		true,
		debug,
	)
	if err != nil {
		return fmt.Errorf("failed to create argocd-secret: %w", err)
	}

	err = kubernetes.RestartPods(kubeconfig, argocdNamespace,
		[]string{
			argocdServerDeploymentName,
			argocdDexServerDeploymentName,
			argocdRepoServerDeploymentName,
			argocdRedisDeploymentName,
			argocdApplicationsetControllerDeploymentName,
			argocdNotificationsControllerDeploymentName,
		},
		debug)
	if err != nil {
		return fmt.Errorf("failed to restart argocd server pod: %w", err)
	}

	// wait for argocd server to be ready
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	err = kubernetes.IsDeploymentReady(
		ctxWithTimeout,
		kubeconfig,
		argocdNamespace,
		[]string{
			argocdServerDeploymentName,
		},
		debug,
	)
	if err != nil {
		return fmt.Errorf("failed to wait for argocd server to be ready \n %w", err)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func validateValues(values Values) error {
	if values.ChartVersion == "" {
		values.ChartVersion = argocdChartVersion
	}
	return nil
}

func patchConfigMap(kubeconfig string, namespace string, name string, patch map[string]string, debug bool) error {
	cm, err := kubernetes.GetConfigMap(kubeconfig, namespace, name)
	if err != nil {
		return fmt.Errorf("failed to get configmap %s: %w", name, err)
	}

	for k, v := range patch {
		cm.Data[k] = v
	}

	err = kubernetes.UpdateConfigMap(kubeconfig, cm, namespace)
	if err != nil {
		return fmt.Errorf("failed to update configmap %s: %w", name, err)
	}

	return nil
}

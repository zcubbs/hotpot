package argocd

import (
	"context"
	"fmt"
	"github.com/zcubbs/x/helm"
	"github.com/zcubbs/x/kubernetes"
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

	vals := map[string]interface{}{
		"configs.params.server\\.insecure": values.Insecure,
	}

	err := helm.Install(helm.Chart{
		Name:            argocdChartName,
		Repo:            argocdHelmRepoName,
		URL:             argocdHelmRepoURL,
		Version:         values.ChartVersion,
		Values:          vals,
		ValuesFiles:     nil,
		Namespace:       argocdNamespace,
		Upgrade:         true,
		CreateNamespace: true,
	}, kubeconfig, debug)
	if err != nil {
		return err
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

	return nil
}

func Uninstall(kubeconfig string, debug bool) error {
	return helm.Uninstall(helm.Chart{
		Name:      argocdChartName,
		Namespace: argocdNamespace,
	}, kubeconfig, debug)
}

type Values struct {
	Insecure      bool
	ChartVersion  string
	AdminPassword string
}

func PatchPassword(values Values, kubeconfig string, debug bool) error {
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

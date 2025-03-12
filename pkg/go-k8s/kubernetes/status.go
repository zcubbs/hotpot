package kubernetes

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

// IsDeploymentReady checks if a kubernetes deployment is ready
func IsDeploymentReady(ctx context.Context,
	kubeconfig, namespace string,
	deploymentNames []string,
	debug bool) error {

	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return err
	}
	for _, deploymentName := range deploymentNames {
		err := isDeploymentUnitReady(ctx, namespace, deploymentName, cs, debug)
		if err != nil {
			return err
		}
	}

	return nil
}

func isDeploymentUnitReady(ctx context.Context,
	namespace, deploymentName string,
	clientSet *kubernetes.Clientset,
	debug bool) error {
	for {
		deployment, err := clientSet.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
			if debug {
				fmt.Printf("Deployment %s is ready\n", deploymentName)
			}
			break
		} else {
			if debug {
				fmt.Printf("Deployment %s is not ready yet, ready replicas: %v\n", deploymentName, deployment.Status.ReadyReplicas)
			}
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

// IsClusterReady checks if a kubernetes cluster is ready
func IsClusterReady(ctx context.Context, kubeconfig string) (bool, error) {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return false, err
	}
	nodes, err := cs.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, node := range nodes.Items {
		isNodeReady := false
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" && condition.Status == "True" {
				isNodeReady = true
			}
		}

		if !isNodeReady {
			return false, nil
		}
	}

	return true, nil
}

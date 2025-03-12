package kubernetes

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func RestartPods(kubeconfig, namespace string, podNames []string, debug bool) error {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return err
	}

	pods, err := GetPodsInNamespace(kubeconfig, namespace, debug)
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground

	for _, pod := range pods {
		for _, podName := range podNames {
			if strings.Contains(pod, podName) {
				err := cs.CoreV1().
					Pods(namespace).
					Delete(context.TODO(), pod, metav1.DeleteOptions{
						PropagationPolicy: &deletePolicy,
					})
				if err != nil {
					return fmt.Errorf("failed to delete pod: %s", err)
				}

				if debug {
					fmt.Printf("Restarted pod: %s\n", pod)
				}
			}
		}
	}
	return nil
}

func GetPodsInNamespace(kubeconfig string, namespace string, debug bool) ([]string, error) {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return nil, err
	}

	pods, err := cs.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	if debug {
		fmt.Printf("Found pods: %v\n", podNames)
	}
	return podNames, nil
}

func GetPodsInDeployment(kubeconfig, namespace, deploymentName string, debug bool) ([]string, error) {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return nil, err
	}

	deploy, err := cs.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	labelSelector := metav1.FormatLabelSelector(deploy.Spec.Selector)
	pods, err := cs.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	if debug {
		fmt.Printf("Found pods: %v\n", podNames)
	}
	return podNames, nil
}

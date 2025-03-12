package kubernetes

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetServiceCIDR returns the service CIDR from the Kubernetes cluster
func GetServiceCIDR(ctx context.Context, kubeconfig string) (string, error) {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return "", err
	}

	nodeList, err := cs.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	for _, annotation := range nodeList.Items[0].Annotations {
		if annotation == "kubeadm.alpha.kubernetes.io/cri-socket" {
			return annotation, nil
		}
	}

	return "", nil
}

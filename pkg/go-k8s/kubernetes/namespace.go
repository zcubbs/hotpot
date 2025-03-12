package kubernetes

import (
	"context"
	apiv1 "k8s.io/api/core/v1"
	errosv1 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNamespace(kubeconfig string, namespace []string) error {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return err
	}

	for _, ns := range namespace {
		ns := &apiv1.Namespace{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Namespace",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: ns,
				Labels: map[string]string{
					"name": ns,
				},
			},
		}

		_, err := cs.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
		if err != nil && !errosv1.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

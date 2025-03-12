package kubernetes

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfigMap(kubeconfig, namespace, name string) (*v1.ConfigMap, error) {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return nil, err
	}
	cm, err := cs.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func CreateConfigMap(kubeconfig string, cm *v1.ConfigMap, namespace string) error {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return err
	}
	_, err = cs.CoreV1().ConfigMaps(namespace).Create(context.Background(), cm, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func UpdateConfigMap(kubeconfig string, cm *v1.ConfigMap, namespace string) error {
	cs, err := GetClientSet(kubeconfig)
	if err != nil {
		return err
	}
	_, err = cs.CoreV1().ConfigMaps(namespace).Update(context.Background(), cm, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

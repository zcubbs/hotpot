package certmanager

// DefaultManager is the default implementation of CertManager
type DefaultManager struct{}

func (d DefaultManager) Install(values Values, kubeconfig string, debug bool) error {
	return Install(values, kubeconfig, debug)
}

func (d DefaultManager) Uninstall(kubeconfig string, debug bool) error {
	return Uninstall(kubeconfig, debug)
}

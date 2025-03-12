package argocd

// DefaultManager is the default implementation of ArgoCDManager
type DefaultManager struct{}

func (d DefaultManager) Install(values Values, kubeconfig string, debug bool) error {
	return Install(values, kubeconfig, debug)
}

func (d DefaultManager) Uninstall(kubeconfig string, debug bool) error {
	return Uninstall(kubeconfig, debug)
}

func (d DefaultManager) CreateProject(project Project, kubeconfig string, debug bool) error {
	return CreateProject(project, kubeconfig, debug)
}

func (d DefaultManager) CreateApplication(app Application, kubeconfig string, debug bool) error {
	return CreateApplication(app, kubeconfig, debug)
}

func (d DefaultManager) CreateRepository(repo Repository, kubeconfig string, debug bool) error {
	return CreateRepository(repo, kubeconfig, debug)
}

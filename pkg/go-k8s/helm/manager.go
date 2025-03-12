package helm

// DefaultManager is the default implementation of HelmManager
type DefaultManager struct{}

func (d DefaultManager) Install(debug bool) error {
	return Install(debug)
}

func (d DefaultManager) InstallCli(debug bool) error {
	return Install(debug)
}

func (d DefaultManager) Uninstall(debug bool) error {
	return Uninstall(debug)
}

func (d DefaultManager) IsHelmInstalled() (bool, error) {
	return IsHelmInstalled()
}

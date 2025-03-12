package k3s

// DefaultManager is the default implementation of K3sManager
type DefaultManager struct{}

func (d DefaultManager) Install(cfg Config, debug bool) error {
	return Install(cfg, debug)
}

func (d DefaultManager) Uninstall(debug bool) error {
	return Uninstall(debug)
}

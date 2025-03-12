package k9s

// DefaultManager is the default implementation of K9sManager
type DefaultManager struct{}

func (d DefaultManager) Install(debug bool) error {
	return Install(debug)
}

func (d DefaultManager) Uninstall(debug bool) error {
	return Uninstall(debug)
}

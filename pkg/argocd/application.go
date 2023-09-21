package argocd

type Application struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func CreateApplication(app Application, _ string, debug bool) error {
	return nil
}

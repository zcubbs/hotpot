package argocd

type Project struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
}

func CreateProject(project Project, _ string, debug bool) error {
	return nil
}

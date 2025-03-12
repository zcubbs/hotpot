package argocd

// Config represents ArgoCD configuration
type Config struct {
	Values map[string]interface{} `yaml:"values"`
}

// ProjectConfig represents ArgoCD project configuration
type ProjectConfig = Project

// ApplicationConfig represents ArgoCD application configuration
type ApplicationConfig = Application

// RepositoryConfig represents ArgoCD repository configuration
type RepositoryConfig = Repository

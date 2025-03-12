package traefik

// Config represents Traefik configuration
type Config struct {
	Values map[string]interface{} `yaml:"values"`
}

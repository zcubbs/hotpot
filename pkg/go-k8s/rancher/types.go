package rancher

// Config represents Rancher configuration
type Config struct {
	Values map[string]interface{} `yaml:"values"`
}

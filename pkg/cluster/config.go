package cluster

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Configuration struct {
	Kubeconfig string    `mapstructure:"kubeconfig" json:"kubeconfig" yaml:"kubeconfig"`
	Clusters   []Cluster `mapstructure:"clusters" json:"clusters" yaml:"clusters"`
}

func LoadConfig(path string) (*Configuration, error) {
	var config *Configuration
	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into the provided interface
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *Configuration) error {
	return nil
}

package recipe

import (
	"gopkg.in/yaml.v2"
	"os"
)

func Load(path string) (*Recipe, error) {
	var config *Recipe
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

func validate(_ *Recipe) error {
	return nil
}

package recipe

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

func Load(path string) (*Recipe, error) {
	var recipe Recipe

	initViperPresets(path)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load recipe file path=%s err=%s", path, err)
	}

	err = viper.Unmarshal(&recipe)
	if err != nil {
		return nil, fmt.Errorf("could not decode recipe into struct err=%s", err)
	}
	return &recipe, nil
}

func initViperPresets(path string) {
	dir := filepath.Dir(path)
	file := filepath.Base(path)
	viper.AddConfigPath(dir)
	viper.SetConfigName(file)
	viper.SetConfigType("yaml")
}

func validate(r *Recipe) error {
	if r.Kubeconfig == "" {
		r.Kubeconfig = "/etc/rancher/k3s/k3s.yaml"
	}
	return nil
}

func printRecipe(recipe *Recipe) {
	jsonConfig, err := json.MarshalIndent(recipe, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("recipe path: %s\n", viper.ConfigFileUsed())
	fmt.Printf("%v\n", string(jsonConfig))
}

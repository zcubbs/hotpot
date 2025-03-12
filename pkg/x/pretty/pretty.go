package pretty

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

func PrintJson(v interface{}) {
	fmt.Println(string(toJson(v)))
}

func PrintYaml(v interface{}) {
	fmt.Println(string(toYaml(v)))
}

func toJson(v interface{}) []byte {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return b
}

func toYaml(v interface{}) []byte {
	b, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

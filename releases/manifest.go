package releases

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Manifest struct {
	Name       string   `yaml:"name"`
	Cluster    string   `yaml:"cluster"`
	Namespace  string   `yaml:"namespace"`
	Chart      string   `yaml:"chart"`
	ValueFiles []string `yaml:"values"`
	Values     []string `yaml:"set"`
}

func ReadManifestFromFile(filePath string) (*Manifest, error) {
	var manifest Manifest

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(content, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

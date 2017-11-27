package cmd

import (
	"testing"

	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
	"github.com/imdario/mergo"
	yaml "gopkg.in/yaml.v2"
)

func TestLoadValues(t *testing.T) {
	fileA, _ := ioutil.ReadFile("samples/a.yml")
	fileB, _ := ioutil.ReadFile("samples/b.yml")

	var aMap map[string]interface{}
	yaml.Unmarshal(fileA, &aMap)

	var bMap map[string]interface{}
	yaml.Unmarshal(fileB, &bMap)

	err := mergo.Merge(&aMap, bMap)

	spew.Dump(err)
	spew.Dump(aMap)
}

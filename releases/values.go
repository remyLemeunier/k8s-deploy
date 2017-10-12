package releases

import (
	"fmt"
	"io/ioutil"

	"github.com/peterbourgon/mergemap"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/strvals"
)

// taken from helm's cmd/helm/install.go

func loadValues(valueFiles []string, values []string) ([]byte, error) {
	base := map[string]interface{}{}

	for _, filePath := range valueFiles {
		currentMap := map[string]interface{}{}

		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return []byte{}, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return []byte{}, fmt.Errorf("failed to parse %s : %s", filePath, err)
		}
		base = mergemap.Merge(base, currentMap)
	}

	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return []byte{}, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}

	return yaml.Marshal(base)
}

package releases

import (
	"fmt"
	"io/ioutil"

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
		base = mergeValues(base, currentMap)
	}

	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return []byte{}, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}

	return yaml.Marshal(base)
}

func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		if !ok {
			dest[k] = v
			continue
		}
		if _, exists := dest[k]; !exists {
			dest[k] = nextMap
			continue
		}
		destMap, isMap := dest[k].(map[string]interface{})
		if !isMap {
			dest[k] = v
			continue
		}
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}

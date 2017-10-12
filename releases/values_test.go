package releases

import (
	"testing"
)

func TestLoadValues(t *testing.T) {
	valueFiles := []string{
		"./testdata/values1.yaml",
		"./testdata/values2.yaml",
	}
	_, err := loadValues(valueFiles, []string{"a=value", "c=1"})
	if err != nil {
		t.Errorf("Unexpected err: %q", err)
	}
}

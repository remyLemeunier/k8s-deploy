package releases

import (
	"testing"
)

func TestLoadValues(t *testing.T) {
	valueFiles := []string{
		"./testdata/values1.yaml",
		"./testdata/values2.yaml",
	}
	v, err := LoadValues(valueFiles, []string{"a=value", "c=1"})
	if err != nil {
		t.Errorf("Unexpected err: %q", err)
	}

	if string(v) != "a: value\nattributes:\n  values1.yaml:\n    values1: a\n  values2.yaml:\n    values2: b\nc: 1\n" {
		t.Errorf("Unexpected values: %q", string(v))
	}
}

package releases

import (
	"testing"
)

func TestLoadValues(t *testing.T) {
	valueFiles := []string{
		"./testdata/values1.yaml",
		"./testdata/values2.yaml",
	}
	v, err := loadValues(valueFiles, []string{"a=value", "c=1"})
	if err != nil {
		t.Errorf("Unexpected err: %q", err)
	}

	if string(v) != "a: value\nc: 1\ntata: z\ntoto: a\n" {
		t.Errorf("Unexpected values : %q", v)
	}
}

package releases

import (
	"testing"
)

func TestLoadValues(t *testing.T) {
	v, _ := loadValues([]string{}, []string{"a=value", "c=1"})
	if string(v) != "a: value\nc: 1\n" {
		t.Errorf("Unexpected values : %q", v)
	}
}

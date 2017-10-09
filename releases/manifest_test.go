package releases

import "testing"

func TestReadManifestFromFile(t *testing.T) {
	man, err := ReadManifestFromFile("./testdata/manifest.yaml")
	if err != nil {
		t.Fatalf("Unexpected error : %q", err)
	}
	if man.Name != "worker-payment-monitoring" {
		t.Errorf("Unexpected Name : %q", man.Name)
	}
}

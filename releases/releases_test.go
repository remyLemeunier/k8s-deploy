package releases

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
)

var mockRelease Release = Release{
	name:      "fake",
	chart:     &chart.Chart{},
	cluster:   "test",
	namespace: "test",
	hapiClient: &helm.FakeClient{
		Rels: []*release.Release{
			&release.Release{Name: "fake"},
		},
		Err: nil,
	},
	release: nil,
}

func TestNewRelease(t *testing.T) {
}

func TestNewReleaseFromManifest(t *testing.T) {
	_, err := NewReleaseFromManifest("./testdata/manifest.yaml")
	if err != nil {
		t.Errorf("Unexpected err: %q", err)
	}
}

func TestIsInstalled(t *testing.T) {
	r := mockRelease
	installed := r.isInstalled()
	if installed != true {
		t.Errorf("Unexpected value for isInstalled() : %q", installed)
	}

	r.hapiClient = &helm.FakeClient{
		Err: fmt.Errorf("release: %q not found", r.name),
	}
	installed = r.isInstalled()
	if installed != false {
		t.Errorf("Unexpected value for isInstalled() : %q", installed)
	}

}

func TestReleaseLoadValues(t *testing.T) {
	r, err := NewReleaseFromManifest("./testdata/manifest.yaml")
	if err != nil {
		t.Errorf("Unexpected err: %q", err)
	}
	spew.Dump(r.valueFiles)
	overrides, err := r.LoadValues()
	if err != nil {
		t.Errorf("Unexpected err: %q", err)
	}
	spew.Dump(string(overrides))
}

func TestDeploy(t *testing.T) {
	r := mockRelease
	if err := r.Deploy(); err != nil {
		t.Errorf("Unexpected return for Deploy() : %q", err)
	}
}

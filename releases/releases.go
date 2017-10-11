package releases

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"text/tabwriter"

	"k8s.io/helm/pkg/proto/hapi/services"

	"github.com/remyLemeunier/k8s-deploy/helmclient"
	"github.com/sergi/go-diff/diffmatchpatch"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/timeconv"
)

type Releaser interface {
	Deploy() error
	Content() (*services.GetReleaseContentResponse, error)
	Status() (*services.GetReleaseStatusResponse, error)
	AddValues([]string, []string)
}

type FakeRelease struct {
	Release
}

func (r *FakeRelease) Deploy() error {
	return nil
}

type Release struct {
	name       string
	chart      *chart.Chart
	cluster    string
	namespace  string
	valueFiles []string
	values     []string
	hapiClient helm.Interface
	release    *release.Release
}

func NewRelease(name string, cluster string, namespace string, chartPath string,
	valueFiles []string, values []string) (*Release, error) {

	chart, err := chartutil.Load(chartPath)
	if err != nil {
		log.Debugf("Could no load chart : %q", chartPath)
		return nil, err
	}

	// find out which tiller to talk to
	kubeClient, err := helmclient.NewKubeClient("", cluster)
	if err != nil {
		return nil, err
	}

	tillerEps, err := helmclient.GetTillerHosts(kubeClient, "kube-system")
	if err != nil {
		return nil, err
	}

	return &Release{
		name:       name,
		chart:      chart,
		cluster:    cluster,
		namespace:  namespace,
		valueFiles: valueFiles,
		values:     values,
		hapiClient: helm.NewClient(helm.Host(tillerEps[0])),
	}, nil
}

func NewReleaseFromManifest(manifestPath string) (*Release, error) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	if err := os.Chdir(path.Dir(manifestPath)); err != nil {
		return nil, err
	}

	manifest, err := ReadManifestFromFile(manifestPath)
	if err != nil {
		return nil, err
	}

	release, err := NewRelease(manifest.Name, manifest.Cluster, manifest.Namespace, manifest.Chart, manifest.ValueFiles, manifest.Values)
	if err != nil {
		return nil, err
	}

	if cwd != "" {
		os.Chdir(cwd)
	}
	return release, nil
}

func (r *Release) isInstalled() bool {
	_, err := r.hapiClient.ReleaseHistory(r.name, helm.WithMaxHistory(1))
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("release: %q not found", r.name)) {
		return false
	}
	return true
}

func (r *Release) AddValues(valueFiles []string, values []string) {
	for _, vf := range valueFiles {
		r.valueFiles = append(r.valueFiles, vf)
	}
	for _, v := range values {
		r.values = append(r.values, v)
	}
}

func (r *Release) Deploy() error {
	overrides, err := loadValues(r.valueFiles, r.values)
	if err != nil {
		return err
	}

	if !r.isInstalled() {
		log.Debugf("Installing release %s", r.name)
		response, err := r.hapiClient.InstallReleaseFromChart(
			r.chart,
			r.namespace,
			helm.ReleaseName(r.name),
			helm.InstallDryRun(false),
			helm.ValueOverrides(overrides),
		)
		if err != nil {
			return err
		}
		r.release = response.GetRelease()
	} else {
		log.Debugf("Updating release %s", r.name)
		response, err := r.hapiClient.UpdateReleaseFromChart(
			r.name,
			r.chart,
			helm.UpgradeDryRun(false),
			helm.UpdateValueOverrides(overrides),
		)
		if err != nil {
			return err
		}

		r.release = response.GetRelease()
	}
	log.Debugf("Deployed release %s", r.name)
	return nil
}

func (r *Release) PrintDiff(out io.Writer) error {
	response, err := r.hapiClient.ReleaseContent(r.name)
	if err != nil {
		return err
	}
	orig := response.Release.Manifest

	responseNew, err := r.hapiClient.UpdateReleaseFromChart(
		r.name,
		r.chart,
		helm.UpgradeDryRun(true),
		helm.UpdateValueOverrides([]byte{}),
	)

	if err != nil {
		return err
	}
	new := responseNew.Release.Manifest

	dmp := diffmatchpatch.New()
	diff := dmp.DiffMain(orig, new, false)
	fmt.Fprintf(out, dmp.DiffPrettyText(diff))
	return nil
}

func (r *Release) Content() (*services.GetReleaseContentResponse, error) {
	response, err := r.hapiClient.ReleaseContent(r.name)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *Release) PrintContent(out io.Writer) error {
	response, err := r.hapiClient.ReleaseContent(r.name)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, response.Release.Manifest)
	return nil
}

func (r *Release) Status() (*services.GetReleaseStatusResponse, error) {
	status, err := r.hapiClient.ReleaseStatus(r.name)
	if err != nil {
		return nil, errors.New(grpc.ErrorDesc(err))
	}
	return status, nil
}

// taken from helm's printstatus()
func (r *Release) PrintStatus(out io.Writer) error {
	status, err := r.hapiClient.ReleaseStatus(r.name)
	if err != nil {
		return errors.New(grpc.ErrorDesc(err))
	}

	if status.Info.LastDeployed != nil {
		fmt.Fprintf(out, "LAST DEPLOYED: %s\n", timeconv.String(status.Info.LastDeployed))
	}
	fmt.Fprintf(out, "NAMESPACE: %s\n", status.Namespace)
	fmt.Fprintf(out, "STATUS: %s\n", status.Info.Status.Code)
	fmt.Fprintf(out, "\n")
	if len(status.Info.Status.Resources) > 0 {
		re := regexp.MustCompile("  +")

		w := tabwriter.NewWriter(out, 0, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "RESOURCES:\n%s\n", re.ReplaceAllString(status.Info.Status.Resources, "\t"))
		w.Flush()
	}
	if status.Info.Status.LastTestSuiteRun != nil {
		lastRun := status.Info.Status.LastTestSuiteRun
		fmt.Fprintf(out, "TEST SUITE:\n%s\n%s\n\n",
			fmt.Sprintf("Last Started: %s", timeconv.String(lastRun.StartedAt)),
			fmt.Sprintf("Last Completed: %s", timeconv.String(lastRun.CompletedAt)),
		)
	}

	if len(status.Info.Status.Notes) > 0 {
		fmt.Fprintf(out, "NOTES:\n%s\n", status.Info.Status.Notes)
	}
	return nil
}

package helmclient

import (
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

//type Chart struct {
//	name string
//}

func NewChart(chartPath string) (*chart.Chart, error) {
	chart, err := chartutil.Load(chartPath)
	if err != nil {
		return nil, err
	}

	return chart, nil
}

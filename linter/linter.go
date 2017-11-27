package linter

import (
	"path/filepath"

	"github.com/remyLemeunier/k8s-deploy/linter/rules"
	"github.com/remyLemeunier/k8s-deploy/linter/support"
)

func All(basedir string) support.HelmDeployLinter {
	chartDir, _ := filepath.Abs(basedir)

	linter := support.HelmDeployLinter{ChartDir: chartDir}
	rules.DeploymentRule(&linter)

	return linter
}

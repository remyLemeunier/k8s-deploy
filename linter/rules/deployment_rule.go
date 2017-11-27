package rules

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/remyLemeunier/k8s-deploy/linter/support"
	"gopkg.in/yaml.v2"
)

type Deployment struct {
	Spec struct {
		Template struct {
			Spec struct {
				Containers []struct {
					Name           string                 `yaml:"name"`
					LivenessProbe  map[string]interface{} `yaml:"livenessProbe,omitempty"`
					ReadinessProbe map[string]interface{} `yaml:"readinessProbe,omitempty"`
				} `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

func DeploymentRule(linter *support.HelmDeployLinter) {
	deploymentFileName := "deployment.yaml"
	deploymentFileDir := "templates"
	deploymentPath := filepath.Join(linter.ChartDir, deploymentFileDir, deploymentFileName)

	linter.RunLinterRule(deploymentPath, validateLivenessReadynessProb(deploymentPath))

}

func validateLivenessReadynessProb(deploymentPath string) error {
	if _, err := os.Stat(deploymentPath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("File not found %s", deploymentPath))
	}
	content, err := ioutil.ReadFile(deploymentPath)
	if err != nil {
		return err
	}
	var deployment Deployment
	yaml.Unmarshal(content, &deployment)

	for _, container := range deployment.Spec.Template.Spec.Containers {
		if len(container.LivenessProbe) == 0 && len(container.ReadinessProbe) == 0 {
			return errors.New(fmt.Sprintf("The aci %s need a readiness or a liveness prob.", container.Name))
		}
	}

	return nil
}

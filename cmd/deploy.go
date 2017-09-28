package cmd

import (
	"fmt"

	"github.com/remyLemeunier/k8s-deploy/releases"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "",
	Long:  ``,
	RunE:  deploy,
}

var cluster string
var namespace string
var valueFiles []string
var values []string
var chartPath string

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&cluster, "cluster", "", "cluster name")
	deployCmd.Flags().StringVar(&namespace, "namespace", "", "namespace")
	deployCmd.Flags().StringVar(&chartPath, "chart", "", "chart")
	deployCmd.Flags().StringArrayVar(&valueFiles, "values", []string{}, "values")
	deployCmd.Flags().StringArrayVar(&values, "set", []string{}, "values")
}

func deploy(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Missing release name argument")
	}
	if cluster == "" {
		return fmt.Errorf("Missing cluster name argument")
	}
	if namespace == "" {
		return fmt.Errorf("Missing namespace argument")
	}

	release, err := releases.NewRelease(args[0], cluster, namespace, chartPath, valueFiles, values)
	if err != nil {
		return err
	}

	if err := release.Deploy(); err != nil {
		return err
	}
	return nil
}

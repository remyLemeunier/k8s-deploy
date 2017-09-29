package cmd

import (
	"fmt"
	"os"

	"github.com/remyLemeunier/k8s-deploy/releases"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "",
	Long:  ``,
	RunE:  deploy,
}

func init() {
	RootCmd.AddCommand(deployCmd)
}

func deploy(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Missing release manifest argument")
	}
	release, err := releases.NewReleaseFromManifest(args[0])
	if err != nil {
		return err
	}

	if err := release.Deploy(); err != nil {
		return err
	}

	if err := release.PrintStatus(os.Stdout); err != nil {
		return err
	}

	return nil
}

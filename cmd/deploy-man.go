package cmd

import (
	"fmt"
	"os"

	"github.com/remyLemeunier/k8s-deploy/releases"
	"github.com/spf13/cobra"
)

var deployManCmd = &cobra.Command{
	Use:   "deploy-man",
	Short: "",
	Long:  ``,
	RunE:  deployMan,
}

func init() {
	RootCmd.AddCommand(deployManCmd)
}

func deployMan(cmd *cobra.Command, args []string) error {
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

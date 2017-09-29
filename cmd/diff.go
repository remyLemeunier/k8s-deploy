package cmd

import (
	"fmt"
	"os"

	"github.com/remyLemeunier/k8s-deploy/releases"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "",
	Long:  ``,
	RunE:  diff,
}

func init() {
	RootCmd.AddCommand(diffCmd)
}

func diff(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Missing release manifest argument")
	}
	release, err := releases.NewReleaseFromManifest(args[0])
	if err != nil {
		return err
	}

	if err := release.PrintDiff(os.Stdout); err != nil {
		return err
	}

	return nil
}

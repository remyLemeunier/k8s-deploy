package cmd

import (
	"fmt"
	"os"

	"github.com/remyLemeunier/k8s-deploy/releases"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "",
	Long:  ``,
	RunE:  status,
}

func init() {
	RootCmd.AddCommand(statusCmd)
}

func status(cmd *cobra.Command, args []string) error {
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
	if err := release.PrintContent(os.Stdout); err != nil {
	}
	return nil
}

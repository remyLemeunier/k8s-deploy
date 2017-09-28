package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "k8s-deploy",
	Short: "",
	Long:  ``,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test.yaml)")
}

// // initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := homedir.Dir()
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
//
// 		// Search config in home directory with name ".test" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigName(".test")
// 	}
//
// 	viper.AutomaticEnv() // read in environment variables that match
//
// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }

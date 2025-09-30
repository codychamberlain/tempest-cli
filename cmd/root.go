/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tempest-cli",
	Short: "Command line application for accessing tempest station and forecast data",
	Long: `Application for accessing tempest station and forecast data via the command line.
	
	Use of the data requires an API token which can be obtained from tempestwx.com. `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringP("station", "s", "", "Station ID to pull data from")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output format - JSON")
}

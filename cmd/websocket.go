/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// websocketCmd represents the websocket command
var websocketCmd = &cobra.Command{
	Use:   "websocket",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("websocket called")

	},
}

func init() {
	rootCmd.AddCommand(websocketCmd)

	// websocketCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"chartgen/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "chartgen",
	Short: "Generate Helm-like values.yaml from Kubernetes resources",
	Long:  `chartgen is a kubectl plugin that outputs deployment, service, and ingress data in Helm values.yaml format`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use subcommands like 'generate'")
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(cmd.GenerateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
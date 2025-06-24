package cmd

import (
	"fmt"
	"os"

	"chartgen/internal"

	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Helm values.yaml from Kubernetes resources",
	Long:  `Generate a Helm values.yaml file from existing Kubernetes deployments, services, and ingresses`,
	Run: func(cmd *cobra.Command, args []string) {
		outputFile, _ := cmd.Flags().GetString("output")
		namespace, _ := cmd.Flags().GetString("namespace")
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		insecureSkipTLSVerify, _ := cmd.Flags().GetBool("insecure-skip-tls-verify")
		
		fmt.Printf("Generating Helm values.yaml from Kubernetes resources...\n")
		if namespace != "" {
			fmt.Printf("Target namespace: %s\n", namespace)
		}
		if kubeconfig != "" {
			fmt.Printf("Using kubeconfig: %s\n", kubeconfig)
		}
		
		// Create parser with or without kubeconfig
		var parser *internal.Parser
		if kubeconfig != "" {
			parser = internal.NewParserWithKubeconfig(namespace, kubeconfig, insecureSkipTLSVerify)
		} else {
			parser = internal.NewParser(namespace)
		}
		
		// Get Kubernetes resources
		fmt.Println("Fetching Kubernetes resources...")
		resources, err := parser.GetK8sResources()
		if err != nil {
			fmt.Printf("Error fetching Kubernetes resources: %v\n", err)
			os.Exit(1)
		}
		
		if len(resources) == 0 {
			fmt.Println("No Kubernetes resources found.")
			os.Exit(0)
		}
		
		fmt.Printf("Found %d resources\n", len(resources))
		
		// Parse to multiple Helm values (one per service)
		fmt.Println("Converting to Helm values structure...")
		serviceValues, err := parser.ParseToMultipleHelmValues(resources)
		if err != nil {
			fmt.Printf("Error parsing resources: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Generated values for %d services\n", len(serviceValues))
		
		// Generate YAML
		fmt.Println("Generating YAML output...")
		yamlOutput, err := parser.GenerateMultipleYAML(serviceValues)
		if err != nil {
			fmt.Printf("Error generating YAML: %v\n", err)
			os.Exit(1)
		}
		
		// Write to file or stdout
		if outputFile == "-" {
			fmt.Println(yamlOutput)
		} else {
			err = os.WriteFile(outputFile, []byte(yamlOutput), 0644)
			if err != nil {
				fmt.Printf("Error writing to file %s: %v\n", outputFile, err)
				os.Exit(1)
			}
			fmt.Printf("Helm values written to: %s\n", outputFile)
		}
	},
}

func init() {
	// Add flags here if needed
	GenerateCmd.Flags().StringP("output", "o", "values.yaml", "Output file path (use '-' for stdout)")
	GenerateCmd.Flags().StringP("namespace", "n", "", "Target namespace (default: current)")
	GenerateCmd.Flags().StringP("kubeconfig", "k", "", "Path to kubeconfig file")
	GenerateCmd.Flags().Bool("insecure-skip-tls-verify", false, "Skip TLS certificate verification when connecting to the Kubernetes API server")
} 
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "smart-log-analyser",
	Short: "A high-performance CLI tool for analysing Nginx access logs",
	Long: `Smart Log Analyser is designed to help system administrators and developers 
gain insights from their Nginx access logs. It provides statistical analysis, 
error pattern detection, traffic analysis, and real-time monitoring with 
configurable alerting.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(analyseCmd)
}
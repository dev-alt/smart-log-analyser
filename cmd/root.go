package cmd

import (
	"os"
	
	"github.com/spf13/cobra"
	"smart-log-analyser/pkg/menu"
)

var rootCmd = &cobra.Command{
	Use:   "smart-log-analyser",
	Short: "A high-performance CLI tool for analysing Nginx access logs",
	Long: `Smart Log Analyser is designed to help system administrators and developers 
gain insights from their Nginx access logs. It provides statistical analysis, 
error pattern detection, traffic analysis, and real-time monitoring with 
configurable alerting.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, launch interactive menu
		if len(args) == 0 {
			menuSystem := menu.New()
			if err := menuSystem.Run(); err != nil {
				os.Exit(1)
			}
			return
		}
		
		// Otherwise show help
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(analyseCmd)
	rootCmd.AddCommand(downloadCmd)
}
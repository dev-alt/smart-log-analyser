package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"smart-log-analyser/pkg/analyser"
	"smart-log-analyser/pkg/parser"
)

var (
	since   string
	until   string
	topIPs  int
	topURLs int
)

var analyseCmd = &cobra.Command{
	Use:   "analyse [log-file]",
	Short: "Analyse Nginx access logs",
	Long:  `Parse and analyse Nginx access logs to provide statistical insights.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logFile := args[0]
		
		p := parser.New()
		logs, err := p.ParseFile(logFile)
		if err != nil {
			log.Fatalf("Failed to parse log file: %v", err)
		}

		var sinceTime, untilTime *time.Time
		if since != "" {
			t, err := time.Parse("2006-01-02 15:04:05", since)
			if err != nil {
				log.Fatalf("Invalid since time format: %v", err)
			}
			sinceTime = &t
		}
		if until != "" {
			t, err := time.Parse("2006-01-02 15:04:05", until)
			if err != nil {
				log.Fatalf("Invalid until time format: %v", err)
			}
			untilTime = &t
		}

		a := analyser.New()
		results := a.Analyse(logs, sinceTime, untilTime)
		
		printResults(results)
	},
}

func init() {
	analyseCmd.Flags().StringVar(&since, "since", "", "Start time (YYYY-MM-DD HH:MM:SS)")
	analyseCmd.Flags().StringVar(&until, "until", "", "End time (YYYY-MM-DD HH:MM:SS)")
	analyseCmd.Flags().IntVar(&topIPs, "top-ips", 10, "Number of top IPs to show")
	analyseCmd.Flags().IntVar(&topURLs, "top-urls", 10, "Number of top URLs to show")
}

func printResults(results *analyser.Results) {
	fmt.Printf("=== Smart Log Analyser Results ===\n\n")
	fmt.Printf("Total Requests: %d\n", results.TotalRequests)
	fmt.Printf("Date Range: %s to %s\n\n", 
		results.TimeRange.Start.Format("2006-01-02 15:04:05"),
		results.TimeRange.End.Format("2006-01-02 15:04:05"))

	fmt.Printf("=== Status Code Distribution ===\n")
	for status, count := range results.StatusCodes {
		fmt.Printf("%s: %d\n", status, count)
	}

	fmt.Printf("\n=== Top %d IP Addresses ===\n", topIPs)
	count := 0
	for _, ip := range results.TopIPs {
		if count >= topIPs {
			break
		}
		fmt.Printf("%s: %d requests\n", ip.IP, ip.Count)
		count++
	}

	fmt.Printf("\n=== Top %d URLs ===\n", topURLs)
	count = 0
	for _, url := range results.TopURLs {
		if count >= topURLs {
			break
		}
		fmt.Printf("%s: %d requests\n", url.URL, url.Count)
		count++
	}
}
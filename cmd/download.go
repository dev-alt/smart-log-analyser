package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"smart-log-analyser/pkg/remote"
)

var (
	configFile   string
	serverName   string
	outputDir    string
	testConn     bool
	createConfig bool
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download log files from remote servers via SSH",
	Long: `Download Nginx access logs from remote servers using SSH credentials.
Requires a JSON configuration file with server details.`,
	Run: func(cmd *cobra.Command, args []string) {
		if createConfig {
			handleCreateConfig()
			return
		}

		if testConn {
			handleTestConnection()
			return
		}

		handleDownload()
	},
}

func init() {
	downloadCmd.Flags().StringVar(&configFile, "config", "servers.json", "Path to SSH configuration file")
	downloadCmd.Flags().StringVar(&serverName, "server", "", "Specific server to download from (host name)")
	downloadCmd.Flags().StringVar(&outputDir, "output", "./downloads", "Directory to save downloaded files")
	downloadCmd.Flags().BoolVar(&testConn, "test", false, "Test SSH connection without downloading")
	downloadCmd.Flags().BoolVar(&createConfig, "init", false, "Create a sample configuration file")
}

func handleCreateConfig() {
	if err := remote.CreateSampleConfig(configFile); err != nil {
		log.Fatalf("Failed to create config file: %v", err)
	}
	
	fmt.Printf("Created sample configuration file: %s\n", configFile)
	fmt.Println("Please edit the file with your server details before using.")
	fmt.Println("\nExample usage:")
	fmt.Println("  # Test connection")
	fmt.Println("  smart-log-analyser download --test")
	fmt.Println("  # Download logs")
	fmt.Println("  smart-log-analyser download")
}

func handleTestConnection() {
	config, err := remote.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Servers) == 0 {
		log.Fatal("No servers configured")
	}

	fmt.Println("Testing SSH connections...")

	for _, server := range config.Servers {
		if serverName != "" && server.Host != serverName {
			continue
		}

		fmt.Printf("Testing connection to %s@%s:%d... ", server.Username, server.Host, server.Port)
		
		if err := remote.TestConnection(&server); err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
		} else {
			fmt.Printf("✅ SUCCESS\n")
		}
	}
}

func handleDownload() {
	config, err := remote.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Servers) == 0 {
		log.Fatal("No servers configured")
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fmt.Printf("Downloading log files to: %s\n\n", outputDir)

	for _, server := range config.Servers {
		if serverName != "" && server.Host != serverName {
			continue
		}

		fmt.Printf("Connecting to %s@%s:%d...\n", server.Username, server.Host, server.Port)
		
		client := remote.NewSSHClient(&server)
		
		if err := client.Connect(); err != nil {
			fmt.Printf("❌ Failed to connect: %v\n\n", err)
			continue
		}

		// Generate local filename with timestamp and server name
		timestamp := time.Now().Format("20060102_150405")
		localFilename := fmt.Sprintf("%s_%s_access.log", server.Host, timestamp)
		localPath := filepath.Join(outputDir, localFilename)

		fmt.Printf("Downloading %s -> %s\n", server.LogPath, localFilename)
		
		if err := client.DownloadFile(server.LogPath, localPath); err != nil {
			fmt.Printf("❌ Download failed: %v\n", err)
			client.Close()
			continue
		}

		// Check file size
		if stat, err := os.Stat(localPath); err == nil {
			fmt.Printf("✅ Downloaded successfully (%d bytes)\n", stat.Size())
		} else {
			fmt.Printf("✅ Downloaded successfully\n")
		}

		client.Close()
		fmt.Println()
	}

	fmt.Println("Download completed!")
	fmt.Printf("Files saved to: %s\n", outputDir)
	fmt.Println("\nYou can now analyse the downloaded files:")
	fmt.Printf("  smart-log-analyser analyse %s/*.log\n", outputDir)
}
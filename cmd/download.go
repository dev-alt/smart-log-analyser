package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	downloadAll  bool
	listFiles    bool
	maxFiles     int
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

		if listFiles {
			handleListFiles()
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
	downloadCmd.Flags().BoolVar(&downloadAll, "all", false, "Download all access log files (current + rotated)")
	downloadCmd.Flags().BoolVar(&listFiles, "list", false, "List available log files without downloading")
	downloadCmd.Flags().IntVar(&maxFiles, "max-files", 10, "Maximum number of files to download when using --all")
}

func handleCreateConfig() {
	if err := remote.CreateSampleConfig(configFile); err != nil {
		// Check if it's because file already exists
		if strings.Contains(err.Error(), "already exists") {
			fmt.Printf("Configuration file '%s' already exists.\n", configFile)
			fmt.Println("Use --config flag to specify a different filename if needed.")
			fmt.Println("\nCurrent configuration:")
			
			// Try to load and display current config (safely)
			if config, loadErr := remote.LoadConfig(configFile); loadErr == nil {
				fmt.Printf("  - %d server(s) configured\n", len(config.Servers))
				for i, server := range config.Servers {
					fmt.Printf("  - Server %d: %s@%s:%d\n", i+1, server.Username, server.Host, server.Port)
				}
			}
			
			fmt.Println("\nExample usage:")
			fmt.Println("  # Test existing configuration")
			fmt.Println("  smart-log-analyser download --test")
			fmt.Println("  # Download logs with existing configuration")
			fmt.Println("  smart-log-analyser download")
			return
		}
		
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

func handleListFiles() {
	config, err := remote.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Servers) == 0 {
		log.Fatal("No servers configured")
	}

	fmt.Println("Listing available log files...\n")

	for _, server := range config.Servers {
		if serverName != "" && server.Host != serverName {
			continue
		}

		fmt.Printf("📋 Server: %s@%s:%d\n", server.Username, server.Host, server.Port)
		
		client := remote.NewSSHClient(&server)
		
		if err := client.Connect(); err != nil {
			fmt.Printf("❌ Failed to connect: %v\n\n", err)
			continue
		}

		// Get log directory from configured path
		logDir := filepath.Dir(server.LogPath)
		if logDir == "." {
			logDir = "/var/log/nginx"
		}

		// List access log files
		accessFiles, err := client.ListAccessLogFiles(logDir)
		if err != nil {
			fmt.Printf("❌ Failed to list files: %v\n", err)
			client.Close()
			continue
		}

		if len(accessFiles) > 0 {
			fmt.Printf("📁 Access log files in %s:\n", logDir)
			for i, file := range accessFiles {
				if i >= 20 { // Limit display
					fmt.Printf("   ... and %d more files\n", len(accessFiles)-i)
					break
				}
				fmt.Printf("   • %s\n", file)
			}
			fmt.Printf("   Total: %d files\n", len(accessFiles))
		} else {
			fmt.Printf("   No access log files found in %s\n", logDir)
		}

		client.Close()
		fmt.Println()
	}

	fmt.Println("Use --all flag to download all access log files.")
	fmt.Println("Use --max-files to limit the number of files downloaded.")
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

		var filesToDownload []string
		logDir := filepath.Dir(server.LogPath)
		if logDir == "." {
			logDir = "/var/log/nginx"
		}

		if downloadAll {
			// Download all access log files
			accessFiles, err := client.ListAccessLogFiles(logDir)
			if err != nil {
				fmt.Printf("❌ Failed to list files: %v\n", err)
				client.Close()
				continue
			}
			
			// Limit number of files
			if len(accessFiles) > maxFiles {
				fmt.Printf("⚠️  Found %d files, downloading first %d (use --max-files to change)\n", len(accessFiles), maxFiles)
				accessFiles = accessFiles[:maxFiles]
			}
			
			filesToDownload = accessFiles
			fmt.Printf("📦 Downloading %d access log files...\n", len(filesToDownload))
		} else {
			// Download single file as before
			filesToDownload = []string{server.LogPath}
		}

		timestamp := time.Now().Format("20060102_150405")
		totalBytes := int64(0)
		successCount := 0

		for i, remoteFile := range filesToDownload {
			// Generate local filename
			baseName := filepath.Base(remoteFile)
			localFilename := fmt.Sprintf("%s_%s_%s", server.Host, timestamp, baseName)
			localPath := filepath.Join(outputDir, localFilename)

			fmt.Printf("  [%d/%d] %s -> %s\n", i+1, len(filesToDownload), remoteFile, localFilename)
			
			if err := client.DownloadFile(remoteFile, localPath); err != nil {
				fmt.Printf("    ❌ Failed: %v\n", err)
				continue
			}

			// Check file size
			if stat, err := os.Stat(localPath); err == nil {
				fmt.Printf("    ✅ Downloaded (%d bytes)\n", stat.Size())
				totalBytes += stat.Size()
				successCount++
			} else {
				fmt.Printf("    ✅ Downloaded\n")
				successCount++
			}
		}

		fmt.Printf("📊 Summary: %d/%d files downloaded successfully", successCount, len(filesToDownload))
		if totalBytes > 0 {
			fmt.Printf(" (%d bytes total)", totalBytes)
		}
		fmt.Println()

		client.Close()
		fmt.Println()
	}

	fmt.Println("Download completed!")
	fmt.Printf("Files saved to: %s\n", outputDir)
	fmt.Println("\nYou can now analyse the downloaded files:")
	fmt.Printf("  smart-log-analyser analyse %s/*.log\n", outputDir)
}
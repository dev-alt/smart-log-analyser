package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"smart-log-analyser/pkg/ipc"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start IPC server for dashboard integration",
	Long: `Start the Smart Log Analyser IPC server for communication with external dashboards.

The server automatically uses:
- Named Pipes on Windows (\\.\pipe\SmartLogAnalyser)
- Unix Domain Sockets on Linux/macOS (/tmp/smart-log-analyser.sock)

Supported actions:
- analyze: Perform log analysis
- query: Execute SLAQ queries
- listPresets: List available analysis presets
- runPreset: Execute a specific preset
- getConfig: Retrieve current configuration
- updateConfig: Update configuration
- getStatus: Get server status
- shutdown: Gracefully shutdown server

Example usage:
  smart-log-analyser server

The server will run until interrupted (Ctrl+C) or shutdown via IPC command.`,
	Run: runServer,
}

var (
	serverPort int
	serverHost string
)

func init() {
	rootCmd.AddCommand(serverCmd)
	
	serverCmd.Flags().IntVar(&serverPort, "port", 0, "TCP port for testing (0 = use platform-specific IPC)")
	serverCmd.Flags().StringVar(&serverHost, "host", "127.0.0.1", "Host for TCP testing mode")
}

func runServer(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Starting Smart Log Analyser IPC Server...")
	
	server, err := ipc.NewServer()
	if err != nil {
		log.Fatalf("Failed to create IPC server: %v", err)
	}

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start IPC server: %v", err)
	}
	
	fmt.Println("✅ IPC Server is running")
	fmt.Println("📊 Ready to accept dashboard connections")
	fmt.Println("🔧 Supported actions: analyze, query, listPresets, runPreset, getConfig, updateConfig, getStatus, shutdown")
	fmt.Println("⚡ Use Ctrl+C to shutdown")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	<-sigChan
	fmt.Println("\n🛑 Shutting down IPC server...")
	
	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
	
	fmt.Println("👋 IPC Server stopped")
}
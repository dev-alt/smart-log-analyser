package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"smart-log-analyser/pkg/config"
)

var (
	configDir    string
	configList   string
	configReset  bool
	configBackup bool
	configExport string
	configImport string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration, presets, and templates",
	Long: `Manage Smart Log Analyser configuration including:

- Analysis presets for common scenarios (security, performance, traffic)
- Report templates for consistent output formatting  
- Server connection profiles for remote log access
- User preferences and default settings

Examples:
  # Initialize configuration with built-in presets
  ./smart-log-analyser config --init

  # List available presets
  ./smart-log-analyser config --list presets

  # List available templates
  ./smart-log-analyser config --list templates

  # Show configuration status
  ./smart-log-analyser config --status

  # Backup configuration
  ./smart-log-analyser config --backup

  # Reset to defaults
  ./smart-log-analyser config --reset`,
	Run: runConfig,
}

func init() {
	configCmd.Flags().StringVar(&configDir, "config-dir", "config", "Configuration directory path")
	configCmd.Flags().StringVar(&configList, "list", "", "List items (presets, templates, servers, categories)")
	configCmd.Flags().BoolVar(&configReset, "reset", false, "Reset configuration to defaults")
	configCmd.Flags().BoolVar(&configBackup, "backup", false, "Create configuration backup")
	configCmd.Flags().StringVar(&configExport, "export", "", "Export presets to file")
	configCmd.Flags().StringVar(&configImport, "import", "", "Import presets from file")
	configCmd.Flags().Bool("init", false, "Initialize configuration")
	configCmd.Flags().Bool("status", false, "Show configuration status")

	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) {
	installer := config.NewInstaller(configDir)

	// Handle initialization
	if init, _ := cmd.Flags().GetBool("init"); init {
		fmt.Println("🔧 Initializing Smart Log Analyser configuration...")
		
		if err := installer.Initialize(); err != nil {
			fmt.Printf("❌ Failed to initialize configuration: %v\n", err)
			os.Exit(1)
		}
		
		status, _ := installer.GetStatus()
		fmt.Println("✅ Configuration initialized successfully!")
		fmt.Printf("📊 Installed %d presets, %d templates\n", status.Presets, status.Templates)
		return
	}

	// Handle status
	if status, _ := cmd.Flags().GetBool("status"); status {
		showConfigStatus(installer)
		return
	}

	// Handle reset
	if configReset {
		fmt.Println("⚠️  Resetting configuration to defaults...")
		
		if err := installer.Reset(); err != nil {
			fmt.Printf("❌ Failed to reset configuration: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ Configuration reset successfully!")
		return
	}

	// Handle backup
	if configBackup {
		fmt.Println("💾 Creating configuration backup...")
		
		backupFile, err := installer.Backup()
		if err != nil {
			fmt.Printf("❌ Failed to create backup: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✅ Backup created: %s\n", backupFile)
		return
	}

	// Handle export
	if configExport != "" {
		fmt.Printf("📤 Exporting presets to %s...\n", configExport)
		
		if err := installer.ExportPresets(configExport); err != nil {
			fmt.Printf("❌ Failed to export presets: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ Presets exported successfully!")
		return
	}

	// Handle import
	if configImport != "" {
		fmt.Printf("📥 Importing presets from %s...\n", configImport)
		
		if err := installer.ImportPresets(configImport); err != nil {
			fmt.Printf("❌ Failed to import presets: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ Presets imported successfully!")
		return
	}

	// Handle listing
	if configList != "" {
		handleConfigList(installer, configList)
		return
	}

	// Default: show help
	cmd.Help()
}

func showConfigStatus(installer *config.Installer) {
	status, err := installer.GetStatus()
	if err != nil {
		fmt.Printf("❌ Failed to get configuration status: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("📋 Smart Log Analyser Configuration Status")
	fmt.Println("==========================================")
	fmt.Printf("Configuration Directory: %s\n", status.ConfigDir)
	fmt.Printf("Configuration File:      %s\n", status.ConfigFile)
	fmt.Printf("Initialized:             %v\n", status.Initialized)
	fmt.Printf("Presets:                 %d\n", status.Presets)
	fmt.Printf("Templates:               %d\n", status.Templates)
	fmt.Printf("Server Profiles:         %d\n", status.Servers)
	fmt.Println()

	if !status.Initialized {
		fmt.Println("💡 Run './smart-log-analyser config --init' to initialize configuration")
	}
}

func handleConfigList(installer *config.Installer, listType string) {
	// Ensure configuration is loaded
	configManager := config.NewConfigManager(configDir)
	if err := configManager.Load(); err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		fmt.Println("💡 Run './smart-log-analyser config --init' to initialize configuration")
		os.Exit(1)
	}

	switch listType {
	case "presets":
		listPresets(configManager)
	case "templates":
		listTemplates(configManager)
	case "servers":
		listServerProfiles(configManager)
	case "categories":
		listPresetCategories()
	default:
		fmt.Printf("❌ Unknown list type: %s\n", listType)
		fmt.Println("Available types: presets, templates, servers, categories")
		os.Exit(1)
	}
}

func listPresets(cm *config.ConfigManager) {
	presets := cm.GetConfig().Presets
	
	if len(presets) == 0 {
		fmt.Println("No presets available. Run './smart-log-analyser config --init' to install built-in presets.")
		return
	}

	fmt.Printf("📊 Available Analysis Presets (%d)\n", len(presets))
	fmt.Println("====================================")

	// Group by category
	categories := make(map[string][]config.AnalysisPreset)
	for _, preset := range presets {
		categories[preset.Category] = append(categories[preset.Category], preset)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	
	for category, categoryPresets := range categories {
		fmt.Printf("\n🏷️  %s\n", category)
		fmt.Fprintln(w, "Name\tDescription\tQuery")
		fmt.Fprintln(w, "----\t-----------\t-----")
		
		for _, preset := range categoryPresets {
			query := preset.Query
			if len(query) > 50 {
				query = query[:47] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", preset.Name, preset.Description, query)
		}
		w.Flush()
	}
	fmt.Println()
}

func listTemplates(cm *config.ConfigManager) {
	templates := cm.GetConfig().Templates
	
	if len(templates) == 0 {
		fmt.Println("No templates available. Run './smart-log-analyser config --init' to install built-in templates.")
		return
	}

	fmt.Printf("📄 Available Report Templates (%d)\n", len(templates))
	fmt.Println("==================================")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Name\tCategory\tSections\tDescription")
	fmt.Fprintln(w, "----\t--------\t--------\t-----------")
	
	for _, template := range templates {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", 
			template.Name, 
			template.Category, 
			len(template.Sections), 
			template.Description)
	}
	
	w.Flush()
	fmt.Println()
}

func listServerProfiles(cm *config.ConfigManager) {
	servers := cm.GetConfig().Servers
	
	if len(servers) == 0 {
		fmt.Println("No server profiles configured.")
		fmt.Println("💡 Add server profiles using the interactive menu or by editing config/app.yaml")
		return
	}

	fmt.Printf("🌐 Server Connection Profiles (%d)\n", len(servers))
	fmt.Println("==================================")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Name\tHost\tPort\tUser\tLog Path")
	fmt.Fprintln(w, "----\t----\t----\t----\t--------")
	
	for _, server := range servers {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n", 
			server.Name, 
			server.Host, 
			server.Port, 
			server.Username, 
			server.LogPath)
	}
	
	w.Flush()
	fmt.Println()
}

func listPresetCategories() {
	categories := config.GetPresetCategories()
	
	fmt.Printf("🏷️  Preset Categories (%d)\n", len(categories))
	fmt.Println("=========================")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Icon\tName\tDescription")
	fmt.Fprintln(w, "----\t----\t-----------")
	
	for _, category := range categories {
		fmt.Fprintf(w, "%s\t%s\t%s\n", category.Icon, category.Name, category.Description)
	}
	
	w.Flush()
	fmt.Println()
}
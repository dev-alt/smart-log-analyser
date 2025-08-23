# Configuration Directory

This directory is reserved for configuration files and settings for the Smart Log Analyser.

## Future Configuration Files

This directory is planned to contain:

### SSH Configuration
- `servers.yaml` - YAML-based server configuration (alternative to servers.json)
- `ssh_config` - SSH client configuration options

### Analysis Configuration
- `analysis.yaml` - Custom analysis rules and patterns
- `filters.yaml` - Custom filtering and exclusion rules
- `alerts.yaml` - Alert thresholds and notification settings

### Bot Detection Configuration
- `bots.yaml` - Custom bot detection patterns and rules
- `user_agents.yaml` - User agent classification rules

### Output Configuration
- `templates/` - Custom report templates
- `formats.yaml` - Output format configurations

## Current Configuration

Currently, server configuration is handled through:
- `servers.json` (in root directory)
- Command-line flags and environment variables

## Security Note

Configuration files may contain sensitive information such as:
- Server credentials
- API keys
- SSH keys and certificates

This directory will be properly secured in the .gitignore when configuration files are implemented.
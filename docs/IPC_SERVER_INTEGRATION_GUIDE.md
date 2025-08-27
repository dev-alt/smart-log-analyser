# Smart Log Analyser IPC Server Integration Guide

**Version**: 1.0  
**Date**: August 2025  
**Target Audience**: Avalonia Dashboard Development Team

---

## Table of Contents

1. [Overview](#overview)
2. [Getting Started](#getting-started)
3. [Server Architecture](#server-architecture)
4. [Communication Protocol](#communication-protocol)
5. [Available Operations](#available-operations)
6. [C# Client Implementation](#c-client-implementation)
7. [Data Models](#data-models)
8. [Error Handling](#error-handling)
9. [Examples](#examples)
10. [Testing](#testing)
11. [Deployment Considerations](#deployment-considerations)

---

## Overview

The Smart Log Analyser IPC (Inter-Process Communication) Server provides a robust, cross-platform interface for integrating log analysis capabilities into external applications, specifically designed for C# Avalonia dashboard development.

### Key Features

- **🌐 Cross-Platform**: Automatically uses Named Pipes (Windows) or Unix Domain Sockets (Linux/macOS)
- **📡 JSON Protocol**: Simple, language-agnostic communication format
- **⚡ High Performance**: Direct IPC communication without HTTP overhead
- **🛡️ Secure**: Local-only communication, no network exposure
- **🔄 Multi-Client**: Supports concurrent dashboard connections
- **📊 Full Feature Access**: Complete Smart Log Analyser functionality via IPC

### Use Cases

- **Real-time Dashboard Analytics**: Live log analysis and visualization
- **Custom Reporting**: Execute specialized SLAQ queries for business intelligence
- **Preset Management**: Access to 12 built-in analysis presets
- **Configuration Synchronization**: Manage analysis settings from dashboard UI
- **Interactive Report Generation**: Generate HTML reports with embedded charts

---

## Getting Started

### Prerequisites

- Smart Log Analyser executable
- .NET 6.0 or later for C# client
- Understanding of async/await patterns in C#

### Starting the IPC Server

```bash
# Start the IPC server
./smart-log-analyser server

# Output:
🚀 Starting Smart Log Analyser IPC Server...
✅ IPC Server is running
📊 Ready to accept dashboard connections
🔧 Supported actions: analyze, query, listPresets, runPreset, getConfig, updateConfig, getStatus, shutdown
⚡ Use Ctrl+C to shutdown
```

### Communication Endpoints

| Platform | Communication Method | Endpoint |
|----------|---------------------|-----------|
| **Windows** | Named Pipes | `\\.\pipe\SmartLogAnalyser` |
| **Linux/macOS** | Unix Domain Sockets | `/tmp/smart-log-analyser.sock` |

---

## Server Architecture

### Core Components

```
┌─────────────────────────────────────┐
│         IPC Server                  │
├─────────────────────────────────────┤
│  • Cross-platform connection mgmt  │
│  • JSON protocol handler           │
│  • Multi-client support            │
│  • Request/response correlation    │
└─────────────────────────────────────┘
                    │
┌─────────────────────────────────────┐
│      Smart Log Analyser Engine      │
├─────────────────────────────────────┤
│  • Log parsing & analysis          │
│  • SLAQ query execution            │
│  • HTML report generation          │
│  • Configuration management        │
│  • Preset system                   │
└─────────────────────────────────────┘
```

### Concurrency Model

- **Multi-client support**: Multiple dashboard instances can connect simultaneously
- **Request isolation**: Each client request is handled independently
- **Thread-safe operations**: All analysis operations are concurrent-safe
- **Connection management**: Automatic cleanup on client disconnect

---

## Communication Protocol

### Message Format

All communication uses JSON over the IPC channel with the following structure:

#### Request Format
```json
{
  "id": "unique-request-identifier",
  "action": "operation-name",
  "logFile": "/path/to/log/file",
  "options": { /* operation-specific options */ },
  "query": "SLAQ query string",
  "preset": "preset-name",
  "config": { /* configuration data */ }
}
```

#### Response Format
```json
{
  "id": "matching-request-identifier",
  "success": true,
  "data": { /* operation results */ },
  "error": "error message if success=false"
}
```

### Request Correlation

Each request must include a unique `id` field. The server will return this same ID in the response, allowing clients to correlate requests with responses in concurrent scenarios.

---

## Available Operations

### 1. **analyze** - Comprehensive Log Analysis

Performs complete log analysis with all Smart Log Analyser features.

**Request:**
```json
{
  "id": "req-001",
  "action": "analyze",
  "logFile": "/var/log/nginx/access.log",
  "options": {
    "enableSecurity": true,
    "enablePerformance": true,
    "enableTrends": true,
    "generateHtml": true,
    "interactive": true,
    "htmlTitle": "Dashboard Analysis Report",
    "outputPath": "reports/analysis.html"
  }
}
```

**Response:**
```json
{
  "id": "req-001",
  "success": true,
  "data": {
    "results": {
      "summary": {
        "totalRequests": 15432,
        "uniqueIPs": 1247,
        "errorCount": 234,
        "errorRate": 0.0152
      },
      "security": {
        "threatCount": 12,
        "securityGrade": "Good",
        "securityScore": 78
      },
      "performance": {
        "performanceGrade": "Excellent",
        "performanceScore": 92,
        "recommendations": ["Enable gzip compression", "Optimize image sizes"]
      }
    },
    "htmlPath": "reports/analysis.html"
  }
}
```

### 2. **query** - Execute SLAQ Queries

Execute custom Smart Log Analyser Query Language (SLAQ) queries.

**Request:**
```json
{
  "id": "req-002",
  "action": "query",
  "logFile": "/var/log/nginx/access.log",
  "query": "SELECT ip, COUNT(*) as requests FROM logs WHERE status_code >= 400 GROUP BY ip ORDER BY requests DESC LIMIT 10"
}
```

**Response:**
```json
{
  "id": "req-002",
  "success": true,
  "data": {
    "queryResults": {
      "columns": ["ip", "requests"],
      "rows": [
        ["192.168.1.100", 45],
        ["10.0.0.5", 32],
        ["203.0.113.1", 28]
      ]
    }
  }
}
```

### 3. **listPresets** - Retrieve Analysis Presets

Get all available analysis presets (12 built-in presets available).

**Request:**
```json
{
  "id": "req-003",
  "action": "listPresets"
}
```

**Response:**
```json
{
  "id": "req-003",
  "success": true,
  "data": {
    "presets": [
      {
        "name": "security-overview",
        "description": "Comprehensive security threat analysis",
        "category": "security",
        "query": "SELECT ip, COUNT(*) as attempts FROM logs WHERE status_code = 401 GROUP BY ip"
      },
      {
        "name": "performance-analysis",
        "description": "Performance bottleneck identification",
        "category": "performance",
        "query": "SELECT url, AVG(response_size) as avg_size FROM logs GROUP BY url ORDER BY avg_size DESC"
      }
    ]
  }
}
```

### 4. **runPreset** - Execute Analysis Preset

Execute a specific analysis preset by name.

**Request:**
```json
{
  "id": "req-004",
  "action": "runPreset",
  "logFile": "/var/log/nginx/access.log",
  "preset": "security-overview"
}
```

### 5. **getConfig** - Retrieve Configuration

Get current system configuration including presets and templates.

**Request:**
```json
{
  "id": "req-005",
  "action": "getConfig"
}
```

### 6. **updateConfig** - Update Configuration

Update system configuration (presets, templates, preferences).

**Request:**
```json
{
  "id": "req-006",
  "action": "updateConfig",
  "config": {
    "analysis": {
      "defaultTopIPs": 15,
      "defaultTopURLs": 20
    },
    "preferences": {
      "defaultExportDir": "custom-reports"
    }
  }
}
```

### 7. **getStatus** - Server Status

Get server status and connection information.

**Request:**
```json
{
  "id": "req-007",
  "action": "getStatus"
}
```

**Response:**
```json
{
  "id": "req-007",
  "success": true,
  "data": {
    "status": "Smart Log Analyser IPC Server - Running on linux (Unix Domain Sockets) - 2 active clients"
  }
}
```

### 8. **shutdown** - Graceful Shutdown

Gracefully shutdown the IPC server.

**Request:**
```json
{
  "id": "req-008",
  "action": "shutdown"
}
```

---

## C# Client Implementation

### Complete Client Library

A comprehensive C# client library is provided at `examples/csharp/SmartLogAnalyserClient.cs`. This library handles:

- Automatic platform detection
- Connection management
- Async/await operations
- Error handling
- Request/response correlation

### Basic Usage

```csharp
using var client = new SmartLogAnalyserClient();

// Connect to server
if (!await client.ConnectAsync())
{
    throw new Exception("Failed to connect to Smart Log Analyser IPC server");
}

// Perform analysis
var analysisResult = await client.AnalyzeAsync("access.log", new AnalysisOptions
{
    EnableSecurity = true,
    EnablePerformance = true,
    GenerateHtml = true,
    Interactive = true
});

Console.WriteLine($"Total requests: {analysisResult.Results.Summary.TotalRequests}");
Console.WriteLine($"Security grade: {analysisResult.Results.Security.SecurityGrade}");
```

### Connection Handling

```csharp
public class SmartLogAnalyserClient : IDisposable
{
    private Stream communicationStream;
    private StreamReader reader;
    private StreamWriter writer;
    
    public async Task<bool> ConnectAsync()
    {
        try
        {
            if (RuntimeInformation.IsOSPlatform(OSPlatform.Windows))
            {
                await ConnectNamedPipeAsync();  // Named Pipes
            }
            else
            {
                await ConnectUnixSocketAsync(); // Unix Domain Sockets
            }
            
            isConnected = true;
            return true;
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Connection failed: {ex.Message}");
            return false;
        }
    }
}
```

---

## Data Models

### Core Data Structures

```csharp
public class AnalysisResult
{
    public LogAnalysisResults Results { get; set; }
    public string HtmlPath { get; set; }
}

public class LogAnalysisResults
{
    public LogEntry[] Entries { get; set; }
    public StatsSummary Summary { get; set; }
    public SecurityAnalysis Security { get; set; }
    public PerformanceAnalysis Performance { get; set; }
}

public class StatsSummary
{
    public int TotalRequests { get; set; }
    public int UniqueIPs { get; set; }
    public int ErrorCount { get; set; }
    public double ErrorRate { get; set; }
}

public class SecurityAnalysis
{
    public SecurityThreat[] Threats { get; set; }
    public string SecurityGrade { get; set; }  // Excellent/Good/Fair/Poor/Critical
    public double SecurityScore { get; set; }   // 0-100
}

public class PerformanceAnalysis
{
    public string PerformanceGrade { get; set; } // Excellent/Good/Fair/Poor/Critical
    public double PerformanceScore { get; set; }  // 0-100
    public string[] Recommendations { get; set; }
}
```

### Log Entry Structure

```csharp
public class LogEntry
{
    public string IP { get; set; }
    public DateTime Time { get; set; }
    public string Method { get; set; }       // GET, POST, etc.
    public string URL { get; set; }
    public int StatusCode { get; set; }      // 200, 404, 500, etc.
    public long Size { get; set; }           // Response size in bytes
    public string UserAgent { get; set; }
    public string Referer { get; set; }
}
```

---

## Error Handling

### Common Error Scenarios

| Error Type | Cause | Handling Strategy |
|------------|-------|-------------------|
| **Connection Failed** | Server not running | Retry with exponential backoff |
| **File Not Found** | Invalid log file path | Validate file paths before sending |
| **Invalid Query** | SLAQ syntax error | Provide query validation feedback |
| **Preset Not Found** | Preset name doesn't exist | List available presets first |
| **Configuration Error** | Invalid config format | Validate config structure |

### Error Response Format

```json
{
  "id": "req-001",
  "success": false,
  "error": "Log file not found: /invalid/path/access.log"
}
```

### C# Error Handling Pattern

```csharp
try
{
    var result = await client.AnalyzeAsync(logPath, options);
    // Handle success
}
catch (FileNotFoundException)
{
    // Handle file not found
    ShowFileSelectionDialog();
}
catch (InvalidOperationException ex) when (ex.Message.Contains("Not connected"))
{
    // Handle connection issues
    await ReconnectAsync();
}
catch (Exception ex)
{
    // Handle general errors
    LogError($"Analysis failed: {ex.Message}");
}
```

---

## Examples

### Example 1: Basic Dashboard Analytics

```csharp
public async Task<DashboardData> LoadDashboardDataAsync(string logFile)
{
    using var client = new SmartLogAnalyserClient();
    await client.ConnectAsync();
    
    // Get basic analysis
    var analysis = await client.AnalyzeAsync(logFile, new AnalysisOptions
    {
        EnableSecurity = true,
        EnablePerformance = true
    });
    
    // Get top error IPs
    var errorQuery = await client.QueryAsync(logFile,
        "SELECT ip, COUNT(*) as errors FROM logs WHERE status_code >= 400 GROUP BY ip ORDER BY errors DESC LIMIT 5");
    
    return new DashboardData
    {
        TotalRequests = analysis.Results.Summary.TotalRequests,
        ErrorRate = analysis.Results.Summary.ErrorRate,
        SecurityGrade = analysis.Results.Security.SecurityGrade,
        PerformanceScore = analysis.Results.Performance.PerformanceScore,
        TopErrorIPs = ProcessQueryResults(errorQuery.QueryResults)
    };
}
```

### Example 2: Preset Management UI

```csharp
public async Task LoadPresetsAsync()
{
    var presets = await client.ListPresetsAsync();
    
    foreach (var preset in presets.Presets)
    {
        var button = new Button
        {
            Content = $"{preset.Name}\n{preset.Description}",
            Tag = preset.Name,
            Command = new RelayCommand(() => ExecutePreset(preset.Name))
        };
        
        PresetPanel.Children.Add(button);
    }
}

private async Task ExecutePreset(string presetName)
{
    try
    {
        var result = await client.RunPresetAsync(CurrentLogFile, presetName);
        DisplayResults(result.QueryResults);
    }
    catch (Exception ex)
    {
        MessageBox.Show($"Failed to execute preset: {ex.Message}");
    }
}
```

### Example 3: Real-time Status Monitoring

```csharp
public async Task StartStatusMonitoring()
{
    var timer = new DispatcherTimer { Interval = TimeSpan.FromSeconds(30) };
    timer.Tick += async (s, e) =>
    {
        try
        {
            var status = await client.GetStatusAsync();
            StatusLabel.Content = status.Status;
            ConnectionIndicator.Fill = Brushes.Green;
        }
        catch
        {
            ConnectionIndicator.Fill = Brushes.Red;
            StatusLabel.Content = "Connection Lost";
        }
    };
    timer.Start();
}
```

---

## Testing

### Unit Testing with Mock IPC

```csharp
[Test]
public async Task AnalyzeAsync_ValidLogFile_ReturnsResults()
{
    // Arrange
    var mockClient = new Mock<ISmartLogAnalyserClient>();
    mockClient.Setup(x => x.AnalyzeAsync(It.IsAny<string>(), It.IsAny<AnalysisOptions>()))
           .ReturnsAsync(new AnalysisResult 
           {
               Results = new LogAnalysisResults 
               {
                   Summary = new StatsSummary { TotalRequests = 100 }
               }
           });
    
    // Act
    var result = await mockClient.Object.AnalyzeAsync("test.log", new AnalysisOptions());
    
    // Assert
    Assert.AreEqual(100, result.Results.Summary.TotalRequests);
}
```

### Integration Testing

```bash
# Test script provided at examples/test_ipc.sh
#!/bin/bash
echo "Testing IPC server..."
../smart-log-analyser server &
SERVER_PID=$!
sleep 2

echo '{"id":"test-1","action":"getStatus"}' | nc -U /tmp/smart-log-analyser.sock

kill $SERVER_PID
```

---

## Deployment Considerations

### Development Environment

1. **Server Process**: Run `smart-log-analyser server` during development
2. **Auto-start**: Consider adding server startup to IDE launch configuration
3. **Debugging**: Both server and client can be debugged independently

### Production Environment

1. **Service Installation**: Consider running server as system service
2. **Process Management**: Implement server health checks and auto-restart
3. **Log File Access**: Ensure server has appropriate file permissions
4. **Resource Management**: Monitor memory usage for large log files

### Security Considerations

- **Local Only**: IPC communication is restricted to local machine
- **File Access**: Server inherits file system permissions of launching user
- **Process Isolation**: Each client connection runs in isolated context
- **Input Validation**: All log file paths and queries are validated

### Performance Guidelines

- **Concurrent Connections**: Server supports multiple simultaneous clients
- **Large Files**: Consider chunking analysis for files > 1GB
- **Memory Usage**: Monitor memory consumption during analysis
- **Query Optimization**: Use appropriate WHERE clauses to limit result sets

---

## Support and Documentation

### Additional Resources

- **Main Documentation**: `/docs/README.md`
- **SLAQ Query Reference**: `/docs/SLAQ_REFERENCE.md`
- **Example Projects**: `/examples/csharp/`
- **Configuration Guide**: `/docs/CONFIGURATION.md`

### Troubleshooting

| Issue | Solution |
|-------|----------|
| Server won't start | Check if port/socket is already in use |
| Connection refused | Verify server is running and permissions are correct |
| Large file analysis slow | Consider using query filters to reduce dataset size |
| Memory issues | Implement pagination for large result sets |

---

**Document Version**: 1.0  
**Last Updated**: August 2025  
**Contact**: Smart Log Analyser Development Team  

This document provides comprehensive guidance for integrating Smart Log Analyser's IPC server into your Avalonia dashboard application. For technical support or additional questions, please refer to the project documentation or contact the development team.
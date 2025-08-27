using System;
using System.IO;
using System.IO.Pipes;
using System.Net.Sockets;
using System.Runtime.InteropServices;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;

namespace SmartLogAnalyser.Client
{
    /// <summary>
    /// Cross-platform client for Smart Log Analyser IPC server
    /// Automatically uses Named Pipes on Windows, Unix Domain Sockets on other platforms
    /// </summary>
    public class SmartLogAnalyserClient : IDisposable
    {
        private Stream communicationStream;
        private StreamReader reader;
        private StreamWriter writer;
        private bool isConnected;
        
        /// <summary>
        /// Connect to the Smart Log Analyser IPC server
        /// </summary>
        public async Task<bool> ConnectAsync()
        {
            try
            {
                if (RuntimeInformation.IsOSPlatform(OSPlatform.Windows))
                {
                    await ConnectNamedPipeAsync();
                }
                else
                {
                    await ConnectUnixSocketAsync();
                }
                
                isConnected = true;
                Console.WriteLine("✅ Connected to Smart Log Analyser IPC server");
                return true;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"❌ Connection failed: {ex.Message}");
                return false;
            }
        }
        
        /// <summary>
        /// Connect using Named Pipes (Windows)
        /// </summary>
        private async Task ConnectNamedPipeAsync()
        {
            var pipeClient = new NamedPipeClientStream(".", "SmartLogAnalyser", PipeDirection.InOut);
            await pipeClient.ConnectAsync(5000); // 5 second timeout
            
            communicationStream = pipeClient;
            reader = new StreamReader(communicationStream, Encoding.UTF8);
            writer = new StreamWriter(communicationStream, Encoding.UTF8) { AutoFlush = true };
        }
        
        /// <summary>
        /// Connect using Unix Domain Sockets (Linux/macOS)
        /// </summary>
        private async Task ConnectUnixSocketAsync()
        {
            var socket = new Socket(AddressFamily.Unix, SocketType.Stream, ProtocolType.Unspecified);
            var endpoint = new UnixDomainSocketEndPoint("/tmp/smart-log-analyser.sock");
            
            await socket.ConnectAsync(endpoint);
            
            communicationStream = new NetworkStream(socket, true);
            reader = new StreamReader(communicationStream, Encoding.UTF8);
            writer = new StreamWriter(communicationStream, Encoding.UTF8) { AutoFlush = true };
        }
        
        /// <summary>
        /// Analyze a log file
        /// </summary>
        public async Task<AnalysisResult> AnalyzeAsync(string logFilePath, AnalysisOptions options = null)
        {
            if (!isConnected) throw new InvalidOperationException("Not connected to server");
            
            options ??= new AnalysisOptions();
            
            var request = new IPCRequest
            {
                Id = Guid.NewGuid().ToString(),
                Action = "analyze",
                LogFile = logFilePath,
                Options = options
            };
            
            return await SendRequestAsync<AnalysisResult>(request);
        }
        
        /// <summary>
        /// Execute a SLAQ query
        /// </summary>
        public async Task<QueryResult> QueryAsync(string logFilePath, string query)
        {
            if (!isConnected) throw new InvalidOperationException("Not connected to server");
            
            var request = new IPCRequest
            {
                Id = Guid.NewGuid().ToString(),
                Action = "query",
                LogFile = logFilePath,
                Query = query
            };
            
            return await SendRequestAsync<QueryResult>(request);
        }
        
        /// <summary>
        /// List available presets
        /// </summary>
        public async Task<PresetListResult> ListPresetsAsync()
        {
            if (!isConnected) throw new InvalidOperationException("Not connected to server");
            
            var request = new IPCRequest
            {
                Id = Guid.NewGuid().ToString(),
                Action = "listPresets"
            };
            
            return await SendRequestAsync<PresetListResult>(request);
        }
        
        /// <summary>
        /// Run a specific preset
        /// </summary>
        public async Task<QueryResult> RunPresetAsync(string logFilePath, string presetName)
        {
            if (!isConnected) throw new InvalidOperationException("Not connected to server");
            
            var request = new IPCRequest
            {
                Id = Guid.NewGuid().ToString(),
                Action = "runPreset",
                LogFile = logFilePath,
                Preset = presetName
            };
            
            return await SendRequestAsync<QueryResult>(request);
        }
        
        /// <summary>
        /// Get server status
        /// </summary>
        public async Task<StatusResult> GetStatusAsync()
        {
            if (!isConnected) throw new InvalidOperationException("Not connected to server");
            
            var request = new IPCRequest
            {
                Id = Guid.NewGuid().ToString(),
                Action = "getStatus"
            };
            
            return await SendRequestAsync<StatusResult>(request);
        }
        
        /// <summary>
        /// Send request and receive response
        /// </summary>
        private async Task<T> SendRequestAsync<T>(IPCRequest request) where T : class
        {
            var json = JsonSerializer.Serialize(request, new JsonSerializerOptions
            {
                PropertyNamingPolicy = JsonNamingPolicy.CamelCase
            });
            
            await writer.WriteLineAsync(json);
            
            var responseJson = await reader.ReadLineAsync();
            if (string.IsNullOrEmpty(responseJson))
            {
                throw new Exception("No response received from server");
            }
            
            var response = JsonSerializer.Deserialize<IPCResponse>(responseJson, new JsonSerializerOptions
            {
                PropertyNamingPolicy = JsonNamingPolicy.CamelCase
            });
            
            if (!response.Success)
            {
                throw new Exception($"Server error: {response.Error}");
            }
            
            return JsonSerializer.Deserialize<T>(response.Data.ToString(), new JsonSerializerOptions
            {
                PropertyNamingPolicy = JsonNamingPolicy.CamelCase
            });
        }
        
        /// <summary>
        /// Disconnect from server
        /// </summary>
        public void Disconnect()
        {
            writer?.Close();
            reader?.Close();
            communicationStream?.Close();
            isConnected = false;
            Console.WriteLine("👋 Disconnected from Smart Log Analyser IPC server");
        }
        
        public void Dispose()
        {
            Disconnect();
        }
    }
    
    // Data classes for requests and responses
    public class IPCRequest
    {
        public string Id { get; set; }
        public string Action { get; set; }
        public string LogFile { get; set; }
        public AnalysisOptions Options { get; set; }
        public string Query { get; set; }
        public string Preset { get; set; }
    }
    
    public class IPCResponse
    {
        public string Id { get; set; }
        public bool Success { get; set; }
        public object Data { get; set; }
        public string Error { get; set; }
    }
    
    public class AnalysisOptions
    {
        public bool EnableSecurity { get; set; } = true;
        public bool EnablePerformance { get; set; } = true;
        public bool EnableTrends { get; set; } = true;
        public bool GenerateHtml { get; set; } = true;
        public bool Interactive { get; set; } = true;
        public string HtmlTitle { get; set; }
        public string OutputPath { get; set; }
    }
    
    public class AnalysisResult
    {
        public LogAnalysisResults Results { get; set; }
        public string HtmlPath { get; set; }
    }
    
    public class QueryResult
    {
        public object QueryResults { get; set; }
    }
    
    public class PresetListResult
    {
        public PresetInfo[] Presets { get; set; }
    }
    
    public class StatusResult
    {
        public string Status { get; set; }
    }
    
    public class LogAnalysisResults
    {
        public LogEntry[] Entries { get; set; }
        public StatsSummary Summary { get; set; }
        public SecurityAnalysis Security { get; set; }
        public PerformanceAnalysis Performance { get; set; }
    }
    
    public class LogEntry
    {
        public string IP { get; set; }
        public DateTime Time { get; set; }
        public string Method { get; set; }
        public string URL { get; set; }
        public int StatusCode { get; set; }
        public long Size { get; set; }
        public string UserAgent { get; set; }
        public string Referer { get; set; }
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
        public string SecurityGrade { get; set; }
        public double SecurityScore { get; set; }
    }
    
    public class PerformanceAnalysis
    {
        public string PerformanceGrade { get; set; }
        public double PerformanceScore { get; set; }
        public string[] Recommendations { get; set; }
    }
    
    public class SecurityThreat
    {
        public string Type { get; set; }
        public string Description { get; set; }
        public string Severity { get; set; }
        public int Count { get; set; }
    }
    
    public class PresetInfo
    {
        public string Name { get; set; }
        public string Description { get; set; }
        public string Category { get; set; }
        public string Query { get; set; }
    }
}
using System;
using System.Threading.Tasks;
using SmartLogAnalyser.Client;

namespace SmartLogAnalyser.Example
{
    class Program
    {
        static async Task Main(string[] args)
        {
            Console.WriteLine("🔧 Smart Log Analyser C# Client Example");
            Console.WriteLine("==========================================");
            
            using var client = new SmartLogAnalyserClient();
            
            // Connect to the IPC server
            if (!await client.ConnectAsync())
            {
                Console.WriteLine("Failed to connect to Smart Log Analyser IPC server");
                Console.WriteLine("Make sure the server is running: smart-log-analyser server");
                return;
            }
            
            try
            {
                // Example 1: Get server status
                Console.WriteLine("\n📊 Getting server status...");
                var status = await client.GetStatusAsync();
                Console.WriteLine($"Status: {status.Status}");
                
                // Example 2: List available presets
                Console.WriteLine("\n📋 Listing available presets...");
                var presets = await client.ListPresetsAsync();
                Console.WriteLine($"Found {presets.Presets.Length} presets:");
                foreach (var preset in presets.Presets)
                {
                    Console.WriteLine($"  • {preset.Name} ({preset.Category}): {preset.Description}");
                }
                
                // Example 3: Analyze a log file (replace with your log file path)
                var logFilePath = "/path/to/your/access.log";
                if (args.Length > 0)
                {
                    logFilePath = args[0];
                }
                
                Console.WriteLine($"\n🔍 Analyzing log file: {logFilePath}");
                
                var analysisOptions = new AnalysisOptions
                {
                    EnableSecurity = true,
                    EnablePerformance = true,
                    EnableTrends = true,
                    GenerateHtml = true,
                    Interactive = true,
                    HtmlTitle = "Dashboard Analysis Report",
                    OutputPath = "dashboard-report.html"
                };
                
                var analysisResult = await client.AnalyzeAsync(logFilePath, analysisOptions);
                
                Console.WriteLine("✅ Analysis completed!");
                Console.WriteLine($"📊 Total requests: {analysisResult.Results.Summary.TotalRequests}");
                Console.WriteLine($"🌐 Unique IPs: {analysisResult.Results.Summary.UniqueIPs}");
                Console.WriteLine($"❌ Error rate: {analysisResult.Results.Summary.ErrorRate:P2}");
                Console.WriteLine($"🔒 Security grade: {analysisResult.Results.Security.SecurityGrade}");
                Console.WriteLine($"⚡ Performance grade: {analysisResult.Results.Performance.PerformanceGrade}");
                
                if (!string.IsNullOrEmpty(analysisResult.HtmlPath))
                {
                    Console.WriteLine($"📄 HTML report generated: {analysisResult.HtmlPath}");
                }
                
                // Example 4: Execute a custom query
                Console.WriteLine("\n🔎 Executing custom query...");
                var queryResult = await client.QueryAsync(logFilePath, 
                    "SELECT ip, COUNT(*) as requests FROM logs WHERE status_code >= 400 GROUP BY ip ORDER BY requests DESC LIMIT 10");
                    
                Console.WriteLine("Top IPs with errors:");
                // Process queryResult.QueryResults as needed
                
                // Example 5: Run a preset
                if (presets.Presets.Length > 0)
                {
                    Console.WriteLine($"\n⚙️ Running preset: {presets.Presets[0].Name}");
                    var presetResult = await client.RunPresetAsync(logFilePath, presets.Presets[0].Name);
                    Console.WriteLine("✅ Preset execution completed");
                }
                
            }
            catch (Exception ex)
            {
                Console.WriteLine($"❌ Error: {ex.Message}");
            }
            
            Console.WriteLine("\n👋 Example completed");
        }
    }
}
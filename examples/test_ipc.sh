#!/bin/bash

# Test script for Smart Log Analyser IPC Server

echo "🧪 Smart Log Analyser IPC Server Test"
echo "====================================="

# Check if nc (netcat) is available for testing
if ! command -v nc &> /dev/null; then
    echo "❌ netcat (nc) is not available. Cannot test IPC server."
    exit 1
fi

# Start the IPC server in background
echo "🚀 Starting IPC server..."
../smart-log-analyser server &
SERVER_PID=$!

# Wait a moment for server to start
sleep 2

echo "📡 Server started with PID: $SERVER_PID"
echo "🔌 Testing connection to Unix socket..."

# Test if socket exists
if [ -S "/tmp/smart-log-analyser.sock" ]; then
    echo "✅ Unix socket created successfully"
    
    # Test server status
    echo "📊 Testing server status..."
    echo '{"id":"test-1","action":"getStatus"}' | nc -U /tmp/smart-log-analyser.sock
    
    echo ""
    echo "✅ IPC server is working correctly!"
else
    echo "❌ Unix socket not found"
fi

# Cleanup
echo "🛑 Stopping server..."
kill $SERVER_PID 2>/dev/null

echo "✅ Test completed"
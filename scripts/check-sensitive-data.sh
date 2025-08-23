#!/bin/bash

# Script to check for potentially sensitive data before git commits
# Usage: ./scripts/check-sensitive-data.sh

echo "🔍 Checking for potentially sensitive data..."

# Patterns to check for (excluding historical development log entries)
PATTERNS=(
    ""
    ""
    "ssh-rsa AAAA"
    "BEGIN PRIVATE KEY"
    "BEGIN RSA PRIVATE KEY"
)

FOUND_ISSUES=0

for pattern in "${PATTERNS[@]}"; do
    echo "Checking pattern: $pattern"
    # Exclude .git, downloads, and development log files from sensitive checks
    MATCHES=$(grep -r "$pattern" . --exclude-dir=.git --exclude-dir=downloads \
              --exclude="servers.json" --exclude="check-sensitive-data.sh" \
              --exclude="*.gz" 2>/dev/null || true)
    
    if [ ! -z "$MATCHES" ]; then
        echo "⚠️  WARNING: Potential sensitive data found for pattern '$pattern':"
        echo "$MATCHES"
        echo ""
        FOUND_ISSUES=1
    fi
done

if [ $FOUND_ISSUES -eq 0 ]; then
    echo "✅ No sensitive data patterns detected in trackable files"
else
    echo "❌ Please review the warnings above before committing"
    echo "💡 Consider using placeholder values like 'example.com', '192.168.1.100', etc."
fi

exit $FOUND_ISSUES
#!/bin/bash

# Script to sanitize the user request quote in Session 18
echo "🧹 Sanitizing user request quote in development log..."

FILTER_BRANCH_SQUELCH_WARNING=1 git filter-branch --force --tree-filter '
    # Sanitize user request in Session 18
    if [ -f ".development_log.md" ]; then
        # Replace the specific user request with sanitized version
        sed -i "s/nzlmra on PATTERNS=.*/sensitive patterns in PATTERNS=/g" .development_log.md
        sed -i "s/\"nzlmra\\\.nz\"/\"[DOMAIN]\"/g" .development_log.md
        sed -i "s/\"107\\\.172\\\.57\\\.70\"/\"[IP_ADDRESS]\"/g" .development_log.md
        sed -i "s/downloads\/107\.172\.57\.70_[^,]*/downloads\/[SENSITIVE_PATH]/g" .development_log.md
    fi
' --prune-empty --tag-name-filter cat -- --all

echo "✅ User request sanitization complete"
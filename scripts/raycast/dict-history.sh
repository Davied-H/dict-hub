#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Dict History
# @raycast.mode fullOutput

# Optional parameters:
# @raycast.icon 📜
# @raycast.packageName Dict-Hub

# Documentation:
# @raycast.description View search history from Dict-Hub
# @raycast.author Dict-Hub
# @raycast.authorURL https://github.com/dict-hub

# Configuration
API_BASE="${DICT_HUB_API:-http://localhost:8080}"
FRONTEND_URL="${DICT_HUB_FRONTEND:-http://localhost:3000}"

LIMIT=20

# Get history
RESPONSE=$(curl -s "${API_BASE}/api/v1/history?limit=${LIMIT}")

if [ -z "$RESPONSE" ]; then
    echo "Error: Could not connect to Dict-Hub API at ${API_BASE}"
    echo "Please make sure the Dict-Hub server is running."
    exit 1
fi

# Check if jq is available
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed."
    echo "Install with: brew install jq"
    exit 1
fi

# Parse the response
CODE=$(echo "$RESPONSE" | jq -r '.code // 0')

if [ "$CODE" != "0" ]; then
    MESSAGE=$(echo "$RESPONSE" | jq -r '.message // "Unknown error"')
    echo "Error: $MESSAGE"
    exit 1
fi

# Extract history
HISTORY=$(echo "$RESPONSE" | jq -r '.data.history // .data // []')
COUNT=$(echo "$HISTORY" | jq 'length')

if [ "$COUNT" == "0" ]; then
    echo "# 📜 Search History"
    echo ""
    echo "No search history yet."
    echo ""
    echo "Start searching words to build your history!"
    exit 0
fi

echo "# 📜 Search History (Recent ${COUNT} items)"
echo ""

# Display history
echo "$HISTORY" | jq -c '.[]' | while read -r item; do
    WORD=$(echo "$item" | jq -r '.word // .keyword // "Unknown"')
    TIMESTAMP=$(echo "$item" | jq -r '.created_at // .timestamp // ""')
    SEARCH_COUNT=$(echo "$item" | jq -r '.search_count // 1')
    
    # Format timestamp if available
    if [ -n "$TIMESTAMP" ] && [ "$TIMESTAMP" != "null" ]; then
        # Try to format the timestamp (works on macOS)
        FORMATTED_TIME=$(date -j -f "%Y-%m-%dT%H:%M:%S" "${TIMESTAMP%%.*}" "+%Y-%m-%d %H:%M" 2>/dev/null || echo "$TIMESTAMP")
    else
        FORMATTED_TIME=""
    fi
    
    if [ "$SEARCH_COUNT" -gt 1 ]; then
        echo "- **${WORD}** (x${SEARCH_COUNT}) ${FORMATTED_TIME}"
    else
        echo "- **${WORD}** ${FORMATTED_TIME}"
    fi
done

echo ""
echo "---"
echo "🔗 [View full history](${FRONTEND_URL}/history)"

#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Dict Suggest
# @raycast.mode fullOutput

# Optional parameters:
# @raycast.icon 🔍
# @raycast.argument1 { "type": "text", "placeholder": "Enter prefix to get suggestions" }
# @raycast.packageName Dict-Hub

# Documentation:
# @raycast.description Get word suggestions from Dict-Hub
# @raycast.author Dict-Hub
# @raycast.authorURL https://github.com/dict-hub

# Configuration
API_BASE="${DICT_HUB_API:-http://localhost:8080}"
FRONTEND_URL="${DICT_HUB_FRONTEND:-http://localhost:3000}"

QUERY="$1"
LIMIT=15

if [ -z "$QUERY" ]; then
    echo "Please enter a prefix to get suggestions"
    exit 1
fi

# Get suggestions
RESPONSE=$(curl -s "${API_BASE}/api/v1/search/suggest?q=${QUERY}&limit=${LIMIT}")

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

# Extract suggestions
SUGGESTIONS=$(echo "$RESPONSE" | jq -r '.data.suggestions // []')
COUNT=$(echo "$SUGGESTIONS" | jq 'length')

if [ "$COUNT" == "0" ]; then
    echo "No suggestions found for: $QUERY"
    exit 0
fi

echo "# 🔍 Suggestions for \"$QUERY\""
echo ""

# Display suggestions
INDEX=1
echo "$SUGGESTIONS" | jq -r '.[]' | while read -r word; do
    echo "${INDEX}. **${word}** - [Search](${FRONTEND_URL}/word/${word})"
    INDEX=$((INDEX + 1))
done

echo ""
echo "---"
echo "💡 Use **Dict Search** command to look up any word"

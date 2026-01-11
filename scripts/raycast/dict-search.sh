#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Dict Search
# @raycast.mode fullOutput

# Optional parameters:
# @raycast.icon 📖
# @raycast.argument1 { "type": "text", "placeholder": "Enter word to search" }
# @raycast.packageName Dict-Hub

# Documentation:
# @raycast.description Search word definition in Dict-Hub
# @raycast.author Dict-Hub
# @raycast.authorURL https://github.com/dict-hub

# Configuration
API_BASE="${DICT_HUB_API:-http://localhost:8080}"
FRONTEND_URL="${DICT_HUB_FRONTEND:-http://localhost:3000}"

WORD="$1"

if [ -z "$WORD" ]; then
    echo "Please enter a word to search"
    exit 1
fi

# Search the word
RESPONSE=$(curl -s "${API_BASE}/api/v1/search?word=${WORD}")

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

# Extract results
RESULTS=$(echo "$RESPONSE" | jq -r '.data.results // []')
COUNT=$(echo "$RESULTS" | jq 'length')

if [ "$COUNT" == "0" ]; then
    echo "No results found for: $WORD"
    echo ""
    echo "Try a different word or check your spelling."
    exit 0
fi

echo "# 📖 $WORD"
echo ""

# Display each dictionary result
echo "$RESULTS" | jq -c '.[]' | while read -r result; do
    DICT_NAME=$(echo "$result" | jq -r '.dict_name // "Unknown"')
    DEFINITION=$(echo "$result" | jq -r '.definition // ""')
    
    echo "## 📚 $DICT_NAME"
    echo ""
    
    # Strip HTML tags for cleaner display and limit length
    CLEAN_DEF=$(echo "$DEFINITION" | sed 's/<[^>]*>//g' | head -c 2000)
    echo "$CLEAN_DEF"
    echo ""
    echo "---"
    echo ""
done

echo ""
echo "🔗 Open in browser: ${FRONTEND_URL}/word/${WORD}"

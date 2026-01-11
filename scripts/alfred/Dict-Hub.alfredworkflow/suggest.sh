#!/bin/bash

# Alfred Script Filter for Dict-Hub word suggestions
# Outputs JSON in Alfred Script Filter format

# Configuration
API_BASE="${DICT_HUB_API:-http://localhost:8080}"
FRONTEND_URL="${DICT_HUB_FRONTEND:-http://localhost:3000}"

QUERY="$1"
LIMIT=10

# Function to output Alfred JSON
output_error() {
    cat << EOF
{
    "items": [
        {
            "title": "Error: $1",
            "subtitle": "$2",
            "valid": false,
            "icon": {
                "path": "icon.png"
            }
        }
    ]
}
EOF
}

# Check if jq is available
if ! command -v jq &> /dev/null; then
    output_error "jq is not installed" "Install with: brew install jq"
    exit 0
fi

# If no query, show help
if [ -z "$QUERY" ]; then
    cat << EOF
{
    "items": [
        {
            "title": "Type a word to search",
            "subtitle": "Dict-Hub will show suggestions as you type",
            "valid": false,
            "icon": {
                "path": "icon.png"
            }
        }
    ]
}
EOF
    exit 0
fi

# Get suggestions from API
RESPONSE=$(curl -s --max-time 5 "${API_BASE}/api/v1/search/suggest?q=${QUERY}&limit=${LIMIT}" 2>/dev/null)

if [ -z "$RESPONSE" ]; then
    output_error "Cannot connect to Dict-Hub" "Make sure the server is running at ${API_BASE}"
    exit 0
fi

# Parse response
CODE=$(echo "$RESPONSE" | jq -r '.code // 0')

if [ "$CODE" != "0" ]; then
    MESSAGE=$(echo "$RESPONSE" | jq -r '.message // "Unknown error"')
    output_error "API Error" "$MESSAGE"
    exit 0
fi

# Extract suggestions
SUGGESTIONS=$(echo "$RESPONSE" | jq -r '.data.suggestions // []')
COUNT=$(echo "$SUGGESTIONS" | jq 'length')

# If no suggestions, offer to search the query directly
if [ "$COUNT" == "0" ]; then
    cat << EOF
{
    "items": [
        {
            "title": "Search for \"${QUERY}\"",
            "subtitle": "Press Enter to look up this word",
            "arg": "${FRONTEND_URL}/word/${QUERY}",
            "valid": true,
            "icon": {
                "path": "icon.png"
            }
        }
    ]
}
EOF
    exit 0
fi

# Build Alfred items JSON
ITEMS=""
FIRST=true

while IFS= read -r word; do
    if [ "$FIRST" = true ]; then
        FIRST=false
    else
        ITEMS="${ITEMS},"
    fi
    
    # Escape special characters in word for JSON
    ESCAPED_WORD=$(echo "$word" | jq -R '.')
    
    ITEMS="${ITEMS}
        {
            \"title\": ${ESCAPED_WORD},
            \"subtitle\": \"Press Enter to look up this word\",
            \"arg\": \"${FRONTEND_URL}/word/${word}\",
            \"autocomplete\": ${ESCAPED_WORD},
            \"valid\": true,
            \"mods\": {
                \"cmd\": {
                    \"subtitle\": \"Copy word to clipboard\",
                    \"arg\": ${ESCAPED_WORD},
                    \"valid\": true
                }
            },
            \"icon\": {
                \"path\": \"icon.png\"
            }
        }"
done < <(echo "$SUGGESTIONS" | jq -r '.[]')

cat << EOF
{
    "items": [${ITEMS}
    ]
}
EOF

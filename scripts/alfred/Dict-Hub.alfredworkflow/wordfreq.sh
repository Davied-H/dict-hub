#!/bin/bash

# Alfred Script Filter for Dict-Hub word frequency
# Outputs JSON in Alfred Script Filter format

# Configuration
API_BASE="${DICT_HUB_API:-http://localhost:8080}"
FRONTEND_URL="${DICT_HUB_FRONTEND:-http://localhost:3000}"

LIMIT=20

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

# Get word frequency from API
RESPONSE=$(curl -s --max-time 5 "${API_BASE}/api/v1/wordfreq/top?limit=${LIMIT}" 2>/dev/null)

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

# Extract word frequency data
WORDFREQ=$(echo "$RESPONSE" | jq -r '.data.words // .data // []')
COUNT=$(echo "$WORDFREQ" | jq 'length')

# If no data
if [ "$COUNT" == "0" ]; then
    cat << EOF
{
    "items": [
        {
            "title": "No word frequency data",
            "subtitle": "Start searching words to track frequency",
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

# Build Alfred items JSON
ITEMS=""
FIRST=true
RANK=1

echo "$WORDFREQ" | jq -c '.[]' | while read -r item; do
    WORD=$(echo "$item" | jq -r '.word // .keyword // "Unknown"')
    FREQ=$(echo "$item" | jq -r '.frequency // .count // 0')
    
    SUBTITLE="#${RANK} • Searched ${FREQ} times"
    
    if [ "$FIRST" = true ]; then
        FIRST=false
    else
        echo ","
    fi
    
    # Escape special characters
    ESCAPED_WORD=$(echo "$WORD" | jq -R '.')
    ESCAPED_SUBTITLE=$(echo "$SUBTITLE" | jq -R '.')
    
    cat << ITEM
        {
            "title": ${ESCAPED_WORD},
            "subtitle": ${ESCAPED_SUBTITLE},
            "arg": "${FRONTEND_URL}/word/${WORD}",
            "autocomplete": ${ESCAPED_WORD},
            "valid": true,
            "mods": {
                "cmd": {
                    "subtitle": "Copy word to clipboard",
                    "arg": ${ESCAPED_WORD},
                    "valid": true
                }
            },
            "icon": {
                "path": "icon.png"
            }
        }
ITEM
    RANK=$((RANK + 1))
done | {
    ITEMS=$(cat)
    cat << EOF
{
    "items": [${ITEMS}
    ]
}
EOF
}

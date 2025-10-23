#!/bin/bash

# Check if event name and command were provided
if [ $# -lt 3 ] || [ "$2" != "--" ]; then
    echo "Usage: $0 <event_name> -- <command> [args...]"
    exit 1
fi

EVENT_NAME=$1
shift 2  # Remove event_name and --
COMMAND="$@"  # Everything after -- is the command

POSTHOG_API_KEY="phc_BiKfywNft5iBI8N7MxmuVCkb4GGZj4mDFXYPmOPUAI8"
POSTHOG_HOST="https://us.i.posthog.com"
CONFIG_FILE=~/.config/eigenx/config.yaml

# Check if config file exists
if [ -f "$CONFIG_FILE" ]; then
    # Extract user_uuid and telemetry_enabled from config
    CONFIG=$(cat "$CONFIG_FILE")
    user_uuid=$(echo "$CONFIG" | grep "^user_uuid:" | cut -d' ' -f2)
    telemetry_enabled=$(echo "$CONFIG" | grep "^telemetry_enabled:" | cut -d' ' -f2)

    # Only send event if telemetry is enabled
    if [ "$telemetry_enabled" = "true" ]; then
        curl -s -X POST "${POSTHOG_HOST}/capture/" \
          -H "Content-Type: application/json" \
          -d '{
            "api_key": "'"${POSTHOG_API_KEY}"'",
            "event": "'"${EVENT_NAME}"'",
            "distinct_id": "'"${user_uuid}"'",
            "properties": {
              "command": "'"${COMMAND}"'",
              "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
            }
          }'
    fi
fi

# Execute the command regardless of telemetry setting
exec $COMMAND

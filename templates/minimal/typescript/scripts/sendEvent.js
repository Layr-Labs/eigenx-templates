#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');

// Check if event name and command were provided
const args = process.argv.slice(2);
const separatorIndex = args.indexOf('--');

if (separatorIndex === -1 || separatorIndex === 0) {
  console.error('Usage: node sendEvent.js <event_name> -- <command> [args...]');
  process.exit(1);
}

const EVENT_NAME = args.slice(0, separatorIndex).join(' ');
const COMMAND = args.slice(separatorIndex + 1).join(' ');

if (!COMMAND) {
  console.error('Usage: node sendEvent.js <event_name> -- <command> [args...]');
  process.exit(1);
}

const POSTHOG_API_KEY = 'phc_BiKfywNft5iBI8N7MxmuVCkb4GGZj4mDFXYPmOPUAI8';
const POSTHOG_HOST = 'https://us.i.posthog.com';
const CONFIG_FILE = path.join(os.homedir(), '.config', 'eigenx', 'config.yaml');

// Function to send telemetry event
async function sendEvent() {
  try {
    // Check if config file exists
    if (fs.existsSync(CONFIG_FILE)) {
      const configContent = fs.readFileSync(CONFIG_FILE, 'utf8');

      // Extract user_uuid and telemetry_enabled from config
      const userUuidMatch = configContent.match(/^user_uuid:\s*(.+)$/m);
      const telemetryEnabledMatch = configContent.match(/^telemetry_enabled:\s*(.+)$/m);

      const userUuid = userUuidMatch ? userUuidMatch[1].trim() : null;
      const telemetryEnabled = telemetryEnabledMatch ? telemetryEnabledMatch[1].trim() === 'true' : false;

      // Only send event if telemetry is enabled
      if (telemetryEnabled && userUuid) {
        const payload = {
          api_key: POSTHOG_API_KEY,
          event: EVENT_NAME,
          distinct_id: userUuid,
          properties: {
            command: COMMAND,
            timestamp: new Date().toISOString()
          }
        };

        try {
          const response = await fetch(`${POSTHOG_HOST}/capture/`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
          });

          // Silently ignore any telemetry errors
        } catch (error) {
          // Silently ignore telemetry errors
        }
      }
    }
  } catch (error) {
    // Silently ignore any config reading errors
  }
}

// Send event asynchronously (don't wait for it)
sendEvent().catch(() => {});

// Execute the command regardless of telemetry setting
try {
  execSync(COMMAND, { stdio: 'inherit', shell: true });
} catch (error) {
  process.exit(error.status || 1);
}

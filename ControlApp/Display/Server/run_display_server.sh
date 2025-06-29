#!/bin/bash

# Path to your Python script
PYTHON_SCRIPT="display_server.py"

nmcli con up id Boxi

while true; do
    # Wait until wifi is connected
    until nmcli con up id Boxi; do
        echo "Failed to bring up the connection, retrying in 3 seconds..."
        sleep 3
    done

    # Wait until the endpoint is available
    until curl -s -o /dev/null -w "%{http_code}" http://192.168.4.1:8080/api/ping | grep -q "200"; do
        echo "Control app not available, retrying in 3 seconds..."
        sleep 3
    done

    echo "Control app is up. Starting display server..."
    python3 "$PYTHON_SCRIPT"

    echo "Display server ended. Rechecking endpoint..."
done

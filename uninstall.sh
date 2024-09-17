#!/bin/bash

# Check if prm is installed
if ! command -v prm &> /dev/null; then
    echo "prm is not installed on your system."
    exit 0
fi

# Output the location of the prm binary
PRM_PATH=$(command -v prm)
echo "Found prm at: $PRM_PATH"

# Run the 'prm purge -f' command to purge app data
echo "Purging prm data..."
prm purge -f

# Check if the purge was successful
if [[ $? -ne 0 ]]; then
    echo "Failed to purge prm data. Exiting without uninstalling."
    exit 1
fi

# Remove the prm binary file
echo "Removing prm binary..."
if [[ -w "$PRM_PATH" ]]; then
    rm "$PRM_PATH"
    echo "prm has been successfully uninstalled."
else
    echo "You do not have permission to remove prm. Try with sudo..."
    exit 1
fi

exit 0

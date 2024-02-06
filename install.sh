#!/bin/bash

REPO_URL="https://github.com/IkehAkinyemi/wc-util.git"

read -p "Enter the name for the binary (default: cwc): " BINARY_NAME
BINARY_NAME=$(BINARY_NAME:-cwc) # Use 'cwc' as default if input is empty

git clone "$REPO_URL" || { echo "Failed to clone repository"; exit 1; }
cd "$(basename "$REPO_URL" .git)" || { echo "Failed to enter directory"; exit 1; }

GOBIN=$(go env GOPATH)/bin
go install -o "$GOBIN/$BINARY_NAME" ./cmd || { echo "Build and install failed"; exit 1; }

echo "Installation completed successfully"
echo "Make sure $GOBIN is in your PATH."
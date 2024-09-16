#!/bin/bash

# Check if a release tag is provided as an argument
if [[ -n "$1" ]]; then
  RELEASE_TAG="$1"
else
    echo "Fetching the latest version..."
    LATEST_VERSION=$(curl -s https://api.github.com/repos/dhruv1397/prm/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [[ -z "$LATEST_VERSION" ]]; then
      echo "Failed to fetch the latest version. Please check your internet connection or the repository."
      exit 1
    fi
    echo "Latest version is $LATEST_VERSION"
    echo "$LATEST_VERSION"
  RELEASE_TAG=$LATEST_VERSION
fi

# Remove the 'v' prefix from the version tag, if present
VERSION=${RELEASE_TAG#v}

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Mapping architecture names to match GitHub release names
if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# Check if OS is supported
if [[ "$OS" != "linux" && "$OS" != "darwin" ]]; then
    echo "Unsupported OS: $OS"
    exit 1
fi

# Construct the download URL
REPO_URL="https://github.com/dhruv1397/prm/releases/download"
FILE_NAME="prm-${VERSION}-${OS}-${ARCH}"
DOWNLOAD_URL="${REPO_URL}/${RELEASE_TAG}/${FILE_NAME}"

# Download the binary file
echo "Downloading $FILE_NAME from $DOWNLOAD_URL..."
curl -L -o prm "$DOWNLOAD_URL"

# Check if the file was downloaded successfully
if [[ ! -f "prm" ]]; then
    echo "Failed to download the binary. Check if the file exists at the URL."
    exit 1
fi

# Make the binary executable
chmod +x prm

# Move the binary to /usr/local/bin (or ~/.local/bin if no sudo privileges)
if [[ -w /usr/local/bin ]]; then
    sudo mv prm /usr/local/bin/prm
else
    mkdir -p ~/.local/bin
    mv prm ~/.local/bin/prm
    export PATH=$PATH:~/.local/bin
    echo "Since we couldn't write to /usr/local/bin, we have moved the binary to ~/.local/bin."
    echo "Please add ~/.local/bin to your PATH. You can do this by running:"
    echo "  export PATH=\$PATH:~/.local/bin"
    echo "You may want to add the above line to your ~/.bashrc or ~/.zshrc file to make this change permanent."
fi

# Verify installation
if command -v prm &> /dev/null; then
    echo "prm has been successfully installed!"
else
    echo "prm installation failed. Please check your PATH."
    exit 1
fi

#!/bin/bash

# Get architecture and OS
ARCH=$(uname -m)
OS=$(uname -s)

# Check if the architecture is supported
if [[ $ARCH != "x86_64" && $ARCH != "arm64" ]]; then
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Check if the OS is supported
if [[ $OS != "Linux" && $OS != "Darwin" ]]; then
  echo "Unsupported operating system: $OS"
  exit 1
fi

# Determine the appropriate binary URL
if [[ $OS == "Linux" ]]; then
  URL="https://github.com/zcubbs/hotpot/releases/latest/download/Hotpot_Linux_$ARCH.tar.gz"
elif [[ $OS == "Darwin" ]]; then
  URL="https://github.com/zcubbs/hotpot/releases/latest/download/Hotpot_Darwin_$ARCH.tar.gz"
fi

# Get the file name from the URL
FILE=$(basename $URL)

echo "Installing $FILE for $OS ($ARCH)"

# Download the binary
curl -L -O $URL

# Unpack the binary
if [[ $FILE == *.tar.gz ]]; then
  tar -xzf $FILE
elif [[ $FILE == *.zip ]]; then
  unzip $FILE
fi

# The file that has been unpacked will usually be the binary itself.
BINARY="hotpot"

# Check if the file is executable
if [[ ! -x $BINARY ]]; then
  chmod +x $BINARY
fi

# Move the binary into the PATH, so it can be executed anywhere
sudo mv $BINARY /usr/local/bin/

# Verify if the binary is now in the PATH and executable
if ! command -v $BINARY &>/dev/null; then
  echo "Installation failed."
  exit 1
else
  echo "Installation succeeded."
fi

# Remove the downloaded file
rm $FILE

# Run the installed binary to verify functionality
$BINARY about

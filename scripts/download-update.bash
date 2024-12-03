#!/bin/bash

version=$1

if [ -z "$version" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

# Download the new version
curl "https://github.com/stormaaja/homeautomation/releases/download/v$version/data-store-v$version-linux-amd64" -LO
chmod u+x data-store-v$version-linux-amd64
rm data-store-latest
ln -s data-store-v$version-linux-amd64 data-store-latest
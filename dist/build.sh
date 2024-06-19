#!/bin/bash

# Check if the win directory exists, if not, create it
if [ ! -d "win" ]; then
    mkdir win
fi

# Check if the linux directory exists, if not, create it
if [ ! -d "linux" ]; then
    mkdir linux
fi

# Check if the macos directory exists, if not, create it
if [ ! -d "macos" ]; then
    mkdir macos
fi

# Cross-compile for Windows
export GOOS=windows
export GOARCH=amd64
go build -ldflags="-s -w" -o win/rune.exe

# Cross-compile for Linux
export GOOS=linux
export GOARCH=amd64
go build -ldflags="-s -w" -o linux/rune

# Cross-compile for macOS
export GOOS=darwin
export GOARCH=amd64
go build -ldflags="-s -w" -o macos/rune

echo "Cross-compilation completed."

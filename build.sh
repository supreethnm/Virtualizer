#!/bin/bash

# OSX
# OSX 64 bit
echo "Building for OSX 64 bit ..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o virtualizer_darwin_amd64 main/main.go
echo "Build completed for OSX 64 bit"
# OSX 32 bit
echo "Building for OSX 32 bit ..."
GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -o virtualizer_darwin_386 main/main.go
echo "Build completed for OSX 32 bit"

# Linux
# Linux 64 bit
echo "Building for Linux 64 bit ..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o virtualizer_linux_amd64 main/main.go
echo "Build completed for Linux 64 bit"
# Linux 32 bit
echo "Building for Linux 32 bit ..."
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o virtualizer_linux_386 main/main.go
echo "Build completed for Linux 32 bit"

# Windows
# Windows 64 bit
echo "Building for Windows 64 bit ..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o virtualizer_windows_amd64.exe main/main.go
echo "Build completed for Windows 64 bit"
# Windows 32 bit
echo "Building for Windows 32 bit ..."
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o virtualizer_windows_386.exe main/main.go
echo "Build completed for Windows 32 bit"

echo "Build completed"

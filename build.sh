#!/bin/bash

echo "🔨 Building GoStep..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed or not in PATH"
    echo "   Please install Go from https://golang.org/dl/"
    exit 1
fi

# Check for required build tools
echo "🔍 Checking build dependencies..."
if ! command -v gcc &> /dev/null; then
    echo "❌ Error: gcc is not installed. Please install build-essential:"
    echo "   sudo apt-get update && sudo apt-get install build-essential"
    exit 1
fi

if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    echo "❌ Error: MinGW-w64 is not installed. Required for Windows cross-compilation."
    echo "   Please install with: sudo apt-get update && sudo apt-get install gcc-mingw-w64"
    exit 1
fi

if ! command -v x86_64-w64-mingw32-windres &> /dev/null; then
    echo "❌ Error: windres is not installed. Required for Windows resource compilation."
    echo "   Please install with: sudo apt-get update && sudo apt-get install binutils-mingw-w64"
    exit 1
fi

# Verify icons exist
if [ ! -f "pkg/assets/icons/16.ico" ] || [ ! -f "pkg/assets/icons/32.ico" ] || [ ! -f "pkg/assets/icons/256.ico" ]; then
    echo "❌ Error: Icon files missing in pkg/assets/icons/"
    exit 1
fi

# Verify resource file exists
if [ ! -f "pkg/assets/windows/resource.rc" ]; then
    echo "❌ Error: Windows resource file missing at pkg/assets/windows/resource.rc"
    exit 1
fi

# Download dependencies
echo "📥 Downloading dependencies..."
go mod download
if [ $? -ne 0 ]; then
    echo "❌ Error: Failed to download dependencies"
    exit 1
fi

# Detect WSL environment
if grep -q -E "Microsoft|WSL" /proc/version; then
    echo "✓ WSL environment detected"
    
    echo "🔨 Compiling Windows resources..."
    # Change to the resource directory first
    cd pkg/assets/windows
    x86_64-w64-mingw32-windres resource.rc -O coff -o resource.syso
    if [ $? -ne 0 ]; then
        echo "❌ Error: Failed to compile Windows resources"
        cd ../../..
        exit 1
    fi
    # Move the resource file to the main package directory
    mv resource.syso ../../../cmd/step-recorder/
    cd ../../..
    
    echo "🔄 Cross-compiling for Windows..."
    GOOS=windows \
    GOARCH=amd64 \
    CGO_ENABLED=1 \
    CC=x86_64-w64-mingw32-gcc \
    CXX=x86_64-w64-mingw32-g++ \
    CGO_CFLAGS="-g -O2 -D_WIN32_WINNT=0x0A00 -DWINVER=0x0A00" \
    CGO_LDFLAGS="-lcomctl32 -luser32 -lgdi32 -lole32 -lshell32 -ladvapi32 -lmsimg32 -lopengl32 -lwinmm" \
    TAGS="windows" \
    go build -v -tags "windows" -trimpath -ldflags="-H windowsgui -extldflags '-static'" -o gostep.exe ./cmd/step-recorder

    if [ $? -ne 0 ]; then
        echo "❌ Error: Windows cross-compilation failed"
        rm -f cmd/step-recorder/resource.syso
        exit 1
    fi

    # Clean up the resource file
    rm -f cmd/step-recorder/resource.syso

    echo "✅ Build completed successfully!"
    echo "📦 Output: gostep.exe"
else
    echo "❌ Error: This script must be run from WSL for proper Windows compatibility"
    exit 1
fi 
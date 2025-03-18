
# Build for all supported platforms and architectures
build-all:
    @echo "Building for all platforms and architectures..."
    just build-linux-amd64
    just build-linux-arm64
    just build-darwin-amd64
    just build-darwin-arm64
    just build-windows-amd64
    just build-windows-arm64
    @echo "All builds completed"

# Linux builds
build-linux-amd64:
    @echo "Building for Linux/amd64..."
    GOOS=linux GOARCH=amd64 go build -o bin/renamedit-linux-amd64

build-linux-arm64:
    @echo "Building for Linux/arm64..."
    GOOS=linux GOARCH=arm64 go build -o bin/renamedit-linux-arm64

# macOS builds
build-darwin-amd64:
    @echo "Building for macOS/amd64..."
    GOOS=darwin GOARCH=amd64 go build -o bin/renamedit-darwin-amd64

build-darwin-arm64:
    @echo "Building for macOS/arm64..."
    GOOS=darwin GOARCH=arm64 go build -o bin/renamedit-darwin-arm64

# Windows builds
build-windows-amd64:
    @echo "Building for Windows/amd64..."
    GOOS=windows GOARCH=amd64 go build -o bin/renamedit-windows-amd64.exe

build-windows-arm64:
    @echo "Building for Windows/arm64..."
    GOOS=windows GOARCH=arm64 go build -o bin/renamedit-windows-arm64.exe

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf bin/*
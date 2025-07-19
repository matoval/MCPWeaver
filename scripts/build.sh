#!/bin/bash

# MCPWeaver Cross-Platform Build Script
# Builds MCPWeaver for multiple platforms with proper configuration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
VERSION=${VERSION:-"1.0.0"}
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_COMMIT=${BUILD_COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}
BUILD_DIR="./build/bin"
DIST_DIR="./dist"

# Build flags
LDFLAGS="-s -w -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.buildCommit=${BUILD_COMMIT}"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    if ! command_exists wails; then
        print_error "Wails CLI not found. Please install it first:"
        print_error "go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        exit 1
    fi
    
    if ! command_exists npm; then
        print_error "npm not found. Please install Node.js first."
        exit 1
    fi
    
    if ! command_exists go; then
        print_error "Go not found. Please install Go first."
        exit 1
    fi
    
    print_success "All prerequisites found"
}

# Function to setup build environment
setup_build_env() {
    print_status "Setting up build environment..."
    
    # Create build directories
    mkdir -p "$BUILD_DIR"
    mkdir -p "$DIST_DIR"
    
    # Install frontend dependencies
    print_status "Installing frontend dependencies..."
    cd frontend && npm install && cd ..
    
    # Update go dependencies
    print_status "Updating Go dependencies..."
    go mod tidy
    
    print_success "Build environment ready"
}

# Function to build for a specific platform
build_platform() {
    local platform=$1
    local arch=$2
    local output_name=$3
    
    print_status "Building for ${platform}/${arch}..."
    
    case $platform in
        "windows")
            wails build -platform "${platform}/${arch}" \
                -ldflags "$LDFLAGS" \
                -o "${BUILD_DIR}/${output_name}" \
                -clean
            ;;
        "darwin")
            wails build -platform "${platform}/${arch}" \
                -ldflags "$LDFLAGS" \
                -o "${BUILD_DIR}/${output_name}" \
                -clean
            ;;
        "linux")
            wails build -platform "${platform}/${arch}" \
                -ldflags "$LDFLAGS" \
                -o "${BUILD_DIR}/${output_name}" \
                -clean
            ;;
        *)
            print_error "Unknown platform: $platform"
            return 1
            ;;
    esac
    
    if [ $? -eq 0 ]; then
        print_success "Build completed for ${platform}/${arch}"
    else
        print_error "Build failed for ${platform}/${arch}"
        return 1
    fi
}

# Function to package builds
package_builds() {
    print_status "Packaging builds..."
    
    cd "$BUILD_DIR"
    
    # Package Windows builds
    if [ -f "MCPWeaver-windows-amd64.exe" ]; then
        print_status "Creating Windows package..."
        zip -r "../${DIST_DIR}/MCPWeaver-${VERSION}-windows-amd64.zip" "MCPWeaver-windows-amd64.exe"
    fi
    
    # Package macOS builds
    if [ -d "MCPWeaver-darwin-amd64.app" ]; then
        print_status "Creating macOS Intel package..."
        tar -czf "../${DIST_DIR}/MCPWeaver-${VERSION}-darwin-amd64.tar.gz" "MCPWeaver-darwin-amd64.app"
    fi
    
    if [ -d "MCPWeaver-darwin-arm64.app" ]; then
        print_status "Creating macOS ARM package..."
        tar -czf "../${DIST_DIR}/MCPWeaver-${VERSION}-darwin-arm64.tar.gz" "MCPWeaver-darwin-arm64.app"
    fi
    
    # Package Linux builds
    if [ -f "MCPWeaver-linux-amd64" ]; then
        print_status "Creating Linux package..."
        tar -czf "../${DIST_DIR}/MCPWeaver-${VERSION}-linux-amd64.tar.gz" "MCPWeaver-linux-amd64"
    fi
    
    cd ..
    print_success "Packaging completed"
}

# Function to create checksums
create_checksums() {
    print_status "Creating checksums..."
    
    cd "$DIST_DIR"
    
    # Create SHA256 checksums
    if command_exists sha256sum; then
        sha256sum *.zip *.tar.gz > "checksums.sha256" 2>/dev/null || true
    elif command_exists shasum; then
        shasum -a 256 *.zip *.tar.gz > "checksums.sha256" 2>/dev/null || true
    fi
    
    cd ..
    print_success "Checksums created"
}

# Function to display build summary
build_summary() {
    print_status "Build Summary:"
    echo "=================="
    echo "Version: $VERSION"
    echo "Build Date: $BUILD_DATE"
    echo "Build Commit: $BUILD_COMMIT"
    echo "=================="
    
    if [ -d "$DIST_DIR" ]; then
        echo "Distribution files:"
        ls -la "$DIST_DIR/"
    fi
}

# Main build function
main() {
    local platforms=()
    local all_platforms=false
    local clean_build=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --all)
                all_platforms=true
                shift
                ;;
            --windows)
                platforms+=("windows/amd64")
                shift
                ;;
            --macos)
                platforms+=("darwin/amd64" "darwin/arm64")
                shift
                ;;
            --linux)
                platforms+=("linux/amd64")
                shift
                ;;
            --clean)
                clean_build=true
                shift
                ;;
            --version)
                VERSION="$2"
                shift 2
                ;;
            --help)
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --all          Build for all platforms"
                echo "  --windows      Build for Windows (amd64)"
                echo "  --macos        Build for macOS (Intel and ARM)"
                echo "  --linux        Build for Linux (amd64)"
                echo "  --clean        Clean build directory before building"
                echo "  --version VER  Set version number (default: 1.0.0)"
                echo "  --help         Show this help message"
                echo ""
                echo "Examples:"
                echo "  $0 --all                    # Build for all platforms"
                echo "  $0 --windows --linux        # Build for Windows and Linux"
                echo "  $0 --macos --version 1.1.0  # Build macOS with version 1.1.0"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Set default platforms if none specified
    if [ ${#platforms[@]} -eq 0 ] && [ "$all_platforms" = false ]; then
        print_warning "No platforms specified. Building for current platform only."
        platforms+=("current")
    fi
    
    # Build all platforms if requested
    if [ "$all_platforms" = true ]; then
        platforms=("windows/amd64" "darwin/amd64" "darwin/arm64" "linux/amd64")
    fi
    
    print_status "Starting MCPWeaver build process..."
    
    # Clean build directory if requested
    if [ "$clean_build" = true ]; then
        print_status "Cleaning build directory..."
        rm -rf "$BUILD_DIR"
        rm -rf "$DIST_DIR"
    fi
    
    # Run build steps
    check_prerequisites
    setup_build_env
    
    # Build for each platform
    for platform in "${platforms[@]}"; do
        if [ "$platform" = "current" ]; then
            print_status "Building for current platform..."
            wails build -ldflags "$LDFLAGS" -clean
        else
            IFS='/' read -r os arch <<< "$platform"
            output_name="MCPWeaver-${os}-${arch}"
            if [ "$os" = "windows" ]; then
                output_name="${output_name}.exe"
            elif [ "$os" = "darwin" ]; then
                output_name="${output_name}.app"
            fi
            
            build_platform "$os" "$arch" "$output_name"
        fi
    done
    
    # Package and create checksums only if building for multiple platforms
    if [ ${#platforms[@]} -gt 1 ] || [ "$all_platforms" = true ]; then
        package_builds
        create_checksums
    fi
    
    build_summary
    print_success "Build process completed successfully!"
}

# Run main function with all arguments
main "$@"
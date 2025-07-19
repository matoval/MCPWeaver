#!/bin/bash

# MCPWeaver Universal Packaging Script
# Handles packaging for all supported platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Configuration
VERSION=${VERSION:-"1.0.0"}
PLATFORMS=()
SIGN_CODE=false
NOTARIZE_MACOS=false
CREATE_INSTALLER=false
OUTPUT_DIR="./dist"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --version)
            VERSION="$2"
            shift 2
            ;;
        --windows)
            PLATFORMS+=("windows")
            shift
            ;;
        --macos)
            PLATFORMS+=("macos")
            shift
            ;;
        --linux)
            PLATFORMS+=("linux")
            shift
            ;;
        --all)
            PLATFORMS=("windows" "macos" "linux")
            shift
            ;;
        --sign)
            SIGN_CODE=true
            shift
            ;;
        --notarize)
            NOTARIZE_MACOS=true
            shift
            ;;
        --installer)
            CREATE_INSTALLER=true
            shift
            ;;
        --output-dir)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        --help)
            echo "MCPWeaver Packaging Script"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --version VER     Set version number (default: 1.0.0)"
            echo "  --windows         Package for Windows"
            echo "  --macos           Package for macOS"
            echo "  --linux           Package for Linux"
            echo "  --all             Package for all platforms"
            echo "  --sign            Enable code signing"
            echo "  --notarize        Enable macOS notarization"
            echo "  --installer       Create platform-specific installers"
            echo "  --output-dir DIR  Output directory (default: ./dist)"
            echo "  --help            Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 --all --version 1.2.0"
            echo "  $0 --macos --sign --notarize"
            echo "  $0 --windows --installer"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Default to current platform if none specified
if [ ${#PLATFORMS[@]} -eq 0 ]; then
    case "$(uname -s)" in
        Darwin*)    PLATFORMS=("macos") ;;
        Linux*)     PLATFORMS=("linux") ;;
        MINGW*|CYGWIN*|MSYS*) PLATFORMS=("windows") ;;
        *)          print_error "Unsupported platform"; exit 1 ;;
    esac
    print_warning "No platform specified, defaulting to current platform: ${PLATFORMS[0]}"
fi

print_status "MCPWeaver Packaging Script"
print_status "Version: $VERSION"
print_status "Platforms: ${PLATFORMS[*]}"
print_status "Output Directory: $OUTPUT_DIR"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Function to package for Windows
package_windows() {
    print_status "Packaging for Windows..."
    
    local WINDOWS_BINARY="./build/bin/MCPWeaver-windows-amd64.exe"
    local PACKAGE_NAME="MCPWeaver-${VERSION}-windows-amd64"
    
    if [ ! -f "$WINDOWS_BINARY" ]; then
        print_error "Windows binary not found at $WINDOWS_BINARY"
        print_error "Please build for Windows first: wails build -platform windows/amd64"
        return 1
    fi
    
    # Create ZIP package
    print_status "Creating Windows ZIP package..."
    local ZIP_DIR="${OUTPUT_DIR}/${PACKAGE_NAME}"
    mkdir -p "$ZIP_DIR"
    
    cp "$WINDOWS_BINARY" "${ZIP_DIR}/MCPWeaver.exe"
    [ -f "./LICENSE" ] && cp "./LICENSE" "${ZIP_DIR}/"
    [ -f "./README.md" ] && cp "./README.md" "${ZIP_DIR}/"
    
    # Create a simple batch launcher
    cat > "${ZIP_DIR}/MCPWeaver.bat" << 'EOF'
@echo off
cd /d "%~dp0"
start "" "MCPWeaver.exe"
EOF
    
    cd "$OUTPUT_DIR"
    zip -r "${PACKAGE_NAME}.zip" "$(basename "$ZIP_DIR")"
    rm -rf "$(basename "$ZIP_DIR")"
    cd - > /dev/null
    
    # Create NSIS installer if requested
    if [ "$CREATE_INSTALLER" = true ]; then
        print_status "Creating Windows installer..."
        if command -v makensis >/dev/null 2>&1; then
            # Copy the binary to the expected location for NSIS
            mkdir -p "./build/bin"
            cp "$WINDOWS_BINARY" "./build/bin/"
            
            # Run NSIS to create installer
            makensis -DVERSION="$VERSION" ./build/windows/installer/install.nsi
            
            # Move installer to output directory
            if [ -f "MCPWeaver-*.exe" ]; then
                mv MCPWeaver-*.exe "$OUTPUT_DIR/"
            fi
        else
            print_warning "NSIS not found, skipping installer creation"
        fi
    fi
    
    print_success "Windows packaging completed"
}

# Function to package for macOS
package_macos() {
    print_status "Packaging for macOS..."
    
    local MACOS_APP="./build/bin/MCPWeaver.app"
    local PACKAGE_NAME="MCPWeaver-${VERSION}-macos"
    
    if [ ! -d "$MACOS_APP" ]; then
        print_error "macOS app bundle not found at $MACOS_APP"
        print_error "Please build for macOS first: wails build -platform darwin/universal"
        return 1
    fi
    
    # Code signing
    if [ "$SIGN_CODE" = true ]; then
        print_status "Code signing macOS application..."
        if [ -n "${APPLE_DEVELOPER_ID:-}" ]; then
            codesign --force --deep --sign "$APPLE_DEVELOPER_ID" --options runtime "$MACOS_APP"
            print_success "Code signing completed"
        else
            print_warning "APPLE_DEVELOPER_ID not set, skipping code signing"
        fi
    fi
    
    # Create tar.gz package
    print_status "Creating macOS tar.gz package..."
    cd "./build/bin"
    tar -czf "../../${OUTPUT_DIR}/${PACKAGE_NAME}.tar.gz" "$(basename "$MACOS_APP")"
    cd - > /dev/null
    
    # Create DMG if on macOS and create-dmg is available
    if [ "$(uname -s)" = "Darwin" ] && [ "$CREATE_INSTALLER" = true ]; then
        print_status "Creating macOS DMG..."
        if [ -x "./scripts/create-dmg.sh" ]; then
            ./scripts/create-dmg.sh "$VERSION"
            if [ -f "MCPWeaver-${VERSION}.dmg" ]; then
                mv "MCPWeaver-${VERSION}.dmg" "$OUTPUT_DIR/"
                [ -f "MCPWeaver-${VERSION}.dmg.sha256" ] && mv "MCPWeaver-${VERSION}.dmg.sha256" "$OUTPUT_DIR/"
            fi
        else
            print_warning "DMG creation script not found or not executable"
        fi
    fi
    
    # Notarization
    if [ "$NOTARIZE_MACOS" = true ]; then
        print_status "Submitting for notarization..."
        if [ -n "${APPLE_ID:-}" ] && [ -n "${APPLE_PASSWORD:-}" ]; then
            # Create a zip for notarization
            cd "./build/bin"
            zip -r "../../${OUTPUT_DIR}/${PACKAGE_NAME}-notarization.zip" "$(basename "$MACOS_APP")"
            cd - > /dev/null
            
            # Submit for notarization
            xcrun altool --notarize-app \
                --primary-bundle-id "com.mcpweaver.app" \
                --username "$APPLE_ID" \
                --password "$APPLE_PASSWORD" \
                --file "${OUTPUT_DIR}/${PACKAGE_NAME}-notarization.zip"
            
            print_status "Notarization submitted (this is an async process)"
        else
            print_warning "APPLE_ID or APPLE_PASSWORD not set, skipping notarization"
        fi
    fi
    
    print_success "macOS packaging completed"
}

# Function to package for Linux
package_linux() {
    print_status "Packaging for Linux..."
    
    local LINUX_BINARY="./build/bin/MCPWeaver-linux-amd64"
    local PACKAGE_NAME="MCPWeaver-${VERSION}-linux-amd64"
    
    if [ ! -f "$LINUX_BINARY" ]; then
        print_error "Linux binary not found at $LINUX_BINARY"
        print_error "Please build for Linux first: wails build -platform linux/amd64"
        return 1
    fi
    
    # Create tar.gz package
    print_status "Creating Linux tar.gz package..."
    local TAR_DIR="${OUTPUT_DIR}/${PACKAGE_NAME}"
    mkdir -p "$TAR_DIR"
    
    cp "$LINUX_BINARY" "${TAR_DIR}/MCPWeaver"
    chmod +x "${TAR_DIR}/MCPWeaver"
    [ -f "./LICENSE" ] && cp "./LICENSE" "${TAR_DIR}/"
    [ -f "./README.md" ] && cp "./README.md" "${TAR_DIR}/"
    
    # Create a simple launcher script
    cat > "${TAR_DIR}/mcpweaver.sh" << 'EOF'
#!/bin/bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec "$DIR/MCPWeaver" "$@"
EOF
    chmod +x "${TAR_DIR}/mcpweaver.sh"
    
    cd "$OUTPUT_DIR"
    tar -czf "${PACKAGE_NAME}.tar.gz" "$(basename "$TAR_DIR")"
    rm -rf "$(basename "$TAR_DIR")"
    cd - > /dev/null
    
    # Create AppImage if requested
    if [ "$CREATE_INSTALLER" = true ]; then
        print_status "Creating Linux AppImage..."
        if [ -x "./scripts/create-appimage.sh" ]; then
            ./scripts/create-appimage.sh "$VERSION"
            if [ -f "MCPWeaver-${VERSION}-x86_64.AppImage" ]; then
                mv "MCPWeaver-${VERSION}-x86_64.AppImage" "$OUTPUT_DIR/"
                [ -f "MCPWeaver-${VERSION}-x86_64.AppImage.sha256" ] && mv "MCPWeaver-${VERSION}-x86_64.AppImage.sha256" "$OUTPUT_DIR/"
            fi
        else
            print_warning "AppImage creation script not found or not executable"
        fi
    fi
    
    print_success "Linux packaging completed"
}

# Function to create checksums
create_checksums() {
    print_status "Creating checksums for all packages..."
    
    cd "$OUTPUT_DIR"
    
    # Remove any existing checksums file
    rm -f checksums.sha256
    
    # Create checksums for all package files
    for file in *.zip *.tar.gz *.dmg *.AppImage *.exe; do
        if [ -f "$file" ]; then
            if command -v sha256sum >/dev/null 2>&1; then
                sha256sum "$file" >> checksums.sha256
            elif command -v shasum >/dev/null 2>&1; then
                shasum -a 256 "$file" >> checksums.sha256
            fi
        fi
    done
    
    cd - > /dev/null
    
    if [ -f "${OUTPUT_DIR}/checksums.sha256" ]; then
        print_success "Checksums created in ${OUTPUT_DIR}/checksums.sha256"
    fi
}

# Function to display summary
display_summary() {
    print_status "Packaging Summary"
    echo "=================="
    echo "Version: $VERSION"
    echo "Output Directory: $OUTPUT_DIR"
    echo ""
    
    if [ -d "$OUTPUT_DIR" ]; then
        echo "Created packages:"
        ls -la "$OUTPUT_DIR/" | grep -E '\.(zip|tar\.gz|dmg|AppImage|exe)$' || echo "No packages found"
        echo ""
        
        if [ -f "${OUTPUT_DIR}/checksums.sha256" ]; then
            echo "Checksums:"
            cat "${OUTPUT_DIR}/checksums.sha256"
        fi
    fi
}

# Main execution
print_status "Starting packaging process..."

# Package for each requested platform
for platform in "${PLATFORMS[@]}"; do
    case $platform in
        "windows")
            package_windows
            ;;
        "macos")
            package_macos
            ;;
        "linux")
            package_linux
            ;;
        *)
            print_error "Unknown platform: $platform"
            ;;
    esac
done

# Create checksums and display summary
create_checksums
display_summary

print_success "All packaging completed successfully!"
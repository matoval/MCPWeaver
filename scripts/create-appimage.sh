#!/bin/bash

# MCPWeaver Linux AppImage Creation Script
# Creates a portable AppImage for Linux distribution

set -e

# Configuration
APP_NAME="MCPWeaver"
APP_BINARY="./build/bin/${APP_NAME}-linux-amd64"
APPDIR_NAME="${APP_NAME}.AppDir"
VERSION=${1:-"1.0.0"}
ARCH="x86_64"

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

# Check dependencies
check_dependencies() {
    print_status "Checking dependencies..."
    
    if ! command -v wget >/dev/null 2>&1; then
        print_error "wget is required but not installed."
        exit 1
    fi
    
    if ! command -v file >/dev/null 2>&1; then
        print_error "file is required but not installed."
        exit 1
    fi
    
    print_success "All dependencies found"
}

# Download appimagetool if not present
download_appimagetool() {
    if [ ! -f "appimagetool-x86_64.AppImage" ]; then
        print_status "Downloading appimagetool..."
        wget -q "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage"
        chmod +x appimagetool-x86_64.AppImage
        print_success "appimagetool downloaded"
    fi
}

# Check if binary exists
if [ ! -f "$APP_BINARY" ]; then
    print_error "Binary not found at $APP_BINARY"
    print_error "Please build the application first with: wails build -platform linux/amd64"
    exit 1
fi

print_status "Creating AppImage for ${APP_NAME} version ${VERSION}..."

# Clean up any existing AppDir
rm -rf "$APPDIR_NAME"

# Check dependencies and download tools
check_dependencies
download_appimagetool

# Create AppDir structure
print_status "Creating AppDir structure..."
mkdir -p "$APPDIR_NAME/usr/bin"
mkdir -p "$APPDIR_NAME/usr/share/applications"
mkdir -p "$APPDIR_NAME/usr/share/icons/hicolor/256x256/apps"
mkdir -p "$APPDIR_NAME/usr/share/icons/hicolor/scalable/apps"
mkdir -p "$APPDIR_NAME/usr/share/pixmaps"
mkdir -p "$APPDIR_NAME/usr/share/doc/${APP_NAME}"

# Copy the binary
print_status "Copying application binary..."
cp "$APP_BINARY" "$APPDIR_NAME/usr/bin/${APP_NAME}"
chmod +x "$APPDIR_NAME/usr/bin/${APP_NAME}"

# Create desktop file
print_status "Creating desktop file..."
cat > "$APPDIR_NAME/usr/share/applications/${APP_NAME}.desktop" << EOF
[Desktop Entry]
Type=Application
Name=${APP_NAME}
Comment=Transform OpenAPI specifications into Model Context Protocol servers
Exec=${APP_NAME} %F
Icon=${APP_NAME}
Categories=Development;IDE;Programming;
StartupNotify=true
NoDisplay=false
MimeType=application/x-openapi;application/json;application/yaml;text/yaml;
Keywords=OpenAPI;MCP;API;Development;Tools;Generator;
GenericName=API Transformer
StartupWMClass=${APP_NAME}
EOF

# Copy desktop file to AppDir root
cp "$APPDIR_NAME/usr/share/applications/${APP_NAME}.desktop" "$APPDIR_NAME/"

# Create/copy icon
print_status "Setting up application icon..."
if [ -f "./build/appicon.png" ]; then
    cp "./build/appicon.png" "$APPDIR_NAME/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png"
    cp "./build/appicon.png" "$APPDIR_NAME/usr/share/pixmaps/${APP_NAME}.png"
    cp "./build/appicon.png" "$APPDIR_NAME/${APP_NAME}.png"
else
    print_warning "App icon not found, creating a simple text-based icon..."
    # Create a simple PNG icon using ImageMagick if available
    if command -v convert >/dev/null 2>&1; then
        convert -size 256x256 xc:lightblue -gravity center -pointsize 48 -fill darkblue \
            -annotate +0+0 "MCP" "$APPDIR_NAME/${APP_NAME}.png"
        cp "$APPDIR_NAME/${APP_NAME}.png" "$APPDIR_NAME/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png"
        cp "$APPDIR_NAME/${APP_NAME}.png" "$APPDIR_NAME/usr/share/pixmaps/${APP_NAME}.png"
    else
        print_warning "ImageMagick not available, skipping icon creation"
    fi
fi

# Create AppRun script
print_status "Creating AppRun script..."
cat > "$APPDIR_NAME/AppRun" << 'EOF'
#!/bin/bash
HERE="$(dirname "$(readlink -f "${0}")")"
export APPDIR="$HERE"
export PATH="${HERE}/usr/bin/:$PATH"
export LD_LIBRARY_PATH="${HERE}/usr/lib/:$LD_LIBRARY_PATH"
export XDG_DATA_DIRS="${HERE}/usr/share/:$XDG_DATA_DIRS"

# Set up environment for the application
export APPIMAGE_RUNTIME_DIR="$HERE"

# Run the application
exec "${HERE}/usr/bin/MCPWeaver" "$@"
EOF

chmod +x "$APPDIR_NAME/AppRun"

# Copy documentation
print_status "Copying documentation..."
if [ -f "./LICENSE" ]; then
    cp "./LICENSE" "$APPDIR_NAME/usr/share/doc/${APP_NAME}/"
fi
if [ -f "./README.md" ]; then
    cp "./README.md" "$APPDIR_NAME/usr/share/doc/${APP_NAME}/"
fi

# Create AppStream metadata
print_status "Creating AppStream metadata..."
mkdir -p "$APPDIR_NAME/usr/share/metainfo"
cat > "$APPDIR_NAME/usr/share/metainfo/${APP_NAME}.appdata.xml" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<component type="desktop-application">
    <id>com.mcpweaver.MCPWeaver</id>
    <metadata_license>CC0-1.0</metadata_license>
    <project_license>AGPL-3.0</project_license>
    <name>MCPWeaver</name>
    <summary>Transform OpenAPI specifications into Model Context Protocol servers</summary>
    <description>
        <p>
            MCPWeaver is a powerful desktop application that transforms OpenAPI specifications 
            into Model Context Protocol (MCP) servers. It provides an intuitive interface for 
            developers to convert their API specifications into MCP-compatible servers.
        </p>
        <p>Features:</p>
        <ul>
            <li>Import OpenAPI specifications from files or URLs</li>
            <li>Generate MCP server code with customizable templates</li>
            <li>Validate and test generated servers</li>
            <li>Cross-platform support for Windows, macOS, and Linux</li>
            <li>Template marketplace for community templates</li>
        </ul>
    </description>
    <launchable type="desktop-id">MCPWeaver.desktop</launchable>
    <provides>
        <binary>MCPWeaver</binary>
    </provides>
    <screenshots>
        <screenshot type="default">
            <caption>Main application window</caption>
        </screenshot>
    </screenshots>
    <url type="homepage">https://github.com/matoval/MCPWeaver</url>
    <url type="bugtracker">https://github.com/matoval/MCPWeaver/issues</url>
    <url type="help">https://github.com/matoval/MCPWeaver/wiki</url>
    <developer_name>MCPWeaver Team</developer_name>
    <releases>
        <release version="$VERSION" date="$(date +%Y-%m-%d)">
            <description>
                <p>Latest release of MCPWeaver</p>
            </description>
        </release>
    </releases>
    <content_rating type="oars-1.1" />
</component>
EOF

# Verify the binary works
print_status "Verifying binary..."
if ! file "$APPDIR_NAME/usr/bin/${APP_NAME}" | grep -q "ELF 64-bit"; then
    print_error "Binary verification failed - not a valid 64-bit ELF executable"
    exit 1
fi

# Check for required libraries
print_status "Checking library dependencies..."
LDD_OUTPUT=$(ldd "$APPDIR_NAME/usr/bin/${APP_NAME}" 2>/dev/null || echo "ldd failed")
if echo "$LDD_OUTPUT" | grep -q "not found"; then
    print_warning "Some libraries are missing:"
    echo "$LDD_OUTPUT" | grep "not found"
    print_warning "The AppImage may not work on systems without these libraries"
fi

# Create version info
echo "$VERSION" > "$APPDIR_NAME/VERSION"

# Build the AppImage
print_status "Building AppImage..."
export ARCH="$ARCH"
./appimagetool-x86_64.AppImage "$APPDIR_NAME" "${APP_NAME}-${VERSION}-${ARCH}.AppImage"

if [ $? -eq 0 ]; then
    # Calculate size and checksum
    APPIMAGE_FILE="${APP_NAME}-${VERSION}-${ARCH}.AppImage"
    APPIMAGE_SIZE=$(du -h "$APPIMAGE_FILE" | cut -f1)
    APPIMAGE_CHECKSUM=$(sha256sum "$APPIMAGE_FILE" | cut -d' ' -f1)
    
    print_success "AppImage created successfully!"
    print_success "File: $APPIMAGE_FILE"
    print_success "Size: $APPIMAGE_SIZE"
    print_success "SHA256: $APPIMAGE_CHECKSUM"
    
    # Create checksums file
    echo "$APPIMAGE_CHECKSUM  $APPIMAGE_FILE" > "${APPIMAGE_FILE}.sha256"
    
    # Make AppImage executable
    chmod +x "$APPIMAGE_FILE"
    
    # Clean up
    if [ "$2" != "--keep-appdir" ]; then
        rm -rf "$APPDIR_NAME"
        print_status "Cleaned up temporary AppDir"
    fi
    
    print_status "AppImage creation completed successfully!"
    print_status "You can now distribute $APPIMAGE_FILE"
    
    # Test the AppImage
    if [ "$2" = "--test" ]; then
        print_status "Testing AppImage..."
        timeout 10s "./$APPIMAGE_FILE" --version || print_warning "AppImage test timed out or failed"
    fi
else
    print_error "AppImage creation failed!"
    exit 1
fi
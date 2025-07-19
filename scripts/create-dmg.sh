#!/bin/bash

# MCPWeaver macOS DMG Creation Script
# Creates a professional DMG installer for macOS

set -e

# Configuration
APP_NAME="MCPWeaver"
APP_PATH="./build/bin/${APP_NAME}.app"
DMG_NAME="${APP_NAME}"
DMG_BACKGROUND="./build/darwin/dmg-background.png"
DMG_ICON="./build/appicon.icns"
DMG_WINDOW_SIZE="600x400"
DMG_ICON_SIZE=100
DMG_TEXT_SIZE=16

# Version and build info
VERSION=${1:-"1.0.0"}
VOLUME_NAME="${APP_NAME} ${VERSION}"
DMG_TEMP_NAME="${APP_NAME}-temp.dmg"
DMG_FINAL_NAME="${APP_NAME}-${VERSION}.dmg"

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

# Check if app bundle exists
if [ ! -d "$APP_PATH" ]; then
    print_error "App bundle not found at $APP_PATH"
    print_error "Please build the application first with: wails build"
    exit 1
fi

print_status "Creating DMG for ${APP_NAME} version ${VERSION}..."

# Clean up any existing DMG files
rm -f "${DMG_TEMP_NAME}" "${DMG_FINAL_NAME}"

# Create a temporary DMG
print_status "Creating temporary DMG..."
SIZE=$(du -sh "$APP_PATH" | sed 's/\([0-9\.]*\)M.*/\1/')
SIZE=$(echo "${SIZE} + 50" | bc)

hdiutil create -srcfolder "$APP_PATH" -volname "$VOLUME_NAME" -fs HFS+ \
    -fsargs "-c c=64,a=16,e=16" -format UDRW -size ${SIZE}M "$DMG_TEMP_NAME"

# Mount the DMG
print_status "Mounting temporary DMG..."
DEVICE=$(hdiutil attach -readwrite -noverify "$DMG_TEMP_NAME" | \
    egrep '^/dev/' | sed 1q | awk '{print $1}')

sleep 2

# Get the volume path
VOLUME_PATH="/Volumes/$VOLUME_NAME"

# Create Applications shortcut
print_status "Creating Applications shortcut..."
ln -sf /Applications "$VOLUME_PATH/Applications"

# Create a custom background if it doesn't exist
if [ ! -f "$DMG_BACKGROUND" ]; then
    print_warning "DMG background not found, creating a simple one..."
    mkdir -p "$(dirname "$DMG_BACKGROUND")"
    
    # Create a simple gradient background using sips (if available)
    if command -v sips > /dev/null; then
        sips -s format png --out "$DMG_BACKGROUND" /System/Library/Desktop\ Pictures/Solid\ Colors/Space\ Gray.png 2>/dev/null || \
        cp /System/Library/Desktop\ Pictures/Solid\ Colors/Stone.heic "$DMG_BACKGROUND" 2>/dev/null || \
        echo "Could not create background image"
    fi
fi

# Set up the DMG appearance using AppleScript
print_status "Configuring DMG appearance..."
cat > dmg_setup.applescript << EOF
tell application "Finder"
    tell disk "$VOLUME_NAME"
        open
        set current view of container window to icon view
        set toolbar visible of container window to false
        set statusbar visible of container window to false
        set the bounds of container window to {100, 100, 700, 500}
        set viewOptions to the icon view options of container window
        set arrangement of viewOptions to not arranged
        set icon size of viewOptions to $DMG_ICON_SIZE
        set background picture of viewOptions to file ".background:dmg-background.png"
        
        -- Position icons
        set position of item "$APP_NAME.app" of container window to {150, 200}
        set position of item "Applications" of container window to {450, 200}
        
        -- Set icon properties
        if exists file ".VolumeIcon.icns" then
            set icon of disk "$VOLUME_NAME" to file ".VolumeIcon.icns"
        end if
        
        close
        open
        update without registering applications
        delay 2
    end tell
end tell
EOF

# Copy background image if it exists
if [ -f "$DMG_BACKGROUND" ]; then
    mkdir -p "$VOLUME_PATH/.background"
    cp "$DMG_BACKGROUND" "$VOLUME_PATH/.background/dmg-background.png"
fi

# Copy volume icon if it exists
if [ -f "$DMG_ICON" ]; then
    cp "$DMG_ICON" "$VOLUME_PATH/.VolumeIcon.icns"
    SetFile -c icnC "$VOLUME_PATH/.VolumeIcon.icns"
fi

# Execute the AppleScript
osascript dmg_setup.applescript

# Clean up
rm dmg_setup.applescript

# Set custom attributes
print_status "Setting DMG attributes..."
SetFile -a C "$VOLUME_PATH"

# Make .background and .VolumeIcon.icns hidden
if [ -d "$VOLUME_PATH/.background" ]; then
    SetFile -a V "$VOLUME_PATH/.background"
fi
if [ -f "$VOLUME_PATH/.VolumeIcon.icns" ]; then
    SetFile -a V "$VOLUME_PATH/.VolumeIcon.icns"
fi

# Force a sync
sync

# Unmount the DMG
print_status "Unmounting temporary DMG..."
hdiutil detach "$DEVICE"

# Convert to final compressed DMG
print_status "Creating final compressed DMG..."
hdiutil convert "$DMG_TEMP_NAME" -format UDZO -imagekey zlib-level=9 -o "$DMG_FINAL_NAME"

# Clean up temporary DMG
rm -f "$DMG_TEMP_NAME"

# Verify the DMG
print_status "Verifying DMG..."
hdiutil verify "$DMG_FINAL_NAME"

# Calculate size and checksum
DMG_SIZE=$(du -h "$DMG_FINAL_NAME" | cut -f1)
DMG_CHECKSUM=$(shasum -a 256 "$DMG_FINAL_NAME" | cut -d' ' -f1)

print_success "DMG created successfully!"
print_success "File: $DMG_FINAL_NAME"
print_success "Size: $DMG_SIZE"
print_success "SHA256: $DMG_CHECKSUM"

# Create a checksums file
echo "$DMG_CHECKSUM  $DMG_FINAL_NAME" > "${DMG_FINAL_NAME}.sha256"

print_status "DMG creation completed. Test the DMG before distribution!"

# Optional: Open the DMG for inspection
if [ "$2" = "--open" ]; then
    open "$DMG_FINAL_NAME"
fi
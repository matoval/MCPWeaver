#!/bin/bash

# MCPWeaver Code Signing Script
# Handles code signing for macOS and Windows builds

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
PLATFORM=""
BINARY_PATH=""
SIGN_IDENTITY=""
ENTITLEMENTS_PATH=""
KEYCHAIN_PATH=""
KEYCHAIN_PASSWORD=""
CERTIFICATE_PATH=""
CERTIFICATE_PASSWORD=""
TIMESTAMP_URL=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --platform)
            PLATFORM="$2"
            shift 2
            ;;
        --binary)
            BINARY_PATH="$2"
            shift 2
            ;;
        --identity)
            SIGN_IDENTITY="$2"
            shift 2
            ;;
        --entitlements)
            ENTITLEMENTS_PATH="$2"
            shift 2
            ;;
        --keychain)
            KEYCHAIN_PATH="$2"
            shift 2
            ;;
        --keychain-password)
            KEYCHAIN_PASSWORD="$2"
            shift 2
            ;;
        --certificate)
            CERTIFICATE_PATH="$2"
            shift 2
            ;;
        --certificate-password)
            CERTIFICATE_PASSWORD="$2"
            shift 2
            ;;
        --timestamp-url)
            TIMESTAMP_URL="$2"
            shift 2
            ;;
        --help)
            echo "MCPWeaver Code Signing Script"
            echo ""
            echo "Usage: $0 --platform PLATFORM --binary PATH [OPTIONS]"
            echo ""
            echo "Required:"
            echo "  --platform PLATFORM     Platform to sign for (macos|windows)"
            echo "  --binary PATH           Path to binary/app to sign"
            echo ""
            echo "macOS Options:"
            echo "  --identity IDENTITY     Code signing identity"
            echo "  --entitlements PATH     Path to entitlements plist"
            echo "  --keychain PATH         Path to keychain"
            echo "  --keychain-password PWD Keychain password"
            echo ""
            echo "Windows Options:"
            echo "  --certificate PATH      Path to certificate file (.pfx/.p12)"
            echo "  --certificate-password PWD Certificate password"
            echo "  --timestamp-url URL     Timestamp server URL"
            echo ""
            echo "Environment Variables (alternatives to command line):"
            echo "  APPLE_DEVELOPER_ID      macOS signing identity"
            echo "  APPLE_KEYCHAIN_PATH     macOS keychain path"
            echo "  APPLE_KEYCHAIN_PASSWORD macOS keychain password"
            echo "  WINDOWS_CERTIFICATE     Windows certificate path"
            echo "  WINDOWS_CERT_PASSWORD   Windows certificate password"
            echo "  TIMESTAMP_SERVER        Windows timestamp server"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Validate required parameters
if [ -z "$PLATFORM" ] || [ -z "$BINARY_PATH" ]; then
    print_error "Platform and binary path are required"
    echo "Use --help for usage information"
    exit 1
fi

# Check if binary exists
if [ ! -e "$BINARY_PATH" ]; then
    print_error "Binary not found: $BINARY_PATH"
    exit 1
fi

print_status "Code signing for $PLATFORM: $BINARY_PATH"

# Function to sign macOS application
sign_macos() {
    print_status "Signing macOS application..."
    
    # Use environment variables if command line options not provided
    local identity="${SIGN_IDENTITY:-${APPLE_DEVELOPER_ID:-}}"
    local keychain="${KEYCHAIN_PATH:-${APPLE_KEYCHAIN_PATH:-}}"
    local keychain_pwd="${KEYCHAIN_PASSWORD:-${APPLE_KEYCHAIN_PASSWORD:-}}"
    local entitlements="${ENTITLEMENTS_PATH:-./build/darwin/entitlements.plist}"
    
    if [ -z "$identity" ]; then
        print_error "No signing identity provided"
        print_error "Set APPLE_DEVELOPER_ID environment variable or use --identity"
        return 1
    fi
    
    # Setup keychain if provided
    if [ -n "$keychain" ] && [ -n "$keychain_pwd" ]; then
        print_status "Setting up keychain..."
        security unlock-keychain -p "$keychain_pwd" "$keychain"
        security set-keychain-settings -t 3600 -l "$keychain"
    fi
    
    # List available identities for debugging
    print_status "Available signing identities:"
    security find-identity -p codesigning -v || true
    
    # Sign the application
    print_status "Signing with identity: $identity"
    
    local sign_args=(
        --force
        --sign "$identity"
        --options runtime
        --timestamp
        --verbose
    )
    
    # Add entitlements if they exist
    if [ -f "$entitlements" ]; then
        sign_args+=(--entitlements "$entitlements")
        print_status "Using entitlements: $entitlements"
    else
        print_warning "Entitlements file not found: $entitlements"
    fi
    
    # Sign all binaries within the app bundle
    if [[ "$BINARY_PATH" == *.app ]]; then
        print_status "Signing app bundle and all contents..."
        
        # Sign frameworks and libraries first
        find "$BINARY_PATH" -type f \( -name "*.dylib" -o -name "*.framework" \) -exec \
            codesign "${sign_args[@]}" {} \; 2>/dev/null || true
        
        # Sign the main app bundle
        codesign "${sign_args[@]}" "$BINARY_PATH"
    else
        # Sign a single binary
        codesign "${sign_args[@]}" "$BINARY_PATH"
    fi
    
    # Verify the signature
    print_status "Verifying signature..."
    codesign --verify --deep --strict --verbose=2 "$BINARY_PATH"
    
    # Check if the signature is valid
    if codesign --verify --deep "$BINARY_PATH" 2>/dev/null; then
        print_success "macOS code signing completed successfully"
        
        # Display signature information
        print_status "Signature information:"
        codesign -dv --verbose=4 "$BINARY_PATH" 2>&1 | head -10
    else
        print_error "Code signature verification failed"
        return 1
    fi
}

# Function to sign Windows executable
sign_windows() {
    print_status "Signing Windows executable..."
    
    # Use environment variables if command line options not provided
    local cert_path="${CERTIFICATE_PATH:-${WINDOWS_CERTIFICATE:-}}"
    local cert_pwd="${CERTIFICATE_PASSWORD:-${WINDOWS_CERT_PASSWORD:-}}"
    local timestamp="${TIMESTAMP_URL:-${TIMESTAMP_SERVER:-http://timestamp.digicert.com}}"
    
    if [ -z "$cert_path" ]; then
        print_error "No certificate path provided"
        print_error "Set WINDOWS_CERTIFICATE environment variable or use --certificate"
        return 1
    fi
    
    if [ ! -f "$cert_path" ]; then
        print_error "Certificate file not found: $cert_path"
        return 1
    fi
    
    # Check if we're on Windows or have wine available
    local signtool_path=""
    
    if command -v signtool.exe >/dev/null 2>&1; then
        signtool_path="signtool.exe"
    elif command -v wine >/dev/null 2>&1; then
        # Try to find signtool via wine
        local wine_signtool=$(wine cmd /c 'where signtool.exe' 2>/dev/null | tr -d '\r' | head -1)
        if [ -n "$wine_signtool" ]; then
            signtool_path="wine $wine_signtool"
        fi
    fi
    
    if [ -z "$signtool_path" ]; then
        print_error "signtool.exe not found"
        print_error "Please install Windows SDK or use Windows environment"
        return 1
    fi
    
    print_status "Using signtool: $signtool_path"
    
    # Prepare signing command
    local sign_cmd="$signtool_path sign"
    sign_cmd="$sign_cmd /f \"$cert_path\""
    
    if [ -n "$cert_pwd" ]; then
        sign_cmd="$sign_cmd /p \"$cert_pwd\""
    fi
    
    sign_cmd="$sign_cmd /tr \"$timestamp\""
    sign_cmd="$sign_cmd /td sha256"
    sign_cmd="$sign_cmd /fd sha256"
    sign_cmd="$sign_cmd /v"
    sign_cmd="$sign_cmd \"$BINARY_PATH\""
    
    print_status "Signing command: $sign_cmd"
    
    # Execute signing
    if eval "$sign_cmd"; then
        print_success "Windows code signing completed successfully"
        
        # Verify the signature
        print_status "Verifying signature..."
        local verify_cmd="$signtool_path verify /pa /v \"$BINARY_PATH\""
        if eval "$verify_cmd"; then
            print_success "Signature verification passed"
        else
            print_warning "Signature verification failed, but signing may still be valid"
        fi
    else
        print_error "Code signing failed"
        return 1
    fi
}

# Function to check signing prerequisites
check_prerequisites() {
    print_status "Checking prerequisites for $PLATFORM..."
    
    case $PLATFORM in
        "macos")
            if ! command -v codesign >/dev/null 2>&1; then
                print_error "codesign not found - macOS development tools required"
                return 1
            fi
            
            if ! command -v security >/dev/null 2>&1; then
                print_error "security not found - macOS security framework required"
                return 1
            fi
            ;;
        "windows")
            # Check will be done in sign_windows function
            ;;
        *)
            print_error "Unsupported platform: $PLATFORM"
            return 1
            ;;
    esac
    
    print_success "Prerequisites check passed"
}

# Main execution
check_prerequisites

case $PLATFORM in
    "macos")
        sign_macos
        ;;
    "windows")
        sign_windows
        ;;
    *)
        print_error "Unsupported platform: $PLATFORM"
        exit 1
        ;;
esac

print_success "Code signing process completed for $PLATFORM"
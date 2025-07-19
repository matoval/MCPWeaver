#!/bin/bash

# MCPWeaver macOS Notarization Script
# Handles Apple notarization process for macOS applications

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
APP_PATH=""
APPLE_ID=""
APPLE_PASSWORD=""
TEAM_ID=""
BUNDLE_ID="com.mcpweaver.app"
WAIT_FOR_COMPLETION=false
ZIP_PATH=""
DMG_PATH=""
TIMEOUT=3600  # 1 hour timeout

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --app-path)
            APP_PATH="$2"
            shift 2
            ;;
        --apple-id)
            APPLE_ID="$2"
            shift 2
            ;;
        --password)
            APPLE_PASSWORD="$2"
            shift 2
            ;;
        --team-id)
            TEAM_ID="$2"
            shift 2
            ;;
        --bundle-id)
            BUNDLE_ID="$2"
            shift 2
            ;;
        --wait)
            WAIT_FOR_COMPLETION=true
            shift
            ;;
        --timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        --help)
            echo "MCPWeaver macOS Notarization Script"
            echo ""
            echo "Usage: $0 --app-path PATH [OPTIONS]"
            echo ""
            echo "Required:"
            echo "  --app-path PATH         Path to .app bundle or .dmg file"
            echo ""
            echo "Options:"
            echo "  --apple-id EMAIL        Apple ID email (or set APPLE_ID env var)"
            echo "  --password PASSWORD     App-specific password (or set APPLE_PASSWORD env var)"
            echo "  --team-id TEAM          Team ID (or set APPLE_TEAM_ID env var)"
            echo "  --bundle-id ID          Bundle identifier (default: com.mcpweaver.app)"
            echo "  --wait                  Wait for notarization to complete"
            echo "  --timeout SECONDS       Timeout for waiting (default: 3600)"
            echo "  --help                  Show this help"
            echo ""
            echo "Environment Variables:"
            echo "  APPLE_ID               Apple ID email"
            echo "  APPLE_PASSWORD         App-specific password"
            echo "  APPLE_TEAM_ID          Team ID"
            echo ""
            echo "Note: You must create an app-specific password at appleid.apple.com"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Use environment variables if not provided via command line
APPLE_ID="${APPLE_ID:-${APPLE_ID}}"
APPLE_PASSWORD="${APPLE_PASSWORD:-${APPLE_PASSWORD}}"
TEAM_ID="${TEAM_ID:-${APPLE_TEAM_ID}}"

# Validate required parameters
if [ -z "$APP_PATH" ]; then
    print_error "App path is required"
    echo "Use --help for usage information"
    exit 1
fi

if [ ! -e "$APP_PATH" ]; then
    print_error "App path not found: $APP_PATH"
    exit 1
fi

if [ -z "$APPLE_ID" ]; then
    print_error "Apple ID is required"
    print_error "Set APPLE_ID environment variable or use --apple-id"
    exit 1
fi

if [ -z "$APPLE_PASSWORD" ]; then
    print_error "Apple password is required"
    print_error "Set APPLE_PASSWORD environment variable or use --password"
    print_error "Create an app-specific password at https://appleid.apple.com"
    exit 1
fi

print_status "Starting notarization process..."
print_status "App Path: $APP_PATH"
print_status "Apple ID: $APPLE_ID"
print_status "Bundle ID: $BUNDLE_ID"

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if we're on macOS
    if [ "$(uname -s)" != "Darwin" ]; then
        print_error "This script must be run on macOS"
        exit 1
    fi
    
    # Check for required tools
    if ! command -v xcrun >/dev/null 2>&1; then
        print_error "xcrun not found - Xcode command line tools required"
        exit 1
    fi
    
    # Check if notarytool is available (Xcode 13+)
    if xcrun notarytool --help >/dev/null 2>&1; then
        NOTARY_TOOL="notarytool"
        print_status "Using notarytool (recommended)"
    elif xcrun altool --help >/dev/null 2>&1; then
        NOTARY_TOOL="altool"
        print_warning "Using altool (deprecated, upgrade to Xcode 13+ recommended)"
    else
        print_error "Neither notarytool nor altool found"
        exit 1
    fi
    
    print_success "Prerequisites check passed"
}

# Function to create zip for notarization
prepare_for_notarization() {
    print_status "Preparing for notarization..."
    
    local base_name=$(basename "$APP_PATH" | sed 's/\.[^.]*$//')
    
    if [[ "$APP_PATH" == *.app ]]; then
        # Create zip from app bundle
        ZIP_PATH="${base_name}-notarization.zip"
        print_status "Creating zip archive: $ZIP_PATH"
        ditto -c -k --keepParent "$APP_PATH" "$ZIP_PATH"
    elif [[ "$APP_PATH" == *.dmg ]]; then
        # Use DMG directly
        DMG_PATH="$APP_PATH"
        print_status "Using DMG for notarization: $DMG_PATH"
    else
        print_error "Unsupported file type: $APP_PATH"
        print_error "Only .app bundles and .dmg files are supported"
        exit 1
    fi
}

# Function to submit for notarization using notarytool
submit_notarytool() {
    print_status "Submitting for notarization using notarytool..."
    
    local file_to_notarize="${ZIP_PATH:-$DMG_PATH}"
    local submit_args=(
        "submit"
        "$file_to_notarize"
        "--apple-id" "$APPLE_ID"
        "--password" "$APPLE_PASSWORD"
    )
    
    if [ -n "$TEAM_ID" ]; then
        submit_args+=("--team-id" "$TEAM_ID")
    fi
    
    if [ "$WAIT_FOR_COMPLETION" = true ]; then
        submit_args+=("--wait" "--timeout" "$TIMEOUT")
    fi
    
    print_status "Submitting to Apple..."
    local result
    if result=$(xcrun notarytool "${submit_args[@]}" 2>&1); then
        echo "$result"
        
        # Extract submission ID from output
        local submission_id=$(echo "$result" | grep -E "id: [a-f0-9-]+" | head -1 | sed 's/.*id: //')
        
        if [ -n "$submission_id" ]; then
            print_success "Submission ID: $submission_id"
            
            if [ "$WAIT_FOR_COMPLETION" = false ]; then
                print_status "Notarization submitted. Check status with:"
                echo "xcrun notarytool info \"$submission_id\" --apple-id \"$APPLE_ID\" --password \"$APPLE_PASSWORD\""
                
                if [ -n "$TEAM_ID" ]; then
                    echo "  --team-id \"$TEAM_ID\""
                fi
            else
                print_status "Waiting for notarization to complete..."
                check_notarization_status "$submission_id"
            fi
        else
            print_error "Could not extract submission ID from response"
            return 1
        fi
    else
        print_error "Notarization submission failed:"
        echo "$result"
        return 1
    fi
}

# Function to submit for notarization using altool (legacy)
submit_altool() {
    print_status "Submitting for notarization using altool..."
    
    local file_to_notarize="${ZIP_PATH:-$DMG_PATH}"
    local submit_args=(
        "--notarize-app"
        "--primary-bundle-id" "$BUNDLE_ID"
        "--username" "$APPLE_ID"
        "--password" "$APPLE_PASSWORD"
        "--file" "$file_to_notarize"
    )
    
    if [ -n "$TEAM_ID" ]; then
        submit_args+=("--asc-provider" "$TEAM_ID")
    fi
    
    print_status "Submitting to Apple..."
    local result
    if result=$(xcrun altool "${submit_args[@]}" 2>&1); then
        echo "$result"
        
        # Extract request UUID from output
        local request_uuid=$(echo "$result" | grep "RequestUUID" | sed 's/.*RequestUUID = //')
        
        if [ -n "$request_uuid" ]; then
            print_success "Request UUID: $request_uuid"
            
            if [ "$WAIT_FOR_COMPLETION" = true ]; then
                print_status "Waiting for notarization to complete..."
                check_altool_status "$request_uuid"
            else
                print_status "Notarization submitted. Check status with:"
                echo "xcrun altool --notarization-info \"$request_uuid\" --username \"$APPLE_ID\" --password \"$APPLE_PASSWORD\""
            fi
        else
            print_error "Could not extract request UUID from response"
            return 1
        fi
    else
        print_error "Notarization submission failed:"
        echo "$result"
        return 1
    fi
}

# Function to check notarization status (notarytool)
check_notarization_status() {
    local submission_id="$1"
    local elapsed=0
    
    while [ $elapsed -lt $TIMEOUT ]; do
        sleep 30
        elapsed=$((elapsed + 30))
        
        print_status "Checking status... (${elapsed}s elapsed)"
        
        local status_result
        local status_args=(
            "info"
            "$submission_id"
            "--apple-id" "$APPLE_ID"
            "--password" "$APPLE_PASSWORD"
        )
        
        if [ -n "$TEAM_ID" ]; then
            status_args+=("--team-id" "$TEAM_ID")
        fi
        
        if status_result=$(xcrun notarytool "${status_args[@]}" 2>&1); then
            echo "$status_result"
            
            if echo "$status_result" | grep -q "status: Accepted"; then
                print_success "Notarization completed successfully!"
                staple_app
                return 0
            elif echo "$status_result" | grep -q "status: Invalid"; then
                print_error "Notarization failed!"
                get_notarization_log "$submission_id"
                return 1
            elif echo "$status_result" | grep -q "status: In Progress"; then
                print_status "Notarization still in progress..."
            else
                print_warning "Unknown status, continuing to wait..."
            fi
        else
            print_warning "Failed to check status, retrying..."
        fi
    done
    
    print_error "Timeout waiting for notarization to complete"
    return 1
}

# Function to check altool status (legacy)
check_altool_status() {
    local request_uuid="$1"
    local elapsed=0
    
    while [ $elapsed -lt $TIMEOUT ]; do
        sleep 30
        elapsed=$((elapsed + 30))
        
        print_status "Checking status... (${elapsed}s elapsed)"
        
        local status_result
        local status_args=(
            "--notarization-info" "$request_uuid"
            "--username" "$APPLE_ID"
            "--password" "$APPLE_PASSWORD"
        )
        
        if status_result=$(xcrun altool "${status_args[@]}" 2>&1); then
            echo "$status_result"
            
            if echo "$status_result" | grep -q "Status: success"; then
                print_success "Notarization completed successfully!"
                staple_app
                return 0
            elif echo "$status_result" | grep -q "Status: invalid"; then
                print_error "Notarization failed!"
                return 1
            elif echo "$status_result" | grep -q "Status: in progress"; then
                print_status "Notarization still in progress..."
            else
                print_warning "Unknown status, continuing to wait..."
            fi
        else
            print_warning "Failed to check status, retrying..."
        fi
    done
    
    print_error "Timeout waiting for notarization to complete"
    return 1
}

# Function to get notarization log
get_notarization_log() {
    local submission_id="$1"
    print_status "Retrieving notarization log..."
    
    local log_args=(
        "log"
        "$submission_id"
        "--apple-id" "$APPLE_ID"
        "--password" "$APPLE_PASSWORD"
    )
    
    if [ -n "$TEAM_ID" ]; then
        log_args+=("--team-id" "$TEAM_ID")
    fi
    
    xcrun notarytool "${log_args[@]}" || print_warning "Could not retrieve log"
}

# Function to staple the notarization ticket
staple_app() {
    if [[ "$APP_PATH" == *.app ]]; then
        print_status "Stapling notarization ticket to app..."
        if xcrun stapler staple "$APP_PATH"; then
            print_success "Stapling completed successfully"
            
            # Verify stapling
            print_status "Verifying stapled ticket..."
            if xcrun stapler validate "$APP_PATH"; then
                print_success "Stapled ticket verification passed"
            else
                print_warning "Stapled ticket verification failed"
            fi
        else
            print_error "Stapling failed"
        fi
    elif [[ "$APP_PATH" == *.dmg ]]; then
        print_status "Stapling notarization ticket to DMG..."
        if xcrun stapler staple "$APP_PATH"; then
            print_success "DMG stapling completed successfully"
        else
            print_error "DMG stapling failed"
        fi
    fi
}

# Function to cleanup temporary files
cleanup() {
    if [ -n "$ZIP_PATH" ] && [ -f "$ZIP_PATH" ]; then
        print_status "Cleaning up temporary zip file..."
        rm -f "$ZIP_PATH"
    fi
}

# Set trap for cleanup
trap cleanup EXIT

# Main execution
check_prerequisites
prepare_for_notarization

# Submit based on available tool
if [ "$NOTARY_TOOL" = "notarytool" ]; then
    submit_notarytool
else
    submit_altool
fi

print_success "Notarization process completed!"

# Final verification
if [[ "$APP_PATH" == *.app ]]; then
    print_status "Final verification..."
    if spctl -a -t exec -vv "$APP_PATH" 2>&1 | grep -q "accepted"; then
        print_success "App will be accepted by Gatekeeper"
    else
        print_warning "App may not be accepted by Gatekeeper"
    fi
fi
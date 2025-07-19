#!/bin/bash

# MCPWeaver Release Automation Script
# Handles complete release process including versioning, building, signing, and publishing

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
VERSION=""
RELEASE_TYPE="patch"  # major, minor, patch
PRERELEASE=false
DRY_RUN=false
SKIP_TESTS=false
SKIP_BUILD=false
SKIP_SIGN=false
SKIP_PUBLISH=false
PLATFORMS=("windows" "macos" "linux")
OUTPUT_DIR="./dist"
CHANGELOG_FILE="CHANGELOG.md"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --version)
            VERSION="$2"
            shift 2
            ;;
        --type)
            RELEASE_TYPE="$2"
            shift 2
            ;;
        --prerelease)
            PRERELEASE=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --skip-sign)
            SKIP_SIGN=true
            shift
            ;;
        --skip-publish)
            SKIP_PUBLISH=true
            shift
            ;;
        --platforms)
            IFS=',' read -ra PLATFORMS <<< "$2"
            shift 2
            ;;
        --output-dir)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        --help)
            echo "MCPWeaver Release Automation Script"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --version VERSION     Specific version to release (e.g., 1.2.0)"
            echo "  --type TYPE           Release type: major, minor, patch (default: patch)"
            echo "  --prerelease          Mark as prerelease"
            echo "  --dry-run             Simulate release without making changes"
            echo "  --skip-tests          Skip running tests"
            echo "  --skip-build          Skip building applications"
            echo "  --skip-sign           Skip code signing"
            echo "  --skip-publish        Skip publishing to GitHub"
            echo "  --platforms LIST      Comma-separated list of platforms (windows,macos,linux)"
            echo "  --output-dir DIR      Output directory for builds (default: ./dist)"
            echo "  --help                Show this help"
            echo ""
            echo "Environment Variables:"
            echo "  GITHUB_TOKEN          GitHub token for publishing releases"
            echo "  APPLE_ID              Apple ID for notarization"
            echo "  APPLE_PASSWORD        Apple app-specific password"
            echo "  APPLE_TEAM_ID         Apple team ID"
            echo "  WINDOWS_CERTIFICATE   Windows signing certificate"
            echo "  WINDOWS_CERT_PASSWORD Windows certificate password"
            echo ""
            echo "Examples:"
            echo "  $0 --type minor                    # Bump minor version and release"
            echo "  $0 --version 2.0.0 --prerelease    # Release specific version as prerelease"
            echo "  $0 --dry-run --type major          # Simulate major version release"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

print_status "MCPWeaver Release Automation"
print_status "=============================="

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check required tools
    local required_tools=("git" "go" "npm" "wails")
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            print_error "$tool is required but not installed"
            exit 1
        fi
    done
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        print_error "Not in a git repository"
        exit 1
    fi
    
    # Check for uncommitted changes
    if ! git diff-index --quiet HEAD --; then
        print_error "Uncommitted changes detected. Please commit or stash changes first."
        exit 1
    fi
    
    # Check if we're on main branch (unless dry run)
    local current_branch=$(git branch --show-current)
    if [ "$current_branch" != "main" ] && [ "$DRY_RUN" = false ]; then
        print_warning "Not on main branch (current: $current_branch)"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    print_success "Prerequisites check passed"
}

# Function to determine next version
determine_version() {
    if [ -n "$VERSION" ]; then
        print_status "Using specified version: $VERSION"
        return
    fi
    
    print_status "Determining next version..."
    
    # Get current version from git tags
    local current_version=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0")
    print_status "Current version: $current_version"
    
    # Parse version components
    IFS='.' read -ra VERSION_PARTS <<< "$current_version"
    local major=${VERSION_PARTS[0]:-0}
    local minor=${VERSION_PARTS[1]:-0}
    local patch=${VERSION_PARTS[2]:-0}
    
    # Bump version based on type
    case $RELEASE_TYPE in
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "patch")
            patch=$((patch + 1))
            ;;
        *)
            print_error "Invalid release type: $RELEASE_TYPE"
            exit 1
            ;;
    esac
    
    VERSION="${major}.${minor}.${patch}"
    
    if [ "$PRERELEASE" = true ]; then
        local timestamp=$(date +%Y%m%d%H%M%S)
        VERSION="${VERSION}-rc.${timestamp}"
    fi
    
    print_status "Next version: $VERSION"
}

# Function to update version in files
update_version_files() {
    print_status "Updating version in project files..."
    
    # Update wails.json
    if [ -f "wails.json" ]; then
        local temp_file=$(mktemp)
        jq --arg version "$VERSION" '.info.productVersion = $version' wails.json > "$temp_file"
        mv "$temp_file" wails.json
        print_status "Updated wails.json"
    fi
    
    # Update package.json if it exists
    if [ -f "package.json" ]; then
        local temp_file=$(mktemp)
        jq --arg version "$VERSION" '.version = $version' package.json > "$temp_file"
        mv "$temp_file" package.json
        print_status "Updated package.json"
    fi
    
    # Update frontend package.json
    if [ -f "frontend/package.json" ]; then
        local temp_file=$(mktemp)
        jq --arg version "$VERSION" '.version = $version' frontend/package.json > "$temp_file"
        mv "$temp_file" frontend/package.json
        print_status "Updated frontend/package.json"
    fi
}

# Function to update changelog
update_changelog() {
    print_status "Updating changelog..."
    
    if [ ! -f "$CHANGELOG_FILE" ]; then
        print_status "Creating changelog file..."
        cat > "$CHANGELOG_FILE" << EOF
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [${VERSION}] - $(date +%Y-%m-%d)

### Added
- Initial release

### Changed
- 

### Fixed
- 

### Removed
- 

EOF
    else
        # Add new version entry at the top
        local temp_file=$(mktemp)
        {
            head -n 5 "$CHANGELOG_FILE"
            echo ""
            echo "## [${VERSION}] - $(date +%Y-%m-%d)"
            echo ""
            echo "### Added"
            echo "- "
            echo ""
            echo "### Changed"
            echo "- "
            echo ""
            echo "### Fixed"
            echo "- "
            echo ""
            echo "### Removed"
            echo "- "
            echo ""
            tail -n +6 "$CHANGELOG_FILE"
        } > "$temp_file"
        mv "$temp_file" "$CHANGELOG_FILE"
    fi
    
    print_status "Please update the changelog with release notes"
    if [ "$DRY_RUN" = false ]; then
        read -p "Open changelog in editor? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            ${EDITOR:-nano} "$CHANGELOG_FILE"
        fi
    fi
}

# Function to run tests
run_tests() {
    if [ "$SKIP_TESTS" = true ]; then
        print_warning "Skipping tests"
        return
    fi
    
    print_status "Running tests..."
    
    # Run Go tests
    go test ./... -v -race -coverprofile=coverage.out
    
    # Run frontend tests if available
    if [ -f "frontend/package.json" ] && grep -q '"test"' frontend/package.json; then
        cd frontend
        npm test || print_warning "Frontend tests failed or not available"
        cd ..
    fi
    
    print_success "Tests completed"
}

# Function to build applications
build_applications() {
    if [ "$SKIP_BUILD" = true ]; then
        print_warning "Skipping build"
        return
    fi
    
    print_status "Building applications..."
    
    # Clean previous builds
    rm -rf "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR"
    
    # Build for each platform
    for platform in "${PLATFORMS[@]}"; do
        print_status "Building for $platform..."
        
        case $platform in
            "windows")
                wails build -platform windows/amd64 -clean \
                    -ldflags "-s -w -X main.version=$VERSION -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.buildCommit=$(git rev-parse --short HEAD)"
                ;;
            "macos")
                wails build -platform darwin/universal -clean \
                    -ldflags "-s -w -X main.version=$VERSION -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.buildCommit=$(git rev-parse --short HEAD)"
                ;;
            "linux")
                wails build -platform linux/amd64 -clean \
                    -ldflags "-s -w -X main.version=$VERSION -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.buildCommit=$(git rev-parse --short HEAD)"
                ;;
            *)
                print_warning "Unknown platform: $platform"
                ;;
        esac
    done
    
    print_success "Build completed"
}

# Function to sign applications
sign_applications() {
    if [ "$SKIP_SIGN" = true ]; then
        print_warning "Skipping code signing"
        return
    fi
    
    print_status "Signing applications..."
    
    # Sign macOS app if built and credentials available
    if [[ " ${PLATFORMS[*]} " =~ " macos " ]] && [ -n "${APPLE_DEVELOPER_ID:-}" ]; then
        if [ -d "./build/bin/MCPWeaver.app" ]; then
            ./scripts/sign-code.sh --platform macos --binary "./build/bin/MCPWeaver.app"
            
            # Notarize if credentials available
            if [ -n "${APPLE_ID:-}" ] && [ -n "${APPLE_PASSWORD:-}" ]; then
                ./scripts/notarize-macos.sh --app-path "./build/bin/MCPWeaver.app" --wait
            fi
        fi
    fi
    
    # Sign Windows executable if built and certificate available
    if [[ " ${PLATFORMS[*]} " =~ " windows " ]] && [ -n "${WINDOWS_CERTIFICATE:-}" ]; then
        if [ -f "./build/bin/MCPWeaver-windows-amd64.exe" ]; then
            ./scripts/sign-code.sh --platform windows --binary "./build/bin/MCPWeaver-windows-amd64.exe"
        fi
    fi
    
    print_success "Code signing completed"
}

# Function to package applications
package_applications() {
    print_status "Packaging applications..."
    
    # Use the packaging script
    local package_args=("--version" "$VERSION" "--output-dir" "$OUTPUT_DIR")
    
    for platform in "${PLATFORMS[@]}"; do
        package_args+=("--$platform")
    done
    
    if [ "$SKIP_SIGN" = false ]; then
        package_args+=("--sign")
    fi
    
    package_args+=("--installer")
    
    ./scripts/package.sh "${package_args[@]}"
    
    print_success "Packaging completed"
}

# Function to create git tag
create_git_tag() {
    if [ "$DRY_RUN" = true ]; then
        print_status "Would create git tag: v$VERSION"
        return
    fi
    
    print_status "Creating git tag..."
    
    # Commit version changes
    git add .
    git commit -m "chore: bump version to $VERSION" || print_warning "No changes to commit"
    
    # Create tag
    git tag -a "v$VERSION" -m "Release version $VERSION"
    
    print_success "Git tag created: v$VERSION"
}

# Function to publish release
publish_release() {
    if [ "$SKIP_PUBLISH" = true ]; then
        print_warning "Skipping publish"
        return
    fi
    
    if [ "$DRY_RUN" = true ]; then
        print_status "Would publish release to GitHub"
        return
    fi
    
    print_status "Publishing release to GitHub..."
    
    if [ -z "${GITHUB_TOKEN:-}" ]; then
        print_error "GITHUB_TOKEN environment variable is required for publishing"
        exit 1
    fi
    
    # Push changes and tags
    git push origin main
    git push origin "v$VERSION"
    
    # Extract changelog for this version
    local release_notes=""
    if [ -f "$CHANGELOG_FILE" ]; then
        release_notes=$(awk "/## \[${VERSION}\]/,/## \[.*\]/{if(/## \[.*\]/ && !/## \[${VERSION}\]/) exit; print}" "$CHANGELOG_FILE" | head -n -1)
    fi
    
    # Create GitHub release
    local prerelease_flag=""
    if [ "$PRERELEASE" = true ]; then
        prerelease_flag="--prerelease"
    fi
    
    gh release create "v$VERSION" \
        --title "MCPWeaver $VERSION" \
        --notes "$release_notes" \
        $prerelease_flag \
        "$OUTPUT_DIR"/*
    
    print_success "Release published to GitHub"
}

# Function to display summary
display_summary() {
    print_status "Release Summary"
    echo "==============="
    echo "Version: $VERSION"
    echo "Type: $RELEASE_TYPE"
    echo "Prerelease: $PRERELEASE"
    echo "Platforms: ${PLATFORMS[*]}"
    echo "Output Directory: $OUTPUT_DIR"
    echo ""
    
    if [ -d "$OUTPUT_DIR" ]; then
        echo "Created artifacts:"
        ls -la "$OUTPUT_DIR/"
    fi
    
    if [ "$DRY_RUN" = true ]; then
        print_warning "This was a dry run - no changes were made"
    else
        print_success "Release $VERSION completed successfully!"
    fi
}

# Main execution
print_status "Starting release process..."

if [ "$DRY_RUN" = true ]; then
    print_warning "DRY RUN MODE - No changes will be made"
fi

check_prerequisites
determine_version
update_version_files
update_changelog
run_tests
build_applications
sign_applications
package_applications
create_git_tag
publish_release
display_summary
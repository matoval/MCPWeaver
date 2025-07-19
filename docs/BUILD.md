# MCPWeaver Build System Documentation

This document provides comprehensive information about the MCPWeaver build system, including cross-platform builds, packaging, code signing, and release automation.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Build Configuration](#build-configuration)
- [Building Applications](#building-applications)
- [Platform-Specific Packaging](#platform-specific-packaging)
- [Code Signing](#code-signing)
- [Release Automation](#release-automation)
- [CI/CD Pipeline](#cicd-pipeline)
- [Troubleshooting](#troubleshooting)

## Overview

MCPWeaver uses a comprehensive build system based on Wails v2 that supports:

- **Cross-platform builds**: Windows, macOS (Intel/ARM), and Linux
- **Automated packaging**: Platform-specific installers (.exe, .dmg, .AppImage)
- **Code signing**: Windows Authenticode and macOS Developer ID
- **Notarization**: macOS app notarization for Gatekeeper compatibility
- **CI/CD integration**: GitHub Actions for automated builds and releases
- **Quality gates**: Testing, linting, and security scanning

## Prerequisites

### Required Tools

- **Go 1.21+**: Backend compilation
- **Node.js 18+**: Frontend build system
- **Wails v2.10.1+**: Application framework
- **Git**: Version control

### Platform-Specific Requirements

#### Windows
- **Windows 10/11** (for building Windows binaries)
- **Windows SDK** (for code signing)
- **NSIS** (optional, for advanced installers)

#### macOS
- **macOS 10.15+** (for building macOS binaries)
- **Xcode Command Line Tools**
- **Apple Developer Account** (for code signing and notarization)

#### Linux
- **Ubuntu 18.04+** or equivalent
- **GTK 3 development libraries**:
  ```bash
  sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev
  ```

### Installation

1. **Install Wails CLI**:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

2. **Verify installation**:
   ```bash
   wails doctor
   ```

3. **Install frontend dependencies**:
   ```bash
   cd frontend && npm install
   ```

## Build Configuration

### Wails Configuration

The build system is configured through `wails.json`:

```json
{
  "name": "MCPWeaver",
  "outputfilename": "MCPWeaver",
  "info": {
    "productName": "MCPWeaver",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2025 MCPWeaver. All rights reserved."
  },
  "build": {
    "ldflags": ["-s", "-w", "-X main.version={{.Info.ProductVersion}}"],
    "tags": ["production"]
  },
  "macos": {
    "bundleID": "com.mcpweaver.app",
    "enableCodeSigning": true,
    "enableNotarization": true
  },
  "windows": {
    "wixVersion": "v4",
    "webview2InstallMode": "downloadBootstrapper"
  },
  "linux": {
    "packageType": "appimage"
  }
}
```

### Build Scripts

The build system includes several scripts in the `scripts/` directory:

- `build.sh` / `build.ps1`: Cross-platform build automation
- `package.sh`: Platform-specific packaging
- `sign-code.sh`: Code signing automation
- `notarize-macos.sh`: macOS notarization
- `release.sh`: Complete release automation

## Building Applications

### Quick Build Commands

```bash
# Build for current platform
wails build

# Build for specific platform
wails build -platform windows/amd64
wails build -platform darwin/universal
wails build -platform linux/amd64

# Development build with hot reload
wails dev
```

### Using Build Scripts

#### Cross-Platform Build

```bash
# Build for all platforms
./scripts/build.sh --all

# Build for specific platforms
./scripts/build.sh --windows --linux

# Build with custom version
./scripts/build.sh --all --version 1.2.0

# Clean build
./scripts/build.sh --all --clean
```

#### Windows (PowerShell)

```powershell
# Build for all platforms
.\scripts\build.ps1 -All

# Build for Windows only
.\scripts\build.ps1 -Windows -Version "1.2.0"
```

### Build Options

- `--all`: Build for all supported platforms
- `--windows`: Build for Windows (amd64)
- `--macos`: Build for macOS (Intel and ARM)
- `--linux`: Build for Linux (amd64)
- `--clean`: Clean build directory first
- `--version VER`: Set custom version number

## Platform-Specific Packaging

### Windows Packaging

#### ZIP Package
```bash
./scripts/package.sh --windows --version 1.0.0
```

#### NSIS Installer
```bash
./scripts/package.sh --windows --installer --version 1.0.0
```

**Requirements**:
- NSIS installed and in PATH
- Certificate for code signing (optional)

### macOS Packaging

#### DMG Creation
```bash
./scripts/create-dmg.sh 1.0.0
```

#### Tar.gz Package
```bash
./scripts/package.sh --macos --version 1.0.0
```

**Requirements**:
- macOS development environment
- Code signing certificate (for distribution)
- Apple Developer account (for notarization)

### Linux Packaging

#### AppImage Creation
```bash
./scripts/create-appimage.sh 1.0.0
```

#### Tar.gz Package
```bash
./scripts/package.sh --linux --version 1.0.0
```

**Requirements**:
- AppImageTool (automatically downloaded)
- Desktop file and icon

## Code Signing

### macOS Code Signing

#### Setup
1. **Install certificate** in Keychain Access
2. **Set environment variables**:
   ```bash
   export APPLE_DEVELOPER_ID="Developer ID Application: Your Name (TEAM_ID)"
   export APPLE_KEYCHAIN_PATH="/path/to/keychain"
   export APPLE_KEYCHAIN_PASSWORD="keychain_password"
   ```

#### Sign Application
```bash
./scripts/sign-code.sh --platform macos --binary ./build/bin/MCPWeaver.app
```

#### Notarization
```bash
export APPLE_ID="your@apple.id"
export APPLE_PASSWORD="app-specific-password"
export APPLE_TEAM_ID="TEAM_ID"

./scripts/notarize-macos.sh --app-path ./build/bin/MCPWeaver.app --wait
```

### Windows Code Signing

#### Setup
1. **Obtain code signing certificate** (.pfx or .p12 file)
2. **Set environment variables**:
   ```bash
   export WINDOWS_CERTIFICATE="/path/to/certificate.pfx"
   export WINDOWS_CERT_PASSWORD="certificate_password"
   ```

#### Sign Executable
```bash
./scripts/sign-code.sh --platform windows --binary ./build/bin/MCPWeaver.exe
```

## Release Automation

### Complete Release Process

```bash
# Patch release
./scripts/release.sh --type patch

# Minor release with prerelease
./scripts/release.sh --type minor --prerelease

# Specific version
./scripts/release.sh --version 2.0.0

# Dry run (simulate without changes)
./scripts/release.sh --dry-run --type major
```

### Release Steps

The release script automates:

1. **Version determination** (semantic versioning)
2. **File updates** (wails.json, package.json)
3. **Changelog updates** (CHANGELOG.md)
4. **Test execution** (Go and frontend tests)
5. **Cross-platform builds**
6. **Code signing** (if certificates available)
7. **Packaging** (platform-specific installers)
8. **Git tagging**
9. **GitHub release creation**

### Environment Variables

```bash
# GitHub
export GITHUB_TOKEN="github_token"

# Apple (macOS)
export APPLE_ID="your@apple.id"
export APPLE_PASSWORD="app-specific-password"
export APPLE_TEAM_ID="TEAM_ID"
export APPLE_DEVELOPER_ID="Developer ID Application: Your Name (TEAM_ID)"

# Microsoft (Windows)
export WINDOWS_CERTIFICATE="/path/to/certificate.pfx"
export WINDOWS_CERT_PASSWORD="certificate_password"
```

## CI/CD Pipeline

### GitHub Actions Workflows

#### 1. Build and Release (`build.yml`)
- **Triggers**: Push to main, tags, manual dispatch
- **Platforms**: Windows, macOS, Linux
- **Features**:
  - Cross-platform matrix builds
  - Code signing and notarization
  - Automated releases on tags
  - Artifact management

#### 2. Pull Request Validation (`pr-validation.yml`)
- **Triggers**: Pull requests
- **Checks**:
  - Code quality (linting, formatting)
  - Tests (unit, integration)
  - Security scanning
  - Cross-platform build verification
  - Performance regression testing

#### 3. Nightly Testing (`nightly.yml`)
- **Schedule**: Daily at 2 AM UTC
- **Features**:
  - Extended testing across Go versions
  - Performance monitoring
  - Security audits
  - Code quality analysis

### Workflow Configuration

```yaml
# Example build matrix
strategy:
  matrix:
    os: [windows-latest, macos-latest, ubuntu-latest]
    include:
      - os: windows-latest
        platform: windows
        arch: amd64
      - os: macos-latest
        platform: darwin
        arch: universal
      - os: ubuntu-latest
        platform: linux
        arch: amd64
```

### Secrets Configuration

Required GitHub repository secrets:

```
# Code Signing
MACOS_CERTIFICATE          # Base64 encoded .p12
MACOS_CERTIFICATE_PWD       # Certificate password
MACOS_KEYCHAIN_PWD         # Keychain password
WINDOWS_CERTIFICATE         # Base64 encoded .pfx
WINDOWS_CERTIFICATE_PWD     # Certificate password

# Notarization
APPLE_ID                   # Apple ID email
APPLE_PASSWORD             # App-specific password
APPLE_TEAM_ID              # Apple Team ID

# Release
GITHUB_TOKEN               # GitHub token (auto-provided)
```

## Troubleshooting

### Common Build Issues

#### 1. Wails Build Failures

**Problem**: Build fails with WebView2 errors on Windows
**Solution**:
```bash
# Ensure WebView2 is installed
winget install Microsoft.WebView2
```

**Problem**: CGO compilation errors on Linux
**Solution**:
```bash
sudo apt-get install build-essential libgtk-3-dev libwebkit2gtk-4.0-dev
```

#### 2. Code Signing Issues

**Problem**: macOS code signing fails with "no identity found"
**Solution**:
```bash
# List available identities
security find-identity -p codesigning -v

# Unlock keychain
security unlock-keychain ~/Library/Keychains/login.keychain-db
```

**Problem**: Windows signing fails with "signtool not found"
**Solution**:
- Install Windows SDK
- Add signtool.exe to PATH
- Use full path to signtool.exe

#### 3. Notarization Issues

**Problem**: macOS notarization fails with "invalid binary"
**Solution**:
- Ensure app is properly code signed first
- Check entitlements.plist configuration
- Verify hardened runtime is enabled

#### 4. CI/CD Issues

**Problem**: GitHub Actions build timeouts
**Solution**:
- Increase timeout values
- Optimize build caching
- Use matrix builds for parallelization

**Problem**: Artifact upload failures
**Solution**:
- Check artifact size limits
- Verify file paths
- Use compressed formats

### Performance Optimization

#### Build Speed

1. **Use build caching**:
   ```yaml
   - uses: actions/cache@v4
     with:
       path: |
         ~/.cache/go-build
         ~/go/pkg/mod
       key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
   ```

2. **Parallel builds**:
   ```bash
   # Build platforms in parallel
   ./scripts/build.sh --windows &
   ./scripts/build.sh --linux &
   ./scripts/build.sh --macos &
   wait
   ```

3. **Optimize frontend builds**:
   ```bash
   # Use production build optimizations
   cd frontend && npm run build
   ```

#### Resource Usage

1. **Monitor build resources**:
   ```bash
   # Check build size
   du -sh build/bin/*
   
   # Analyze binary size
   go tool nm -size build/bin/MCPWeaver | head -20
   ```

2. **Optimize binary size**:
   ```bash
   # Use build flags for smaller binaries
   wails build -ldflags "-s -w" -trimpath
   ```

### Debug Information

#### Build Logs

Enable verbose logging:
```bash
# Wails verbose build
wails build -v

# Go verbose build
go build -v

# Script debugging
DEBUG=1 ./scripts/build.sh --all
```

#### Verification Commands

```bash
# Verify macOS app bundle
codesign -dv --verbose=4 MCPWeaver.app
spctl -a -t exec -vv MCPWeaver.app

# Verify Windows executable
signtool.exe verify /pa /v MCPWeaver.exe

# Check Linux binary
file MCPWeaver
ldd MCPWeaver
```

## Support

For build system issues:

1. **Check the troubleshooting section** above
2. **Review GitHub Actions logs** for CI/CD issues
3. **Consult Wails documentation**: https://wails.io/docs/
4. **Create an issue**: Include build logs and environment details

## Contributing

When contributing to the build system:

1. **Test changes** across all supported platforms
2. **Update documentation** for any new features
3. **Maintain backward compatibility** where possible
4. **Follow security best practices** for signing and credentials
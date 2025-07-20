# Code Signing Guide for MCPWeaver

This document provides comprehensive instructions for code signing MCPWeaver binaries for Windows and macOS platforms.

## Overview

Code signing is a security mechanism that:
- Verifies the authenticity of the software publisher
- Ensures the binary hasn't been tampered with since signing
- Enables users to trust the application
- Allows the application to pass operating system security checks

## Prerequisites

### macOS Code Signing

1. **Apple Developer Account**
   - Paid Apple Developer Program membership
   - Developer ID Application certificate

2. **Development Environment**
   - macOS with Xcode command line tools
   - `codesign` and `security` tools available

3. **Certificates**
   - Developer ID Application certificate installed in Keychain
   - Intermediate certificates (usually automatic)

4. **Notarization (Recommended)**
   - Apple ID with app-specific password
   - Team ID from Apple Developer account

### Windows Code Signing

1. **Code Signing Certificate**
   - Extended Validation (EV) certificate (recommended)
   - Standard code signing certificate
   - Certificate from trusted CA (DigiCert, Sectigo, etc.)

2. **Development Environment**
   - Windows with Windows SDK
   - `signtool.exe` available in PATH
   - Alternatively: Cross-signing with wine on Linux/macOS

3. **Hardware Token (for EV certificates)**
   - USB token or HSM for EV certificates
   - Token drivers installed

## Configuration

### Environment Variables

Create a `.env` file or export these variables:

#### macOS Variables
```bash
# Required
export APPLE_DEVELOPER_ID="Developer ID Application: Your Name (TEAM_ID)"

# Optional but recommended
export APPLE_KEYCHAIN_PATH="/path/to/signing.keychain"
export APPLE_KEYCHAIN_PASSWORD="keychain_password"

# For notarization
export APPLE_ID="your-apple-id@example.com"
export APPLE_ID_PASSWORD="app-specific-password"
export APPLE_TEAM_ID="YOUR_TEAM_ID"
```

#### Windows Variables
```bash
# Option 1: Certificate file
export WINDOWS_CERTIFICATE="/path/to/certificate.pfx"
export WINDOWS_CERT_PASSWORD="certificate_password"

# Option 2: Certificate thumbprint (for certificates in Windows store)
export WINDOWS_CERT_THUMBPRINT="SHA1_THUMBPRINT_HERE"

# Optional
export TIMESTAMP_SERVER="http://timestamp.sectigo.com"
```

### Certificate Setup

#### macOS Certificate Installation

1. **Download Certificate from Apple Developer**
   ```bash
   # Download from developer.apple.com and double-click to install
   # Or import via Keychain Access
   ```

2. **Verify Certificate Installation**
   ```bash
   security find-identity -p codesigning -v
   ```

3. **Create Dedicated Keychain (Optional)**
   ```bash
   security create-keychain -p "password" signing.keychain
   security import certificate.p12 -k signing.keychain -P "cert_password"
   security set-keychain-settings -t 3600 -l signing.keychain
   ```

#### Windows Certificate Installation

1. **Install Certificate**
   ```cmd
   # For .pfx/.p12 files
   certlm.msc  # Open Certificate Manager
   # Import into Personal > Certificates
   ```

2. **Verify Certificate**
   ```cmd
   certlm.msc
   # Check Personal > Certificates for your code signing cert
   ```

3. **Get Certificate Thumbprint**
   ```cmd
   certutil -store My
   # Find your certificate and copy the SHA1 thumbprint
   ```

## Usage

### Basic Signing

#### macOS
```bash
# Sign an application bundle
./scripts/sign-code.sh --platform macos --binary ./build/MCPWeaver.app

# Sign a single binary
./scripts/sign-code.sh --platform macos --binary ./build/mcpweaver-darwin-amd64
```

#### Windows
```bash
# Sign an executable
./scripts/sign-code.sh --platform windows --binary ./build/mcpweaver-windows-amd64.exe

# With specific certificate
./scripts/sign-code.sh --platform windows --binary ./build/mcpweaver.exe \
    --certificate /path/to/cert.pfx --certificate-password "password"
```

### Advanced Options

#### Custom Entitlements (macOS)
```bash
./scripts/sign-code.sh --platform macos --binary ./build/MCPWeaver.app \
    --entitlements ./custom-entitlements.plist
```

#### Custom Timestamp Server
```bash
./scripts/sign-code.sh --platform windows --binary ./build/mcpweaver.exe \
    --timestamp-url "http://timestamp.digicert.com"
```

#### Using Keychain (macOS)
```bash
./scripts/sign-code.sh --platform macos --binary ./build/MCPWeaver.app \
    --keychain /path/to/signing.keychain --keychain-password "password"
```

### Verification

#### Verify macOS Signature
```bash
# Basic verification
codesign --verify --deep ./build/MCPWeaver.app

# Detailed verification
codesign --verify --deep --strict --verbose=2 ./build/MCPWeaver.app

# Check signature details
codesign -dv --verbose=4 ./build/MCPWeaver.app

# Gatekeeper assessment
spctl --assess --type execute ./build/MCPWeaver.app
```

#### Verify Windows Signature
```bash
# Basic verification
signtool.exe verify /pa ./build/mcpweaver.exe

# Detailed verification
signtool.exe verify /pa /v ./build/mcpweaver.exe

# Check signature details
signtool.exe verify /pa /v /d ./build/mcpweaver.exe
```

## Troubleshooting

### Common macOS Issues

#### "No identity found" Error
```bash
# List available identities
security find-identity -p codesigning -v

# Check keychain access
security list-keychains

# Unlock keychain
security unlock-keychain ~/Library/Keychains/login.keychain
```

#### "Resource fork, Finder information, or similar detritus not allowed"
```bash
# Clean extended attributes
xattr -cr ./build/MCPWeaver.app

# Re-sign after cleaning
./scripts/sign-code.sh --platform macos --binary ./build/MCPWeaver.app
```

#### Notarization Issues
```bash
# Check notarization status
xcrun altool --notarization-info REQUEST_UUID \
    --username "apple-id@example.com" \
    --password "app-specific-password"

# Get detailed notarization log
xcrun altool --notarization-info REQUEST_UUID \
    --username "apple-id@example.com" \
    --password "app-specific-password" --verbose
```

### Common Windows Issues

#### "SignTool Error: No certificates were found that met all the given criteria"
```cmd
# List available certificates
certutil -store My

# Verify certificate thumbprint
certutil -store My "THUMBPRINT_HERE"

# Check certificate validity
certutil -verify certificate.pfx
```

#### Timestamp Server Issues
```bash
# Try alternative timestamp servers
--timestamp-url "http://timestamp.digicert.com"
--timestamp-url "http://timestamp.globalsign.com/scripts/timstamp.dll"
--timestamp-url "http://timestamp.sectigo.com"
```

#### Cross-Platform Signing (Linux/macOS to Windows)
```bash
# Install wine
sudo apt-get install wine  # Ubuntu
brew install wine          # macOS

# Install Windows SDK in wine
# Download and install Windows 10 SDK
```

## Security Best Practices

### Certificate Security

1. **Store Certificates Securely**
   - Use dedicated keychains/stores
   - Strong passwords
   - Limited access permissions

2. **Certificate Backup**
   - Backup certificates and private keys
   - Store backups securely offline
   - Document recovery procedures

3. **Access Control**
   - Limit who has access to signing certificates
   - Use dedicated signing machines
   - Audit certificate usage

### Signing Process Security

1. **Verify Binaries Before Signing**
   - Check binary integrity
   - Scan for malware
   - Verify build reproducibility

2. **Secure Build Environment**
   - Clean build machines
   - Updated development tools
   - Isolated signing environment

3. **Post-Signing Validation**
   - Always verify signatures
   - Test signed binaries
   - Check with antivirus scanners

## Automation

### CI/CD Integration

#### GitHub Actions Example
```yaml
name: Code Signing
on:
  release:
    types: [published]

jobs:
  sign-macos:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v3
      - name: Import Certificates
        run: |
          echo "${{ secrets.MACOS_CERTIFICATE }}" | base64 --decode > certificate.p12
          security create-keychain -p "${{ secrets.KEYCHAIN_PASSWORD }}" build.keychain
          security import certificate.p12 -k build.keychain -P "${{ secrets.CERTIFICATE_PASSWORD }}"
      - name: Sign Application
        run: |
          export APPLE_DEVELOPER_ID="${{ secrets.APPLE_DEVELOPER_ID }}"
          ./scripts/sign-code.sh --platform macos --binary ./build/MCPWeaver.app

  sign-windows:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v3
      - name: Sign Executable
        run: |
          $env:WINDOWS_CERTIFICATE = "${{ secrets.WINDOWS_CERTIFICATE }}"
          $env:WINDOWS_CERT_PASSWORD = "${{ secrets.WINDOWS_CERT_PASSWORD }}"
          ./scripts/sign-code.sh --platform windows --binary ./build/mcpweaver.exe
```

### Build Script Integration

Add to your build scripts:
```bash
# After building
if [ "$ENABLE_CODESIGNING" = "true" ]; then
    case "$GOOS" in
        "darwin")
            ./scripts/sign-code.sh --platform macos --binary "$OUTPUT_BINARY"
            ;;
        "windows")
            ./scripts/sign-code.sh --platform windows --binary "$OUTPUT_BINARY"
            ;;
    esac
fi
```

## Certificate Management

### Certificate Renewal

1. **Monitor Expiry Dates**
   ```bash
   # macOS
   security find-certificate -c "Developer ID Application" -p | openssl x509 -text -noout

   # Windows
   certutil -store My "THUMBPRINT" | findstr "NotAfter"
   ```

2. **Renewal Process**
   - Start renewal 30 days before expiry
   - Test new certificate thoroughly
   - Update build systems
   - Archive old certificates securely

### Multiple Certificate Management

For organizations with multiple certificates:

1. **Certificate Inventory**
   - Maintain certificate database
   - Track expiry dates
   - Document usage

2. **Role-Based Access**
   - Different certificates for different teams
   - Separate development and production certificates
   - Audit trails

## Performance Considerations

### Signing Performance

1. **Parallel Signing**
   - Sign multiple binaries in parallel
   - Use build matrix for different platforms

2. **Caching**
   - Cache intermediate certificates
   - Reuse keychain sessions

3. **Optimization**
   - Sign only final binaries
   - Skip signing for development builds
   - Use local timestamp servers when possible

### Build Time Impact

- macOS notarization can take 5-15 minutes
- Windows signing typically takes seconds
- Plan build times accordingly
- Consider async notarization for faster builds

## Legal and Compliance

### Code Signing Policies

1. **Internal Policies**
   - Who can sign code
   - When to sign (all releases vs. major releases)
   - Certificate management procedures

2. **Compliance Requirements**
   - Industry-specific requirements
   - Security standards (SOC 2, ISO 27001)
   - Audit requirements

3. **Distribution Policies**
   - App store requirements
   - Enterprise distribution policies
   - Open source considerations

## Resources

### Official Documentation

- [Apple Code Signing Guide](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/)
- [Apple Notarization Guide](https://developer.apple.com/documentation/xcode/notarizing_macos_software_before_distribution)
- [Microsoft Authenticode](https://docs.microsoft.com/en-us/windows-hardware/drivers/install/authenticode)
- [Windows Code Signing](https://docs.microsoft.com/en-us/windows/win32/seccrypto/cryptography-tools)

### Tools and Utilities

- [Keychain Access](https://support.apple.com/guide/keychain-access/welcome/mac) (macOS)
- [Certificate Manager](https://docs.microsoft.com/en-us/dotnet/framework/tools/certmgr-exe-certificate-manager-tool) (Windows)
- [SignTool](https://docs.microsoft.com/en-us/windows/win32/seccrypto/signtool) (Windows)
- [Codesign](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html) (macOS)

### Certificate Authorities

- [DigiCert](https://www.digicert.com/code-signing/)
- [Sectigo](https://sectigo.com/ssl-certificates-tls/code-signing)
- [GlobalSign](https://www.globalsign.com/en/code-signing-certificate/)
- [Entrust](https://www.entrust.com/digital-security/certificate-solutions/products/digital-certificates/code-signing-certificates)

## Support

For issues with code signing:

1. Check this documentation first
2. Review error messages carefully
3. Verify certificate validity and permissions
4. Test with sample applications
5. Contact certificate authority support if needed

## Changelog

### Version 1.0.0
- Initial code signing documentation
- Support for macOS and Windows platforms
- Basic automation examples
- Troubleshooting guide
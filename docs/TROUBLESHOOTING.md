# MCPWeaver Troubleshooting Guide

This guide helps you diagnose and resolve common issues with MCPWeaver. If you can't find a solution here, check our [FAQ](FAQ.md) or [report an issue](https://github.com/matoval/MCPWeaver/issues).

## Table of Contents

- [Installation Issues](#installation-issues)
- [Application Startup Problems](#application-startup-problems)
- [OpenAPI Import and Validation Issues](#openapi-import-and-validation-issues)
- [Generation Failures](#generation-failures)
- [Performance Issues](#performance-issues)
- [Network and Connectivity Problems](#network-and-connectivity-problems)
- [File System Issues](#file-system-issues)
- [Platform-Specific Issues](#platform-specific-issues)
- [Database Issues](#database-issues)
- [UI and Display Problems](#ui-and-display-problems)
- [Advanced Diagnostics](#advanced-diagnostics)
- [Getting Additional Help](#getting-additional-help)

## Installation Issues

### Windows Installation Problems

#### "Windows protected your PC" Warning

**Problem**: Windows Defender SmartScreen blocks MCPWeaver installation.

**Solution**:
1. Click "More info" in the SmartScreen dialog
2. Click "Run anyway" to proceed with installation
3. Alternatively, download from the official GitHub releases page and verify checksums

**Prevention**: Future releases will have extended validation (EV) code signing to avoid this warning.

#### "The application was unable to start correctly (0xc000007b)"

**Problem**: Missing Visual C++ Redistributable or corrupted system files.

**Solution**:
```powershell
# Install Visual C++ Redistributable
# Download from: https://aka.ms/vs/17/release/vc_redist.x64.exe

# Run System File Checker
sfc /scannow

# Update Windows
# Go to Settings > Update & Security > Windows Update
```

#### Installation Directory Permission Issues

**Problem**: Access denied when installing to Program Files.

**Solution**:
```powershell
# Option 1: Run installer as Administrator
# Right-click installer > "Run as administrator"

# Option 2: Install to user directory
# Choose a directory like C:\Users\YourName\MCPWeaver during installation

# Option 3: Use portable installation
# Download the ZIP version instead of the installer
```

#### WebView2 Runtime Missing

**Problem**: Application fails to start due to missing WebView2.

**Solution**:
```powershell
# Download WebView2 Evergreen Standalone Installer
# From: https://developer.microsoft.com/microsoft-edge/webview2/

# Or install via PowerShell
winget install Microsoft.EdgeWebView2Runtime
```

### macOS Installation Problems

#### "MCPWeaver cannot be opened because it is from an unidentified developer"

**Problem**: Gatekeeper prevents unsigned applications from running.

**Solution**:
```bash
# Option 1: Remove quarantine attribute
sudo xattr -rd com.apple.quarantine /Applications/MCPWeaver.app

# Option 2: Override in System Preferences
# Go to System Preferences > Security & Privacy > General
# Click "Open Anyway" next to the MCPWeaver message
```

#### "MCPWeaver is damaged and can't be opened"

**Problem**: File corruption during download or false positive from Gatekeeper.

**Solution**:
```bash
# 1. Re-download from official source
curl -L -o MCPWeaver.dmg "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-macos-universal.dmg"

# 2. Verify checksum
shasum -a 256 MCPWeaver.dmg
# Compare with official checksums.sha256

# 3. Clear quarantine after verification
sudo xattr -rd com.apple.quarantine MCPWeaver.dmg
```

#### DMG Won't Mount

**Problem**: "No mountable filesystems" error when opening DMG.

**Solution**:
```bash
# Check if DMG is corrupted
hdiutil verify MCPWeaver.dmg

# Force mount if verification passes
hdiutil attach MCPWeaver.dmg -force

# If still failing, re-download the DMG
```

### Linux Installation Problems

#### AppImage Permission Denied

**Problem**: Cannot execute AppImage file.

**Solution**:
```bash
# Make AppImage executable
chmod +x MCPWeaver-linux-amd64.AppImage

# Check file ownership
ls -la MCPWeaver-linux-amd64.AppImage

# If owned by root, change ownership
sudo chown $USER:$USER MCPWeaver-linux-amd64.AppImage
```

#### "No such file or directory" for AppImage

**Problem**: FUSE not available for AppImage support.

**Solution**:
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install fuse

# Fedora/CentOS/RHEL
sudo dnf install fuse

# Arch Linux
sudo pacman -S fuse2

# If FUSE installation fails, extract and run manually
./MCPWeaver-linux-amd64.AppImage --appimage-extract
./squashfs-root/AppRun
```

#### Missing GTK/WebKit Libraries

**Problem**: Application fails to start due to missing system libraries.

**Solution**:
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37 libglib2.0-0

# Fedora/CentOS/RHEL
sudo dnf install gtk3 webkit2gtk3 glib2

# Arch Linux
sudo pacman -S gtk3 webkit2gtk glib2

# openSUSE
sudo zypper install gtk3-devel webkit2gtk3-devel glib2-devel
```

#### Library Version Conflicts

**Problem**: Incompatible library versions causing crashes.

**Solution**:
```bash
# Check library versions
ldd MCPWeaver-linux-amd64.AppImage

# Install compatible versions
# Ubuntu 18.04+
sudo apt install libgtk-3-0=3.22* libwebkit2gtk-4.0-37=2.24*

# If issues persist, use the portable tar.gz version
wget "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-linux-amd64.tar.gz"
tar -xzf MCPWeaver-linux-amd64.tar.gz
cd MCPWeaver
./mcpweaver
```

## Application Startup Problems

### Application Won't Start

#### Database Initialization Failure

**Problem**: MCPWeaver fails to start with database-related errors.

**Solution**:
```bash
# Check if database file is corrupted
# Location varies by platform:
# Windows: %APPDATA%\MCPWeaver\mcpweaver.db
# macOS: ~/Library/Application Support/MCPWeaver/mcpweaver.db  
# Linux: ~/.config/MCPWeaver/mcpweaver.db

# Backup and reset database
mv mcpweaver.db mcpweaver.db.backup
# Restart MCPWeaver to create new database
```

#### Insufficient Permissions

**Problem**: MCPWeaver can't create necessary files or directories.

**Solution**:
```bash
# Check application data directory permissions
# Windows
icacls "%APPDATA%\MCPWeaver" /grant %USERNAME%:F

# macOS/Linux
chmod -R 755 ~/.config/MCPWeaver  # Linux
chmod -R 755 "~/Library/Application Support/MCPWeaver"  # macOS
```

#### Port Conflict

**Problem**: Development server port already in use.

**Solution**:
```bash
# Find process using port (usually 34115)
# Windows
netstat -ano | findstr :34115

# macOS/Linux  
lsof -i :34115

# Kill the conflicting process or restart MCPWeaver
```

### Slow Startup

#### Performance Optimization

**Problem**: MCPWeaver takes longer than 5 seconds to start.

**Solution**:
```bash
# Clear application cache
# Windows: %APPDATA%\MCPWeaver\cache\
# macOS: ~/Library/Caches/MCPWeaver/
# Linux: ~/.cache/MCPWeaver/

# Reduce recent files list
# Open MCPWeaver > Settings > Clear Recent Files

# Check available system resources
# Ensure at least 2GB RAM available
# Close unnecessary applications
```

#### Antivirus Interference

**Problem**: Antivirus software slowing down startup.

**Solution**:
1. Add MCPWeaver installation directory to antivirus exclusions
2. Add MCPWeaver process to exclusions
3. Temporarily disable real-time scanning to test

**Directories to exclude**:
- Installation directory (e.g., `C:\Program Files\MCPWeaver`)
- Data directory (e.g., `%APPDATA%\MCPWeaver`)
- Temp directories used during generation

## OpenAPI Import and Validation Issues

### Import Failures

#### "Invalid OpenAPI Specification" Error

**Problem**: Specification fails to import with validation errors.

**Diagnostic Steps**:
```bash
# 1. Check if file is valid JSON/YAML
# For JSON files:
python -m json.tool your-spec.json

# For YAML files:
python -c "import yaml; yaml.safe_load(open('your-spec.yaml'))"

# 2. Validate against OpenAPI schema
# Use online validator: https://editor.swagger.io/
# Or install swagger-codegen-cli locally
```

**Common Issues and Fixes**:

**Missing Required Fields**:
```yaml
# ❌ Invalid - missing required fields
{
  "paths": {}
}

# ✅ Valid - includes required fields
{
  "openapi": "3.0.0",
  "info": {
    "title": "My API",
    "version": "1.0.0"
  },
  "paths": {}
}
```

**Invalid Path Definitions**:
```yaml
# ❌ Invalid - path doesn't start with /
paths:
  "users":
    get: {}

# ✅ Valid - path starts with /
paths:
  "/users":
    get: {}
```

**Unsupported OpenAPI Version**:
```yaml
# ❌ Unsupported - OpenAPI 2.0 (Swagger)
swagger: "2.0"
info:
  title: "My API"
  version: "1.0.0"

# ✅ Supported - OpenAPI 3.0+
openapi: "3.0.0"
info:
  title: "My API"
  version: "1.0.0"
```

#### Network Import Issues

**Problem**: Cannot import from URLs.

**Solution**:
```bash
# Test URL accessibility
curl -I "https://your-api.com/openapi.json"

# Check for SSL certificate issues
curl -k "https://your-api.com/openapi.json"

# If behind corporate firewall/proxy
# Configure proxy settings in MCPWeaver:
# Settings > Network > Proxy Configuration
```

**Common Network Issues**:
- **SSL Certificate Errors**: Update system certificates or use HTTP for testing
- **Proxy/Firewall**: Configure proxy settings or request firewall rules
- **Rate Limiting**: Wait and retry, or download file manually
- **Authentication Required**: Download file manually with credentials

#### Large File Import

**Problem**: Import fails for large OpenAPI specifications.

**Solution**:
```bash
# Check file size (MCPWeaver limit: 10MB)
ls -lh your-spec.json

# If too large, try to optimize:
# 1. Remove examples and descriptions
# 2. Split into multiple smaller specs
# 3. Use $ref to external schemas

# Temporary workaround: Use command line tools
# openapi-generator-cli generate -i large-spec.yaml -g go-server
```

### Validation Problems

#### False Positive Validation Errors

**Problem**: Valid specification shows validation errors.

**Solution**:
```bash
# 1. Clear validation cache
# Settings > Advanced > Clear Validation Cache

# 2. Update MCPWeaver to latest version
# Check for updates in Help > About

# 3. Test with external validator
# Upload to https://editor.swagger.io/
# Use swagger-codegen validate command
```

#### Slow Validation

**Problem**: Validation takes longer than expected.

**Solution**:
```bash
# 1. Enable validation caching
# Settings > Validation > Enable Cache

# 2. Reduce specification complexity
# Remove unnecessary examples and descriptions

# 3. Check system resources
# Ensure sufficient RAM and CPU available
```

## Generation Failures

### MCP Server Generation Issues

#### Generation Process Hangs

**Problem**: Generation starts but never completes.

**Diagnostic Steps**:
```bash
# 1. Check generation logs
# Windows: %APPDATA%\MCPWeaver\logs\generation.log
# macOS: ~/Library/Logs/MCPWeaver/generation.log
# Linux: ~/.local/share/MCPWeaver/logs/generation.log

# 2. Monitor system resources
# Task Manager (Windows) / Activity Monitor (macOS) / htop (Linux)

# 3. Check output directory permissions
ls -la /path/to/output/directory
```

**Solutions**:
```bash
# 1. Cancel and retry generation
# Use "Cancel Generation" button in UI

# 2. Clear generation cache
# Settings > Advanced > Clear Generation Cache

# 3. Use different output directory
# Choose directory with write permissions

# 4. Restart MCPWeaver
# Close and reopen the application
```

#### "Permission Denied" During Generation

**Problem**: Cannot write to output directory.

**Solution**:
```bash
# Check directory permissions
# Windows
icacls "C:\path\to\output" /grant %USERNAME%:F

# macOS/Linux
chmod 755 /path/to/output
sudo chown $USER:$USER /path/to/output

# Use alternative directory
# Choose a directory in your home folder
mkdir ~/mcpweaver-output
```

#### Generated Code Compilation Errors

**Problem**: Generated MCP server has compilation errors.

**Solution**:
```bash
# 1. Check Go installation
go version
# Should be 1.21 or later

# 2. Test with minimal OpenAPI spec
# Use Petstore example to verify generation works

# 3. Check for invalid identifiers in spec
# Ensure operation IDs and schema names are valid Go identifiers

# 4. Review generation template
# Settings > Templates > Reset to Default
```

**Common Spec Issues Causing Generation Errors**:

**Invalid Operation IDs**:
```yaml
# ❌ Invalid - contains spaces and special characters
paths:
  /users:
    get:
      operationId: "get all users"

# ✅ Valid - camelCase identifier
paths:
  /users:
    get:
      operationId: "getAllUsers"
```

**Reserved Keywords**:
```yaml
# ❌ Invalid - uses Go reserved keyword
components:
  schemas:
    type:  # 'type' is reserved in Go
      properties:
        name:
          type: string

# ✅ Valid - avoid reserved keywords
components:
  schemas:
    UserType:
      properties:
        name:
          type: string
```

### Template Issues

#### Custom Template Errors

**Problem**: Custom generation templates cause failures.

**Solution**:
```bash
# 1. Validate template syntax
# Templates use Go text/template syntax

# 2. Reset to default template
# Settings > Templates > Reset to Default

# 3. Test template with simple spec
# Use minimal OpenAPI spec for testing

# 4. Check template permissions
# Ensure template files are readable
```

#### Template Loading Failures

**Problem**: Cannot load or apply templates.

**Solution**:
```bash
# 1. Check template directory
# Settings > Templates > Template Directory

# 2. Verify template file format
# Templates should be .tmpl or .gotmpl files

# 3. Clear template cache
# Settings > Advanced > Clear Template Cache
```

## Performance Issues

### High Memory Usage

#### Memory Leaks

**Problem**: MCPWeaver uses excessive memory over time.

**Solution**:
```bash
# 1. Monitor memory usage
# Task Manager > MCPWeaver process

# 2. Force garbage collection
# Settings > Performance > Force Garbage Collection

# 3. Restart application periodically
# Close and reopen MCPWeaver

# 4. Reduce concurrent operations
# Avoid running multiple generations simultaneously
```

#### Large Specification Handling

**Problem**: High memory usage with large OpenAPI specs.

**Solution**:
```bash
# 1. Increase available memory
# Close other applications
# Ensure at least 4GB RAM available

# 2. Split large specifications
# Break into smaller, focused APIs

# 3. Remove unnecessary content
# Delete examples, descriptions, and unused schemas

# 4. Use streaming validation
# Settings > Validation > Enable Streaming Mode
```

### Slow Generation

#### Performance Optimization

**Problem**: Generation takes longer than expected.

**Benchmarks**:
- Small APIs (< 10 endpoints): < 1 second
- Medium APIs (10-100 endpoints): < 3 seconds  
- Large APIs (100+ endpoints): < 10 seconds

**Solutions**:
```bash
# 1. Check system specifications
# CPU: Multi-core recommended
# RAM: 4GB+ recommended
# Storage: SSD recommended

# 2. Close unnecessary applications
# Free up CPU and memory resources

# 3. Disable antivirus real-time scanning
# Temporarily disable for MCPWeaver directory

# 4. Use local specifications
# Avoid network imports during generation

# 5. Enable generation caching
# Settings > Generation > Enable Cache
```

#### Disk I/O Bottlenecks

**Problem**: Slow file operations during generation.

**Solution**:
```bash
# 1. Use SSD storage
# Move output directory to SSD

# 2. Exclude from antivirus scanning
# Add output directory to exclusions

# 3. Close file indexing services
# Windows: Disable Windows Search for output directory
# macOS: Add to Spotlight Privacy exclusions

# 4. Use faster directory
# Choose output directory on fastest available drive
```

## Network and Connectivity Problems

### URL Import Issues

#### Connection Timeouts

**Problem**: URL imports fail with timeout errors.

**Solution**:
```bash
# 1. Test network connectivity
ping 8.8.8.8
curl -I https://www.google.com

# 2. Check firewall settings
# Ensure MCPWeaver can access internet
# Add MCPWeaver to firewall exceptions

# 3. Configure proxy settings
# Settings > Network > Proxy Configuration
# Use system proxy settings

# 4. Increase timeout values
# Settings > Network > Request Timeout (default: 30s)
```

#### SSL/TLS Certificate Errors

**Problem**: HTTPS imports fail with certificate errors.

**Solution**:
```bash
# 1. Update system certificates
# Windows: Windows Update
# macOS: Software Update
# Linux: sudo apt update && sudo apt upgrade ca-certificates

# 2. Test with curl
curl -v "https://your-api.com/openapi.json"

# 3. Temporarily use HTTP for testing
# Replace https:// with http:// if available

# 4. Download file manually
# Save to local file and import
```

#### Corporate Proxy Issues

**Problem**: Cannot access external URLs behind corporate proxy.

**Solution**:
```bash
# 1. Configure proxy in MCPWeaver
# Settings > Network > Proxy Configuration
# Enter proxy server, port, username, password

# 2. Test proxy configuration
curl --proxy proxy.company.com:8080 "https://api.example.com/openapi.json"

# 3. Contact IT department
# Request firewall rules for MCPWeaver
# Get correct proxy settings

# 4. Use VPN if available
# Connect to VPN that bypasses proxy
```

### API Rate Limiting

#### Too Many Requests

**Problem**: API servers reject requests due to rate limiting.

**Solution**:
```bash
# 1. Implement retry with backoff
# MCPWeaver automatically retries with exponential backoff

# 2. Cache responses
# Settings > Network > Enable Response Caching

# 3. Download specifications manually
# Save to local file for repeated use

# 4. Contact API provider
# Request higher rate limits or API key
```

## File System Issues

### File Permission Problems

#### Cannot Read/Write Files

**Problem**: File operations fail with permission errors.

**Solution**:
```bash
# Windows
# Run as Administrator or change file permissions
icacls "path\to\file" /grant %USERNAME%:F

# macOS/Linux
# Change file ownership and permissions
sudo chown $USER:$USER /path/to/file
chmod 644 /path/to/file  # for read/write
chmod 755 /path/to/directory  # for directories
```

#### File Locking Issues

**Problem**: "File in use" errors during operations.

**Solution**:
```bash
# 1. Close other applications using the file
# Check Task Manager / Activity Monitor

# 2. Use different file location
# Copy file to different directory

# 3. Restart system
# Clear all file locks

# 4. Check for background processes
# Antivirus, backup software, file indexers
```

### Storage Space

#### Insufficient Disk Space

**Problem**: Operations fail due to lack of storage.

**Solution**:
```bash
# 1. Check available space
# Windows: dir
# macOS/Linux: df -h

# 2. Clean up temporary files
# Windows: %TEMP%\MCPWeaver\
# macOS: /tmp/MCPWeaver/
# Linux: /tmp/MCPWeaver/

# 3. Change output directory
# Choose location with more space

# 4. Clean old projects
# Remove unused projects from MCPWeaver
```

#### Path Length Limitations

**Problem**: File paths too long for file system.

**Solution**:
```bash
# Windows: Enable long path support
# Group Policy: Computer Configuration > Administrative Templates 
# > System > Filesystem > Enable Win32 long paths

# Alternative: Use shorter output paths
# Example: C:\Out\ instead of C:\Users\LongUsername\Documents\MCPWeaver\Projects\
```

## Platform-Specific Issues

### Windows-Specific Problems

#### Registry Issues

**Problem**: Application settings not persisting.

**Solution**:
```powershell
# Check registry permissions
# HKEY_CURRENT_USER\Software\MCPWeaver

# Reset registry entries
reg delete "HKEY_CURRENT_USER\Software\MCPWeaver" /f
# Restart MCPWeaver to recreate settings
```

#### COM Component Errors

**Problem**: WebView2 COM errors on startup.

**Solution**:
```powershell
# Re-register WebView2
regsvr32 "C:\Program Files (x86)\Microsoft\EdgeWebView2\Application\*\msedgewebview2.exe"

# Update WebView2 runtime
winget upgrade Microsoft.EdgeWebView2Runtime
```

### macOS-Specific Problems

#### Keychain Access Issues

**Problem**: Cannot save secure settings.

**Solution**:
```bash
# Reset keychain access
security delete-generic-password -s "MCPWeaver" -a "$USER"
# Restart MCPWeaver and re-enter credentials
```

#### App Translocation

**Problem**: Application moved by Gatekeeper and can't access files.

**Solution**:
```bash
# Move app to Applications folder
mv ~/Downloads/MCPWeaver.app /Applications/

# Or clear quarantine attribute
sudo xattr -rd com.apple.quarantine /Applications/MCPWeaver.app
```

### Linux-Specific Problems

#### Desktop Integration

**Problem**: Application doesn't appear in application menu.

**Solution**:
```bash
# Create desktop entry
cat > ~/.local/share/applications/mcpweaver.desktop << EOF
[Desktop Entry]
Type=Application
Name=MCPWeaver
Exec=/path/to/MCPWeaver.AppImage
Icon=mcpweaver
Categories=Development;
EOF

# Update desktop database
update-desktop-database ~/.local/share/applications/
```

#### Font Rendering Issues

**Problem**: Fonts appear blurry or incorrect.

**Solution**:
```bash
# Install better fonts
sudo apt install fonts-noto fonts-liberation

# Configure font rendering
# Install gnome-tweaks or similar tool
# Adjust font hinting and antialiasing
```

## Database Issues

### Database Corruption

#### Corrupt Database File

**Problem**: Application crashes or fails to load projects.

**Solution**:
```bash
# 1. Backup current database
cp mcpweaver.db mcpweaver.db.backup

# 2. Try to repair database
sqlite3 mcpweaver.db "PRAGMA integrity_check;"

# 3. If repair fails, reset database
rm mcpweaver.db
# Restart MCPWeaver to create new database

# 4. Restore projects from backup if available
# Settings > Projects > Import Projects
```

#### Migration Failures

**Problem**: Database schema migration fails during updates.

**Solution**:
```bash
# 1. Backup database before updating
cp mcpweaver.db mcpweaver.db.pre-update

# 2. Check database version
sqlite3 mcpweaver.db "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;"

# 3. Manual migration recovery
# Contact support with database backup
# Restore from backup and try alternative update path
```

### Performance Issues

#### Slow Database Queries

**Problem**: Project loading and operations are slow.

**Solution**:
```bash
# 1. Optimize database
sqlite3 mcpweaver.db "VACUUM; ANALYZE;"

# 2. Clear old data
# Settings > Maintenance > Clean Old Data

# 3. Increase cache size
# Settings > Performance > Database Cache Size

# 4. Monitor database size
ls -lh mcpweaver.db
# If very large (>100MB), consider cleaning old projects
```

## UI and Display Problems

### Layout Issues

#### Interface Elements Overlapping

**Problem**: UI components overlap or are positioned incorrectly.

**Solution**:
```bash
# 1. Reset window settings
# Settings > Window > Reset Window Layout

# 2. Change display scaling
# Windows: Settings > Display > Scale
# macOS: System Preferences > Displays > Resolution
# Linux: Display settings in desktop environment

# 3. Update graphics drivers
# Ensure latest drivers installed

# 4. Try different zoom level
# View > Zoom > Reset to 100%
```

#### Dark/Light Theme Issues

**Problem**: Text unreadable or colors incorrect.

**Solution**:
```bash
# 1. Reset theme settings
# Settings > Appearance > Reset Theme

# 2. Check system theme compatibility
# Ensure system and MCPWeaver themes match

# 3. Update application
# Check for theme fixes in newer versions

# 4. Use contrast mode
# Settings > Accessibility > High Contrast Mode
```

### Performance Issues

#### Slow UI Responsiveness

**Problem**: Interface feels sluggish or unresponsive.

**Solution**:
```bash
# 1. Check system resources
# Ensure sufficient RAM and CPU available

# 2. Disable animations
# Settings > Appearance > Disable Animations

# 3. Reduce UI complexity
# Close unused panels and windows

# 4. Update graphics drivers
# Especially important for integrated graphics

# 5. Check for background processes
# Disable unnecessary startup programs
```

#### High CPU Usage from UI

**Problem**: MCPWeaver uses high CPU when idle.

**Solution**:
```bash
# 1. Disable real-time features
# Settings > General > Disable Auto-refresh

# 2. Check for update loops
# View browser console for repeated requests

# 3. Restart application
# Close and reopen to clear any loops

# 4. Check extensions/plugins
# Disable browser extensions that might interfere
```

## Advanced Diagnostics

### Logging and Debugging

#### Enable Debug Logging

**Windows**:
```powershell
# Set environment variable for debug logging
$env:MCPWEAVER_LOG_LEVEL = "debug"
# Restart MCPWeaver

# View logs
Get-Content "$env:APPDATA\MCPWeaver\logs\mcpweaver.log" -Tail 50 -Wait
```

**macOS**:
```bash
# Enable debug logging
export MCPWEAVER_LOG_LEVEL=debug
open -a MCPWeaver

# View logs
tail -f ~/Library/Logs/MCPWeaver/mcpweaver.log
```

**Linux**:
```bash
# Enable debug logging
export MCPWEAVER_LOG_LEVEL=debug
./MCPWeaver.AppImage

# View logs
tail -f ~/.local/share/MCPWeaver/logs/mcpweaver.log
```

#### Common Log Messages

**Normal Operation**:
```
INFO: Application started successfully
INFO: Database initialized
INFO: Project loaded: project-123
INFO: Validation completed: spec.yaml (0 errors)
INFO: Generation started: project-123
INFO: Generation completed: /path/to/output
```

**Warning Signs**:
```
WARN: Validation took longer than expected: 5.2s
WARN: High memory usage detected: 85%
WARN: Network request retry: attempt 3/5
WARN: Template cache miss for: custom-template
```

**Error Indicators**:
```
ERROR: Failed to connect to database: database locked
ERROR: Generation failed: template compilation error
ERROR: File system error: permission denied
ERROR: Network timeout: request exceeded 30s
```

### Performance Monitoring

#### Built-in Performance Monitor

**Access**: Settings > Performance > Monitor

**Key Metrics**:
- **Memory Usage**: Should stay under 100MB for normal use
- **CPU Usage**: Should be <5% when idle
- **Database Operations**: Response times should be <100ms
- **Generation Times**: Within expected ranges per API size

#### System Resource Monitoring

**Windows**:
```powershell
# Monitor MCPWeaver process
Get-Process MCPWeaver | Select-Object ProcessName, CPU, WorkingSet, PagedMemorySize

# Performance counters
Get-Counter "\Process(MCPWeaver)\% Processor Time"
Get-Counter "\Process(MCPWeaver)\Working Set"
```

**macOS**:
```bash
# Monitor MCPWeaver process
top -pid $(pgrep MCPWeaver)

# Detailed process info
ps -p $(pgrep MCPWeaver) -o pid,ppid,cpu,mem,time,comm
```

**Linux**:
```bash
# Monitor MCPWeaver process
htop -p $(pgrep mcpweaver)

# Resource usage
cat /proc/$(pgrep mcpweaver)/status
```

### Network Diagnostics

#### Test Network Connectivity

```bash
# Test basic connectivity
ping -c 4 8.8.8.8

# Test HTTPS connectivity
curl -I https://api.github.com

# Test specific API endpoint
curl -v "https://petstore3.swagger.io/api/v3/openapi.json"

# Test with proxy (if applicable)
curl --proxy proxy.company.com:8080 "https://api.example.com"
```

#### DNS Resolution Issues

```bash
# Test DNS resolution
nslookup api.example.com
dig api.example.com

# Try alternative DNS servers
# Windows: netsh interface ip set dns "Wi-Fi" static 8.8.8.8
# macOS/Linux: Add nameserver 8.8.8.8 to /etc/resolv.conf
```

### File System Diagnostics

#### Check File Permissions

**Windows**:
```powershell
# Check directory permissions
icacls "C:\path\to\directory"

# Check specific file
Get-Acl "C:\path\to\file.json" | Format-List
```

**macOS/Linux**:
```bash
# Check directory permissions
ls -la /path/to/directory

# Check file permissions
stat /path/to/file.json

# Check disk space
df -h /path/to/directory
```

#### Test File Operations

```bash
# Test read access
# Try opening the file in a text editor

# Test write access
echo "test" > /path/to/test-file.txt
rm /path/to/test-file.txt

# Test directory creation
mkdir /path/to/test-directory
rmdir /path/to/test-directory
```

## Getting Additional Help

### Collecting Diagnostic Information

When reporting issues, please include:

1. **System Information**:
   - Operating system and version
   - MCPWeaver version
   - Available RAM and storage
   - Processor information

2. **Error Details**:
   - Exact error message
   - Steps to reproduce
   - Expected vs actual behavior
   - Screenshots or screen recordings

3. **Log Files**:
   - Application logs (last 100 lines)
   - System error logs
   - Network logs if relevant

4. **Configuration**:
   - Settings that differ from defaults
   - Custom templates in use
   - Proxy/network configuration

### Automated Diagnostic Collection

MCPWeaver includes a diagnostic collection tool:

**Access**: Help > Generate Diagnostic Report

**Includes**:
- System information
- Application logs (last 24 hours)
- Performance metrics
- Configuration summary
- Recent error events

**Privacy**: Personal data like project names and API content are not included.

### Support Channels

1. **GitHub Issues**: [Report bugs and feature requests](https://github.com/matoval/MCPWeaver/issues)
   - Use issue templates
   - Include diagnostic information
   - Search existing issues first

2. **GitHub Discussions**: [Community support](https://github.com/matoval/MCPWeaver/discussions)
   - General questions
   - Usage tips
   - Feature discussions

3. **Documentation**: 
   - [User Guide](USER_GUIDE.md)
   - [API Documentation](API.md)
   - [Installation Guide](INSTALLATION.md)
   - [FAQ](FAQ.md)

### Community Resources

- **Stack Overflow**: Tag questions with `mcpweaver`
- **Reddit**: r/MCPWeaver community
- **Discord**: [Join development discussions](https://discord.gg/mcpweaver)

### Professional Support

For enterprise users requiring dedicated support:
- Priority issue resolution
- Custom feature development
- Integration assistance
- Training and consultation

Contact: support@mcpweaver.dev

---

## Quick Reference

### Emergency Fixes

**Application Won't Start**:
1. Delete database file
2. Clear application cache
3. Reset settings to defaults
4. Reinstall application

**Generation Fails**:
1. Check output directory permissions
2. Validate OpenAPI specification
3. Reset generation templates
4. Try minimal test specification

**Performance Issues**:
1. Force garbage collection
2. Clear all caches
3. Restart application
4. Check system resources

**Network Issues**:
1. Test with curl/wget
2. Check firewall settings
3. Configure proxy if needed
4. Try alternative DNS servers

Remember: When in doubt, try the simplest solution first - restart MCPWeaver and try again with a fresh start! 🔄
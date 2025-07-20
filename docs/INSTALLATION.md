# MCPWeaver Installation Guide

This guide provides detailed instructions for installing MCPWeaver on Windows, macOS, and Linux systems.

## Table of Contents

- [System Requirements](#system-requirements)
- [Download Options](#download-options)
- [Windows Installation](#windows-installation)
- [macOS Installation](#macos-installation)
- [Linux Installation](#linux-installation)
- [Portable Installation](#portable-installation)
- [Building from Source](#building-from-source)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)
- [Uninstallation](#uninstallation)

## System Requirements

### Minimum Requirements

| Component | Requirement |
|-----------|-------------|
| **Operating System** | Windows 10/11, macOS 10.15+, or Linux (Ubuntu 18.04+/equivalent) |
| **Memory (RAM)** | 2 GB available RAM |
| **Storage** | 100 MB free disk space |
| **Architecture** | x64 (Intel/AMD) or ARM64 (Apple Silicon/ARM) |
| **Network** | Internet connection for downloads and URL imports |

### Recommended Requirements

| Component | Recommendation |
|-----------|----------------|
| **Memory (RAM)** | 4 GB or more for large OpenAPI specifications |
| **Storage** | 500 MB for application, templates, and project data |
| **CPU** | Multi-core processor for faster generation |
| **Display** | 1920x1080 or higher resolution |

### Additional Requirements by Platform

**Windows:**
- Windows 10 version 1903 or later
- WebView2 runtime (automatically installed if needed)
- Visual C++ 2019 Redistributable (usually pre-installed)

**macOS:**
- macOS Catalina (10.15) or later
- Gatekeeper must allow the application (handled automatically for signed releases)

**Linux:**
- GTK 3.20 or later
- WebKitGTK 2.24 or later
- glibc 2.28 or later

## Download Options

### Official Releases

Download the latest stable release from the [GitHub Releases page](https://github.com/matoval/MCPWeaver/releases/latest).

**Available Downloads:**
- `MCPWeaver-windows-amd64.exe` - Windows installer
- `MCPWeaver-macos-universal.dmg` - macOS disk image (Intel + Apple Silicon)
- `MCPWeaver-linux-amd64.AppImage` - Linux AppImage
- `MCPWeaver-linux-amd64.tar.gz` - Linux portable archive
- Source code archives for building from source

### Pre-release Versions

For testing the latest features, pre-release versions are available:
- Navigate to [GitHub Releases](https://github.com/matoval/MCPWeaver/releases)
- Look for releases marked as "Pre-release"
- Download the appropriate package for your platform

⚠️ **Warning:** Pre-release versions may contain bugs and are not recommended for production use.

### Checksums

All release files include SHA256 checksums for verification:
- Download the `checksums.sha256` file from the release page
- Verify your download using the appropriate command for your platform

## Windows Installation

### Method 1: Installer (Recommended)

1. **Download the Installer**
   ```powershell
   # Using PowerShell
   Invoke-WebRequest -Uri "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-windows-amd64.exe" -OutFile "MCPWeaver-installer.exe"
   ```
   
   Or download manually from the [releases page](https://github.com/matoval/MCPWeaver/releases/latest).

2. **Run the Installer**
   - Double-click `MCPWeaver-installer.exe`
   - If Windows Defender SmartScreen appears, click "More info" then "Run anyway"
   - Follow the installation wizard prompts
   - Choose installation directory (default: `C:\Program Files\MCPWeaver`)
   - Select additional options (desktop shortcut, start menu entry)
   - Click "Install" to complete the installation

3. **Launch MCPWeaver**
   - Use the desktop shortcut or start menu entry
   - Or run from Command Prompt: `mcpweaver`

### Method 2: Portable Installation

1. **Download Archive**
   ```powershell
   Invoke-WebRequest -Uri "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-windows-amd64.zip" -OutFile "MCPWeaver.zip"
   ```

2. **Extract and Run**
   ```powershell
   Expand-Archive -Path "MCPWeaver.zip" -DestinationPath "C:\MCPWeaver"
   cd C:\MCPWeaver
   .\MCPWeaver.exe
   ```

### Method 3: Package Manager

**Chocolatey:**
```powershell
# Coming soon - package under review
choco install mcpweaver
```

**Winget:**
```powershell
# Coming soon - package under review  
winget install MCPWeaver.MCPWeaver
```

### Windows-Specific Configuration

**File Associations:**
The installer can optionally associate `.mcpweaver` project files with MCPWeaver.

**PATH Environment Variable:**
Add MCPWeaver to your PATH to run from anywhere:
1. Open System Properties > Environment Variables
2. Add `C:\Program Files\MCPWeaver` to your PATH
3. Restart your command prompt

**Windows Defender:**
If Windows Defender blocks the application:
1. Open Windows Security > Virus & threat protection
2. Add MCPWeaver installation directory to exclusions
3. Or submit the file to Microsoft for analysis

## macOS Installation

### Method 1: DMG Package (Recommended)

1. **Download the DMG**
   ```bash
   curl -L -o MCPWeaver.dmg "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-macos-universal.dmg"
   ```

2. **Mount and Install**
   ```bash
   # Mount the DMG
   open MCPWeaver.dmg
   
   # Or mount via command line
   hdiutil attach MCPWeaver.dmg
   ```

3. **Install Application**
   - Drag MCPWeaver.app to the Applications folder
   - Eject the DMG when installation is complete
   - Launch MCPWeaver from Applications or Launchpad

### Method 2: Homebrew

```bash
# Add the tap (repository)
brew tap matoval/mcpweaver

# Install MCPWeaver
brew install --cask mcpweaver

# Launch MCPWeaver
open -a MCPWeaver
```

### Method 3: Portable Installation

```bash
# Download and extract
curl -L -o MCPWeaver.tar.gz "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-macos-universal.tar.gz"
tar -xzf MCPWeaver.tar.gz
mv MCPWeaver.app /Applications/
```

### macOS-Specific Setup

**Gatekeeper Permission:**
If macOS prevents MCPWeaver from running:
1. Go to System Preferences > Security & Privacy > General
2. Click "Open Anyway" next to the MCPWeaver message
3. Or run: `sudo xattr -rd com.apple.quarantine /Applications/MCPWeaver.app`

**CLI Access:**
To run MCPWeaver from the command line:
```bash
# Create symlink
sudo ln -s /Applications/MCPWeaver.app/Contents/MacOS/MCPWeaver /usr/local/bin/mcpweaver

# Or add to PATH in your shell profile
echo 'export PATH="/Applications/MCPWeaver.app/Contents/MacOS:$PATH"' >> ~/.zshrc
```

**Apple Silicon vs Intel:**
The universal binary works on both Intel and Apple Silicon Macs. No special configuration needed.

## Linux Installation

### Method 1: AppImage (Recommended)

AppImage provides a portable, dependency-free installation:

1. **Download AppImage**
   ```bash
   wget "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-linux-amd64.AppImage"
   ```

2. **Make Executable and Run**
   ```bash
   chmod +x MCPWeaver-linux-amd64.AppImage
   ./MCPWeaver-linux-amd64.AppImage
   ```

3. **Optional: Desktop Integration**
   ```bash
   # Install AppImageLauncher for better desktop integration
   # Ubuntu/Debian:
   sudo apt install appimagelauncher
   
   # Or manually create desktop entry
   ./MCPWeaver-linux-amd64.AppImage --appimage-extract
   cp squashfs-root/mcpweaver.desktop ~/.local/share/applications/
   ```

### Method 2: Portable Archive

1. **Download and Extract**
   ```bash
   wget "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-linux-amd64.tar.gz"
   tar -xzf MCPWeaver-linux-amd64.tar.gz
   cd MCPWeaver
   ```

2. **Install Dependencies (if needed)**
   ```bash
   # Ubuntu/Debian
   sudo apt update
   sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37 libglib2.0-0
   
   # CentOS/RHEL/Fedora
   sudo dnf install gtk3 webkit2gtk3 glib2
   
   # Arch Linux
   sudo pacman -S gtk3 webkit2gtk glib2
   ```

3. **Run MCPWeaver**
   ```bash
   ./mcpweaver
   ```

### Method 3: Package Managers

**Snap Package:**
```bash
# Coming soon
sudo snap install mcpweaver
```

**Flatpak:**
```bash
# Coming soon
flatpak install flathub com.github.matoval.MCPWeaver
```

**AUR (Arch Linux):**
```bash
# Coming soon
yay -S mcpweaver-bin
```

### Distribution-Specific Instructions

**Ubuntu/Debian:**
```bash
# Install dependencies
sudo apt update
sudo apt install wget

# Download and install
wget "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-linux-amd64.AppImage"
chmod +x MCPWeaver-linux-amd64.AppImage
sudo mv MCPWeaver-linux-amd64.AppImage /opt/mcpweaver
sudo ln -s /opt/mcpweaver /usr/local/bin/mcpweaver
```

**CentOS/RHEL/Fedora:**
```bash
# Install dependencies  
sudo dnf install wget

# Download and install
wget "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-linux-amd64.AppImage"
chmod +x MCPWeaver-linux-amd64.AppImage
sudo mv MCPWeaver-linux-amd64.AppImage /opt/mcpweaver
sudo ln -s /opt/mcpweaver /usr/bin/mcpweaver
```

**Arch Linux:**
```bash
# Using AUR helper
yay -S mcpweaver-bin

# Or manual installation
wget "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-linux-amd64.AppImage"
chmod +x MCPWeaver-linux-amd64.AppImage
sudo mv MCPWeaver-linux-amd64.AppImage /usr/local/bin/mcpweaver
```

## Portable Installation

For situations where you can't install software system-wide, MCPWeaver supports portable installations:

### Creating a Portable Installation

1. **Download Portable Package**
   - Windows: Download the ZIP archive instead of the installer
   - macOS: Extract the app from the DMG to any folder
   - Linux: Use the AppImage or tar.gz archive

2. **Setup Portable Directory**
   ```
   MCPWeaver-Portable/
   ├── MCPWeaver(.exe)      # Application binary
   ├── templates/           # Custom templates (optional)
   ├── projects/           # Project storage
   └── config/            # Configuration files
   ```

3. **Configure Portable Mode**
   Create a `portable.txt` file in the same directory as the MCPWeaver binary to enable portable mode.

### Portable Mode Features

- All data stored in application directory
- No registry/system configuration changes
- Can run from USB drives or network locations
- Settings and projects travel with the application

### USB Drive Installation

1. **Format USB Drive** (recommended: exFAT for cross-platform compatibility)
2. **Copy MCPWeaver** to the USB drive
3. **Create portable.txt** file
4. **Test on Different Systems** to ensure compatibility

## Building from Source

For developers or users who want to build MCPWeaver from source:

### Prerequisites

- **Go 1.21 or later**
- **Node.js 18 or later** 
- **Wails CLI v2.10.1 or later**
- **Git**

### Platform-Specific Build Requirements

**Windows:**
- WebView2 SDK
- CGO-compatible C compiler (TDM-GCC or Visual Studio)

**macOS:**
- Xcode Command Line Tools
- macOS SDK 10.15+

**Linux:**
- GTK 3 development libraries
- WebKit2GTK development libraries
- GCC/Clang compiler

### Build Instructions

1. **Clone Repository**
   ```bash
   git clone https://github.com/matoval/MCPWeaver.git
   cd MCPWeaver
   ```

2. **Install Dependencies**
   ```bash
   # Install Wails CLI
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   
   # Install frontend dependencies
   cd frontend
   npm install
   cd ..
   ```

3. **Build Application**
   ```bash
   # Development build
   wails build
   
   # Production build with optimizations
   wails build -clean -ldflags "-s -w"
   
   # Cross-platform build
   wails build -platform windows/amd64,darwin/universal,linux/amd64
   ```

4. **Build Output**
   Built applications are in the `build/bin/` directory.

### Custom Build Options

**Build with Custom Features:**
```bash
# Build with specific features
wails build -tags "feature1,feature2"

# Build with custom version
wails build -ldflags "-X main.version=1.0.0-custom"
```

**Development Build:**
```bash
# Run in development mode
wails dev

# Build for development (faster, with debug info)
wails build -devMode
```

## Verification

### Verify Installation

After installation, verify MCPWeaver is working correctly:

1. **Check Version**
   ```bash
   mcpweaver --version
   ```

2. **Run Basic Test**
   ```bash
   mcpweaver --help
   ```

3. **Launch GUI**
   - Start MCPWeaver from your applications menu
   - Verify the interface loads correctly
   - Check for any error messages

### Verify Download Integrity

**Windows (PowerShell):**
```powershell
# Calculate SHA256 hash
Get-FileHash -Algorithm SHA256 MCPWeaver-installer.exe
# Compare with checksums.sha256 file
```

**macOS/Linux:**
```bash
# Calculate SHA256 hash
shasum -a 256 MCPWeaver-*.dmg
# Or for Linux
sha256sum MCPWeaver-*.AppImage

# Compare with checksums.sha256 file
```

### Test Installation

1. **Create Test Project**
   - Launch MCPWeaver
   - Create a new project
   - Import a sample OpenAPI specification
   - Verify validation works

2. **Generate Test Server**
   - Use the built-in Petstore API example
   - Generate an MCP server
   - Check that files are created correctly

## Troubleshooting

### Common Installation Issues

#### Windows Issues

**"Windows protected your PC" message:**
- Click "More info" then "Run anyway"
- Or download from official sources and verify checksums

**"The application was unable to start correctly" error:**
- Install Visual C++ 2019 Redistributable
- Update Windows to the latest version
- Check Windows Event Viewer for specific error details

**Installation directory permission issues:**
- Run installer as Administrator
- Or install to user directory instead of Program Files

#### macOS Issues

**"MCPWeaver cannot be opened because it is from an unidentified developer":**
```bash
# Remove quarantine attribute
sudo xattr -rd com.apple.quarantine /Applications/MCPWeaver.app
```

**"MCPWeaver is damaged and can't be opened":**
- Download a fresh copy from the official releases page
- Verify the DMG file integrity with checksums

**Application crashes on startup:**
- Check Console.app for crash logs
- Ensure macOS version compatibility
- Try running from Terminal to see error messages

#### Linux Issues

**"No such file or directory" when running AppImage:**
```bash
# Install FUSE for AppImage support
sudo apt install fuse  # Ubuntu/Debian
sudo dnf install fuse  # Fedora
```

**Missing GTK/WebKit libraries:**
```bash
# Install required libraries
sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37  # Ubuntu/Debian
sudo dnf install gtk3 webkit2gtk3                   # Fedora
```

**Permission denied errors:**
```bash
# Ensure AppImage is executable
chmod +x MCPWeaver-*.AppImage

# Check file ownership
ls -la MCPWeaver-*.AppImage
```

### Performance Issues

**Slow startup:**
- Check available system memory
- Disable unnecessary startup programs
- Clear MCPWeaver cache: delete `~/.mcpweaver/cache/`

**High memory usage:**
- Adjust memory limits in settings
- Close unused projects
- Restart MCPWeaver periodically

### Network Issues

**Cannot download updates or import from URLs:**
- Check firewall settings
- Verify internet connection
- Configure proxy settings if needed

**SSL/TLS certificate errors:**
- Update your operating system
- Check system time and date
- Verify corporate proxy/firewall settings

### Getting Help

If you encounter issues not covered here:

1. **Check Documentation**
   - [User Guide](USER_GUIDE.md)
   - [Troubleshooting Guide](TROUBLESHOOTING.md)
   - [FAQ](FAQ.md)

2. **Search Existing Issues**
   - [GitHub Issues](https://github.com/matoval/MCPWeaver/issues)
   - [GitHub Discussions](https://github.com/matoval/MCPWeaver/discussions)

3. **Report New Issues**
   - [Bug Report Template](https://github.com/matoval/MCPWeaver/issues/new?template=bug_report.md)
   - Include: OS version, MCPWeaver version, error messages, steps to reproduce

4. **Community Support**
   - [GitHub Discussions](https://github.com/matoval/MCPWeaver/discussions)
   - Stack Overflow (tag: mcpweaver)

## Uninstallation

### Windows Uninstallation

**Using the Uninstaller:**
1. Go to Settings > Apps > Apps & features
2. Find MCPWeaver in the list
3. Click "Uninstall" and follow prompts

**Manual Uninstallation:**
```powershell
# Remove application files
Remove-Item -Recurse "C:\Program Files\MCPWeaver"

# Remove user data
Remove-Item -Recurse "$env:APPDATA\MCPWeaver"

# Remove from PATH (if added)
# Edit Environment Variables to remove MCPWeaver paths
```

### macOS Uninstallation

```bash
# Remove application
sudo rm -rf /Applications/MCPWeaver.app

# Remove user data
rm -rf ~/Library/Application\ Support/MCPWeaver
rm -rf ~/Library/Preferences/com.mcpweaver.app.plist
rm -rf ~/Library/Caches/com.mcpweaver.app

# Remove CLI symlink (if created)
sudo rm /usr/local/bin/mcpweaver
```

### Linux Uninstallation

```bash
# Remove AppImage
rm ~/MCPWeaver-*.AppImage

# Or remove from system location
sudo rm /usr/local/bin/mcpweaver

# Remove user data
rm -rf ~/.config/MCPWeaver
rm -rf ~/.local/share/MCPWeaver
rm -rf ~/.cache/MCPWeaver

# Remove desktop entry (if created)
rm ~/.local/share/applications/mcpweaver.desktop
```

### Data Backup Before Uninstallation

Before uninstalling, consider backing up:

**Projects and Settings:**
- Windows: `%APPDATA%\MCPWeaver\`
- macOS: `~/Library/Application Support/MCPWeaver/`
- Linux: `~/.config/MCPWeaver/`

**Custom Templates:**
- Export templates from Template Manager
- Save template files manually from templates directory

**Generated Servers:**
- Copy any generated servers you want to keep
- They are typically in your chosen output directories

---

## Next Steps

After successful installation:

1. **Read the [User Guide](USER_GUIDE.md)** to learn how to use MCPWeaver
2. **Try the Getting Started tutorial** to create your first MCP server
3. **Explore the [API Documentation](API.md)** for advanced usage
4. **Join the community** on [GitHub Discussions](https://github.com/matoval/MCPWeaver/discussions)

Welcome to MCPWeaver! 🎉
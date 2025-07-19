# MCPWeaver Cross-Platform Build Script for Windows
# PowerShell version of the build script

param(
    [switch]$All,
    [switch]$Windows,
    [switch]$MacOS,
    [switch]$Linux,
    [switch]$Clean,
    [string]$Version = "1.0.0",
    [switch]$Help
)

# Configuration
$BuildDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
$BuildCommit = try { git rev-parse --short HEAD } catch { "unknown" }
$BuildDir = "./build/bin"
$DistDir = "./dist"

# Build flags
$LdFlags = "-s -w -X main.version=$Version -X main.buildDate=$BuildDate -X main.buildCommit=$BuildCommit"

# Colors for output
function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Blue }
function Write-Success { param($Message) Write-Host "[SUCCESS] $Message" -ForegroundColor Green }
function Write-Warning { param($Message) Write-Host "[WARNING] $Message" -ForegroundColor Yellow }
function Write-Error { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }

# Function to check if command exists
function Test-Command {
    param($Command)
    $null -ne (Get-Command $Command -ErrorAction SilentlyContinue)
}

# Function to check prerequisites
function Test-Prerequisites {
    Write-Info "Checking prerequisites..."
    
    if (-not (Test-Command "wails")) {
        Write-Error "Wails CLI not found. Please install it first:"
        Write-Error "go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        exit 1
    }
    
    if (-not (Test-Command "npm")) {
        Write-Error "npm not found. Please install Node.js first."
        exit 1
    }
    
    if (-not (Test-Command "go")) {
        Write-Error "Go not found. Please install Go first."
        exit 1
    }
    
    Write-Success "All prerequisites found"
}

# Function to setup build environment
function Initialize-BuildEnvironment {
    Write-Info "Setting up build environment..."
    
    # Create build directories
    if (-not (Test-Path $BuildDir)) { New-Item -ItemType Directory -Path $BuildDir -Force }
    if (-not (Test-Path $DistDir)) { New-Item -ItemType Directory -Path $DistDir -Force }
    
    # Install frontend dependencies
    Write-Info "Installing frontend dependencies..."
    Push-Location frontend
    npm install
    Pop-Location
    
    # Update go dependencies
    Write-Info "Updating Go dependencies..."
    go mod tidy
    
    Write-Success "Build environment ready"
}

# Function to build for a specific platform
function Build-Platform {
    param($Platform, $Arch, $OutputName)
    
    Write-Info "Building for $Platform/$Arch..."
    
    $BuildCommand = "wails build -platform $Platform/$Arch -ldflags `"$LdFlags`" -o `"$BuildDir/$OutputName`" -clean"
    
    try {
        Invoke-Expression $BuildCommand
        Write-Success "Build completed for $Platform/$Arch"
    }
    catch {
        Write-Error "Build failed for $Platform/$Arch"
        Write-Error $_.Exception.Message
        return $false
    }
    
    return $true
}

# Function to package builds
function New-Packages {
    Write-Info "Packaging builds..."
    
    Push-Location $BuildDir
    
    # Package Windows builds
    $WindowsExe = "MCPWeaver-windows-amd64.exe"
    if (Test-Path $WindowsExe) {
        Write-Info "Creating Windows package..."
        Compress-Archive -Path $WindowsExe -DestinationPath "../$DistDir/MCPWeaver-$Version-windows-amd64.zip" -Force
    }
    
    # Package macOS builds
    $MacOSIntel = "MCPWeaver-darwin-amd64.app"
    if (Test-Path $MacOSIntel) {
        Write-Info "Creating macOS Intel package..."
        tar -czf "../$DistDir/MCPWeaver-$Version-darwin-amd64.tar.gz" $MacOSIntel
    }
    
    $MacOSARM = "MCPWeaver-darwin-arm64.app"
    if (Test-Path $MacOSARM) {
        Write-Info "Creating macOS ARM package..."
        tar -czf "../$DistDir/MCPWeaver-$Version-darwin-arm64.tar.gz" $MacOSARM
    }
    
    # Package Linux builds
    $LinuxBin = "MCPWeaver-linux-amd64"
    if (Test-Path $LinuxBin) {
        Write-Info "Creating Linux package..."
        tar -czf "../$DistDir/MCPWeaver-$Version-linux-amd64.tar.gz" $LinuxBin
    }
    
    Pop-Location
    Write-Success "Packaging completed"
}

# Function to create checksums
function New-Checksums {
    Write-Info "Creating checksums..."
    
    Push-Location $DistDir
    
    $Files = Get-ChildItem -Filter "*.zip", "*.tar.gz"
    if ($Files.Count -gt 0) {
        $Checksums = foreach ($File in $Files) {
            $Hash = Get-FileHash -Path $File.FullName -Algorithm SHA256
            "$($Hash.Hash.ToLower())  $($File.Name)"
        }
        $Checksums | Out-File -FilePath "checksums.sha256" -Encoding UTF8
    }
    
    Pop-Location
    Write-Success "Checksums created"
}

# Function to display build summary
function Show-BuildSummary {
    Write-Info "Build Summary:"
    Write-Host "=================="
    Write-Host "Version: $Version"
    Write-Host "Build Date: $BuildDate"
    Write-Host "Build Commit: $BuildCommit"
    Write-Host "=================="
    
    if (Test-Path $DistDir) {
        Write-Host "Distribution files:"
        Get-ChildItem $DistDir | Format-Table Name, Length, LastWriteTime
    }
}

# Show help
if ($Help) {
    Write-Host "MCPWeaver Build Script for Windows"
    Write-Host ""
    Write-Host "Usage: .\build.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -All          Build for all platforms"
    Write-Host "  -Windows      Build for Windows (amd64)"
    Write-Host "  -MacOS        Build for macOS (Intel and ARM)"
    Write-Host "  -Linux        Build for Linux (amd64)"
    Write-Host "  -Clean        Clean build directory before building"
    Write-Host "  -Version VER  Set version number (default: 1.0.0)"
    Write-Host "  -Help         Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\build.ps1 -All                    # Build for all platforms"
    Write-Host "  .\build.ps1 -Windows -Linux         # Build for Windows and Linux"
    Write-Host "  .\build.ps1 -MacOS -Version 1.1.0   # Build macOS with version 1.1.0"
    exit 0
}

# Main execution
Write-Info "Starting MCPWeaver build process..."

# Determine platforms to build
$Platforms = @()
if ($All) {
    $Platforms = @(
        @{OS="windows"; Arch="amd64"; Output="MCPWeaver-windows-amd64.exe"},
        @{OS="darwin"; Arch="amd64"; Output="MCPWeaver-darwin-amd64.app"},
        @{OS="darwin"; Arch="arm64"; Output="MCPWeaver-darwin-arm64.app"},
        @{OS="linux"; Arch="amd64"; Output="MCPWeaver-linux-amd64"}
    )
} else {
    if ($Windows) {
        $Platforms += @{OS="windows"; Arch="amd64"; Output="MCPWeaver-windows-amd64.exe"}
    }
    if ($MacOS) {
        $Platforms += @{OS="darwin"; Arch="amd64"; Output="MCPWeaver-darwin-amd64.app"}
        $Platforms += @{OS="darwin"; Arch="arm64"; Output="MCPWeaver-darwin-arm64.app"}
    }
    if ($Linux) {
        $Platforms += @{OS="linux"; Arch="amd64"; Output="MCPWeaver-linux-amd64"}
    }
}

# Default to current platform if none specified
if ($Platforms.Count -eq 0) {
    Write-Warning "No platforms specified. Building for current platform only."
    wails build -ldflags $LdFlags -clean
    Show-BuildSummary
    Write-Success "Build process completed successfully!"
    exit 0
}

# Clean build directory if requested
if ($Clean) {
    Write-Info "Cleaning build directory..."
    if (Test-Path $BuildDir) { Remove-Item $BuildDir -Recurse -Force }
    if (Test-Path $DistDir) { Remove-Item $DistDir -Recurse -Force }
}

# Run build steps
Test-Prerequisites
Initialize-BuildEnvironment

# Build for each platform
$BuildSuccess = $true
foreach ($Platform in $Platforms) {
    $Result = Build-Platform -Platform $Platform.OS -Arch $Platform.Arch -OutputName $Platform.Output
    if (-not $Result) {
        $BuildSuccess = $false
    }
}

# Package and create checksums if building for multiple platforms
if ($Platforms.Count -gt 1) {
    New-Packages
    New-Checksums
}

Show-BuildSummary

if ($BuildSuccess) {
    Write-Success "Build process completed successfully!"
} else {
    Write-Error "Build process completed with errors!"
    exit 1
}
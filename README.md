# MCPWeaver

[![Build Status](https://github.com/matoval/MCPWeaver/workflows/Build%20and%20Release/badge.svg)](https://github.com/matoval/MCPWeaver/actions)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)
[![Wails](https://img.shields.io/badge/Wails-v2.10.1-red?logo=wails)](https://wails.io/)

**Transform OpenAPI specifications into Model Context Protocol (MCP) servers with a simple, fast desktop application.**

MCPWeaver is an open-source desktop application that converts OpenAPI 3.0 specifications into fully functional MCP servers. Built with Go and React using the Wails framework, it provides a lightweight, user-controlled solution that runs entirely on your machine.

![MCPWeaver Interface](assets/screenshot.png)

## 🚀 Quick Start

1. **Download** the latest release for your platform
2. **Install** MCPWeaver on your system
3. **Import** your OpenAPI specification
4. **Generate** your MCP server
5. **Deploy** and use your generated server

[📖 **View Full Getting Started Guide**](docs/USER_GUIDE.md) | [⬇️ **Download Latest Release**](https://github.com/matoval/MCPWeaver/releases/latest)

## ✨ Key Features

- **🔄 One-Click Generation**: Transform OpenAPI specs to MCP servers instantly
- **⚡ Lightning Fast**: Generate servers in under 5 seconds
- **🖥️ Desktop Native**: No cloud dependencies, works offline
- **🎯 Real-time Validation**: Instant feedback on OpenAPI specifications
- **📁 Project Management**: Track and manage multiple projects
- **🔧 Template Customization**: Customize generation templates
- **🌐 Cross-Platform**: Windows, macOS, and Linux support
- **💾 Lightweight**: < 50MB memory footprint

## 🛠️ Installation

### Quick Install

**Windows:**
```powershell
# Download and run the installer
Invoke-WebRequest -Uri "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-windows-amd64.exe" -OutFile "MCPWeaver-installer.exe"
.\MCPWeaver-installer.exe
```

**macOS:**
```bash
# Download and install
curl -L -o MCPWeaver.dmg "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-macos-universal.dmg"
open MCPWeaver.dmg
```

**Linux:**
```bash
# Download AppImage
wget "https://github.com/matoval/MCPWeaver/releases/latest/download/MCPWeaver-linux-amd64.AppImage"
chmod +x MCPWeaver-linux-amd64.AppImage
./MCPWeaver-linux-amd64.AppImage
```

[📋 **Detailed Installation Instructions**](docs/INSTALLATION.md)

## 🎯 Use Cases

- **API Integration**: Convert REST APIs to MCP for LLM integration
- **Development Tools**: Generate MCP servers for testing and development
- **Legacy System Integration**: Modernize old APIs with MCP
- **Microservices**: Create MCP interfaces for microservice architectures
- **AI/LLM Workflows**: Enable LLMs to interact with existing APIs

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [🚀 User Guide](docs/USER_GUIDE.md) | Complete guide for using MCPWeaver |
| [⚙️ Installation](docs/INSTALLATION.md) | Platform-specific installation instructions |
| [🔧 API Reference](docs/API.md) | Complete API documentation |
| [👩‍💻 Developer Guide](docs/DEVELOPER.md) | Contributing and development setup |
| [🆘 Troubleshooting](docs/TROUBLESHOOTING.md) | Common issues and solutions |
| [❓ FAQ](docs/FAQ.md) | Frequently asked questions |
| [🔄 Migration Guide](docs/MIGRATION.md) | Migrating from openapi2mcp |

## 🏗️ Architecture

MCPWeaver is built using modern, reliable technologies:

- **Backend**: Go 1.23+ with robust error handling and performance monitoring
- **Frontend**: React 18 with TypeScript for type safety
- **Framework**: Wails v2 for native desktop integration
- **Database**: SQLite for local project storage
- **Build System**: Cross-platform automated builds with code signing

[📋 **View Detailed Architecture**](specs/ARCHITECTURE-SPECIFICATION.md)

## 🤝 Contributing

We welcome contributions! MCPWeaver is built by the community, for the community.

### Quick Contributing Guide

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'Add amazing feature'`
4. **Push** to the branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

[👩‍💻 **Full Developer Guide**](docs/DEVELOPER.md) | [🐛 **Report Issues**](https://github.com/matoval/MCPWeaver/issues)

## 📊 Performance

MCPWeaver is designed for speed and efficiency:

| Metric | Target | Typical |
|--------|--------|---------|
| Startup Time | < 2s | ~1.2s |
| Memory Usage | < 50MB | ~35MB |
| Small API (< 10 endpoints) | < 1s | ~0.5s |
| Medium API (10-100 endpoints) | < 3s | ~2.1s |
| Large API (100+ endpoints) | < 10s | ~7.2s |

## 🔒 Security

- **Local Processing**: All data processed locally, never sent to external servers
- **Code Signing**: All releases are code-signed for authenticity
- **Security Scanning**: Automated vulnerability scanning in CI/CD
- **Sandbox**: Template execution runs in a secure sandbox environment

[🔐 **Security Policy**](SECURITY.md)

## 📄 License

MCPWeaver is licensed under the [GNU Affero General Public License v3.0 (AGPL-3.0)](LICENSE).

This means:
- ✅ **Free to use** for personal and commercial projects
- ✅ **Open source** - you can view and modify the code
- ✅ **Share improvements** - contributions benefit everyone
- ⚠️ **Copyleft** - derivative works must also be open source

## 🌟 Related Projects

- **[openapi2mcp](https://github.com/modelcontextprotocol/openapi2mcp)**: The original CLI tool that inspired MCPWeaver
- **[Model Context Protocol](https://github.com/modelcontextprotocol/specification)**: The MCP specification
- **[Wails](https://wails.io/)**: The framework powering MCPWeaver's desktop interface

## 🙋 Support

Need help? We're here for you:

- 📖 **Documentation**: Start with our [User Guide](docs/USER_GUIDE.md)
- 🐛 **Bug Reports**: [Create an issue](https://github.com/matoval/MCPWeaver/issues/new?template=bug_report.md)
- 💡 **Feature Requests**: [Suggest a feature](https://github.com/matoval/MCPWeaver/issues/new?template=feature_request.md)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/matoval/MCPWeaver/discussions)
- ❓ **FAQ**: [Frequently Asked Questions](docs/FAQ.md)

## 🎉 Acknowledgments

- The **Model Context Protocol** team for creating the MCP specification
- The **openapi2mcp** contributors for the foundational work
- The **Wails** community for the excellent desktop framework
- All our **contributors** who make MCPWeaver better every day

---

<div align="center">

**⭐ Star this repository if MCPWeaver helps you build better integrations! ⭐**

[🚀 Get Started](docs/USER_GUIDE.md) • [📥 Download](https://github.com/matoval/MCPWeaver/releases/latest) • [🤝 Contribute](docs/DEVELOPER.md)

</div>
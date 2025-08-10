# MCPWeaver

[![Build Status](https://github.com/matoval/MCPWeaver/workflows/Build%20and%20Release/badge.svg)](https://github.com/matoval/MCPWeaver/actions)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)

**Transform OpenAPI specifications into Model Context Protocol (MCP) servers with a simple, fast CLI tool.**

MCPWeaver is an open-source CLI tool that converts OpenAPI specifications into fully functional MCP servers. Built with Go, it provides a lightweight, user-controlled solution that runs entirely on your machine with no external dependencies.

## 🚀 Quick Start

```bash
# Install MCPWeaver
curl -sf https://gobinaries.com/matoval/mcpweaver | sh

# Generate MCP server from OpenAPI spec
mcpweaver generate api.yaml --output ./my-server

# Run your generated server
cd my-server && python server.py
```

[📖 **View Full Getting Started Guide**](docs/USER_GUIDE.md) | [⬇️ **Download Latest Release**](https://github.com/matoval/MCPWeaver/releases/latest)

## ✨ Key Features

- **🚀 Simple CLI**: Single command to generate complete MCP servers
- **⚡ Lightning Fast**: Generate servers in seconds
- **🖥️ Pure Go**: No external dependencies, works completely offline
- **🎯 Interactive Selection**: Choose which endpoints to include
- **🐍 Python FastMCP**: Generates ready-to-run Python servers
- **🧪 Tests Included**: Generated servers come with complete test suites
- **🌐 Cross-Platform**: Single binary for Windows, macOS, and Linux
- **💾 Lightweight**: Minimal resource usage

## 🛠️ Installation

### Quick Install

**Using gobinaries (recommended):**

```bash
curl -sf https://gobinaries.com/matoval/mcpweaver | sh
```

**Direct download:**

```bash
# Linux/macOS
wget https://github.com/matoval/MCPWeaver/releases/latest/download/mcpweaver-linux-amd64
chmod +x mcpweaver-linux-amd64
sudo mv mcpweaver-linux-amd64 /usr/local/bin/mcpweaver

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/matoval/MCPWeaver/releases/latest/download/mcpweaver-windows-amd64.exe" -OutFile "mcpweaver.exe"
```

**Package managers:**

```bash
# Homebrew (macOS/Linux)
brew install matoval/tap/mcpweaver

# APT (Ubuntu/Debian)
echo "deb [trusted=yes] https://apt.fury.io/matoval/ /" | sudo tee /etc/apt/sources.list.d/matoval.list
sudo apt update && sudo apt install mcpweaver
```

[📋 **Detailed Installation Instructions**](docs/INSTALLATION.md)

## 🎯 Use Cases

- **API Integration**: Convert REST APIs to MCP for LLM integration
- **Development Tools**: Generate MCP servers for testing and development
- **Legacy System Integration**: Modernize old APIs with MCP
- **Microservices**: Create MCP interfaces for microservice architectures
- **AI/LLM Integration**: Convert REST APIs to MCP servers for LLM workflows
- **Rapid Prototyping**: Quickly generate MCP servers from existing API specs
- **API Modernization**: Bridge legacy OpenAPI specs with modern MCP protocol
- **Development Tools**: Create MCP servers for testing and development environments

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [📋 Requirements](Docs/REQUIREMENTS.md) | Project requirements and scope |
| [🏗️ Architecture](Docs/ARCHITECTURE.md) | Technical architecture and design |
| [⌨️ CLI Design](Docs/CLI-DESIGN.md) | Command interface and user experience |
| [📝 Examples](Docs/EXAMPLES.md) | Input/output examples and use cases |

## 🏗️ Architecture

MCPWeaver is built using modern, reliable technologies:

- **Core**: Pure Go 1.21+ with zero external dependencies
- **Parser**: OpenAPI 2.0/3.0+ support with comprehensive validation
- **Generator**: Template-based Python FastMCP server generation
- **CLI**: Cobra-based interface with interactive endpoint selection
- **Build**: Single binary cross-platform distribution

## 🤝 Contributing

We welcome contributions! MCPWeaver is built by the community, for the community.

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for your changes
5. Open a Pull Request

[👩‍💻 **Full Developer Guide**](docs/DEVELOPER.md) | [🐛 **Report Issues**](https://github.com/matoval/MCPWeaver/issues)

## 🔒 Security

- **Local Processing**: All data processed locally, no network connections required
- **No Dependencies**: Pure Go binary with no external runtime dependencies
- **Offline Operation**: Works completely offline, no data leaves your machine

## 📄 License

MCPWeaver is licensed under the [GNU Affero General Public License v3.0 (AGPL-3.0)](LICENSE).

This means:

- ✅ **Free to use** for personal and commercial projects
- ✅ **Open source** - you can view and modify the code
- ✅ **Share improvements** - contributions benefit everyone
- ⚠️ **Copyleft** - derivative works must also be open source

## 🌟 Related Projects

- **[openapi2mcp](https://github.com/modelcontextprotocol/openapi2mcp)**: Similar CLI tool for OpenAPI to MCP conversion
- **[Model Context Protocol](https://github.com/modelcontextprotocol/specification)**: The MCP specification
- **[Wails](https://wails.io/)**: The framework powering MCPWeaver's desktop interface

## 🙋 Support

Need help? We're here for you:

- 📖 **Documentation**: Check the planning docs in [/Docs](Docs/)
- 🐛 **Bug Reports**: [Create an issue](https://github.com/matoval/MCPWeaver/issues/new)
- 💡 **Feature Requests**: [Suggest a feature](https://github.com/matoval/MCPWeaver/issues/new)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/matoval/MCPWeaver/discussions)

## 🎉 Acknowledgments

- The **Model Context Protocol** team for creating the MCP specification
- The **OpenAPI** and **MCP** communities for their excellent specifications
- The **Go** community for the powerful standard library
- All our **contributors** who make MCPWeaver better every day

---

## **⭐ Star this repository if MCPWeaver helps you build better integrations! ⭐**

[📋 View Planning](Docs/) • [📥 Download](https://github.com/matoval/MCPWeaver/releases/latest) • [🤝 Contribute](https://github.com/matoval/MCPWeaver/blob/main/CONTRIBUTING.md)

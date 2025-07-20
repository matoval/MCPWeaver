# MCPWeaver Security Documentation

This document provides comprehensive security information for the MCPWeaver project, covering security architecture, implemented protections, testing procedures, and security best practices.

## Table of Contents

- [Security Overview](#security-overview)
- [Security Architecture](#security-architecture)
- [Implemented Security Controls](#implemented-security-controls)
- [Security Testing](#security-testing)
- [Vulnerability Management](#vulnerability-management)
- [Code Signing](#code-signing)
- [Security Configuration](#security-configuration)
- [Incident Response](#incident-response)
- [Security Best Practices](#security-best-practices)
- [Compliance](#compliance)
- [Security Resources](#security-resources)

## Security Overview

MCPWeaver is designed with security as a fundamental principle. As a desktop application that processes OpenAPI specifications and generates MCP servers, the application implements multiple layers of security controls to protect against various threat vectors.

### Security Principles

1. **Defense in Depth**: Multiple layers of security controls
2. **Principle of Least Privilege**: Minimal permissions and access rights
3. **Input Validation**: Comprehensive validation of all user inputs
4. **Secure by Default**: Safe default configurations
5. **Zero Trust**: No implicit trust in any component

### Threat Model

The primary threats MCPWeaver protects against include:

- **Malicious OpenAPI Specifications**: Crafted specs designed to exploit parsing vulnerabilities
- **Path Traversal Attacks**: Attempts to access unauthorized file system locations
- **Server-Side Request Forgery (SSRF)**: Attempts to make requests to internal systems
- **Template Injection**: Malicious template variables in generated code
- **Dependency Vulnerabilities**: Known vulnerabilities in third-party libraries
- **Code Injection**: Attempts to inject malicious code through various inputs

## Security Architecture

### Application Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    MCPWeaver Security Architecture          │
├─────────────────────────────────────────────────────────────┤
│ Frontend (React/TypeScript)                                 │
│ ├─ Input Validation Layer                                   │
│ ├─ Error Boundary & Recovery                                │
│ └─ CSP & XSS Protection                                     │
├─────────────────────────────────────────────────────────────┤
│ Backend (Go)                                                │
│ ├─ Security Validator                                       │
│ ├─ Secure File System                                       │
│ ├─ Secure HTTP Client                                       │
│ ├─ Template Security                                        │
│ └─ Error Handling & Logging                                 │
├─────────────────────────────────────────────────────────────┤
│ System Integration                                          │
│ ├─ Code Signing                                             │
│ ├─ File System Permissions                                  │
│ ├─ Network Restrictions                                     │
│ └─ Secure Storage                                           │
└─────────────────────────────────────────────────────────────┘
```

### Security Components

#### 1. Security Validator (`internal/app/security.go`)
- URL validation with SSRF protection
- File path validation with traversal protection
- Template injection prevention
- Input sanitization and validation
- JSON/YAML security validation
- Regular expression safety checks

#### 2. Secure File System (`internal/app/filesystem.go`)
- Atomic file operations
- Path restriction enforcement
- Secure temporary file creation
- File integrity checks
- Controlled file permissions

#### 3. Secure HTTP Client (`internal/app/httpclient.go`)
- TLS 1.2+ enforcement
- SSRF protection for URLs
- Content type validation
- Response size limits
- Secure headers management

#### 4. Template Security (in Generator Service)
- Template variable sanitization
- Path traversal prevention in templates
- Template injection protection
- Safe template execution environment

## Implemented Security Controls

### Input Validation and Sanitization

#### URL Validation
- **SSRF Protection**: Blocks requests to private IP ranges, localhost, and cloud metadata endpoints
- **Scheme Validation**: Only HTTP and HTTPS protocols allowed
- **Length Limits**: Maximum URL length enforcement
- **Format Validation**: Proper URL format validation

```go
// Example: URL validation with SSRF protection
validator := app.NewSecurityValidator(app)
err := validator.ValidateURL("https://api.example.com/openapi.yaml")
```

#### File Path Validation
- **Path Traversal Protection**: Prevents `../` attacks
- **Whitelist Validation**: Only allows access to designated directories
- **File Name Validation**: Restricts characters in file names
- **Length Limits**: Maximum path and filename length enforcement

#### Template Security
- **Variable Sanitization**: Escapes template injection characters
- **Control Character Removal**: Strips dangerous control characters
- **Length Limits**: Enforces maximum template variable length

### Network Security

#### HTTPS Enforcement
- **TLS 1.2+ Required**: Minimum TLS version enforcement
- **Strong Cipher Suites**: Only secure encryption algorithms
- **Certificate Validation**: Proper certificate chain validation
- **HSTS Support**: HTTP Strict Transport Security headers

#### Request Security
- **User-Agent Control**: Consistent and identifiable user agent
- **Header Sanitization**: Removes dangerous headers
- **Timeout Enforcement**: Prevents hanging connections
- **Size Limits**: Response size restrictions

### File System Security

#### Access Control
- **Restricted Paths**: Only allows access to approved directories
- **Atomic Operations**: Prevents race conditions
- **Secure Permissions**: Appropriate file and directory permissions
- **Temporary File Security**: Secure creation and cleanup

#### File Operations
- **Integrity Checks**: Validates file integrity before processing
- **Size Limits**: Prevents resource exhaustion
- **Type Validation**: Ensures files are of expected types
- **Secure Deletion**: Overwrites sensitive files before deletion

### Code Generation Security

#### Template Security
- **Sandboxed Execution**: Safe template processing environment
- **Input Sanitization**: All template variables sanitized
- **Path Validation**: Template paths validated against traversal
- **Safe Functions**: Only safe template functions available

#### Output Validation
- **Generated Code Review**: Basic validation of generated code
- **Syntax Checking**: Ensures generated code is syntactically correct
- **Security Scanning**: Generated code can be scanned for issues

## Security Testing

### Automated Security Testing

The project includes comprehensive automated security testing:

#### Static Application Security Testing (SAST)
- **gosec**: Go security checker for common vulnerabilities
- **semgrep**: Multi-language static analysis for security patterns
- **golangci-lint**: Code quality with security rules

#### Dependency Scanning
- **govulncheck**: Go vulnerability database checking
- **nancy**: Dependency vulnerability scanner
- **trivy**: Container and filesystem vulnerability scanning

#### Dynamic Testing
- **Input Fuzzing**: Automated testing with malformed inputs
- **Boundary Testing**: Edge case validation
- **Integration Testing**: End-to-end security testing

### Running Security Tests

```bash
# Run all security tests
./scripts/security-test.sh

# Run quick security scan
./scripts/security-test.sh --quick

# Run only SAST tools
./scripts/security-test.sh --sast-only

# Run comprehensive vulnerability scan
./scripts/vulnerability-scan.sh

# Run specific security test suite
go test -v ./tests/security/...
```

### Continuous Security Monitoring

#### GitHub Actions Integration
- **Automated Scanning**: Security scans on every push and PR
- **SARIF Upload**: Results uploaded to GitHub Security tab
- **Workflow Artifacts**: Detailed reports available for download
- **Failure Notifications**: Automatic alerts for security issues

#### Scheduled Scans
- **Daily Vulnerability Scans**: Automated daily security checks
- **Dependency Updates**: Regular checks for vulnerable dependencies
- **Security Alerts**: Immediate notifications for critical issues

## Vulnerability Management

### Vulnerability Response Process

1. **Detection**: Automated scans detect vulnerabilities
2. **Assessment**: Security team evaluates severity and impact
3. **Prioritization**: Issues prioritized based on risk level
4. **Remediation**: Patches and fixes developed and tested
5. **Deployment**: Fixes deployed and verified
6. **Communication**: Stakeholders notified of resolution

### Severity Classification

#### High Severity
- Remote code execution vulnerabilities
- Authentication bypass issues
- Data exposure vulnerabilities
- **SLA**: Fix within 24 hours

#### Medium Severity
- Privilege escalation issues
- Information disclosure
- Denial of service vulnerabilities
- **SLA**: Fix within 1 week

#### Low Severity
- Configuration issues
- Information leaks
- Non-exploitable vulnerabilities
- **SLA**: Fix within 1 month

### Vulnerability Scanning Tools

The project uses multiple vulnerability scanning tools:

1. **gosec**: Static analysis for Go security issues
2. **govulncheck**: Go vulnerability database checking
3. **nancy**: Dependency vulnerability scanning
4. **trivy**: Multi-purpose vulnerability scanner
5. **semgrep**: Pattern-based security analysis

## Code Signing

MCPWeaver implements code signing for both Windows and macOS platforms to ensure binary integrity and authenticity.

### Supported Platforms

#### macOS Code Signing
- **Developer ID Application Certificate**: Required for distribution
- **Notarization**: Apple notarization for enhanced security
- **Hardened Runtime**: Enhanced security features enabled
- **Entitlements**: Minimal required permissions

#### Windows Code Signing
- **Authenticode Signing**: Standard Windows code signing
- **Timestamp Servers**: Ensures signature validity beyond certificate expiry
- **Certificate Validation**: Proper certificate chain validation

### Code Signing Process

```bash
# Sign macOS binary
./scripts/sign-code.sh --platform macos --binary ./build/MCPWeaver.app

# Sign Windows binary
./scripts/sign-code.sh --platform windows --binary ./build/mcpweaver.exe

# Verify signatures
./scripts/sign-code.sh --validate ./dist/signed-binary
```

See [CODE_SIGNING.md](CODE_SIGNING.md) for detailed instructions.

## Security Configuration

### Environment Variables

#### Security Settings
```bash
# Enable security features
ENABLE_SECURITY_VALIDATION=true
ENABLE_SSRF_PROTECTION=true
ENABLE_PATH_VALIDATION=true

# File system security
MAX_FILE_SIZE=10485760  # 10MB
ALLOWED_PATHS="./,projects/,output/,temp/"

# Network security
ENABLE_TLS_VALIDATION=true
MIN_TLS_VERSION=1.2
HTTPS_ONLY=true
```

#### Development vs Production
```bash
# Development (relaxed security for testing)
ENVIRONMENT=development
SECURITY_MODE=development
ALLOW_INSECURE_REQUESTS=true

# Production (strict security)
ENVIRONMENT=production
SECURITY_MODE=production
ALLOW_INSECURE_REQUESTS=false
```

### Configuration Files

#### Security Policy (`security-policy.json`)
```json
{
  "version": "1.0",
  "policies": {
    "input_validation": {
      "max_url_length": 2048,
      "max_file_size": 10485760,
      "allowed_schemes": ["http", "https"]
    },
    "network_security": {
      "block_private_ips": true,
      "block_localhost": true,
      "require_tls": true
    },
    "file_system": {
      "allowed_paths": ["./", "projects/", "output/"],
      "max_path_length": 255,
      "secure_temp_files": true
    }
  }
}
```

## Incident Response

### Security Incident Classification

#### Critical Incidents
- Active exploitation of vulnerabilities
- Data breach or exposure
- System compromise
- **Response Time**: Immediate (< 1 hour)

#### High Priority Incidents
- Newly discovered high-severity vulnerabilities
- Security control failures
- Suspicious activity detected
- **Response Time**: < 4 hours

#### Medium Priority Incidents
- Security misconfigurations
- Failed security scans
- Policy violations
- **Response Time**: < 24 hours

### Incident Response Process

1. **Detection and Reporting**
   - Automated monitoring alerts
   - User reports
   - Security scan findings

2. **Initial Response**
   - Incident classification
   - Initial assessment
   - Stakeholder notification

3. **Investigation**
   - Evidence collection
   - Impact assessment
   - Root cause analysis

4. **Containment**
   - Isolate affected systems
   - Prevent further damage
   - Implement temporary fixes

5. **Eradication**
   - Remove threats
   - Patch vulnerabilities
   - Update security controls

6. **Recovery**
   - Restore normal operations
   - Verify system integrity
   - Monitor for recurrence

7. **Post-Incident**
   - Document lessons learned
   - Update procedures
   - Improve security controls

### Contact Information

#### Security Team
- **Security Email**: security@mcpweaver.com
- **Emergency Contact**: +1-XXX-XXX-XXXX
- **PGP Key**: Available on project website

#### Reporting Security Issues
- **GitHub Security Advisory**: Use GitHub's private vulnerability reporting
- **Email**: security@mcpweaver.com
- **Bug Bounty**: Details available on project website

## Security Best Practices

### For Developers

#### Secure Coding Practices
1. **Input Validation**: Validate all inputs at boundaries
2. **Output Encoding**: Properly encode outputs
3. **Error Handling**: Don't expose sensitive information in errors
4. **Logging**: Log security events appropriately
5. **Dependencies**: Keep dependencies updated

#### Code Review Guidelines
1. **Security Focus**: Review for security implications
2. **Input Validation**: Verify proper input handling
3. **Authentication**: Check authentication and authorization
4. **Cryptography**: Verify proper use of crypto functions
5. **Configuration**: Review security configurations

### For Users

#### Safe Usage Guidelines
1. **Trusted Sources**: Only import OpenAPI specs from trusted sources
2. **File Validation**: Verify file integrity before processing
3. **Network Security**: Use secure networks for spec downloads
4. **Updates**: Keep MCPWeaver updated to latest version
5. **Permissions**: Run with minimal required permissions

#### Configuration Recommendations
1. **Restrict Paths**: Limit file system access paths
2. **Network Controls**: Use firewalls and network restrictions
3. **Monitoring**: Enable security logging and monitoring
4. **Backup**: Regularly backup important project data
5. **Training**: Stay informed about security best practices

### For Administrators

#### Deployment Security
1. **Environment Isolation**: Separate development and production
2. **Access Control**: Implement proper access controls
3. **Monitoring**: Deploy comprehensive monitoring
4. **Incident Response**: Have incident response procedures ready
5. **Compliance**: Ensure compliance with relevant standards

## Compliance

### Standards and Frameworks

#### Security Standards
- **OWASP Top 10**: Protection against common web vulnerabilities
- **NIST Cybersecurity Framework**: Comprehensive security approach
- **ISO 27001**: Information security management standards
- **CWE/SANS Top 25**: Common weakness enumeration

#### Development Standards
- **SSDLC**: Secure Software Development Lifecycle
- **DevSecOps**: Security integrated into development process
- **OWASP SAMM**: Software Assurance Maturity Model

### Compliance Checklist

#### Input Validation (OWASP A03:2021)
- ✅ Server-side input validation
- ✅ Parameterized queries (where applicable)
- ✅ Input sanitization
- ✅ File upload restrictions
- ✅ URL validation

#### Security Configuration (OWASP A05:2021)
- ✅ Secure default configurations
- ✅ Security headers implemented
- ✅ Error handling configured
- ✅ Logging and monitoring enabled
- ✅ Regular security updates

#### Vulnerable Components (OWASP A06:2021)
- ✅ Dependency vulnerability scanning
- ✅ Regular dependency updates
- ✅ Component inventory maintained
- ✅ Security patches applied
- ✅ EOL component identification

#### Security Logging (OWASP A09:2021)
- ✅ Comprehensive audit logging
- ✅ Log integrity protection
- ✅ Security event monitoring
- ✅ Incident response procedures
- ✅ Log analysis capabilities

## Security Resources

### Documentation
- [Security Architecture](SECURITY_ARCHITECTURE.md)
- [Code Signing Guide](CODE_SIGNING.md)
- [Vulnerability Management](VULNERABILITY_MANAGEMENT.md)
- [Incident Response Plan](INCIDENT_RESPONSE.md)

### Tools and Scripts
- `scripts/security-test.sh` - Comprehensive security testing
- `scripts/vulnerability-scan.sh` - Vulnerability scanning
- `scripts/sign-code.sh` - Code signing automation
- `tests/security/` - Security test suites

### External Resources
- [OWASP Security Knowledge](https://owasp.org/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [Go Security Checklist](https://github.com/Checkmarx/Go-SCP)
- [Secure Code Review Guide](https://owasp.org/www-project-code-review-guide/)

### Training and Certification
- [OWASP WebGoat](https://owasp.org/www-project-webgoat/) - Security training
- [Secure Code Warrior](https://www.securecodewarrior.com/) - Developer training
- [SANS Security Training](https://www.sans.org/) - Professional security training

## Contributing to Security

### Reporting Security Issues
We take security seriously. If you discover a security issue, please:

1. **Do not** open a public GitHub issue
2. Email security@mcpweaver.com with details
3. Include steps to reproduce the issue
4. Provide proof of concept if possible
5. Allow time for investigation and fix

### Security Contributions
We welcome security contributions including:
- Security feature implementations
- Vulnerability fixes
- Security test improvements
- Documentation enhancements
- Tool integrations

### Recognition
Contributors who help improve MCPWeaver's security will be:
- Credited in release notes
- Listed in security acknowledgments
- Invited to join our security advisory team
- Eligible for bug bounty rewards (when program launches)

## Changelog

### Version 1.0.0 (Current)
- Initial security framework implementation
- Comprehensive input validation
- SSRF protection mechanisms
- Secure file system operations
- Code signing support
- Automated security testing
- Vulnerability scanning integration

### Planned Enhancements
- Enhanced threat modeling
- Additional SAST tool integration
- Runtime security monitoring
- Security metrics dashboard
- Advanced anomaly detection
- Machine learning-based threat detection

---

**Last Updated**: $(date)
**Security Version**: 1.0.0
**Review Schedule**: Quarterly

For questions about this security documentation, contact the security team at security@mcpweaver.com.
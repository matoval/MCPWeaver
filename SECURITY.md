# Security Policy

## Supported Versions

We actively support the following versions of MCPWeaver with security updates:

| Version | Supported          | End of Life |
| ------- | ------------------ | ----------- |
| 1.0.x   | :white_check_mark: | TBD         |
| < 1.0   | :x:                | 2024-12-31  |

## Reporting a Vulnerability

We take the security of MCPWeaver seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please use one of the following methods:

#### Option 1: GitHub Security Advisory (Preferred)
1. Go to the [Security tab](https://github.com/matoval/MCPWeaver/security) of this repository
2. Click "Report a vulnerability"
3. Fill out the security advisory form with detailed information
4. Submit the report

#### Option 2: Email
Send an email to: **security@mcpweaver.com**

Include the following information:
- Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit the issue

### What to Expect

After you submit a vulnerability report, we will:

1. **Acknowledge receipt** within 48 hours
2. **Confirm the vulnerability** and determine its severity within 5 business days
3. **Develop and test a fix** with an estimated timeline
4. **Release the fix** and coordinate disclosure
5. **Credit you** in our security acknowledgments (if desired)

### Severity Guidelines

We use the following severity classifications:

#### Critical
- Remote code execution
- Authentication bypass
- Data exposure affecting all users
- **Response time**: 24 hours
- **Fix timeline**: 1-3 days

#### High
- Privilege escalation
- Significant data exposure
- Authentication/authorization issues
- **Response time**: 48 hours
- **Fix timeline**: 1-2 weeks

#### Medium
- Information disclosure
- Denial of service
- CSRF vulnerabilities
- **Response time**: 5 business days
- **Fix timeline**: 2-4 weeks

#### Low
- Security misconfigurations
- Minor information leaks
- **Response time**: 10 business days
- **Fix timeline**: Next release cycle

## Security Features

MCPWeaver includes the following security features:

### Input Validation
- Comprehensive validation of all user inputs
- Protection against path traversal attacks
- SSRF (Server-Side Request Forgery) protection
- Template injection prevention

### Network Security
- HTTPS enforcement for external requests
- TLS 1.2+ requirement
- Private IP address blocking
- Secure HTTP client configuration

### File System Security
- Restricted file system access
- Atomic file operations
- Secure temporary file handling
- Path traversal protection

### Code Generation Security
- Template variable sanitization
- Safe template execution environment
- Generated code validation

## Security Testing

We maintain comprehensive security testing including:

- **Static Application Security Testing (SAST)**: gosec, semgrep
- **Dependency Vulnerability Scanning**: govulncheck, nancy, trivy
- **Container Security**: trivy container scanning
- **Automated Security Tests**: Custom security test suite
- **Continuous Monitoring**: GitHub Actions integration

## Security Updates

Security updates are released as follows:

1. **Critical vulnerabilities**: Immediate patch release
2. **High severity**: Within 1-2 weeks
3. **Medium/Low severity**: Next regular release

All security updates are:
- Clearly marked in release notes
- Include CVE references when applicable
- Backward compatible when possible
- Include migration guides for breaking changes

## Vulnerability Disclosure

We follow responsible disclosure practices:

1. **Private disclosure** to our security team
2. **Coordinated public disclosure** after fix is available
3. **CVE assignment** for significant vulnerabilities
4. **Public acknowledgment** of reporters (with permission)

### Timeline
- **Day 0**: Vulnerability reported
- **Day 1-5**: Validation and assessment
- **Day 5-30**: Fix development and testing
- **Day 30**: Coordinated disclosure
- **Day 90**: Full public disclosure (if vendor doesn't respond)

## Security Best Practices for Users

### General Guidelines
1. **Keep MCPWeaver updated** to the latest version
2. **Validate OpenAPI specifications** from untrusted sources
3. **Use secure networks** when downloading specifications
4. **Review generated code** before deployment
5. **Follow the principle of least privilege**

### Configuration Security
1. **Restrict file system access** to necessary directories only
2. **Enable all security features** in production environments
3. **Use HTTPS** for all external specification downloads
4. **Monitor logs** for suspicious activity
5. **Regularly backup** project data

### Network Security
1. **Use firewalls** to restrict network access
2. **Disable unnecessary network services**
3. **Monitor network traffic** for anomalies
4. **Use VPN** for remote access
5. **Implement network segmentation**

## Security Resources

### Documentation
- [Comprehensive Security Documentation](docs/SECURITY.md)
- [Code Signing Guide](docs/CODE_SIGNING.md)
- [Security Testing Documentation](docs/SECURITY_TESTING.md)

### Tools and Scripts
- Security testing: `scripts/security-test.sh`
- Vulnerability scanning: `scripts/vulnerability-scan.sh`
- Code signing: `scripts/sign-code.sh`

### External Resources
- [OWASP Application Security](https://owasp.org/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [Go Security Guidelines](https://github.com/Checkmarx/Go-SCP)

## Contact Information

- **Security Team**: security@mcpweaver.com
- **General Inquiries**: support@mcpweaver.com
- **Project Maintainers**: See [MAINTAINERS.md](MAINTAINERS.md)

## Bug Bounty Program

We are currently developing a bug bounty program to reward security researchers who help improve MCPWeaver's security. Details will be announced soon.

### Recognition

Security researchers who responsibly disclose vulnerabilities will be:
- Credited in our security acknowledgments
- Mentioned in release notes (with permission)
- Invited to join our security advisory team
- Eligible for future bug bounty rewards

## Legal

By participating in our vulnerability disclosure program, you agree to:
- Make a good faith effort to avoid privacy violations and destruction of data
- Only interact with accounts you own or have explicit permission to access
- Not access or modify other users' data
- Not perform attacks that could harm the reliability/integrity of our services
- Not use social engineering, physical attacks, or denial of service attacks
- Provide us with reasonable time to resolve vulnerabilities before disclosure

We commit to:
- Respond to your report promptly and keep you informed throughout the process
- Work with you to understand and resolve the issue quickly
- Recognize your contribution to improving our security
- Not pursue legal action against researchers who comply with this policy

---

**Last Updated**: December 2024
**Version**: 1.0

This security policy is reviewed and updated regularly. For the most current version, please check this repository.
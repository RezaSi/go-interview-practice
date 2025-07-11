# Security Policy

## Supported Versions

We actively maintain and provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in this project, please follow these steps:

### 1. **DO NOT** create a public GitHub issue
Security vulnerabilities should be reported privately to prevent potential exploitation.

### 2. **Email us directly**
Send a detailed report to: **rezashiri88@gmail.com**

### 3. **Include the following information**
Your security report should contain:

- **Description**: A clear description of the vulnerability
- **Steps to reproduce**: Detailed steps to reproduce the issue
- **Impact**: Potential impact of the vulnerability
- **Suggested fix**: If you have a solution in mind
- **Affected versions**: Which versions are affected
- **Proof of concept**: If applicable, include a minimal PoC

### 4. **Response timeline**
- **Initial response**: Within 48 hours
- **Status update**: Within 1 week
- **Resolution**: As quickly as possible, typically within 30 days

## Security Measures

### Code Execution Safety
- All user-submitted code runs in isolated environments
- Execution timeouts prevent resource exhaustion
- Memory limits prevent excessive resource usage
- File system access is restricted to challenge directories only

### Input Validation
- All user inputs are validated and sanitized
- SQL injection protection through parameterized queries
- XSS protection through proper output encoding
- File upload restrictions and validation

### Authentication & Authorization
- GitHub OAuth integration for secure authentication
- User submissions are restricted to their own directories
- Pull request validation ensures users can only modify their own code
- Admin-only access to sensitive operations

### Infrastructure Security
- HTTPS enforcement for all web traffic
- Regular dependency updates to patch known vulnerabilities
- Automated security scanning in CI/CD pipeline
- Environment variables for sensitive configuration

## Responsible Disclosure

We follow responsible disclosure practices:

1. **Private reporting**: Vulnerabilities are reported privately first
2. **Timely response**: We respond quickly to security reports
3. **Coordinated disclosure**: We work with reporters to coordinate public disclosure
4. **Credit acknowledgment**: Security researchers are credited in our security advisories

## Security Updates

### Automatic Updates
- Dependencies are automatically updated via Dependabot
- Security patches are applied as soon as they're available
- Automated vulnerability scanning in our CI/CD pipeline

### Manual Updates
- Critical security issues are addressed immediately
- Non-critical issues are scheduled for the next release
- Security advisories are published for significant vulnerabilities

## Security Best Practices for Contributors

### Code Review
- All code changes require security review
- Automated security scanning in pull requests
- Manual review for security-sensitive changes

### Dependency Management
- Regular updates of all dependencies
- Monitoring for known vulnerabilities
- Immediate updates for security patches

### Testing
- Security-focused testing in CI/CD
- Penetration testing for major releases
- Regular security audits

## Contact Information

### Security Team
- **Email**: rezashiri88@gmail.com
- **Response Time**: Within 72 hours
- **PGP Key**: Available upon request

### Project Maintainer
- **GitHub**: [@RezaSi](https://github.com/RezaSi)
- **Email**: rezashiri88@gmail.com

## Security Hall of Fame

We acknowledge security researchers who help improve our security:

| Researcher | Vulnerability | Date |
|------------|---------------|------|
| *Your name could be here* | *Report a vulnerability* | *Help us improve* |

## Bug Bounty

While we don't currently offer a formal bug bounty program, we do:

- Acknowledge security researchers in our documentation
- Provide early access to security patches
- Consider special recognition for significant findings

## Legal

By reporting a security vulnerability, you agree to:

- Keep the vulnerability confidential until we've had time to address it
- Not exploit the vulnerability for malicious purposes
- Work with us to coordinate disclosure
- Follow responsible disclosure practices

---

**Thank you for helping keep our community secure!** 🛡️

*Last updated: July 2025* 

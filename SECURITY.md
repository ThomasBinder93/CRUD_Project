# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| 0.x.x   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in this project, please do **NOT** create a public GitHub issue.

Instead, please report the vulnerability by:

1. **Email**: Send details to the project maintainer
2. **Include**:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Any suggested fixes

Please allow 48 hours for a response. Once we verify the vulnerability, we will:

1. Create a fix on a non-public branch
2. Test the fix thoroughly
3. Create a security release
4. Credit you in the release notes (unless you prefer not to be named)

## Security Best Practices

When using this project:

- Keep Go and dependencies updated
- Don't commit secrets or API keys
- Use environment variables for sensitive configuration
- Enable HTTPS in production
- Run with minimal necessary permissions
- Regularly update dependencies

## Dependencies

This project uses the following key dependencies:

- Go Standard Library
- gorilla/mux (HTTP routing)
- modernc.org/sqlite (SQLite driver)

These are regularly scanned for vulnerabilities using GitHub's Dependabot.

## Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://pkg.go.dev/crypto)
- [Container Security](https://docs.docker.com/engine/security/)

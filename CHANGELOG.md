# Changelog

## [Unreleased]

## [0.2.0] - 2026-11-07

### Breaking

- See repository diff for API changes in this release (secure defaults, renamed packages where applicable, context-aware constructors).

### Security

- Hardened defaults: TLS, timeouts, body limits, PII redaction where logging applies.

### Fixed

- Bug fixes from the v0.2.0 hardening pass.

### Changed

- CI uses go-version-file; added race, staticcheck, govulncheck, golangci-lint (gosec).



- All HTTP subpackages check the status code before JSON parsing (`common.HTTPStatusError`) and cap response bodies (`common.ReadBody`, 4 MiB JSON / 10 MiB images).

## [0.1.0]

- Initial release: mini program login, user, qrcode, subscribe, LBS, and WeChat Pay API v3 subpackages.
- Caller-injected credentials; httptest-friendly HTTP hosts per subpackage.

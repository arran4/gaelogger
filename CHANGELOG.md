# Changelog

## v0.1.1 — CI compatibility maintenance

- Update the Gitleaks GitHub Action from v2 to the Node 24-compatible v3, incorporating Dependabot PR #19.
- Keep the existing Go module dependency set and public logger API unchanged.

## v0.1.0 — Initial release

- Publish the existing Go logger module as its first tagged release.
- Refresh the complete Go dependency graph, including Google Cloud Logging and OpenTelemetry dependencies.
- Verify local logger severity formatting with automated tests and Go race detection.
- Test all pull requests and publish versioned releases only after CI succeeds.

This module retains its existing API. For new integrations, consider the Go standard library's `log/slog` API, as described in the README.

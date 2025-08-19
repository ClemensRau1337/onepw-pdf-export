# CONTRIBUTING

Danke, dass du mithelfen möchtest! / Thanks for contributing!

## Melde ein Problem / Report a Bug
- **Poste niemals echte Passwörter oder sensible Daten.**

## Feature Requests
- Beschreibe das Ziel (Use Case), nicht nur die Lösung.
- Skizziere Benutzerführung (CLI-Flags, Interaktion).

## Code-Richtlinien / Coding Style
- Go >= 1.20, `go fmt`, `go vet`.
- Kleine, fokussierte PRs mit Tests (wo sinnvoll).
- Keine Secrets in Code/Tests/Logs.

## Commit-Messages
- Formatiere prägnant: `feat: ...`, `fix: ...`, `docs: ...`, `chore: ...`.
- Verweise auf Issues: `Fixes #123`.

## Development Setup
```bash
go mod tidy
go build -o onepw-pdf-export ./
```

## Sicherheit
- Melde Sicherheitslücken **vertraulich** (siehe SECURITY.md). Keine Proof-of-Concepts mit echten Daten posten.

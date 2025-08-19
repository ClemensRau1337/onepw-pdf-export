# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2025-08-19
### Added

### Changed
- Removed wrong import package `path/filepath`
- Removed placeholder repo url in english readme
- Removed typo in readme

## [1.0.0] - 2025-08-19
### Added
- Full **interactive mode** for all parameters (output path, risk acceptance, vault picker, password, layout, masking, search).
- **UTF‑8** PDF rendering with automatic DejaVuSans font download.
- Progress spinner and periodic progress updates.
- Bilingual README (DE/EN).
- Interactive vault selection.
- Progress indicator during export.
- Add Changelog file

### Changed
- Removed `live` subcommand. Default mode uses 1Password CLI (`op`).

### Security
- Mandatory PDF encryption with user password (gofpdf SetProtection).

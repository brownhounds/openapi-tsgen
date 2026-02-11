[![Go Reference](https://pkg.go.dev/badge/github.com/brownhounds/openapi-tsgen.svg)](https://pkg.go.dev/github.com/brownhounds/openapi-tsgen)
[![CI](https://github.com/brownhounds/openapi-tsgen/actions/workflows/ci.yml/badge.svg)](https://github.com/brownhounds/openapi-tsgen/actions/workflows/ci.yml)
[![Release](https://github.com/brownhounds/openapi-tsgen/actions/workflows/release.yml/badge.svg)](https://github.com/brownhounds/openapi-tsgen/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/brownhounds/openapi-tsgen)](https://goreportcard.com/report/github.com/brownhounds/openapi-tsgen)
[![Latest Release](https://img.shields.io/github/v/release/brownhounds/openapi-tsgen)](https://github.com/brownhounds/openapi-tsgen/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/brownhounds/openapi-tsgen)](https://github.com/brownhounds/openapi-tsgen/blob/main/go.mod)
[![License](https://img.shields.io/github/license/brownhounds/openapi-tsgen)](https://github.com/brownhounds/openapi-tsgen/blob/main/LICENCE)

## OpenAPI TSGEN

Generate TypeScript types from OpenAPI schemas (YAML or JSON).

## Usage

Basic usage:

```bash
openapi-tsgen -s schema.yml -o type.ts
```

JSON input:

```bash
openapi-tsgen -s schema.json -o type.ts --input-json
```

## Install

### Build From Source

Dependencies: go >= 1.25.6

1. Clone repository
2. Run: `go get`
3. Build:

```bash
go generate ./...
go build -ldflags="-s -w" -o ./bin/openapi-tsgen main.go
```

### Install With Go

```bash
go install github.com/brownhounds/openapi-tsgen@v0.1.1
```

### Install Script LINUX/amd64

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/brownhounds/openapi-tsgen/v0.1.1/tools/install-linux-amd64.sh)"
```

### Install Script LINUX/arm64

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/brownhounds/openapi-tsgen/v0.1.1/tools/install-linux-arm64.sh)"
```

### Post Install/LINUX

Add following to `.bashrc` or `.zshrc`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Source `profile` file or restart your terminal session.

## Other Platforms

Refer to latest release page: [Release](https://github.com/brownhounds/openapi-tsgen/releases)

# WatchDoc

A lightweight development file server with live reload. WatchDoc watches your files for changes, optionally runs a build command, and automatically reloads connected browsers via WebSocket.

## Features

- **Live reload** - Automatically reloads the browser when served files change
- **Build command** - Runs a specified command when source files change (e.g., a doc generator)
- **Multi-directory watch** - Watch additional source directories alongside the served directory
- **Recursive watching** - Monitors all subdirectories (skips `.git`, `node_modules`, `vendor`)
- **Debounced updates** - Coalesces rapid file changes into a single reload
- **Zero config** - Works out of the box with sensible defaults

## Installation

```bash
go install github.com/absmach/watchdoc@latest
```

Or build from source:

```bash
git clone https://github.com/absmach/watchdoc.git
cd watchdoc
make build
```

## Usage

```bash
# Serve current directory with live reload
watchdoc

# Serve a specific directory
watchdoc -dir ./build/output

# Run a build command when source files change
watchdoc -dir ./output -watch-dirs ./src -cmd "make build-docs"

# Custom port
watchdoc -port 3000

# Full example: watch source, run build, serve output
watchdoc -port 8080 -dir ./site -watch-dirs ./docs,./templates -cmd "make generate"
```

### Flags

| Flag          | Default | Description                                     |
| ------------- | ------- | ----------------------------------------------- |
| `-port`       | `8080`  | Port to run the file server on                  |
| `-dir`        | `.`     | Directory to serve files from                   |
| `-watch-dirs` |         | Additional comma-separated directories to watch |
| `-cmd`        |         | Command to execute when source files change     |

## How It Works

1. WatchDoc starts an HTTP file server for the specified directory
2. HTML responses are automatically injected with a WebSocket client script
3. A file watcher monitors the served directory and any additional watch directories
4. When files in a **watch directory** change and a `-cmd` is set, the command runs (e.g., rebuild docs)
5. When files in the **served directory** change, all connected browsers receive a reload signal

This two-stage approach (source → build → output → reload) works well for static site generators, documentation tools, and similar workflows.

## Development

```bash
# Build
make build

# Run tests
make test

# Lint
make lint

# Clean build artifacts
make clean
```

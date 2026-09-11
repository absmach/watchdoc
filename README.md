# WatchDoc 📄👀

[![CI](https://github.com/absmach/watchdoc/actions/workflows/build.yml/badge.svg)](https://github.com/absmach/watchdoc/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/absmach/watchdoc)](https://goreportcard.com/report/github.com/absmach/watchdoc)
[![Go Reference](https://pkg.go.dev/badge/github.com/absmach/watchdoc.svg)](https://pkg.go.dev/github.com/absmach/watchdoc)

![Alt text](img.png)

**WatchDoc** is a lightweight file server with live reload,
built for fast documentation and static-site workflows.

It watches your files for changes, optionally runs a build command, and
automatically reloads all connected browsers via WebSocket, so you can
focus on writing, not refreshing.

This started as an internal tool used for building our website at https://absmach.eu/,
but it can be used as a general-purpose file server in similar scenarios.
If you have custom build scripts but also want to watch changes, trigger
rebuilds between source and compiled files, and observe updates in the browser,
this is the tool for you.

All you need to do is to configure three parameters:
- source dir to watch for changes
- destination dir to serve and watch for browser refresh
- command used to trigger rebuild

The only difference between the source and destination dirs is that destination
dir changes do not trigger the command.
Or, you can just use it as a fast and simple zero-dependency file server (no npm, no heavy dependencies).
You can install it via:
```bash
curl -fsSL https://raw.githubusercontent.com/absmach/watchdoc/main/install.sh | sh
```
or with a custom location:
```bash
INSTALL_DIR=$HOME/.local/bin \
curl -fsSL https://raw.githubusercontent.com/absmach/watchdoc/main/install.sh | sh
```

If you are using Go, you can also install it with:
```bash
go install github.com/absmach/watchdoc@main
```

------------------------------------------------------------------------

## ✨ Features

-   **Live reload**\
    Automatically refreshes the browser when served files change.

-   **Build on change**\
    Runs a custom command when source files update (perfect for doc
    generators and static site builds).

-   **Multiple watch directories**\
    Watch additional source folders alongside the served output
    directory.

-   **Recursive file watching**\
    Monitors all subdirectories while skipping common noise (`.git`,
    `node_modules`, `vendor`).

-   **Debounced updates**\
    Groups rapid file changes into a single rebuild and reload.

-   **Zero configuration**\
    Sensible defaults. Just run it and go.

------------------------------------------------------------------------

## 📦 Installation

Install the latest version with Go:

``` bash
go install github.com/absmach/watchdoc@latest
```

Or build from source:

``` bash
git clone https://github.com/absmach/watchdoc.git
cd watchdoc
make build
```

------------------------------------------------------------------------

## 🚀 Usage

``` bash
# Serve the current directory with live reload
watchdoc

# Serve a specific directory
watchdoc -serve-dir ./build/output

# Run a build command when source files change
watchdoc -serve-dir ./output -watch-dirs ./src -cmd "make build-docs"

# Use a custom port
watchdoc -port 3000

# Disable automatic browser opening (boolean flags take no spaced value)
watchdoc -no-browser

# Full example:
# watch sources → run build → serve output → live reload
watchdoc \
  -port 8080 \
  -serve-dir ./site \
  -watch-dirs ./docs,./templates \
  -cmd "make generate"
```

We use WatchDoc as a local file server and also to watch our blogs directory and 
trigger rebuild on changes [here](https://github.com/absmach/website/blob/main/Makefile#L23)
so we can focus on writing blogs and editing the website, rather than refreshing and rebuilding.

---

## ⚙️ Flags

 
| Flag          | Default | Description                                                          |
| ------------- | ------- | -------------------------------------------------------------------- |
| `-port`       | `8080`  | Port to run the file server on                                       |
| `-serve-dir`  | `.`     | Directory to serve files from, this is usually your built HTML files |
| `-watch-dirs` |         | Comma-separated directories for files to watch for changes           |
| `-cmd`        |         | Command to execute when source files change                          |
| `-no-browser` | `false` | Disable automatic browser opening                                    |

Boolean flags use `-no-browser` or `-no-browser=true`, not
`-no-browser true`.

On Windows, build commands are run with `cmd.exe /C`, so batch files can be
used directly:

``` powershell
watchdoc.exe -no-browser -watch-dirs ./content -serve-dir ./out -cmd "./build.bat"
```

## ❓ How It Works

1.  WatchDoc starts an HTTP file server for the selected output
    directory.
2.  HTML responses are automatically injected with a WebSocket client.
3.  A file watcher monitors:
    -   the served directory
    -   any additional watch directories
4.  When files in a **watch directory** change and `-cmd` is set, the
    command is executed (e.g. regenerate docs).
5.  When files in the **served directory** change, all connected
    browsers receive a reload signal.

This flow:

    source → build → output → reload

is ideal for static site generators, documentation pipelines, and
similar development setups.

------------------------------------------------------------------------

## 🛠 Development

``` bash
# Build
make build

# Run tests
make test

# Lint code
make lint

# Clean build artifacts
make clean
```

# WatchDoc 📄👀
![Alt text](img.png)

**WatchDoc** is a lightweight development file server with live reload,
built for fast documentation and static-site workflows.

It watches your files for changes, optionally runs a build command, and
automatically reloads all connected browsers via WebSocket---so you can
focus on writing, not refreshing.

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
    Sensible defaults---just run it and go.

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

# Full example:
# watch sources → run build → serve output → live reload
watchdoc \
  -port 8080 \
  -serve-dir ./site \
  -watch-dirs ./docs,./templates \
  -cmd "make generate"
```

We use WatchDoc as a local file server and also to watch another source director and 
trigger rebuild on changes here https://github.com/absmach/website/blob/main/Makefile#L23

---

## ⚙️ Flags

 
| Flag          | Default | Description                                     |
| ------------- | ------- | ----------------------------------------------- |
| `-port`       | `8080`  | Port to run the file server on                  |
| `-serve-dir`  | `.`     | Directory to serve files from                   |
| `-watch-dirs` |         | Additional comma-separated directories to watch |
| `-cmd`        |         | Command to execute when source files change     |

## 🧠 How It Works

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

This two-stage flow:

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

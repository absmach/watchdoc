#!/usr/bin/env sh
set -e

REPO="absmach/watchdoc"
BINARY_NAME="watchdoc"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

info()  { printf "[INFO] %s\n" "$1"; }
error() { printf "[ERROR] %s\n" "$1" >&2; exit 1; }

# ---- OS detection -----------------------------------------------------------
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux|darwin) ;;
  *)
    error "Unsupported OS: $OS"
    ;;
esac

# ---- ARCH detection ---------------------------------------------------------
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    error "Unsupported architecture: $ARCH"
    ;;
esac

ASSET="${BINARY_NAME}-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"

info "Installing ${BINARY_NAME}"
info "OS:   ${OS}"
info "ARCH: ${ARCH}"
info "From: ${URL}"

# ---- Download ---------------------------------------------------------------
TMP="$(mktemp)"
curl -fL "$URL" -o "$TMP" || error "Download failed"

chmod +x "$TMP"

# ---- Install ---------------------------------------------------------------
if [ ! -w "$INSTALL_DIR" ]; then
  info "Installing to ${INSTALL_DIR} (requires sudo)"
  sudo mv "$TMP" "${INSTALL_DIR}/${BINARY_NAME}"
else
  mv "$TMP" "${INSTALL_DIR}/${BINARY_NAME}"
fi

info "Installed to ${INSTALL_DIR}/${BINARY_NAME}"
info "Run '${BINARY_NAME} --help' to get started"

#!/usr/bin/env bash
# Verify generated binders/writers under TinyGo.
# TinyGo 0.40 requires Go 1.19–1.25 (not 1.26+).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if ! command -v tinygo >/dev/null; then
  echo "tinygo not installed" >&2
  exit 1
fi

# Prefer a supported Go toolchain when the system Go is too new.
if go version | grep -qE 'go1\.(2[6-9]|[3-9][0-9])'; then
  export GOTOOLCHAIN="${GOTOOLCHAIN:-go1.25.4}"
  GOROOT="$(GOTOOLCHAIN=$GOTOOLCHAIN go env GOROOT)"
  export GOROOT
  export PATH="$GOROOT/bin:$PATH"
  echo "using GOROOT=$GOROOT for TinyGo"
fi

echo "==> tinygo version"
tinygo version

echo "==> tinygo test (runtime + generated mapping)"
tinygo test ./internal/mappingfixture ./internal/tinycheck .

echo "==> tinygo run smoke (Bind/Write via generated code)"
tinygo run ./testdata/cmd/tinygo-bind-smoke

echo "OK: TinyGo checks passed"

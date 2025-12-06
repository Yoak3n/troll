#!/usr/bin/env bash
set -euo pipefail

# Build script for Linux/macOS (bash)
# Steps:
# 1) Build frontend in viewer/viewer-frontend
# 2) Ensure output is placed in viewer/service/app/dist, also copy to troll/dist
# 3) Build viewer cmd program (outputs viewer)
# 4) Build troll main program (outputs troll)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

FRONTEND_DIR="$REPO_ROOT/viewer/viewer-frontend"
ROUTER_DIST_DIR="$REPO_ROOT/viewer/service/app/dist"
TROLL_DIST_DIR="$REPO_ROOT/troll/dist"

echo "[1/4] Building frontend (viewer/viewer-frontend)"

pushd "$FRONTEND_DIR" >/dev/null
if command -v pnpm >/dev/null 2>&1; then
  echo "Using pnpm"
  pnpm install --frozen-lockfile
  pnpm run build
elif command -v npm >/dev/null 2>&1; then
  echo "Using npm"
  npm ci
  npm run build
else
  echo "ERROR: Neither 'pnpm' nor 'npm' found. Please install one of them." >&2
  exit 1
fi
popd >/dev/null

echo "[2/4] Sync dist to viewer/service/app and troll/dist"

mkdir -p "$ROUTER_DIST_DIR"

# Fallback sync from frontend/dist in case Vite didn't write to router/dist
FRONTEND_DIST_DIR="$FRONTEND_DIR/dist"
if [ -d "$FRONTEND_DIST_DIR" ]; then
  rm -rf "$ROUTER_DIST_DIR"
  mkdir -p "$ROUTER_DIST_DIR"
  cp -a "$FRONTEND_DIST_DIR/." "$ROUTER_DIST_DIR/"
fi

rm -rf "$TROLL_DIST_DIR"
mkdir -p "$TROLL_DIST_DIR"
cp -a "$ROUTER_DIST_DIR/." "$TROLL_DIST_DIR/"

echo "[3/4] Building viewer cmd program"

VIEWER_DIR="$REPO_ROOT/viewer"
mkdir -p "$VIEWER_DIR/bin"
pushd "$VIEWER_DIR" >/dev/null
go build -o "$VIEWER_DIR/bin/viewer" ./cmd
popd >/dev/null

echo "[4/4] Building troll main program"

TROLL_DIR="$REPO_ROOT/troll"
pushd "$TROLL_DIR" >/dev/null
go build -o "$TROLL_DIR/troll"
popd >/dev/null

echo "Build complete."
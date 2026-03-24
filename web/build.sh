#!/bin/bash
# Build all Go WASM binaries for claudefun web
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🔨 Building WASM binaries..."

# Constellation
echo "  ✨ Constellation..."
cd "$SCRIPT_DIR/constellation/wasm"
GOOS=js GOARCH=wasm go build -o ../constellation.wasm .
echo "     → constellation.wasm"

# Pizza
echo "  🍕 Pizza..."
cd "$SCRIPT_DIR/pizza/wasm"
GOOS=js GOARCH=wasm go build -o ../pizza.wasm .
echo "     → pizza.wasm"

# Octopus
echo "  🐙 Octopus..."
cd "$SCRIPT_DIR/octopus/wasm"
GOOS=js GOARCH=wasm go build -o ../octopus.wasm .
echo "     → octopus.wasm"

# Kingdom
echo "  🏰 Kingdom..."
cd "$SCRIPT_DIR/kingdom/wasm"
GOOS=js GOARCH=wasm go build -o ../kingdom.wasm .
echo "     → kingdom.wasm"

echo ""
echo "✅ All WASM binaries built!"
ls -lh "$SCRIPT_DIR"/constellation/constellation.wasm \
       "$SCRIPT_DIR"/pizza/pizza.wasm \
       "$SCRIPT_DIR"/octopus/octopus.wasm \
       "$SCRIPT_DIR"/kingdom/kingdom.wasm 2>/dev/null

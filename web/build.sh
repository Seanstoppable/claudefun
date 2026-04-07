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

# Cartographer
echo "  🗺️ Cartographer..."
cd "$SCRIPT_DIR/cartographer/wasm"
GOOS=js GOARCH=wasm go build -o ../cartographer.wasm .
echo "     → cartographer.wasm"

# Elevator
echo "  🎵 Elevator..."
cd "$SCRIPT_DIR/elevator/wasm"
GOOS=js GOARCH=wasm go build -o ../elevator.wasm .
echo "     → elevator.wasm"

echo ""
echo "✅ All WASM binaries built!"
ls -lh "$SCRIPT_DIR"/constellation/constellation.wasm \
       "$SCRIPT_DIR"/pizza/pizza.wasm \
       "$SCRIPT_DIR"/octopus/octopus.wasm \
       "$SCRIPT_DIR"/kingdom/kingdom.wasm \
       "$SCRIPT_DIR"/cartographer/cartographer.wasm \
       "$SCRIPT_DIR"/elevator/elevator.wasm 2>/dev/null

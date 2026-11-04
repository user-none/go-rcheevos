#!/bin/bash
#
# Updates the embedded rcheevos C library to a specific version.
#
# Usage: ./scripts/update-rcheevos.sh <version>
# Example: ./scripts/update-rcheevos.sh v12.2.1
#

set -e

VERSION=${1:-}

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v12.2.1"
    echo ""
    echo "Available versions:"
    git ls-remote --tags https://github.com/RetroAchievements/rcheevos.git | grep -v '\^{}' | awk -F/ '{print $3}' | sort -V | tail -10
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(dirname "$SCRIPT_DIR")"
RCHEEVOS_DIR="$REPO_DIR/rcheevos"

echo "Updating rcheevos to $VERSION..."

# Clone to temp directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

git clone --depth 1 --branch "$VERSION" https://github.com/RetroAchievements/rcheevos.git "$TEMP_DIR/rcheevos"

# Remove git metadata and unnecessary files
rm -rf "$TEMP_DIR/rcheevos/.git"
rm -rf "$TEMP_DIR/rcheevos/.github"
rm -rf "$TEMP_DIR/rcheevos/test"
rm -rf "$TEMP_DIR/rcheevos/validator"
rm -f "$TEMP_DIR/rcheevos/.gitignore"
rm -f "$TEMP_DIR/rcheevos/.editorconfig"
rm -f "$TEMP_DIR/rcheevos/Package.swift"

# Write version file
echo "$VERSION" > "$TEMP_DIR/rcheevos/VERSION"

# Replace existing rcheevos directory
rm -rf "$RCHEEVOS_DIR"
mv "$TEMP_DIR/rcheevos" "$RCHEEVOS_DIR"

echo "Successfully updated rcheevos to $VERSION"
echo ""
echo "Next steps:"
echo "  1. Review changes: git diff"
echo "  2. Test build: go build ./..."
echo "  3. Run tests: go test ./..."
echo "  4. Commit: git add -A && git commit -m 'Update rcheevos to $VERSION'"

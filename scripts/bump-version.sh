#!/bin/bash
# Bump version across all files
# Usage: ./scripts/bump-version.sh 0.2.0

set -e

NEW_VERSION=$1

if [ -z "$NEW_VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 0.2.0"
    exit 1
fi

echo "🔄 Bumping version to $NEW_VERSION..."

# Update main.go
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/version = \"[^\"]*\"/version = \"$NEW_VERSION\"/" cmd/machpay/main.go
else
    # Linux
    sed -i "s/version = \"[^\"]*\"/version = \"$NEW_VERSION\"/" cmd/machpay/main.go
fi

# Verify the change
echo ""
echo "📝 Updated files:"
grep -n "version = " cmd/machpay/main.go

# Build and verify
echo ""
echo "🔨 Building to verify..."
go build -o /tmp/machpay-test ./cmd/machpay
VERSION_OUTPUT=$(/tmp/machpay-test version 2>&1 || true)
echo "Version output: $VERSION_OUTPUT"
rm -f /tmp/machpay-test

echo ""
echo "✅ Version bumped to $NEW_VERSION"
echo ""
echo "Next steps:"
echo "  git add ."
echo "  git commit -m 'chore: bump version to $NEW_VERSION'"
echo "  git tag -a v$NEW_VERSION -m 'Release v$NEW_VERSION'"
echo "  git push origin main"
echo "  git push origin v$NEW_VERSION"


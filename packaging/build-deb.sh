#!/bin/bash
set -e

VERSION="0.1.1"
ARCH="amd64"
PKG_NAME="mcp-is_${VERSION}_${ARCH}"
BUILD_DIR="packaging/${PKG_NAME}"
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

echo "=== Building MCP IS .deb package ==="

# Clean previous build
rm -rf "${PROJECT_DIR}/packaging/${PKG_NAME}"
rm -f "${PROJECT_DIR}/packaging/${PKG_NAME}.deb"

# Create directory structure
mkdir -p "${BUILD_DIR}/DEBIAN"
mkdir -p "${BUILD_DIR}/opt/mcp-is"
mkdir -p "${BUILD_DIR}/etc/mcp-is"
mkdir -p "${BUILD_DIR}/etc/systemd/system"
mkdir -p "${BUILD_DIR}/usr/share/doc/mcp-is"

# Build Go binary
echo "Building Go binary..."
cd "${PROJECT_DIR}"
CGO_ENABLED=0 go build -tags netgo -o "${BUILD_DIR}/opt/mcp-is/mcp-is" ./cmd/server

# Copy project files
echo "Copying project files..."
cp "${PROJECT_DIR}/docker-compose.yaml" "${BUILD_DIR}/opt/mcp-is/"
cp "${PROJECT_DIR}/Dockerfile" "${BUILD_DIR}/opt/mcp-is/"
cp "${PROJECT_DIR}/.env.example" "${BUILD_DIR}/opt/mcp-is/"
cp "${PROJECT_DIR}/go.mod" "${BUILD_DIR}/opt/mcp-is/"
cp "${PROJECT_DIR}/go.sum" "${BUILD_DIR}/opt/mcp-is/"
cp -r "${PROJECT_DIR}/internal" "${BUILD_DIR}/opt/mcp-is/"
cp -r "${PROJECT_DIR}/cmd" "${BUILD_DIR}/opt/mcp-is/"
cp -r "${PROJECT_DIR}/web" "${BUILD_DIR}/opt/mcp-is/"
cp "${PROJECT_DIR}/packaging/update.sh" "${BUILD_DIR}/opt/mcp-is/" 2>/dev/null || true

# Copy DEBIAN files
cp "${PROJECT_DIR}/packaging/DEBIAN/control" "${BUILD_DIR}/DEBIAN/"
cp "${PROJECT_DIR}/packaging/DEBIAN/conffiles" "${BUILD_DIR}/DEBIAN/"
cp "${PROJECT_DIR}/packaging/DEBIAN/postinst" "${BUILD_DIR}/DEBIAN/"
cp "${PROJECT_DIR}/packaging/DEBIAN/preinst" "${BUILD_DIR}/DEBIAN/"
cp "${PROJECT_DIR}/packaging/DEBIAN/prerm" "${BUILD_DIR}/DEBIAN/"

# Make scripts executable
chmod 755 "${BUILD_DIR}/DEBIAN/postinst"
chmod 755 "${BUILD_DIR}/DEBIAN/preinst"
chmod 755 "${BUILD_DIR}/DEBIAN/prerm"

# Copy systemd service
cp "${PROJECT_DIR}/packaging/mcp-is.service" "${BUILD_DIR}/etc/systemd/system/"

# Copy documentation
cp "${PROJECT_DIR}/README.md" "${BUILD_DIR}/usr/share/doc/mcp-is/"
cp "${PROJECT_DIR}/docs/setup/INSTALL.md" "${BUILD_DIR}/usr/share/doc/mcp-is/"
cp "${PROJECT_DIR}/docs/setup/USAGE.md" "${BUILD_DIR}/usr/share/doc/mcp-is/"
cp "${PROJECT_DIR}/docs/setup/API.md" "${BUILD_DIR}/usr/share/doc/mcp-is/"

# Create changelog
echo "mcp-is (${VERSION}) stable; urgency=low" > "${BUILD_DIR}/usr/share/doc/mcp-is/changelog"
echo "" >> "${BUILD_DIR}/usr/share/doc/mcp-is/changelog"
echo "  * Initial release" >> "${BUILD_DIR}/usr/share/doc/mcp-is/changelog"
echo "  * MCP protocol support" >> "${BUILD_DIR}/usr/share/doc/mcp-is/changelog"
echo "  * InSales API integration" >> "${BUILD_DIR}/usr/share/doc/mcp-is/changelog"
echo "" >> "${BUILD_DIR}/usr/share/doc/mcp-is/changelog"
echo " -- MCP IS Team <team@mcp-is.ru>  $(date -R)" >> "${BUILD_DIR}/usr/share/doc/mcp-is/changelog"

gzip -9 "${BUILD_DIR}/usr/share/doc/mcp-is/changelog"

# Set permissions
chmod -R 755 "${BUILD_DIR}"
chmod 644 "${BUILD_DIR}/DEBIAN/control"
chmod 644 "${BUILD_DIR}/DEBIAN/conffiles"
chmod 644 "${BUILD_DIR}/etc/systemd/system/mcp-is.service"

# Build .deb package
echo "Building .deb package..."
cd "${PROJECT_DIR}/packaging"
dpkg-deb --build "${PKG_NAME}"

echo ""
echo "=== Package built successfully ==="
echo "File: packaging/${PKG_NAME}.deb"
echo ""
echo "Install with: sudo dpkg -i packaging/${PKG_NAME}.deb"
echo "Configure:    sudo nano /etc/mcp-is/.env"
echo "Start:        sudo systemctl start mcp-is"

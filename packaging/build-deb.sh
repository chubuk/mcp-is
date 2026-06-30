#!/bin/bash
set -e

VERSION="0.1.0"
ARCH="amd64"
PKG_NAME="mcp-insales_${VERSION}_${ARCH}"
BUILD_DIR="packaging/${PKG_NAME}"
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

echo "=== Building MCP InSales .deb package ==="

# Clean previous build
rm -rf "${PROJECT_DIR}/packaging/${PKG_NAME}"
rm -f "${PROJECT_DIR}/packaging/${PKG_NAME}.deb"

# Create directory structure
mkdir -p "${BUILD_DIR}/DEBIAN"
mkdir -p "${BUILD_DIR}/opt/mcp-insales"
mkdir -p "${BUILD_DIR}/etc/mcp-insales"
mkdir -p "${BUILD_DIR}/etc/systemd/system"
mkdir -p "${BUILD_DIR}/usr/share/doc/mcp-insales"

# Build Go binary
echo "Building Go binary..."
cd "${PROJECT_DIR}"
CGO_ENABLED=0 go build -tags netgo -o "${BUILD_DIR}/opt/mcp-insales/mcp-insales" ./cmd/server

# Copy project files
echo "Copying project files..."
cp "${PROJECT_DIR}/docker-compose.yaml" "${BUILD_DIR}/opt/mcp-insales/"
cp "${PROJECT_DIR}/Dockerfile" "${BUILD_DIR}/opt/mcp-insales/"
cp "${PROJECT_DIR}/.env.example" "${BUILD_DIR}/opt/mcp-insales/"
cp "${PROJECT_DIR}/go.mod" "${BUILD_DIR}/opt/mcp-insales/"
cp "${PROJECT_DIR}/go.sum" "${BUILD_DIR}/opt/mcp-insales/"
cp -r "${PROJECT_DIR}/internal" "${BUILD_DIR}/opt/mcp-insales/"
cp -r "${PROJECT_DIR}/cmd" "${BUILD_DIR}/opt/mcp-insales/"
cp -r "${PROJECT_DIR}/web" "${BUILD_DIR}/opt/mcp-insales/"
cp "${PROJECT_DIR}/packaging/update.sh" "${BUILD_DIR}/opt/mcp-insales/" 2>/dev/null || true

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
cp "${PROJECT_DIR}/packaging/mcp-insales.service" "${BUILD_DIR}/etc/systemd/system/"

# Copy documentation
cp "${PROJECT_DIR}/README.md" "${BUILD_DIR}/usr/share/doc/mcp-insales/"
cp "${PROJECT_DIR}/docs/setup/INSTALL.md" "${BUILD_DIR}/usr/share/doc/mcp-insales/"
cp "${PROJECT_DIR}/docs/setup/USAGE.md" "${BUILD_DIR}/usr/share/doc/mcp-insales/"
cp "${PROJECT_DIR}/docs/setup/API.md" "${BUILD_DIR}/usr/share/doc/mcp-insales/"

# Create changelog
echo "mcp-insales (${VERSION}) stable; urgency=low" > "${BUILD_DIR}/usr/share/doc/mcp-insales/changelog"
echo "" >> "${BUILD_DIR}/usr/share/doc/mcp-insales/changelog"
echo "  * Initial release" >> "${BUILD_DIR}/usr/share/doc/mcp-insales/changelog"
echo "  * MCP protocol support" >> "${BUILD_DIR}/usr/share/doc/mcp-insales/changelog"
echo "  * InSales API integration" >> "${BUILD_DIR}/usr/share/doc/mcp-insales/changelog"
echo "" >> "${BUILD_DIR}/usr/share/doc/mcp-insales/changelog"
echo " -- MCP InSales Team <team@mcp-insales.ru>  $(date -R)" >> "${BUILD_DIR}/usr/share/doc/mcp-insales/changelog"

gzip -9 "${BUILD_DIR}/usr/share/doc/mcp-insales/changelog"

# Set permissions
chmod -R 755 "${BUILD_DIR}"
chmod 644 "${BUILD_DIR}/DEBIAN/control"
chmod 644 "${BUILD_DIR}/DEBIAN/conffiles"
chmod 644 "${BUILD_DIR}/etc/systemd/system/mcp-insales.service"

# Build .deb package
echo "Building .deb package..."
cd "${PROJECT_DIR}/packaging"
dpkg-deb --build "${PKG_NAME}"

echo ""
echo "=== Package built successfully ==="
echo "File: packaging/${PKG_NAME}.deb"
echo ""
echo "Install with: sudo dpkg -i packaging/${PKG_NAME}.deb"
echo "Configure:    sudo nano /etc/mcp-insales/.env"
echo "Start:        sudo systemctl start mcp-insales"

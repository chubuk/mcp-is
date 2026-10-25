#!/bin/bash
set -e

# Single source of truth: internal/mcp/server.go version + CHANGELOG.md
VERSION="${VERSION:-0.7.5}"
ARCH="amd64"
PKG_NAME="mcp-is_${VERSION}_${ARCH}"
BUILD_DIR="packaging/${PKG_NAME}"
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

echo "=== Building MCP IS .deb package ${VERSION} ==="

# Clean previous build
rm -rf "${PROJECT_DIR}/packaging/${PKG_NAME}"
rm -f "${PROJECT_DIR}/packaging/${PKG_NAME}.deb"

# Create directory structure
mkdir -p "${BUILD_DIR}/DEBIAN"
mkdir -p "${BUILD_DIR}/opt/mcp-is"
mkdir -p "${BUILD_DIR}/etc/mcp-is"
mkdir -p "${BUILD_DIR}/etc/systemd/system"
mkdir -p "${BUILD_DIR}/etc/logrotate.d"
mkdir -p "${BUILD_DIR}/usr/share/doc/mcp-is"

# Build Go binary
echo "Building Go binary..."
cd "${PROJECT_DIR}"
CGO_ENABLED=0 go build -tags netgo -ldflags="-s -w" -o "${BUILD_DIR}/opt/mcp-is/mcp-is" ./cmd/server
echo "${VERSION}" > "${BUILD_DIR}/opt/mcp-is/VERSION"

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
cp -r "${PROJECT_DIR}/scripts" "${BUILD_DIR}/opt/mcp-is/"
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
chmod 755 "${BUILD_DIR}/opt/mcp-is/scripts/"*.sh 2>/dev/null || true

# Copy systemd service + logrotate
cp "${PROJECT_DIR}/packaging/mcp-is.service" "${BUILD_DIR}/etc/systemd/system/"
cp "${PROJECT_DIR}/packaging/logrotate-mcp-is-backup" "${BUILD_DIR}/etc/logrotate.d/mcp-is-backup"

# Copy documentation (all current docs, including InSales app mode)
cp "${PROJECT_DIR}/README.md" "${BUILD_DIR}/usr/share/doc/mcp-is/"
cp "${PROJECT_DIR}/LICENSE" "${BUILD_DIR}/usr/share/doc/mcp-is/"
cp "${PROJECT_DIR}/CHANGELOG.md" "${BUILD_DIR}/usr/share/doc/mcp-is/"
cp -r "${PROJECT_DIR}/docs" "${BUILD_DIR}/usr/share/doc/mcp-is/"

# Create changelog (first entries from CHANGELOG.md header)
{
    echo "mcp-is (${VERSION}) stable; urgency=low"
    echo ""
    echo "  * See CHANGELOG.md for full history"
    echo "  * InSales external app mode (install/uninstall/autologin)"
    echo "  * REGISTRATION_ENABLED default false"
    echo ""
    echo " -- Vladislav Chubuk <chubuk@mail.ru>  $(date -R)"
} > "${BUILD_DIR}/usr/share/doc/mcp-is/changelog.Debian"

gzip -9 "${BUILD_DIR}/usr/share/doc/mcp-is/changelog.Debian"

# Set permissions
chmod -R 755 "${BUILD_DIR}"
chmod 644 "${BUILD_DIR}/DEBIAN/control"
chmod 644 "${BUILD_DIR}/DEBIAN/conffiles"
chmod 644 "${BUILD_DIR}/etc/systemd/system/mcp-is.service"
chmod 644 "${BUILD_DIR}/etc/logrotate.d/mcp-is-backup"

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

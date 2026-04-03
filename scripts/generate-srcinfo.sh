#!/usr/bin/env bash
# Usage: generate-srcinfo.sh <version> <sha256_x86_64> <sha256_aarch64>
set -euo pipefail

VERSION="$1"
SHA_X86="$2"
SHA_AARCH64="$3"
BASE_URL="https://github.com/jakobnielsen/tztui/releases/download/v${VERSION}"

cat > .SRCINFO <<EOF
pkgbase = tztui-bin
	pkgdesc = A terminal UI for managing and browsing system timezones
	pkgver = ${VERSION}
	pkgrel = 1
	url = https://github.com/jakobnielsen/tztui
	arch = x86_64
	arch = aarch64
	license = MIT
	provides = tztui
	conflicts = tztui
	source_x86_64 = ${BASE_URL}/tztui_linux_amd64.tar.gz
	sha256sums_x86_64 = ${SHA_X86}
	source_aarch64 = ${BASE_URL}/tztui_linux_arm64.tar.gz
	sha256sums_aarch64 = ${SHA_AARCH64}

pkgname = tztui-bin
	provides = tztui
	conflicts = tztui
EOF

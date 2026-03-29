# Maintainer: Your Name <you@example.com>
pkgname=tztui
pkgver=0.1.0
pkgrel=1
pkgdesc="A terminal UI for managing and browsing system timezones"
arch=('x86_64' 'aarch64')
url="https://github.com/jakobnielsen/tztui"
license=('MIT')
makedepends=('go')
source=("${pkgname}-${pkgver}.tar.gz::${url}/archive/v${pkgver}.tar.gz")
sha256sums=('SKIP')

build() {
    cd "${pkgname}-${pkgver}"
    go build \
        -trimpath \
        -buildmode=pie \
        -ldflags "-linkmode external -extldflags \"${LDFLAGS}\" -X main.version=${pkgver}" \
        -o "${pkgname}" \
        ./cmd/tztui
}

package() {
    cd "${pkgname}-${pkgver}"
    install -Dm755 "${pkgname}" "${pkgdir}/usr/bin/${pkgname}"
    install -Dm644 README.md "${pkgdir}/usr/share/doc/${pkgname}/README.md"
}

class Tztui < Formula
  desc "A terminal UI for managing and browsing system timezones"
  homepage "https://github.com/jakobnielsen/tztui"
  version "0.1.0"

  on_macos do
    on_arm do
      url "https://github.com/jakobnielsen/tztui/releases/download/v#{version}/tztui_darwin_arm64.tar.gz"
      sha256 "REPLACE_WITH_GORELEASER_SHA256_darwin_arm64"
    end
    on_intel do
      url "https://github.com/jakobnielsen/tztui/releases/download/v#{version}/tztui_darwin_amd64.tar.gz"
      sha256 "REPLACE_WITH_GORELEASER_SHA256_darwin_amd64"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/jakobnielsen/tztui/releases/download/v#{version}/tztui_linux_arm64.tar.gz"
      sha256 "REPLACE_WITH_GORELEASER_SHA256_linux_arm64"
    end
    on_intel do
      url "https://github.com/jakobnielsen/tztui/releases/download/v#{version}/tztui_linux_amd64.tar.gz"
      sha256 "REPLACE_WITH_GORELEASER_SHA256_linux_amd64"
    end
  end

  def install
    bin.install "tztui"
  end

  test do
    system "#{bin}/tztui", "--version"
  end
end

class Cfpkg < Formula
  desc "K8s installation tool"
  homepage "https://github.com/wonderix/cfpkg"
  version "{{version}}"
  if OS.mac?
    url "https://github.com/wonderix/cfpkg/releases/download/{{version}}/cfpkg-binary-darwin.tgz"
    sha256 "{{sha256-darwin}}"
  elsif OS.linux?
    url "https://github.com/wonderix/cfpkg/releases/download/{{version}}/cfpkg-binary-linux.tgz"
    sha256 "{{sha256-linux}}"
  end

  depends_on :arch => :x86_64

  def install
    bin.install "cfpkg" => "cfpkg"
  end

  test do
    system "#{bin}/cfpkg", "version"
  end
end
#!/usr/bin/env sh
# Install drizz: download the latest release for this operating system and processor from GitHub, verify its checksum,
# and place it on PATH. Intended for: curl -fsSL https://get.drizz.dev | sh
set -eu

repo="DrizzDev/platform"
dest="${DRIZZ_INSTALL:-$HOME/.local/bin}"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64 | amd64) arch="amd64" ;;
  arm64 | aarch64) arch="arm64" ;;
  *) echo "drizz: unsupported architecture $arch" >&2; exit 1 ;;
esac

tag="$(curl -fsSL "https://api.github.com/repos/$repo/releases/latest" |
  sed -n 's/.*"tag_name" *: *"\([^"]*\)".*/\1/p' | head -1)"
[ -n "$tag" ] || { echo "drizz: could not resolve the latest release" >&2; exit 1; }

archive="drizz-$os-$arch.tar.gz"
base="https://github.com/$repo/releases/download/$tag"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

curl -fsSL "$base/$archive" -o "$tmp/$archive"
curl -fsSL "$base/checksums.txt" -o "$tmp/checksums.txt"

# Verify the download against the published checksum before installing it.
( cd "$tmp" && grep " $archive\$" checksums.txt | shasum -a 256 -c - >/dev/null ) ||
  { echo "drizz: checksum verification failed" >&2; exit 1; }

tar -xzf "$tmp/$archive" -C "$tmp"
mkdir -p "$dest"
install -m 0755 "$tmp/drizz" "$dest/drizz"

echo "drizz: installed to $dest/drizz ($tag)"
case ":$PATH:" in
  *":$dest:"*) ;;
  *) echo "drizz: add $dest to your PATH" ;;
esac

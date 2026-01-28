#!/usr/bin/env bash
set -euo pipefail

force=false
if [[ "${1:-}" == "--force" ]]; then
  force=true
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/../../.." && pwd)"
bin_dir="${repo_root}/bin"
target_bin="${bin_dir}/ezorm"

if [[ -x "${target_bin}" && "${force}" == "false" ]]; then
  printf 'ezorm already present: %s\n' "${target_bin}"
  exit 0
fi

mkdir -p "${bin_dir}"

try_go_install() {
  if ! command -v go >/dev/null 2>&1; then
    return 1
  fi

  printf 'installing ezorm via go install...\n'
  if ! go install github.com/ezbuy/ezorm/v2@latest; then
    return 1
  fi

  local gopath
  gopath="$(go env GOPATH)"
  if [[ -z "${gopath}" ]]; then
    return 1
  fi

  local installed_bin="${gopath}/bin/ezorm"
  if [[ ! -x "${installed_bin}" ]]; then
    return 1
  fi

  cp "${installed_bin}" "${target_bin}"
  chmod +x "${target_bin}"
  printf 'installed ezorm to %s\n' "${target_bin}"
  return 0
}

try_release_download() {
  if ! command -v curl >/dev/null 2>&1; then
    printf 'curl is required to download ezorm release\n'
    return 1
  fi

  local os arch
  case "$(uname -s)" in
    Linux) os="linux" ;;
    Darwin) os="darwin" ;;
    *)
      printf 'unsupported OS: %s\n' "$(uname -s)"
      return 1
      ;;
  esac

  case "$(uname -m)" in
    x86_64) arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *)
      printf 'unsupported arch: %s\n' "$(uname -m)"
      return 1
      ;;
  esac

  local py
  if command -v python3 >/dev/null 2>&1; then
    py="python3"
  elif command -v python >/dev/null 2>&1; then
    py="python"
  else
    printf 'python is required to parse release metadata\n'
    return 1
  fi

  printf 'downloading ezorm release for %s/%s...\n' "${os}" "${arch}"

  local release_json url
  release_json="$(curl -fsSL https://api.github.com/repos/ezbuy/ezorm/releases/latest)"
  url="$(${py} - "${os}" "${arch}" <<'PY'
import json
import sys

target_os = sys.argv[1]
target_arch = sys.argv[2]

data = json.load(sys.stdin)
assets = data.get("assets", [])
for asset in assets:
    name = asset.get("name", "")
    lower = name.lower()
    if "ezorm" not in lower:
        continue
    if target_os in lower and target_arch in lower:
        print(asset.get("browser_download_url", ""))
        raise SystemExit(0)
print("")
raise SystemExit(1)
PY
)"

  if [[ -z "${url}" ]]; then
    printf 'no release asset found for %s/%s\n' "${os}" "${arch}"
    return 1
  fi

  local tmpdir tmpfile
  tmpdir="$(mktemp -d)"
  tmpfile="${tmpdir}/ezorm_download"
  curl -fsSL "${url}" -o "${tmpfile}"

  case "${url}" in
    *.tar.gz)
      tar -xzf "${tmpfile}" -C "${tmpdir}"
      ;;
    *.zip)
      unzip -q "${tmpfile}" -d "${tmpdir}"
      ;;
    *)
      cp "${tmpfile}" "${target_bin}"
      chmod +x "${target_bin}"
      printf 'installed ezorm to %s\n' "${target_bin}"
      rm -rf "${tmpdir}"
      return 0
      ;;
  esac

  local found
  found="$(find "${tmpdir}" -type f -name ezorm -maxdepth 3 | head -n 1)"
  if [[ -z "${found}" ]]; then
    printf 'ezorm binary not found in release archive\n'
    rm -rf "${tmpdir}"
    return 1
  fi

  cp "${found}" "${target_bin}"
  chmod +x "${target_bin}"
  rm -rf "${tmpdir}"
  printf 'installed ezorm to %s\n' "${target_bin}"
  return 0
}

if try_go_install; then
  exit 0
fi

if try_release_download; then
  exit 0
fi

printf 'failed to install ezorm; try go install github.com/ezbuy/ezorm/v2@latest\n'
exit 1

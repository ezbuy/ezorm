#!/usr/bin/env bash
set -euo pipefail

force=false
if [[ "${1:-}" == "--force" ]]; then
  force=true
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
assets_dir="${script_dir}/../assets"
refs_dir="${script_dir}/../references"

mkdir -p "${assets_dir}" "${refs_dir}"

download() {
  local src_path="$1"
  local dst_path="$2"

  if [[ -f "${dst_path}" && "${force}" == "false" ]]; then
    return 0
  fi

  curl -fsSL "https://raw.githubusercontent.com/ezbuy/ezorm/main/${src_path}" -o "${dst_path}"
}

# References

download "Makefile" "${refs_dir}/Makefile"

# Assets (examples)

download "e2e/mysql/user.yaml" "${assets_dir}/mysql_user.yaml"
download "e2e/mysql/blog.yaml" "${assets_dir}/mysql_blog.yaml"
download "e2e/mysqlr/user.yaml" "${assets_dir}/mysqlr_user.yaml"
download "e2e/mysqlr/blog.yaml" "${assets_dir}/mysqlr_blog.yaml"

printf 'init complete (gen mysql/mysqlr)\n'

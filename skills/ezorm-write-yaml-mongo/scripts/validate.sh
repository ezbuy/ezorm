#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
assets_dir="${script_dir}/../assets"
refs_dir="${script_dir}/../references"

required_files=(
  "${refs_dir}/yaml.md"
  "${refs_dir}/yaml_mongo.md"
  "${assets_dir}/mongo_user.yaml"
  "${assets_dir}/mongo_nested.yaml"
)

missing=0
for f in "${required_files[@]}"; do
  if [[ ! -f "${f}" ]]; then
    printf 'missing: %s\n' "${f}"
    missing=1
  fi
done

if [[ "${missing}" -ne 0 ]]; then
  printf 'run scripts/init.sh to download missing files\n'
  exit 1
fi

printf 'validate ok (mongo)\n'

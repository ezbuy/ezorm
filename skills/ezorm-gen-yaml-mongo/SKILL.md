---
name: ezorm-gen-yaml-mongo
description: Generate Go code from ezorm YAML (mongo). Use when asked to produce Go output from mongo YAML schemas; provide the exact ezorm CLI command(s).
---

# Ezorm Gen Yaml (Mongo)

## Scope

Generate Go code from ezorm YAML schema files for mongo using the `ezorm` CLI.
Read these files when command details or examples are needed:
- `Makefile` (target: `gene2e`)
- `e2e/mongo/`
If the local `references/` or `assets/` directories are missing, run:
`scripts/init.sh` (or `scripts/init.sh --force` to refresh).
Validate with `scripts/validate.sh`.
If `bin/ezorm` is missing, run `scripts/ensure-ezorm.sh` first.

## Workflow

1. Confirm input YAML path(s) and output directory.
2. Choose goPackage and namespace (required for `ezorm gen`).
3. Use the mongo command pattern:
   - `bin/ezorm gen -i <yaml-or-dir> -o <out-dir> --goPackage <pkg> --namespace <ns>`
4. Output only the command(s) needed. Do not generate Go code in the response.

## Output Rules

- Emit shell command(s) only.
- Do not invent paths; use the user-provided locations.
- Prefer a single command per directory when input is a folder.

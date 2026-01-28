---
name: ezorm-write-yaml-mongo
description: Write ezorm YAML schema files for this repo (mongo only). Use when asked to define or update ezorm YAML for mongo collections, fields, indexes, flags, embeds, or relations; output YAML only (no Go code).
---

# Ezorm Write Yaml

## Scope

Write ezorm YAML schemas from user requirements for mongo and emit YAML only.
Read these files when schema details or examples are needed:
- `doc/schema/yaml.md`
- `doc/schema/yaml_mongo.md`
- `e2e/mongo/*.yaml`
If the local `references/` or `assets/` directories are missing, run:
`scripts/init.sh` (or `scripts/init.sh --force` to refresh).
Validate with `scripts/validate.sh`.

## Workflow

1. Restate the domain model in simple entities and relationships. Ask one clarification only if critical (missing entities, relationship direction, or required constraints).
2. Choose representation:
   - Embedded structure: define a separate entity with `embed: true` and reference it with `list<EmbedType>` or `EmbedType` fields.
   - Separate collection: define separate entities and connect via ID fields.
3. Map fields:
   - Use mongo field types from `doc/schema/yaml_mongo.md`.
   - Add `flags` like `index`, `unique`, `sort`, `nullable` as needed.
   - Use `attrs` with `bsonTag`/`jsonTag` only when a specific storage name is requested.
4. Set collection metadata:
   - Always include `table` and `dbname`.
   - Default `table` to snake_case of the entity name.
   - Default `dbname` to `default` unless the user provides one.
5. Add constraints:
   - Use `indexes`, `uniques`, `primary` as required.
6. Output only YAML (no Go code, no prose). Use `---` to separate multiple entities.

## Output Rules

- Emit a valid YAML schema file (or multiple YAML documents) and nothing else.
- Keep names consistent with user domain language.
- Prefer minimal fields; do not invent behavior.

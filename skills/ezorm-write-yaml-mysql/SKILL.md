---
name: ezorm-write-yaml-mysql
description: Write ezorm YAML schema files for this repo (mysql/mysqlr). Prefer mysqlr unless the user explicitly asks for mysql. Output YAML only (no Go code).
---

# Ezorm Write Yaml (MySQL/MySQLr)

## Scope

Write ezorm YAML schemas from user requirements for mysql/mysqlr and emit YAML only.
Prefer mysqlr unless the user explicitly requests mysql.
Read these files when schema details or examples are needed:
- `doc/schema/yaml.md`
- `doc/schema/yaml_mysql.md`
- `doc/schema/yaml_mysqlr.md`
- `e2e/mysql/*.yaml`
- `e2e/mysqlr/*.yaml`
If the local `references/` or `assets/` directories are missing, run:
`scripts/init.sh` (or `scripts/init.sh --force` to refresh).
Validate with `scripts/validate.sh`.

## Workflow

1. Restate the domain model in simple entities and relationships. Ask one clarification only if critical (missing entities, relationship direction, or required constraints).
2. Choose db driver:
   - Use `db: mysql` for classic mysql schemas (see `doc/schema/yaml_mysql.md`).
   - Use `db: mysqlr` for mysqlr schemas (see `doc/schema/yaml_mysqlr.md`).
3. Map fields:
   - Use mysql/mysqlr field types from the matching doc.
   - Add `flags` like `index`, `unique`, `nullable`, `primary`, `autoinc`, `noinc` as needed.
   - Use `fk` only when an explicit foreign key is requested.
4. Set table metadata:
   - Always include `dbname`.
   - Use `table` for mysql, `dbtable` for mysqlr.
   - Default table names to snake_case of the entity name unless the user provides one.
5. Add constraints:
   - Use `indexes`, `uniques`, `primary` as required.
6. Output only YAML (no Go code, no prose). Use `---` to separate multiple entities.

## Output Rules

- Emit a valid YAML schema file (or multiple YAML documents) and nothing else.
- Keep names consistent with user domain language.
- Prefer minimal fields; do not invent behavior.

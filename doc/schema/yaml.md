# YAML Schema

EZORM schemas are defined in YAML; generated code is derived from YAML.
Do not infer extra fields or constraints not requested by the user.

## Required vs Optional

Required:
- `db`
- `fields`

Recommended:
- `dbname`
- `table` (mysql) or `dbtable` (mysqlr)

Optional:
- `comment`
- `indexes`
- `uniques`
- `primary`

## Example

```yaml
Blog:
  db: mongo
  comment: my blog
  fields:
    - Title: string
    - Hits: int32
    - Slug: string
      flags: [unique]
    - Body: string
    - User: int32
    - CreateDate: datetime
      flags: [sort]
    - IsPublished: bool
      flags: [index]
  indexes: [[User, IsPublished]]
```

## Schema Definition

Top-level YAML contains one or more entities. Each entity name is the
logical model name used by generated code.

> details: [db: mysql](yaml_mysql.md) | [db: mysqlr](yaml_mysqlr.md) | [db: mongo](yaml_mongo.md)

| Component | Remark | Definition in Example | Other Properties |
|---|---|---|---|
| Entity Name | Logical model name | `Blog` | / |
| `db` | Database driver | `mongo` | `mysql` / `mysqlr` |
| `fields` | Field definitions | `Title`, `Slug`, `Body`, `CreateDate`, `IsPublished` | / |
| Field Type | Field type | `string`, `int32`, `bool`, `datetime` | / |
| Field `flags` | Field properties | `unique`, `sort`, `index` | / |
| Constraints | Table or collection constraints | `indexes` | `uniques` / `primary` |
| `comment` | Entity comment | `comment` | / |

## Driver-specific Pointers

- mongo: [yaml_mongo.md](yaml_mongo.md), examples in `e2e/mongo/`
- mysql: [yaml_mysql.md](yaml_mysql.md), examples in `e2e/mysql/`
- mysqlr: [yaml_mysqlr.md](yaml_mysqlr.md), examples in `e2e/mysqlr/`

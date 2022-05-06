# YAML Schema

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

> details: [db: mysql](yaml_mysql.md) | [db: mysqlr](yaml_mysqlr.md) | [db: mongo](yaml_mongo.md)

| Component | Remark | Definition in Example | Other Properties |
|---|---|---|---|
| Database Name | The Name of database Object , the generated code will use this name as Access Manager|`Blog`| / |
| `db` | The name of database driver |[`mongo`](../../e2e/mongo/)| [`mysql`](../../e2e/mysql/) / [`mysqlr`](../../e2e/mysqlr/) |
| `fields` | The definition of fields | `Title`,`Slug`, `Body`, `CreateDate`,`IsPublished`| / |
| Field Type | The type of field | `string`, `int32`, `bool`, `datetime`| / |
| Field `flags` | The properties of fields | `unique`,`sort`,`index` | /|
|Database Constraint | Describe the table constraint |  `indexes` | `uniques`/ `primary` |
|Database Comment|Describe the table comment | `comment` | / |

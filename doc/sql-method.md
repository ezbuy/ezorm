# generate sql-method

ezorm's basic generated code can provide some methods, such as  `Insert`, `Find`, `FindByID`, etc. But for RDB, these generated sqls are too simple, they only involve one table, but sometimes we want to write more complex sql like adding `JOIN` statement to combine multi-tables.

So for the sql database, we allow users to customize some complex sql statements. The generator can generate the Go method to call these sqls.

There are some restrictions on the sql written by user:

- All query tables and fields must be appeared in ezorm's model definition.
- The `"SELECT *"` is forbidden.
- The `JOIN` clause cannot appear too many times.

We implement these restrictions to make sure users to write more standardized sql statements. ezorm will check whether the statements written by the user meet the above specifications, otherwise the generation will fail.

In addition, we recommend:

- Use GoTemplate to express all table and field names.
- If a sql snippet appears too frequently, it should be reused.

## Define models

The use of sql-method is based on ezorm models, we must first define models. Here I define a series of simple models for demonstration:

**For brevity, I omitted the definition of the indexes. Please do not forget to define indexes in real project.**

```yaml
User:
  db: mysql
  fields:
    - Id: int64
    - Name: string
    - Phone: string
    - Password: string

UserDetail:
  db: mysql
  fields:
    - Id: int64
    - UserId: int64
    - Email: string
    - Introduction: string
    - Age: int32
    - Avatar: string


RoleUser:
  db: mysql
  fields:
    - UserId: int64
    - RoleId: int64

Role:
  db: mysql
  fields:
    - Id: int64
    - Name: string
```

The relationships between these tables are:

- `User` to `UserDetail`: one to one.
- `User` to `Role`: many to many.
- `User` to `UserRole`: one to many.
- `Role` to `UserRole`: one to many.

You can use `ezorm -i models.yaml -o /path/to/generate` to generate basic orm code like before.

## Define methods

If we want to write some complex sql, we need to create a new yaml file, It is recommended to be placed in the same level as `model.yaml`, named `sql.yaml`.

Let's write a simple method: query all roles by user:

```yaml
methods:
  FindUserRoles:
    args: [userId int64]
    ret: list<*Role>
    sql: |
      SELECT
        {{ .Role.Id }},
        {{ .Role.Name }}
      FROM {{ .Role }}
      JOIN {{ .RoleUser }} ON {{ .RoleUser.RoleId }}={{ .Role.Id }}
      WHERE {{ .RoleUser.UserId }}={{ .args.userId }}
```

When defining methods, we must set `args`, `ret`, `sql` options:

- args: The args for method, it is an array.
- ret: The return type for method, 

## Execute generation

## Use in Golang

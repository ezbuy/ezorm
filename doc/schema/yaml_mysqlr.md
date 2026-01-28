## mysqlr

### Required vs Optional

Required:
- `db`
- `fields`

Recommended:
- `dbname`
- `dbtable`

Optional:
- `comment`
- `indexes`
- `uniques`
- `primary`

### Field Properties

* size
* sqltype
* sqlcolumn
* comment
* validator
* attrs (for custom attributes, such as `proto`, `bson`)
* flags (for custom flags, such as `unique`, `sort`, `index`)


### Field Types

* bool
* uint8
* uint16
* uint32
* uint64
* int8
* int16
* int32 (int)
* int64
* time.Time
* timeint / timestamp / datatime

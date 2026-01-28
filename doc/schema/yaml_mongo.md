## mongo

### Required vs Optional

Required:
- `db`
- `fields`

Recommended:
- `dbname`
- `table`

Optional:
- `comment`
- `indexes`
- `uniques`
- `primary`

### Field Properties

* label
* fk
* widget
* remark
* comment
* size
* flags
* attrs
* embed: support the embed-only structure , if this field sets to be true, the structure's orm code will not be generated anymore.

### Field Types
* []byte
* bool
* datetime
* float64
* int
* int32
* int64
* string
* time.Time / datetime
* timestamp

package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const query = `
SELECT
	id
FROM
	user
WHERE name IN ('me') AND id =1 AND phone ='123'
LIMIT 0,10
`

const queryOneLimit = `
SELECT
	id
FROM
	user
WHERE name IN ('me') AND id =1 AND phone ='123'
LIMIT 10
`

const queryNoLimit = `
SELECT
	id
FROM
	user
WHERE name IN ('me') AND id =1 AND phone ='123'
`

const queryWithTableJoin = `
SELECT
	u.id,
	b.id
FROM
	user  u
INNER JOIN
	blog  b
ON
	u.id = b.user_id
WHERE
	u.name = 'me'
`

const queryWithSubquery = `
SELECT
	id
FROM
	user
WHERE name IN (
	SELECT
		name
	FROM
		user
	WHERE id = 1
)
`

func TestTiDBParserParseQuery(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		expect string
	}{
		{"query", query, "SELECT `id` FROM `user` WHERE `name` IN (?) AND `id`=? AND `phone`=? LIMIT ?,?"},
		{"queryOneLimit", queryOneLimit, "SELECT `id` FROM `user` WHERE `name` IN (?) AND `id`=? AND `phone`=? LIMIT ?"},
		{"queryNoLimit", queryNoLimit, "SELECT `id` FROM `user` WHERE `name` IN (?) AND `id`=? AND `phone`=?"},
		{"queryWithTableJoin", queryWithTableJoin, "SELECT `u`.`id`,`b`.`id` FROM `user` AS `u` JOIN `blog` AS `b` ON `u`.`id`=`b`.`user_id` WHERE `u`.`name`=?"},
		{"queryWithSubquery", queryWithSubquery, "SELECT `id` FROM `user` WHERE `name` IN (SELECT `name` FROM `user` WHERE `id`=?)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := NewTiDBParser()
			err := tp.Parse(context.TODO(), tt.query)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, tp.Query())
		})
	}
}

func TestTiDBParserParseMetadata(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		metadata *QueryMetadata
	}{
		{"query", query, &QueryMetadata{
			params: []*QueryField{
				{Name: "col:`name`", Type: T_ARRAY_STRING},
				{Name: "col:`id`", Type: T_INT},
				{Name: "col:`phone`", Type: T_STRING},
				{Name: "limit:count", Type: T_INT},
				{Name: "limit:offset", Type: T_INT},
			},
			result: []*QueryField{
				{Name: "`id`", Type: T_PLACEHOLDER},
			},
		}},
		{"queryOneLimit", queryOneLimit, &QueryMetadata{
			params: []*QueryField{
				{Name: "col:`name`", Type: T_ARRAY_STRING},
				{Name: "col:`id`", Type: T_INT},
				{Name: "col:`phone`", Type: T_STRING},
				{Name: "limit:count", Type: T_INT},
			},

			result: []*QueryField{
				{Name: "`id`", Type: T_PLACEHOLDER},
			},
		}},
		{"queryNoLimit", queryNoLimit, &QueryMetadata{
			params: []*QueryField{
				{Name: "col:`name`", Type: T_ARRAY_STRING},
				{Name: "col:`id`", Type: T_INT},
				{Name: "col:`phone`", Type: T_STRING},
			},
			result: []*QueryField{
				{Name: "`id`", Type: T_PLACEHOLDER},
			},
		}},
		{"queryWithTableJoin", queryWithTableJoin, &QueryMetadata{
			params: []*QueryField{
				{Name: "col:`u`.`name`", Type: T_STRING},
			},
			result: []*QueryField{
				{Name: "`u`.`id`", Type: T_PLACEHOLDER},
				{Name: "`b`.`id`", Type: T_PLACEHOLDER},
			},
		}},
		{"queryWithSubquery", queryWithSubquery, &QueryMetadata{
			params: []*QueryField{
				{Name: "col:`id`", Type: T_INT},
			},
			result: []*QueryField{
				{Name: "`id`", Type: T_PLACEHOLDER},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := NewTiDBParser()
			err := tp.Parse(context.TODO(), tt.query)
			assert.NoError(t, err)
			assert.Equal(t, tt.metadata.String(), tp.Metadata())
		})
	}
}

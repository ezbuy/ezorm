package query

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

const queryIn = `
SELECT
	u.id
FROM
	user AS u
WHERE u.name IN ('me','him')
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

const queryWithColAs = `
SELECT
	u.id as uid
FROM
	user  u
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

const queryLike = `
SELECT
	id
FROM user
WHERE name LIKE 'ezorm%'
`

const queryOrderByDESC = `
SELECT
	id
FROM user
WHERE name LIKE 'ezorm%'
ORDER BY id DESC`

const queryOrderByASC = `
SELECT
	id
FROM user
WHERE name LIKE 'ezorm%'
ORDER BY id ASC`

const queryWithTableJoinOrderBy = `
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
	u.name LIKE 'ezorm%'
ORDER BY u.id DESC
`

func TestTiDBParserParseQuery(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		expect string
	}{
		{"query", query, "SELECT `id` FROM `user` %s"},
		{"queryIn", queryIn, "SELECT `u`.`id` FROM `user` AS `u` %s"},
		{"queryOneLimit", queryOneLimit, "SELECT `id` FROM `user` %s"},
		{"queryNoLimit", queryNoLimit, "SELECT `id` FROM `user` %s"},
		{"queryWithTableJoin", queryWithTableJoin, "SELECT `u`.`id`,`b`.`id` FROM `user` AS `u` JOIN `blog` AS `b` ON `u`.`id`=`b`.`user_id` %s"},
		{"queryWithSubquery", queryWithSubquery, "SELECT `id` FROM `user` %s"},
		{"queryWithLike", queryLike, "SELECT `id` FROM `user` %s"},
		{"queryWithOrderByDESC", queryOrderByDESC, "SELECT `id` FROM `user` %s"},
		{"queryWithOrderByASC", queryOrderByASC, "SELECT `id` FROM `user` %s"},
		{"queryWithTableJoinOrderBy", queryWithTableJoinOrderBy, "SELECT `u`.`id`,`b`.`id` FROM `user` AS `u` JOIN `blog` AS `b` ON `u`.`id`=`b`.`user_id` %s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := NewTiDBParser()
			_, _, err := tp.Parse(context.TODO(), tt.query)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, tp.Query())
		})
	}
}

func TestTiDBParserParseMetadata(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		metadata TableMetadata
	}{
		{
			"query", query, map[Table]*QueryMetadata{
				{Name: "user"}: {
					params: []*QueryField{
						{Name: "col:`name`", Type: T_ARRAY_STRING},
						{Name: "col:`id`", Type: T_INT},
						{Name: "col:`phone`", Type: T_STRING},
						{Name: "limit:offset", Type: T_INT},
						{Name: "limit:count", Type: T_INT},
					},
					result: []*QueryField{
						{Name: "`id`", Type: T_PLACEHOLDER},
					},
				},
			},
		},
		{"queryIn", queryIn, map[Table]*QueryMetadata{
			{Name: "user", Alias: "u"}: {
				params: []*QueryField{
					{Name: "col:`u`.`name`", Type: T_ARRAY_STRING},
				},
				result: []*QueryField{
					{Name: "`u`.`id`", Type: T_PLACEHOLDER},
				},
			},
		}},
		{"queryOneLimit", queryOneLimit, map[Table]*QueryMetadata{
			{Name: "user"}: {
				params: []*QueryField{
					{Name: "col:`name`", Type: T_ARRAY_STRING},
					{Name: "col:`id`", Type: T_INT},
					{Name: "col:`phone`", Type: T_STRING},
					{Name: "limit:count", Type: T_INT},
				},

				result: []*QueryField{
					{Name: "`id`", Type: T_PLACEHOLDER},
				},
			},
		}},
		{"queryNoLimit", queryNoLimit, map[Table]*QueryMetadata{
			{Name: "user"}: {
				params: []*QueryField{
					{Name: "col:`name`", Type: T_ARRAY_STRING},
					{Name: "col:`id`", Type: T_INT},
					{Name: "col:`phone`", Type: T_STRING},
				},
				result: []*QueryField{
					{Name: "`id`", Type: T_PLACEHOLDER},
				},
			},
		}},
		{"queryWithTableJoin", queryWithTableJoin, map[Table]*QueryMetadata{
			{Name: "user", Alias: "u"}: {
				params: []*QueryField{
					{Name: "col:`u`.`name`", Type: T_STRING},
				},
				result: []*QueryField{
					{Name: "`u`.`id`", Type: T_PLACEHOLDER},
				},
			},
			{Name: "blog", Alias: "b"}: {
				result: []*QueryField{
					{Name: "`b`.`id`", Type: T_PLACEHOLDER},
				},
			},
		}},
		{"queryWithSubquery", queryWithSubquery, map[Table]*QueryMetadata{
			{Name: "user"}: {
				params: []*QueryField{
					{Name: "col:`id`", Type: T_INT},
				},
				result: []*QueryField{
					{Name: "`id`", Type: T_PLACEHOLDER},
				},
			},
		}},
		{"queryWithColAs", queryWithColAs, map[Table]*QueryMetadata{
			{Name: "user", Alias: "u"}: {
				params: []*QueryField{
					{Name: "col:`u`.`name`", Type: T_STRING},
				},
				result: []*QueryField{
					{Name: "`u`.`id`", Type: T_PLACEHOLDER, Alias: "uid"},
				},
			},
		}},
		{"queryWithLike", queryLike, map[Table]*QueryMetadata{
			{Name: "user"}: {
				params: []*QueryField{
					{Name: "col:`name`", Type: T_STRING},
				},
				result: []*QueryField{
					{Name: "`id`", Type: T_PLACEHOLDER},
				},
			},
		}},
		{"queryWithOrderByDESC", queryOrderByDESC, map[Table]*QueryMetadata{
			{Name: "user"}: {
				params: []*QueryField{
					{Name: "col:`name`", Type: T_STRING},
					{Name: "orderby-desc:`id`", Type: T_ANY},
				},
				result: []*QueryField{
					{Name: "`id`", Type: T_PLACEHOLDER},
				},
			},
		}},
		{"queryWithOrderByASC", queryOrderByASC, map[Table]*QueryMetadata{
			{Name: "user"}: {
				params: []*QueryField{
					{Name: "col:`name`", Type: T_STRING},
					{Name: "orderby-asc:`id`", Type: T_ANY},
				},
				result: []*QueryField{
					{Name: "`id`", Type: T_PLACEHOLDER},
				},
			},
		}},
		{"queryWithTableJoinOrderBy", queryWithTableJoinOrderBy, map[Table]*QueryMetadata{
			{Name: "user", Alias: "u"}: {
				params: []*QueryField{
					{Name: "col:`u`.`name`", Type: T_STRING},
					{Name: "orderby-desc:`u`.`id`", Type: T_ANY},
				},
				result: []*QueryField{
					{Name: "`u`.`id`", Type: T_PLACEHOLDER},
				},
			},
			{Name: "blog", Alias: "b"}: {
				params: []*QueryField{
					{Name: "orderby-desc:`u`.`id`", Type: T_ANY},
				},
				result: []*QueryField{
					{Name: "`b`.`id`", Type: T_PLACEHOLDER},
				},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := NewTiDBParser()
			_, _, err := tp.Parse(context.TODO(), tt.query)
			assert.NoError(t, err)
			assert.Equal(t, tt.metadata.String(), tp.Metadata())
		})
	}
}

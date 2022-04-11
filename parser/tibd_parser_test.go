package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTiDBParser(t *testing.T) {
	tp := NewTiDBParser()
	err := tp.Parse(context.TODO(), simpleParseSQL)
	assert.NoError(t, err)
	expect := "param: name: `u`.`name`, type: []string\nparam: name: `u`.`id`, type: int64\nparam: name: `u`.`phone`, type: string\nresult: name: `u`.`id`, type: ?\n"
	assert.Equal(t, expect, tp.Metadata())

	query := "SELECT `u`.`id` FROM `user` AS `u` WHERE `u`.`name` IN (?) AND `u`.`id`=? AND `u`.`phone`=? LIMIT ?,?"

	assert.Equal(t, query, tp.Query())
}

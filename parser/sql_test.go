package parser

import (
	"fmt"
	"testing"
)

const parseSQL = `
SELECT
  u.id,
  u.name,
  u.phone,
  u.email,
  ud.desc,
  us.status_code,
  IFNULL(user_status_detail.status_desc, ''),
  IFNULL(usd.status_next, 0) NextStatus
FROM
  user u
JOIN user_detail ud ON u.id=ud.user_id
LEFT JOIN user_status us ON us.user_id=u.id
LEFT JOIN user_status_detail usd ON usd.id=us.detail_id
WHERE u.id=? AND us.status=?
LIMIT ?, ?
`

func TestParseSelect(t *testing.T) {
	stmt, err := ParseSelect(parseSQL)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("====== Fields ======")
	for _, f := range stmt.Fields {
		fmt.Printf("Table = %s, Field = %s, Alias = %s\n",
			f.Table, f.Name, f.Alias)
	}
}

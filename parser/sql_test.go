package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/pingcap/parser/opcode"
	driver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/stretchr/testify/assert"
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
WHERE u.name = "simin"
LIMIT 0,10
`

const simpleParseSQL = `
SELECT
	u.id
FROM
	user u
WHERE u.name IN ('me') AND u.id = 1 AND u.phone = '123'
LIMIT 0,10
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

type condition struct{}

func (c condition) Enter(in ast.Node) (ast.Node, bool) {
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, os.Stdout)
	fmt.Fprintf(os.Stdout, "type: %T\n", in)
	in.Restore(ctx)
	println()
	return in, false
}

func (c condition) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func TestParseSelectRaw(t *testing.T) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(simpleParseSQL, "", "")
	assert.NoError(t, err)
	for _, stmt := range stmtNodes {
		switch x := stmt.(type) {
		case *ast.SelectStmt:
			fmt.Fprintf(os.Stdout, "where type : %T\n", x.Where)
			switch op := x.Where.(type) {
			case *ast.BinaryOperationExpr:
				fmt.Fprintf(os.Stdout, "op type : %v\n", op.Op == opcode.LogicAnd)
				fmt.Fprintf(os.Stdout, "L: op type : %T\n", op.L)
				fmt.Fprintf(os.Stdout, "R: op type : %T\n", op.R)

				ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, os.Stdout)
				op.R.Restore(ctx)
				// assert.Equal(t, op.L.(*ast.ColumnNameExpr).Name.String(), "u.name")
				// assert.Equal(t, op.R.(*driver.ValueExpr).GetType().String(), "var_string(2)")
			case *ast.PatternInExpr:
				assert.Equal(t, op.Expr.(*ast.ColumnNameExpr).Name.String(), "u.name")
				for _, v := range op.List {
					println(v.(*driver.ValueExpr).GetType().String())
				}

			case *ast.PatternLikeExpr:
			case *ast.PatternRegexpExpr:
			}
		}
	}
}

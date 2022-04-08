package parser

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
	"github.com/pingcap/tidb/types"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

type TiDBParser struct {
	meta *QueryMetadata
	b    *QueryBuilder
}

func NewTiDBParser() *TiDBParser {
	return &TiDBParser{
		meta: &QueryMetadata{},
	}
}

func (tp *TiDBParser) Metadata() string {
	return tp.meta.String()
}

func (tp *TiDBParser) parse(node ast.Node, n int) error {
	switch x := node.(type) {
	case *ast.SelectStmt:
		if x.Fields != nil && n == 0 {
			for _, f := range x.Fields.Fields {
				if _, ok := f.Expr.(*ast.ColumnNameExpr); ok {
					field := &QueryField{}
					field.Name = f.Expr.(*ast.ColumnNameExpr).Name.String()
					field.Type = T_PLACEHOLDER
					tp.meta.result = append(tp.meta.result, field)
				}
			}
		}
		if x.Where != nil {
			if err := tp.parse(x.Where, n+1); err != nil {
				return err
			}
		}
	case *ast.BinaryOperationExpr:
		if x.Op == opcode.LogicAnd || x.Op == opcode.LogicOr {
			if err := tp.parse(x.L, n+1); err != nil {
				return err
			}
			if err := tp.parse(x.R, n+1); err != nil {
				return err
			}
		}
		if _, ok := x.L.(*ast.ColumnNameExpr); ok {
			field := &QueryField{
				Name: x.L.(*ast.ColumnNameExpr).Name.String(),
			}
			if _, ook := x.R.(*driver.ValueExpr); ook {
				switch x.R.(*driver.ValueExpr).Kind() {
				case types.KindString:
					field.Type = T_STRING
				case types.KindInt64, types.KindUint64:
					field.Type = T_INT
				default:
					return errors.New("parser: unknown datum type, only support string and int for now")
				}
			}
			tp.meta.params = append(tp.meta.params, field)
		}
	case *ast.PatternInExpr:
		if _, ok := x.Expr.(*ast.ColumnNameExpr); ok {
			field := &QueryField{
				Name: x.Expr.(*ast.ColumnNameExpr).Name.String(),
			}
			if len(x.List) > 0 {
				if _, ok := x.List[0].(*driver.ValueExpr); ok {
					switch x.List[0].(*driver.ValueExpr).Kind() {
					case types.KindString:
						field.Type = T_ARRAY_STRING
					case types.KindInt64, types.KindUint64:
						field.Type = T_ARRAY_INT
					default:
						return errors.New("parser: unknown array datum type, only support string and int for now")
					}
				}
			}
			tp.meta.params = append(tp.meta.params, field)
		}
	case *ast.PatternLikeExpr:
		// TODO
	default:
		return errors.New("parser: unknown node type")
	}
	return nil
}

func (tp *TiDBParser) Parse(ctx context.Context,
	query string) error {
	queries := strings.Split(query, ";")
	for _, q := range queries {
		if err := tp.parseOne(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func (tp *TiDBParser) parseOne(ctx context.Context,
	query string) error {
	node, err := parser.New().ParseOneStmt(query, "", "")
	if err != nil {
		return fmt.Errorf("raw query parser: %w", err)
	}
	return tp.parse(node, 0)
}

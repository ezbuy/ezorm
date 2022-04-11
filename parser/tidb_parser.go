package parser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
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
		b: &QueryBuilder{
			Buffer: bytes.NewBuffer(nil),
			raw: &Raw{
				ins:   map[string]*InBuilder{},
				lo:    map[string]LocationOffset{},
				limit: &LimitOption{},
			},
		},
	}
}

func (tp *TiDBParser) InParams(builders ...*InBuilder) *TiDBParser {
	for _, b := range builders {
		tp.b.raw.ins[b.col] = b
	}
	return tp
}

func (tp *TiDBParser) Metadata() string {
	return tp.meta.String()
}

func (tp *TiDBParser) Query() string {
	return tp.b.rebuild()
}

func (tp *TiDBParser) parse(node ast.Node, n int) error {
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, tp.b)

	buffer := bytes.NewBuffer(nil)
	subCtx := format.NewRestoreCtx(format.DefaultRestoreFlags, buffer)
	if err := node.Restore(subCtx); err != nil {
		return err
	}
	start := strings.Index(tp.b.String(), buffer.String())
	end := start + buffer.Len()

	switch x := node.(type) {
	case *ast.SelectStmt:
		if err := x.Restore(ctx); err != nil {
			return err
		}
		if x.Fields != nil && n == 0 {
			for _, f := range x.Fields.Fields {
				if _, ok := f.Expr.(*ast.ColumnNameExpr); ok {
					field := &QueryField{}
					ff := &strings.Builder{}
					f.Expr.(*ast.ColumnNameExpr).Format(ff)
					field.Name = ff.String()
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
		if x.Limit != nil {
			if _, ok := x.Limit.Count.(*driver.ValueExpr); ok {
				tp.b.raw.limit.count = true
			}
			if _, ok := x.Limit.Offset.(*driver.ValueExpr); ok {
				tp.b.raw.limit.offset = true
			}
			limitBuffer := bytes.NewBuffer(nil)
			subCtx := format.NewRestoreCtx(format.DefaultRestoreFlags, limitBuffer)
			if err := x.Limit.Restore(subCtx); err != nil {
				return err
			}
			start := strings.Index(tp.b.String(), limitBuffer.String())
			end := start + limitBuffer.Len()
			log.Println(limitBuffer.String())
			tp.b.raw.lo["LIMIT"] = LocationOffset{
				// FIXME: solve no well-formated query
				start: start + 5 + 1,
				end:   end,
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
		if l, ok := x.L.(*ast.ColumnNameExpr); ok {
			nameBuilder := &strings.Builder{}
			l.Format(nameBuilder)
			field := &QueryField{
				Name: nameBuilder.String(),
			}
			tp.b.raw.lo[field.Name] = LocationOffset{
				start: start + len(field.Name) + 1,
				end:   end,
			}
			if v, ook := x.R.(*driver.ValueExpr); ook {
				switch v.Kind() {
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
		if expr, ok := x.Expr.(*ast.ColumnNameExpr); ok {
			nameBuilder := &strings.Builder{}
			expr.Format(nameBuilder)
			field := &QueryField{
				Name: nameBuilder.String(),
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
				tp.b.raw.lo[field.Name] = LocationOffset{
					// FIXME: solve no well-formated query
					start: start + len(field.Name) + 5,
					end:   end - 1,
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

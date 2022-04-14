package parser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/pingcap/parser/opcode"
	"github.com/pingcap/tidb/types"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

type TiDBParser struct {
	meta TableMetadata
	b    *QueryBuilder
}

func NewTiDBParser() *TiDBParser {
	return &TiDBParser{
		meta: make(map[Table]*QueryMetadata),
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
	case *ast.Join:
		if x.Left != nil {
			tp.parse(x.Left, n+1)
		}
		if x.Right != nil {
			tp.parse(x.Right, n+1)
		}
	case *ast.TableSource:
		var alias string
		if x.AsName.String() != "" {
			alias = x.AsName.String()
		}
		if table, ok := x.Source.(*ast.TableName); ok {
			tb := Table{
				Alias: alias,
				Name:  table.Name.String(),
			}
			if _, ok := tp.meta[tb]; !ok {
				tp.meta[tb] = &QueryMetadata{
					table: tb.Name,
				}
			}
		}
	// SELECT
	case *ast.SelectStmt:
		if x.From != nil {
			// FROM
			tp.parse(x.From.TableRefs, n+1)
		}
		if n == 0 {
			if err := x.Restore(ctx); err != nil {
				return err
			}
			// Field Ref
			if x.Fields != nil {
				for _, f := range x.Fields.Fields {
					if expr, ok := f.Expr.(*ast.ColumnNameExpr); ok {
						field := &QueryField{}
						ff := &strings.Builder{}
						expr.Format(ff)
						field.Name = ff.String()
						field.Type = T_PLACEHOLDER
						tp.meta.AppendResult(expr.Name.Table.String(), field)
					}
				}
			}
		}
		// WHERE
		if x.Where != nil {
			if err := tp.parse(x.Where, n+1); err != nil {
				return err
			}
		}
		// LIMIT
		if x.Limit != nil {
			if _, ok := x.Limit.Count.(*driver.ValueExpr); ok {
				tp.b.raw.limit.count = true
				// FIXME
				for t := range tp.meta {
					tp.meta[t].params = append(tp.meta[t].params, &QueryField{
						Name: "limit:count",
						Type: T_INT,
					})
				}
			}
			if _, ok := x.Limit.Offset.(*driver.ValueExpr); ok {
				tp.b.raw.limit.offset = true
				// FIXME
				for t := range tp.meta {
					tp.meta[t].params = append(tp.meta[t].params, &QueryField{
						Name: "limit:offset",
						Type: T_INT,
					})
				}
			}
			limitBuffer := bytes.NewBuffer(nil)
			subCtx := format.NewRestoreCtx(format.DefaultRestoreFlags, limitBuffer)
			if err := x.Limit.Restore(subCtx); err != nil {
				return err
			}
			start := strings.Index(tp.b.String(), limitBuffer.String())
			end := start + limitBuffer.Len()
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
			field.Name = fmt.Sprintf("col:%s", field.Name)
			t := l.Name.Table.String()
			tp.meta.AppendParams(t, field)
		}
	case *ast.PatternInExpr:
		switch {
		case x.Sel != nil:
			if err := tp.parse(x.Sel, n+1); err != nil {
				return err
			}
		case len(x.List) > 0:
			if expr, ok := x.Expr.(*ast.ColumnNameExpr); ok {
				nameBuilder := &strings.Builder{}
				expr.Format(nameBuilder)
				field := &QueryField{
					Name: nameBuilder.String(),
				}
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

				tp.b.raw.ins[field.Name] = NewIn(field.Name, len(x.List))
				tp.b.raw.lo[field.Name] = LocationOffset{
					// FIXME: solve no well-formated query
					start: start + len(field.Name) + 5,
					end:   end - 1,
				}
				field.Name = fmt.Sprintf("col:%s", field.Name)
				t := expr.Name.Table.String()
				tp.meta.AppendParams(t, field)
			}

		}
	case *ast.PatternLikeExpr:
		// TODO
	case *ast.SubqueryExpr:
		if err := tp.parse(x.Query, n+1); err != nil {
			return err
		}
	default:
		return fmt.Errorf("parser: unknown node type; %T", node)
	}
	return nil
}

func (tp *TiDBParser) Parse(ctx context.Context,
	query string) (TableMetadata, *QueryBuilder, error) {
	queries := strings.Split(query, ";")
	for _, q := range queries {
		if err := tp.parseOne(ctx, q); err != nil {
			return nil, nil, err
		}
	}

	return tp.meta, tp.b, nil
}

func (tp *TiDBParser) Flush() {
	tp.b.raw = &Raw{
		lo:    make(map[string]LocationOffset),
		ins:   map[string]*InBuilder{},
		limit: &LimitOption{},
	}
	tp.b.Reset()
	tp.meta = make(map[Table]*QueryMetadata)
}

func (tp *TiDBParser) parseOne(ctx context.Context,
	query string) error {
	node, err := parser.New().ParseOneStmt(query, "", "")
	if err != nil {
		return fmt.Errorf("raw query parser: %w", err)
	}
	return tp.parse(node, 0)
}

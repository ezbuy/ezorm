package query

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
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

const (
	LIMIT_COUNT  = "limit:count"
	LIMIT_OFFSET = "limit:offset"
)

func NewTiDBParser() *TiDBParser {
	return &TiDBParser{
		meta: make(map[Table]*QueryMetadata),
		b: &QueryBuilder{
			Buffer: bytes.NewBuffer(nil),
			raw: &Raw{
				ins: map[string]struct{}{},
			},
		},
	}
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
					if expr, ok := f.Expr.(*ast.AggregateFuncExpr); ok {
						field := &QueryField{
							Alias: f.AsName.String(),
						}
						var txt bytes.Buffer
						txt.WriteString(expr.F)
						for _, args := range expr.Args {
							txt.WriteString("_")
							var arg strings.Builder
							args.Format(&arg)
							txt.WriteString(arg.String())
						}
						field.Name = txt.String()
						field.Type = T_ANY
						if len(expr.Args) > 0 {
							if col, ok := expr.Args[0].(*ast.ColumnNameExpr); ok {
								tp.meta.AppendResult(col.Name.Table.String(), field)
								tp.b.resultFields = append(tp.b.resultFields, field)
							}
						}
					}
					if expr, ok := f.Expr.(*ast.FuncCallExpr); ok {
						field := &QueryField{
							Alias: f.AsName.String(),
						}
						var txt bytes.Buffer
						txt.WriteString(expr.FnName.O)
						for _, args := range expr.Args {
							txt.WriteString("_")
							var arg strings.Builder
							args.Format(&arg)
							txt.WriteString(arg.String())
						}
						field.Name = txt.String()
						field.Type = T_ANY
						if len(expr.Args) > 0 {
							for _, arg := range expr.Args {
								if col, ok := arg.(*ast.ColumnNameExpr); ok {
									tp.meta.AppendResult(col.Name.Table.String(), field)
									tp.b.resultFields = append(tp.b.resultFields, field)
								}
							}
						}
					}
					if expr, ok := f.Expr.(*ast.ColumnNameExpr); ok {
						field := &QueryField{
							Alias: f.AsName.String(),
						}
						ff := &strings.Builder{}
						expr.Format(ff)
						field.Name = ff.String()
						field.Type = T_PLACEHOLDER
						tp.meta.AppendResult(expr.Name.Table.String(), field)
						tp.b.resultFields = append(tp.b.resultFields, field)
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
			if _, ok := x.Limit.Offset.(*driver.ValueExpr); ok {
				for t := range tp.meta {
					tp.meta[t].params = append(tp.meta[t].params, &QueryField{
						Name: LIMIT_OFFSET,
						Type: T_INT,
					})
				}
			}
			if _, ok := x.Limit.Count.(*driver.ValueExpr); ok {
				for t := range tp.meta {
					tp.meta[t].params = append(tp.meta[t].params, &QueryField{
						Name: LIMIT_COUNT,
						Type: T_INT,
					})
				}
			}
			limitBuffer := bytes.NewBuffer(nil)
			subCtx := format.NewRestoreCtx(format.DefaultRestoreFlags, limitBuffer)
			if err := x.Limit.Restore(subCtx); err != nil {
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
		if l, ok := x.L.(*ast.ColumnNameExpr); ok {
			nameBuilder := &strings.Builder{}
			l.Format(nameBuilder)
			field := &QueryField{
				Name: nameBuilder.String(),
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
				tp.b.raw.ins[field.Name] = struct{}{}
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
	query string,
) (TableMetadata, *QueryBuilder, error) {
	queries := strings.Split(query, ";")
	for _, q := range queries {
		if len(strings.TrimSpace(q)) == 0 {
			continue
		}
		if err := tp.parseOne(ctx, q); err != nil {
			return nil, nil, err
		}
	}

	for n, meta := range tp.meta {
		sort.SliceStable(meta.params, func(i, j int) bool {
			return meta.params[i].Name < meta.params[j].Name
		})
		tp.meta[n] = meta
	}

	return tp.meta, tp.b, nil
}

func (tp *TiDBParser) Flush() {
	tp.b.raw = &Raw{
		ins: map[string]struct{}{},
	}
	tp.b.resultFields = []*QueryField{}
	tp.b.Reset()
	tp.meta = make(map[Table]*QueryMetadata)
}

func (tp *TiDBParser) parseOne(ctx context.Context,
	query string,
) error {
	node, err := parser.New().ParseOneStmt(query, "", "")
	if err != nil {
		return fmt.Errorf("raw query parser: %w(query: %s)", err, query)
	}
	return tp.parse(node, 0)
}

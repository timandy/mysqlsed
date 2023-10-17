package sqlutil

import (
	"errors"
	"fmt"
	"github.com/xwb1989/sqlparser"
)

// ResolveInsertStatement 解析插入语句, 返回 表名, []列, []值
func ResolveInsertStatement(result string) (tableName sqlparser.TableName, columns []sqlparser.ColIdent, values []sqlparser.Expr, err error) {
	statement, err := sqlparser.Parse(result)
	if err != nil {
		errMsg := fmt.Sprintf("Error parse sql:%v %v", result, err.Error())
		return sqlparser.TableName{}, nil, nil, errors.New(errMsg)
	}
	insertStatement := statement.(*sqlparser.Insert)
	insertStatementValues := insertStatement.Rows.(sqlparser.Values)[0]
	return insertStatement.Table, insertStatement.Columns, insertStatementValues, nil
}

// PrebuildUpdateStatement 预构建更新语句, 返回 set col1=val1, col2=val2 和 id 的值
func PrebuildUpdateStatement(line string, columns []sqlparser.ColIdent, values []sqlparser.Expr) (updateExprs []*sqlparser.UpdateExpr, idValue sqlparser.Expr, err error) {
	idColumnIdx := -1
	for idx, column := range columns {
		col := column
		if col.String() == "id" {
			idColumnIdx = idx
			continue
		}
		updateExpr := &sqlparser.UpdateExpr{
			Name: &sqlparser.ColName{
				Metadata:  nil,
				Name:      col,
				Qualifier: sqlparser.TableName{},
			},
			Expr: values[idx],
		}
		updateExprs = append(updateExprs, updateExpr)
	}
	if idColumnIdx == -1 {
		msg := fmt.Sprintf("sql:%v 没有 id 列", line)
		return nil, nil, errors.New(msg)
	}
	idValue = values[idColumnIdx]
	return updateExprs, idValue, nil
}

// BuildUpdateStatement 构建更新语句
func BuildUpdateStatement(tableName sqlparser.TableName, updateExprs []*sqlparser.UpdateExpr, idValue sqlparser.Expr) *sqlparser.Update {
	updateStatement := &sqlparser.Update{
		Comments: nil,
		TableExprs: []sqlparser.TableExpr{&sqlparser.AliasedTableExpr{
			Expr:       tableName,
			Partitions: nil,
			As:         sqlparser.TableIdent{},
			Hints:      nil,
		}},
		Exprs: updateExprs,
		Where: sqlparser.NewWhere("where", &sqlparser.ComparisonExpr{
			Operator: "=",
			Left: &sqlparser.ColName{
				Metadata:  nil,
				Name:      sqlparser.NewColIdent("id"),
				Qualifier: sqlparser.TableName{},
			},
			Right:  idValue,
			Escape: nil,
		}),
		OrderBy: nil,
		Limit:   nil,
	}
	return updateStatement
}

// SerializeToSQL 更新语句 to SQL
func SerializeToSQL(updateStatement *sqlparser.Update) []byte {
	buffer := sqlparser.NewTrackedBuffer(formatter)
	updateStatement.Format(buffer)
	updateSql := buffer.Bytes()
	return updateSql
}

// 表名, 列名 添加转义符
func formatter(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	switch raw := node.(type) {
	case sqlparser.TableIdent:
		buf.WriteString("`")
		buf.WriteString(raw.String())
		buf.WriteString("`")
	case sqlparser.ColIdent:
		buf.WriteString("`")
		buf.WriteString(raw.String())
		buf.WriteString("`")
	default:
		node.Format(buf)
	}
}

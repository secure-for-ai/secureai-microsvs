//go:build debug
// +build debug

package sqlBuilderV3

import (
	"github.com/secure-for-ai/secureai-microsvs/db"
	"sort"
)

func (m Map) sortedKeys() []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (stmt *Stmt) valuesOneMap(data Map) {
	insertValues := getValExprListWithSize(len(data)) //make([]condExpr, len(data))
	if len(stmt.InsertCols) == 0 {
		for i, col := range data.sortedKeys() {
			stmt.InsertCols = append(stmt.InsertCols, col)
			val := data[col]
			if e, ok := val.(*condExpr); ok {
				insertValues.SetIthWithExpr(i, e)
			} else {
				insertValues.SetIth(i, db.Para, val) //[i].Set(db.Para, val)
			}
		}
	} else {
		for i, col := range data.sortedKeys() {
			val := data[col]
			if e, ok := val.(*condExpr); ok {
				insertValues.SetIthWithExpr(i, e)
			} else {
				insertValues.SetIth(i, db.Para, val) //[i].Set(db.Para, val)
			}
		}
	}

	stmt.InsertValues.append(insertValues)
}

func (stmt *Stmt) buildInsertColsByMap(data Map) {
	stmt.InsertCols = append(stmt.InsertCols, data.sortedKeys()...)
}

func catCondMap(ref *[]Cond, query Map, conds *condAnd) {
	for _, k := range query.sortedKeys() {
		cond := Expr(k, query[k])
		cond.appendSql(" = ")
		cond.appendSql(db.Para)
		// self created cond is stored in the ref
		*ref = append(*ref, cond)
		*conds = append(*conds, cond)
	}
}

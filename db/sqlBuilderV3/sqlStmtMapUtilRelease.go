//go:build !debug
// +build !debug

package sqlBuilderV3

import(
	"github.com/secure-for-ai/secureai-microsvs/db"
)

func (stmt *Stmt) valuesOneMap(data Map) {
	insertValues := getCondExprListWithSize(len(data)) //make([]condExpr, len(curData))
	if len(stmt.InsertCols) == 0 {
		i := 0
		for col, val := range data {
			stmt.InsertCols = append(stmt.InsertCols, col)
			// val := curData[col]
			if e, ok := val.(*condExpr); ok {
				insertValues.SetIthWithExpr(i, e)
			} else {
				insertValues.SetIth(i, db.Para, val) //[i].Set(db.Para, val)
			}
			i++
		}
	} else {
		i := 0
		for _, val := range data {
			// val := curData[col]
			if e, ok := val.(*condExpr); ok {
				insertValues.SetIthWithExpr(i, e)
			} else {
				insertValues.SetIth(i, db.Para, val) //[i].Set(db.Para, val)
			}
			i++
		}
	}

	stmt.InsertValues.append(insertValues)
}

func (stmt *Stmt) buildInsertColsByMap(data Map) {
	for key := range data {
		stmt.InsertCols = append(stmt.InsertCols, key)
	}
}
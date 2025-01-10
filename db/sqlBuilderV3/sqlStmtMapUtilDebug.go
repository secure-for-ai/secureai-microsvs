//go:build debug
// +build debug

package sqlBuilderV3

import(
	"github.com/secure-for-ai/secureai-microsvs/db"
)

// func (m Map) sortedKeys() []string {
// 	keys := make([]string, 0, len(m))
// 	for key := range m {
// 		keys = append(keys, key)
// 	}
// 	sort.Strings(keys)
// 	return keys
// }

func (stmt *Stmt) valuesOneMap(data Map) {
	insertValues := getCondExprListWithSize(len(data)) //make([]condExpr, len(data))
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
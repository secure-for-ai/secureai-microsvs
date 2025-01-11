package sqlBuilderV3

import (
	"sync"
)

type valExpr struct {
	sql  string
	args []interface{}
}

func (expr *valExpr) Set(sql string, args ...interface{}) {
	expr.sql = sql
	expr.args = append(expr.args[:0], args...)
}

func (expr *valExpr) String() string {
	return expr.sql
}

type valExprList []valExpr

var valExprListPool = sync.Pool{
	New: func() interface{} {
		exprList := make(valExprList, 0, 4)
		return &exprList
	},
}

func getValExprListWithSize(n int) *valExprList {
	exprList := valExprListPool.Get().(*valExprList)
	if n <= cap(*exprList) {
		*exprList = (*exprList)[:n]
	} else {
		c := n
		if c < 2*cap(*exprList) {
			c = 2 * cap(*exprList)
		}
		exprList2 := append(valExprList(nil), make(valExprList, c)...)
		copy(exprList2, *exprList)
		*exprList = exprList2[:n]
	}
	return exprList
}

func (exprList *valExprList) Destroy() {
	valExprListPool.Put(exprList)
}

func (exprList *valExprList) SetIth(i int, sql string, args ...interface{}) {
	(*exprList)[i].Set(sql, args...)
}

func (exprList *valExprList) SetIthWithExpr(i int, expr *condExpr) {
	// we deep-copy the memory rather than do (*exprList)[i] = *expr.
	// If we do so, (*exprList)[i].args is indeed expr.args. Then, when we
	// return *exprList to the sync pool, the slice (*exprList)[i].args may be
	// still used by other objects. This can mess up the memory management and
	// cause segmentation faults.
	(*exprList)[i].Set(expr.String(), expr.args...)
}

type valExpr2DList []*valExprList

func newValExpr2DList(n int) valExpr2DList {
	return make(valExpr2DList, 0, n)
}

// Ensure that there are at least n entries left in the list: cap(*list) - len(*list) >= n
func (list *valExpr2DList) grow(n int) {
	c := len(*list) + n

	// grow the cap
	if c > cap(*list) {
		if c < 2*cap(*list) {
			c = 2 * cap(*list)
		}
		list2 := append(valExpr2DList(nil), make(valExpr2DList, c)...)
		copy(list2, *list)
		*list = list2[:len(*list)]
	}
}

func (list *valExpr2DList) reset() {
	for i, il := 0, len(*list); i < il; i++ {
		(*list)[i].Destroy()
	}
	// for _, val := range *list {
	// 	val.Destroy()
	// }
	*list = (*list)[:0]
}

func (list *valExpr2DList) append(val *valExprList) {
	*list = append(*list, val)
}

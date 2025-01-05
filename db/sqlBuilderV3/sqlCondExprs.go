package sqlBuilderV3

import (
	"sync"
)

type condExprList []condExpr

var condExprListPool = sync.Pool{
	New: func() interface{} {
		exprList := make(condExprList, 0, 4)
		return &exprList
	},
}

func getCondExprListWithSize(n int) *condExprList {
	exprList := condExprListPool.Get().(*condExprList)
	if n <= cap(*exprList) {
		*exprList = (*exprList)[:n]
	} else {
		c := n
		if c < 2*cap(*exprList) {
			c = 2 * cap(*exprList)
		}
		exprList2 := append(condExprList(nil), make(condExprList, c)...)
		copy(exprList2, *exprList)
		*exprList = exprList2[:n]
	}
	return exprList
}

func (exprList *condExprList) Destroy() {
	// *exprList = (*exprList)[:0]
	condExprListPool.Put(exprList)
}

func (exprList *condExprList) SetIth(i int, sql string, args ...interface{}) {
	(*exprList)[i].Set(sql, args...)
}

func (exprList *condExprList) SetIthWithExpr(i int, expr *condExpr) {
	(*exprList)[i] = *expr
}

type condExpr2DList []*condExprList

func newCondExpr2DList(n int) condExpr2DList {
	return make(condExpr2DList, 0, n)
}

// Ensure that there are at least n entries left in the list: cap(*list) - len(*list) >= n
func (list *condExpr2DList) grow(n int) {
	c := len(*list) + n

	// grow the cap
	if c > cap(*list) {
		if c < 2*cap(*list) {
			c = 2 * cap(*list)
		}
		list2 := append(condExpr2DList(nil), make(condExpr2DList, c)...)
		copy(list2, *list)
		*list = list2[:len(*list)]
	}
}

func (list *condExpr2DList) reset() {
	for i, il := 0, len(*list); i < il; i++ {
		(*list)[i].Destroy()
	}
	// for _, val := range *list {
	// 	val.Destroy()
	// }
	*list = (*list)[:0]
}

func (list *condExpr2DList) append(val *condExprList) {
	*list = append(*list, val)
}

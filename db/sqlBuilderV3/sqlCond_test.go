package sqlBuilderV3_test

import (
	"testing"

	"github.com/secure-for-ai/secureai-microsvs/db/sqlBuilderV3"
	"github.com/stretchr/testify/assert"
)

var (
	condNull = sqlBuilderV3.Expr("")
	cond1    = sqlBuilderV3.Expr("A < ?", 1)
	cond2    = sqlBuilderV3.Expr("B = ?", "hello")
	cond3    = sqlBuilderV3.Expr("C LIKE ?", "username")
	w        = sqlBuilderV3.NewWriter()
)

func TestExpr(t *testing.T) {
	sql1, args, err := sqlBuilderV3.CondToSQL(cond1, w)
	assert.NoError(t, err)
	assert.EqualValues(t, "A < ?", sql1)
	assert.EqualValues(t, []interface{}{1}, args)
}

func TestExpr_And(t *testing.T) {
	cond1And2 := cond1.And(cond2)
	sql2, args, err := sqlBuilderV3.CondToSQL(cond1And2, w)
	assert.NoError(t, err)
	assert.EqualValues(t, "(A < ?) AND (B = ?)", sql2)
	assert.EqualValues(t, []interface{}{1, "hello"}, args)

	//test And with empty cond
	tmpCond := cond1.And(condNull)
	sql2, args, err = sqlBuilderV3.CondToSQL(tmpCond, w)
	assert.NoError(t, err)
	assert.EqualValues(t, cond1, tmpCond)
	assert.EqualValues(t, "A < ?", sql2)
	assert.EqualValues(t, []interface{}{1}, args)

	//test nest And
	cond1And2And3And1And2 := cond1And2.And(cond3, condNull, condNull, cond1And2)
	sql2, args, err = sqlBuilderV3.CondToSQL(cond1And2And3And1And2, w)
	assert.NoError(t, err)
	assert.EqualValues(t, "(A < ?) AND (B = ?) AND (C LIKE ?) AND (A < ?) AND (B = ?)", sql2)
	assert.EqualValues(t, []interface{}{1, "hello", "username", 1, "hello"}, args)
}

func TestExpr_Or(t *testing.T) {
	cond1Or2 := cond1.Or(cond2)
	sql2, args, err := sqlBuilderV3.CondToSQL(cond1Or2, w)
	assert.NoError(t, err)
	assert.EqualValues(t, "(A < ?) OR (B = ?)", sql2)
	assert.EqualValues(t, []interface{}{1, "hello"}, args)

	//test Or with empty cond
	tmpCond := cond1.Or(condNull)
	sql2, args, err = sqlBuilderV3.CondToSQL(tmpCond, w)
	assert.NoError(t, err)
	assert.EqualValues(t, cond1, tmpCond)
	assert.EqualValues(t, "A < ?", sql2)
	assert.EqualValues(t, []interface{}{1}, args)

	//test nest And
	cond1Or2Or3Or1Or2 := cond1Or2.Or(cond3, condNull, condNull, cond1Or2)
	sql2, args, err = sqlBuilderV3.CondToSQL(cond1Or2Or3Or1Or2, w)
	assert.NoError(t, err)
	assert.EqualValues(t, "(A < ?) OR (B = ?) OR (C LIKE ?) OR (A < ?) OR (B = ?)", sql2)
	assert.EqualValues(t, []interface{}{1, "hello", "username", 1, "hello"}, args)
}

func TestAnd(t *testing.T) {
	assert.EqualValues(t, condNull, sqlBuilderV3.NewCond().And().And(condNull).And(condNull, condNull))
	assert.EqualValues(t, condNull, sqlBuilderV3.And().And(condNull).And(condNull, condNull))
	assert.EqualValues(t, cond1, sqlBuilderV3.And(cond1))
	assert.EqualValues(t, cond1, sqlBuilderV3.And(condNull, cond1))
	assert.EqualValues(t, cond1, sqlBuilderV3.And(cond1, condNull))
	assert.EqualValues(t, cond1, sqlBuilderV3.And(condNull, cond1, condNull))
	assert.EqualValues(t, cond1, condNull.And(condNull, cond1, condNull).And(condNull, condNull))
}

func TestOr(t *testing.T) {
	assert.EqualValues(t, condNull, sqlBuilderV3.NewCond().Or().Or(condNull).Or(condNull, condNull))
	assert.EqualValues(t, condNull, sqlBuilderV3.Or().Or(condNull).Or(condNull, condNull))
	assert.EqualValues(t, cond1, sqlBuilderV3.Or(cond1))
	assert.EqualValues(t, cond1, sqlBuilderV3.Or(condNull, cond1))
	assert.EqualValues(t, cond1, sqlBuilderV3.Or(cond1, condNull))
	assert.EqualValues(t, cond1, sqlBuilderV3.Or(condNull, cond1, condNull))
	assert.EqualValues(t, cond1, condNull.Or(condNull, cond1, condNull).Or(condNull, condNull))
}

func TestNewCond(t *testing.T) {
	assert.EqualValues(t, condNull, sqlBuilderV3.NewCond().Or())
}

func TestComplex(t *testing.T) {
	assert.EqualValues(t, condNull,
		sqlBuilderV3.NewCond().And().And(condNull).And(condNull, condNull).
			Or().Or(condNull).Or(condNull, condNull),
	)
	assert.EqualValues(t, condNull,
		sqlBuilderV3.And().And(condNull).And(condNull, condNull).
			Or().Or(condNull).Or(condNull, condNull),
	)
	assert.EqualValues(t, condNull,
		condNull.
			And().And(condNull).And(condNull, condNull).
			Or().Or(condNull).Or(condNull, condNull),
	)

	cond1And2 := cond1.And(cond2)
	cond1Or3 := cond1.Or(cond3)

	tmpCond := cond1And2.Or(cond1And2).And(cond1Or3)
	sql, args, err := sqlBuilderV3.CondToSQL(tmpCond, w)
	assert.NoError(t, err)
	assert.EqualValues(t, "(((A < ?) AND (B = ?)) OR ((A < ?) AND (B = ?))) AND ((A < ?) OR (C LIKE ?))", sql)
	assert.EqualValues(t, []interface{}{1, "hello", 1, "hello", 1, "username"}, args)
}

func BenchmarkExpr(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		w := sqlBuilderV3.NewWriter()
		for pb.Next() {
			sqlBuilderV3.CondToSQL(cond1, w)
		}
	})
}

func BenchmarkExpr_And(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		w := sqlBuilderV3.NewWriter()
		// dd := make([]sqlBuilderV3.Cond, 3, 4)
		for pb.Next() {
			cond1And2 := sqlBuilderV3.And(cond1, cond2)//cond1.And(cond2)
			// cond1And2 := sqlBuilderV3.AndOne(cond1, cond2)
			// cond1And2 := cond1.And(cond2)
			sqlBuilderV3.CondToSQL(cond1And2, w)

			// //test And with empty cond
			// tmpCond := cond1.And(condNull)
			tmpCond := sqlBuilderV3.And(cond1, condNull)
			sqlBuilderV3.CondToSQL(tmpCond, w)

			// //test nest And
			// cond1And2And3And1And2 := cond1And2.And(cond3, condNull, condNull, cond1And2)
			cond1And2And3And1And2 := sqlBuilderV3.And(cond1And2, cond3, condNull, condNull, cond1And2)
			sqlBuilderV3.CondToSQL(cond1And2And3And1And2, w)

			cond1And2.Destroy()
			// tmpCond.Destroy()
			cond1And2And3And1And2.Destroy()
			// cond1And2.Reset()
			// cond1And2 = cond1And2.And(cond1, cond2, cond2, cond2)
			// dd[0] = cond2
			// dd[1] = cond2
			// dd[2] = cond2
			// cond1And2 := cond1.And(dd...)
			// dd[0] = cond1
			// dd[1] = cond2
			// dd[2] = cond2
			// dd[3] = cond2
			// dd = append(dd, cond1, cond2, cond2, cond2)
			// cond1And2 := condNull.And(dd...)
			// dd = dd[:0]
			// cond1And2 := condNull.And(cond1, cond2, cond2, cond2)
			// cond1And2 := sqlBuilderV3.And(condNull, cond1, cond2, cond2, cond2)
			// // cond1And2 := sqlBuilderV3.And(cond1, cond2, cond2, cond2)
			// // cond1And2 := sqlBuilderV3.And3(cond1, cond2)
			// sqlBuilderV3.CondToSQL(cond1And2, w)
		}
	})
}

func BenchmarkExpr_Or(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		w := sqlBuilderV3.NewWriter()
		// dd := make([]sqlBuilderV3.Cond, 3, 4)
		for pb.Next() {
			// cond1Or2 := cond1.Or(cond2)
			cond1Or2 := sqlBuilderV3.Or(cond1, cond2)
			sqlBuilderV3.CondToSQL(cond1Or2, w)

			//test Or with empty cond
			// tmpCond := cond1.Or(condNull)
			tmpCond := sqlBuilderV3.Or(cond1, condNull)
			sqlBuilderV3.CondToSQL(tmpCond, w)

			//test nest And
			// cond1Or2Or3Or1Or2 := cond1Or2.Or(cond3, condNull, condNull, cond1Or2)
			cond1Or2Or3Or1Or2 := sqlBuilderV3.Or(cond1Or2, cond3, condNull, condNull, cond1Or2)
			sqlBuilderV3.CondToSQL(cond1Or2Or3Or1Or2, w)

			cond1Or2.Destroy()
			cond1Or2Or3Or1Or2.Destroy()
		}
	})
}
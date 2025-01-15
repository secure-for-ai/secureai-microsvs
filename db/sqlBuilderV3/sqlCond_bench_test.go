package sqlBuilderV3_test

import (
	"testing"

	"github.com/secure-for-ai/secureai-microsvs/db/sqlBuilderV3"
	"github.com/stretchr/testify/assert"
)

func BenchmarkExpr(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			sqlBuilderV3.CondToSQL(cond1, w)
			w.Destroy()
		}
	})
}

func BenchmarkExprAnd(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			cond1And2 := sqlBuilderV3.And(cond1, cond2)
			sqlBuilderV3.CondToSQL(cond1And2, w)

			tmpCond := sqlBuilderV3.And(cond1, condNull)
			sqlBuilderV3.CondToSQL(tmpCond, w)

			cond1And2And3And1And2 := sqlBuilderV3.And(cond1And2, cond3, condNull, condNull, cond1And2)
			sqlBuilderV3.CondToSQL(cond1And2And3And1And2, w)

			cond1And2.Destroy()
			cond1And2And3And1And2.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkExprAndAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			cond1And2 := sqlBuilderV3.And(cond1, cond2)
			sql2, args, err := sqlBuilderV3.CondToSQL(cond1And2, w)
			assert.NoError(b, err)
			assert.EqualValues(b, "(A < ?) AND (B = ?)", sql2)
			assert.EqualValues(b, []any{1, "hello"}, args)

			tmpCond := sqlBuilderV3.And(cond1, condNull)
			sql2, args, err = sqlBuilderV3.CondToSQL(tmpCond, w)
			assert.NoError(b, err)
			assert.EqualValues(b, cond1, tmpCond)
			assert.EqualValues(b, "A < ?", sql2)
			assert.EqualValues(b, []any{1}, args)

			cond1And2And3And1And2 := sqlBuilderV3.And(cond1And2, cond3, condNull, condNull, cond1And2)
			sql2, args, err = sqlBuilderV3.CondToSQL(cond1And2And3And1And2, w)
			assert.NoError(b, err)
			assert.EqualValues(b, "(A < ?) AND (B = ?) AND (C LIKE ?) AND (A < ?) AND (B = ?)", sql2)
			assert.EqualValues(b, []any{1, "hello", "username", 1, "hello"}, args)

			cond1And2.Destroy()
			cond1And2And3And1And2.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkExprOr(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			cond1Or2 := sqlBuilderV3.Or(cond1, cond2)
			sqlBuilderV3.CondToSQL(cond1Or2, w)

			tmpCond := sqlBuilderV3.Or(cond1, condNull)
			sqlBuilderV3.CondToSQL(tmpCond, w)

			cond1Or2Or3Or1Or2 := sqlBuilderV3.Or(cond1Or2, cond3, condNull, condNull, cond1Or2)
			sqlBuilderV3.CondToSQL(cond1Or2Or3Or1Or2, w)

			cond1Or2.Destroy()
			cond1Or2Or3Or1Or2.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkExprOrAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			cond1Or2 := sqlBuilderV3.Or(cond1, cond2)
			sql2, args, err := sqlBuilderV3.CondToSQL(cond1Or2, w)
			assert.NoError(b, err)
			assert.EqualValues(b, "(A < ?) OR (B = ?)", sql2)
			assert.EqualValues(b, []any{1, "hello"}, args)

			tmpCond := sqlBuilderV3.Or(cond1, condNull)
			sql2, args, err = sqlBuilderV3.CondToSQL(tmpCond, w)
			assert.NoError(b, err)
			assert.EqualValues(b, cond1, tmpCond)
			assert.EqualValues(b, "A < ?", sql2)
			assert.EqualValues(b, []any{1}, args)

			cond1Or2Or3Or1Or2 := sqlBuilderV3.Or(cond1Or2, cond3, condNull, condNull, cond1Or2)
			sql2, args, err = sqlBuilderV3.CondToSQL(cond1Or2Or3Or1Or2, w)
			assert.NoError(b, err)
			assert.EqualValues(b, "(A < ?) OR (B = ?) OR (C LIKE ?) OR (A < ?) OR (B = ?)", sql2)
			assert.EqualValues(b, []any{1, "hello", "username", 1, "hello"}, args)

			cond1Or2.Destroy()
			cond1Or2Or3Or1Or2.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkCondComplex(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			cond1And2 := sqlBuilderV3.And(cond1, cond2)
			cond1Or3 := sqlBuilderV3.Or(cond1, cond3)
			cond1And2Orcond1And2 := sqlBuilderV3.Or(cond1And2, cond1And2)
			condFin := sqlBuilderV3.And(cond1And2Orcond1And2, cond1Or3)

			sqlBuilderV3.CondToSQL(condFin, w)

			cond1And2.Destroy()
			cond1Or3.Destroy()
			cond1And2Orcond1And2.Destroy()
			condFin.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkCondComplexAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			cond1And2 := sqlBuilderV3.And(cond1, cond2)
			cond1Or3 := sqlBuilderV3.Or(cond1, cond3)
			cond1And2Orcond1And2 := sqlBuilderV3.Or(cond1And2, cond1And2)
			condFin := sqlBuilderV3.And(cond1And2Orcond1And2, cond1Or3)

			sql, args, err := sqlBuilderV3.CondToSQL(condFin, w)
			assert.NoError(b, err)
			assert.EqualValues(b, "(((A < ?) AND (B = ?)) OR ((A < ?) AND (B = ?))) AND ((A < ?) OR (C LIKE ?))", sql)
			assert.EqualValues(b, []any{1, "hello", 1, "hello", 1, "username"}, args)

			cond1And2.Destroy()
			cond1Or3.Destroy()
			cond1And2Orcond1And2.Destroy()
			condFin.Destroy()
			w.Destroy()
		}
	})
}

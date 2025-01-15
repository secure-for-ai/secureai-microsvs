package sqlBuilderV3_test

import (
	"testing"

	"github.com/secure-for-ai/secureai-microsvs/db/sqlBuilderV3"
	"github.com/stretchr/testify/assert"
)

type student2 struct {
	Uid        int64  `db:"uid"`
	Username   string `db:"username"`
	Nickname   string `db:"nickname"`
	Email      string `db:"email"`
	CreateTime int64  `db:"create_time"`
	UpdateTime int64  `db:"update_time"`
}

var stuStruct2 = student2{
	uid,
	"Alice",
	"Ali",
	"ali@gmail.com",
	ts.Unix(),
	ts.Unix(),
}

var stuVal2 = sqlBuilderV3.Map{
	"uid":         uid,
	"username":    "Alice",
	"nickname":    "Ali",
	"email":       "ali@gmail.com",
	"create_time": ts.Unix(),
	"update_time": ts.Unix(),
}

func fastData(data interface{}) interface{} {
	return &data
}
func BenchmarkSQLStmtInsertAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var args []interface{}
		var err error
		var stuInterface interface{} = stuStruct

		evalSingle := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "INSERT INTO student (uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time) VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
			assert.EqualValues(b, stuStructArr, args)
		}

		for pb.Next() {
			// memory alloc is due to converting type to interface{}
			// and Insert() need to use []interface{}
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.InsertOne(stuInterface)
			sql, args, err = stmt.Gen(w)
			evalSingle()
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {

		// This is the most efficient way of inserting one struct.
		// 1. convert it to an interface
		// 2. call InsertOne
		//
		// convert struct to interface to avoid mem alloc
		// cpu: Intel(R) Core(TM) i7-6820HQ CPU @ 2.70GHz
		// === RUN   BenchmarkSQLStmtInsert
		// BenchmarkSQLStmtInsert
		// BenchmarkSQLStmtInsert-8         1539531               745.5 ns/op             0 B/op          0 allocs/op
		var stuInterface interface{} = stuStruct

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.InsertOne(stuInterface)
			stmt.Gen(w)
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertOneMapInterfaceAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {

		var sql string
		var args []interface{}
		var err error
		var stuInterface interface{} = stuVal

		evalSingle := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "INSERT INTO student (age,comp,create_time,email,enrolled,gpa,nickname,tokens,uid,update_time,username) VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
			assert.EqualValues(b, stuStructArrSorted, args)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.Insert().IntoTable("student").ValuesOne(stuInterface)
			sql, args, err = stmt.Gen(w)
			evalSingle()
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertOneMapInterface(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {

		// This is the most efficient way of inserting one struct.
		// 1. convert it to an interface
		// 2. call InsertOne
		//
		// cpu: Intel(R) Core(TM) i7-6820HQ CPU @ 2.70GHz
		// === RUN   BenchmarkSQLStmtInsertOneMapInterface
		// BenchmarkSQLStmtInsertOneMapInterface
		// BenchmarkSQLStmtInsertOneMapInterface-8          1580680               762.5 ns/op             0 B/op          0 allocs/op
		var stuInterface interface{} = stuVal

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.Insert().IntoTable("student").ValuesOne(stuInterface)
			stmt.Gen(w)
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertOneInterfaceAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {

		// This is the most efficient way of inserting one struct.
		// 1. convert it to an interface
		// 2. call InsertOne
		//
		// cpu: Intel(R) Core(TM) i7-6820HQ CPU @ 2.70GHz
		// === RUN   BenchmarkSQLStmtInsertOneInterface
		// BenchmarkSQLStmtInsertOneInterface
		// BenchmarkSQLStmtInsertOneInterface-8     2812692               387.1 ns/op             0 B/op          0 allocs/op
		var sql string
		var args []interface{}
		var err error
		var stuInterface interface{} = stuStructArrExpr

		evalSingle := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "INSERT INTO student VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
			assert.EqualValues(b, stuStructArr, args)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.Insert().IntoTable("student").Values(stuInterface)
			sql, args, err = stmt.Gen(w)
			evalSingle()
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertOneInterface(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {

		var stuInterface interface{} = stuStructArrExpr

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.Insert().IntoTable("student").ValuesOne(stuInterface)
			stmt.Gen(w)
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertBulkStructInterfaceAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var bulkArgs []*[]interface{}
		var err error
		var structInterface = []interface{}{stuStruct, stuStruct}

		evalBulk := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "INSERT INTO student (uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time) VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
			assert.EqualValues(b, []*[]interface{}{&stuStructArr, &stuStructArr}, bulkArgs)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.InsertBulk(&structInterface)
			sql, _, err = stmt.Gen(w)
			bulkArgs = w.BulkArgs()
			evalBulk()
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertBulkStructInterface(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		// but struct in []interface{} is essential the most efficient way to insert.
		// cpu: Intel(R) Core(TM) i7-6820HQ CPU @ 2.70GHz
		// === RUN   BenchmarkSQLStmtInsertBulkStructInterface
		// BenchmarkSQLStmtInsertBulkStructInterface
		// BenchmarkSQLStmtInsertBulkStructInterface-8      1573752               750.7 ns/op             0 B/op          0 allocs/op
		structInterface := []interface{}{stuStruct, stuStruct}
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.InsertBulk(&structInterface)
			stmt.Gen(w)
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertBulkStructInterfaceAssert2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var bulkArgs []*[]interface{}
		var err error
		var structInterface interface{} = stuList

		evalBulk := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "INSERT INTO student (uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time) VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
			assert.EqualValues(b, []*[]interface{}{&stuStructArr, &stuStructArr}, bulkArgs)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.InsertBulk(structInterface)
			sql, _, err = stmt.Gen(w)
			bulkArgs = w.BulkArgs()
			evalBulk()
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertBulkStructInterface2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		// cpu: Intel(R) Core(TM) i7-6820HQ CPU @ 2.70GHz
		// === RUN   BenchmarkSQLStmtInsertBulkStructInterface2
		// BenchmarkSQLStmtInsertBulkStructInterface2
		// BenchmarkSQLStmtInsertBulkStructInterface2-8      600966              1683 ns/op             578 B/op          4 allocs/op
		var structInterface interface{} = stuList
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.InsertBulk(structInterface)
			stmt.Gen(w)
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertBulkMapInterfaceAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var bulkArgs []*[]interface{}
		var err error
		var structInterface = []interface{}{stuVal, stuVal}
		var structInterface2 = []interface{}{stuValExpr, stuValExpr}

		evalBulk := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "INSERT INTO student (age,comp,create_time,email,enrolled,gpa,nickname,tokens,uid,update_time,username) VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
			assert.EqualValues(b, []*[]interface{}{&stuStructArrSorted, &stuStructArrSorted}, bulkArgs)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.InsertBulk(&structInterface).IntoTable("student")
			sql, _, err = stmt.Gen(w)
			bulkArgs = w.BulkArgs()
			evalBulk()
			stmt.Reset()
			stmt.InsertBulk(&structInterface2).IntoTable("student")
			sql, _, err = stmt.Gen(w)
			bulkArgs = w.BulkArgs()
			evalBulk()
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertBulkMapInterface(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var structInterface = []interface{}{stuVal, stuVal}
		var structInterface2 = []interface{}{stuValExpr, stuValExpr}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.InsertBulk(&structInterface).IntoTable("student")
			stmt.Gen(w)
			stmt.Reset()
			stmt.InsertBulk(&structInterface2).IntoTable("student")
			stmt.Gen(w)
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertBulkSliceInterfaceAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var bulkArgs []*[]interface{}
		var err error
		var structInterface = []interface{}{stuStructArr, stuStructArr}
		var structInterface2 = []interface{}{stuStructArrExpr, stuStructArrExpr}

		evalBulk := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "INSERT INTO student VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
			assert.EqualValues(b, []*[]interface{}{&stuStructArr, &stuStructArr}, bulkArgs)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.Insert().IntoTable("student").ValuesBulk(&structInterface)
			sql, _, err = stmt.Gen(w)
			bulkArgs = w.BulkArgs()
			evalBulk()
			stmt.Reset()
			stmt.InsertBulk(&structInterface2).IntoTable("student")
			sql, _, err = stmt.Gen(w)
			bulkArgs = w.BulkArgs()
			evalBulk()
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtInsertBulkSliceInterface(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		// but struct in []interface{} is essential the most efficient way to insert.
		// cpu: Intel(R) Core(TM) i7-6820HQ CPU @ 2.70GHz
		// === RUN   BenchmarkSQLStmtInsertBulkStructInterface
		// BenchmarkSQLStmtInsertBulkStructInterface
		// BenchmarkSQLStmtInsertBulkStructInterface-8      1573752               750.7 ns/op             0 B/op          0 allocs/op
		var structInterface = []interface{}{stuStructArr, stuStructArr}
		var structInterface2 = []interface{}{stuStructArrExpr, stuStructArrExpr}
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.Insert().IntoTable("student").ValuesBulk(&structInterface)
			stmt.Gen(w)
			stmt.Reset()
			stmt.InsertBulk(&structInterface2).IntoTable("student")
			stmt.Gen(w)
			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtDeleteCondAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var args []interface{}
		var err error

		evalStruct := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "DELETE FROM student WHERE uid = ?", sql)
			assert.EqualValues(b, []interface{}{uid}, args)
		}

		evalAnd := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "DELETE FROM student WHERE (uid = ?) AND (username = ?)", sql)
			assert.EqualValues(b, []interface{}{uid, "Alice"}, args)
		}

		evalOr := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "DELETE FROM student WHERE (uid = ?) OR (username = ?)", sql)
			assert.EqualValues(b, []interface{}{uid, "Alice"}, args)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			uidEq100 := sqlBuilderV3.CondExpr("uid = ??", uid)
			usernameEqAlice := sqlBuilderV3.CondExpr("username = ??", "Alice")
			eqCondAnd := sqlBuilderV3.And(uidEq100, usernameEqAlice)
			eqCondOr := sqlBuilderV3.Or(uidEq100, usernameEqAlice)

			stmt := sqlBuilderV3.Delete(&stuStruct, uidEq100)
			sql, args, err = stmt.Gen(w)
			evalStruct()
			stmt.Reset()

			sql, args, err = stmt.Delete(&stuStruct, uidEq100, usernameEqAlice).Gen(w)
			evalAnd()
			stmt.Reset()
			sql, args, err = stmt.Delete(&stuStruct, eqCondAnd).Gen(w)
			evalAnd()
			stmt.Reset()

			sql, args, err = stmt.Delete(&stuStruct).Or(uidEq100, usernameEqAlice).Gen(w)
			evalOr()
			stmt.Reset()
			sql, args, err = stmt.Delete(&stuStruct, eqCondOr).Gen(w)
			evalOr()

			stmt.Destroy()
			uidEq100.Destroy()
			usernameEqAlice.Destroy()
			eqCondAnd.Destroy()
			eqCondOr.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtCondDelete(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			uidEq100 := sqlBuilderV3.CondExpr("uid = ??", uid)
			usernameEqAlice := sqlBuilderV3.CondExpr("username = ??", "Alice")
			eqCondAnd := sqlBuilderV3.And(uidEq100, usernameEqAlice)
			eqCondOr := sqlBuilderV3.Or(uidEq100, usernameEqAlice)

			stmt := sqlBuilderV3.Delete(&stuStruct, uidEq100)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Delete(&stuStruct, uidEq100, usernameEqAlice).Gen(w)
			stmt.Reset()
			stmt.Delete(&stuStruct, eqCondAnd).Gen(w)
			stmt.Reset()

			stmt.Delete(&stuStruct).Or(uidEq100, usernameEqAlice).Gen(w)
			stmt.Reset()
			stmt.Delete(&stuStruct, eqCondOr).Gen(w)

			stmt.Destroy()
			uidEq100.Destroy()
			usernameEqAlice.Destroy()
			eqCondAnd.Destroy()
			eqCondOr.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtDeleteMapAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var args []interface{}
		var err error

		evalStruct := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "DELETE FROM student WHERE uid = ?", sql)
			assert.EqualValues(b, []interface{}{uid}, args)
		}

		evalAnd := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "DELETE FROM student WHERE (uid = ?) AND (username = ?)", sql)
			assert.EqualValues(b, []interface{}{uid, "Alice"}, args)
		}

		evalOr := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "DELETE FROM student WHERE (uid = ?) OR (username = ?)", sql)
			assert.EqualValues(b, []interface{}{uid, "Alice"}, args)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.Delete(&stuStruct, stuMapUid)
			sql, args, err = stmt.Gen(w)
			evalStruct()
			stmt.Reset()

			sql, args, err = stmt.Delete(&stuStruct).Where(stuMapUid).And(stuMapUsername).Gen(w)
			evalAnd()
			stmt.Reset()
			sql, args, err = stmt.Delete(&stuStruct).Where(eqCond).Gen(w)
			evalAnd()
			stmt.Reset()

			sql, args, err = stmt.Delete(&stuStruct).Where(stuMapUid).Or(stuMapUsername).Gen(w)
			evalOr()
			stmt.Reset()
			sql, args, err = stmt.Delete(&stuStruct).Or(eqCond).Gen(w)
			evalOr()

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtDeleteMap(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.Delete(&stuStruct, stuMapUid)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Delete(&stuStruct).Where(stuMapUid).And(stuMapUsername).Gen(w)
			stmt.Reset()
			stmt.Delete(&stuStruct).Where(eqCond).Gen(w)
			stmt.Reset()

			stmt.Delete(&stuStruct).Where(stuMapUid).Or(stuMapUsername).Gen(w)
			stmt.Reset()
			stmt.Delete(&stuStruct).Or(eqCond).Gen(w)

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtDeleteStringAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var args []interface{}
		var err error

		evalStruct := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "DELETE FROM student WHERE uid = ?", sql)
			assert.EqualValues(b, []interface{}{uid}, args)
		}

		evalAnd := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "DELETE FROM student WHERE (uid = ?) AND (username = ?)", sql)
			assert.EqualValues(b, []interface{}{uid, "Alice"}, args)
		}

		evalOr := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "DELETE FROM student WHERE (uid = ?) OR (username = ?)", sql)
			assert.EqualValues(b, []interface{}{uid, "Alice"}, args)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.Delete(&stuStruct, "uid = ??", uid)
			sql, args, err = stmt.Gen(w)
			evalStruct()
			stmt.Reset()

			sql, args, err = stmt.Delete(&stuStruct).Where("uid = ??", uid).And("username = ?", "Alice").Gen(w)
			evalAnd()
			stmt.Reset()

			sql, args, err = stmt.Delete(&stuStruct).Where("uid = ??", uid).Or("username = ?", "Alice").Gen(w)
			evalOr()

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtDeleteString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()
			stmt := sqlBuilderV3.Delete(&stuStruct, "uid = ??", uid)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Delete(&stuStruct).Where("uid = ??", uid).And("username = ?", "Alice").Gen(w)
			stmt.Reset()

			stmt.Delete(&stuStruct).Where("uid = ??", uid).Or("username = ?", "Alice").Gen(w)

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtUpdateAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var args []interface{}
		var err error
		var uidEq100 = sqlBuilderV3.CondExpr("uid = ??", uid)
		var usernameEqAlice = sqlBuilderV3.CondExpr("username = ??", "Alice")
		var stuInterface interface{} = stuStruct
		var createTime = sqlBuilderV3.ExprEq("create_time", ts.Unix())
		var uidExpr = sqlBuilderV3.Expr("??", uid)
		var tokens interface{} = []string{"token1", "token2"}
		var updateTime interface{} = ts.Unix()

		evalStruct := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "UPDATE student SET "+
				"uid = ?,username = ?,nickname = ?,email = ?,age = ?,enrolled = ?,gpa = ?,tokens = ?,comp = ?,create_time = ?,update_time = ? "+
				"WHERE uid = ?", sql)
			assert.EqualValues(b, append(stuStructArr, uid), args)
		}

		evalMap := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "UPDATE student SET uid = ?", sql)
			assert.EqualValues(b, []interface{}{uid}, args)
		}

		evalWhereAnd := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "UPDATE student SET "+
				"uid = ?,username = ?,nickname = ?,email = ?,age = ?,enrolled = ?,gpa = ?,tokens = ?,comp = ?,create_time = ?,update_time = ? "+
				"WHERE (uid = ?) AND (username = ?)", sql)
			assert.EqualValues(b, append(stuStructArr, uid, "Alice"), args)
		}

		evalWhereOr := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "UPDATE student SET "+
				"uid = ?,username = ?,nickname = ?,email = ?,age = ?,enrolled = ?,gpa = ?,tokens = ?,comp = ?,create_time = ?,update_time = ? "+
				"WHERE (uid = ?) OR (username = ?)", sql)
			assert.EqualValues(b, append(stuStructArr, uid, "Alice"), args)
		}

		evalIncr := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "UPDATE student SET username = ?,uid = uid + ? WHERE uid = ?", sql)
			assert.EqualValues(b, []interface{}{"Alice", 10, uid}, args)
		}

		evalDecr := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "UPDATE student SET username = ?,uid = uid - ? WHERE uid = ?", sql)
			assert.EqualValues(b, []interface{}{"Alice", 10, uid}, args)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			// set from struct
			stmt := sqlBuilderV3.SQL().Update(stuInterface, uidEq100)
			sql, args, err = stmt.Gen(w)
			evalStruct()
			stmt.Reset()

			// set from string and *condExpr
			sql, args, err = stmt.Update().From(&stuStruct).
				Set("uid", uidExpr).
				Set("username", "??", "Alice").
				Set("nickname", "Ali").
				Set("email", "ali@gmail.com").
				Set("age", uint64(20)).
				Set("enrolled", "??", true).
				Set("gpa", "??", float64(3.5)).
				Set("tokens", "??", tokens).
				Set("comp", "??", complex(10, 11)).
				Set(createTime).
				Set("update_time", "??", updateTime).
				Where(uidEq100).
				Gen(w)
			evalStruct()
			stmt.Reset()

			// set from Map
			sql, args, err = stmt.Update(stuMapUid).From(stuInterface).Gen(w)
			evalMap()
			stmt.Reset()
			sql, args, err = stmt.Update().From(stuInterface).Set(stuMapUid).Gen(w)
			evalMap()
			stmt.Reset()
			sql, args, err = stmt.Update().From(stuInterface).Set(stuMapExpr).Gen(w)
			evalMap()
			stmt.Reset()

			// update and
			// stmt.Update(stuInterface).Where(uidEq100, usernameEqAlice).Gen(w)
			sql, args, err = stmt.Update(stuInterface).Where(uidEq100, usernameEqAlice).Gen(w)
			evalWhereAnd()
			stmt.Reset()
			// update or
			sql, args, err = stmt.Update(stuInterface).Or(uidEq100, usernameEqAlice).Gen(w)
			evalWhereOr()
			stmt.Reset()
			// update incr
			sql, args, err = stmt.Update("student").Set("username", "??", "Alice").Incr("uid", 10).Where(uidEq100).Gen(w)
			evalIncr()
			stmt.Reset()
			// update decr
			sql, args, err = sqlBuilderV3.Update("student").Set("username", "??", "Alice").Decr("uid", 10).Where(uidEq100).Gen(w)
			evalDecr()

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtUpdate(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {

		var uidEq100 = sqlBuilderV3.CondExpr("uid = ??", uid)
		var usernameEqAlice = sqlBuilderV3.CondExpr("username = ??", "Alice")
		var stuInterface interface{} = stuStruct
		var createTime = sqlBuilderV3.ExprEq("create_time", ts.Unix())
		var uidExpr = sqlBuilderV3.Expr("??", uid)
		var tokens interface{} = []string{"token1", "token2"}
		var updateTime interface{} = ts.Unix()

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			// set from struct
			stmt := sqlBuilderV3.SQL().Update(stuInterface, uidEq100)
			stmt.Gen(w)
			stmt.Reset()

			// set from Map
			stmt.Update().From(&stuStruct).
				Set("uid", uidExpr).
				Set("username", "??", "Alice").
				Set("nickname", "Ali").
				Set("email", "ali@gmail.com").
				Set("age", uint64(20)).
				Set("enrolled", "??", true).
				Set("gpa", "??", float64(3.5)).
				Set("tokens", "??", tokens).
				Set("comp", "??", complex(10, 11)).
				Set(createTime).
				Set("update_time", "??", updateTime).
				Where(uidEq100).
				Gen(w)
			stmt.Reset()

			// set from Map
			stmt.Update(stuMapUid).From(stuInterface).Gen(w)
			stmt.Reset()
			stmt.Update().From(stuInterface).Set(stuMapUid).Gen(w)
			stmt.Reset()
			stmt.Update().From(stuInterface).Set(stuMapExpr).Gen(w)
			stmt.Reset()

			// update and
			stmt.Update(stuInterface).Where(uidEq100, usernameEqAlice).Gen(w)
			stmt.Reset()
			// update or
			stmt.Update(stuInterface).Or(uidEq100, usernameEqAlice).Gen(w)
			stmt.Reset()
			// update incr
			stmt.Update("student").Set("username", "??", "Alice").Incr("uid", 10).Where(uidEq100).Gen(w)
			stmt.Reset()
			// update decr
			stmt.Update("student").Set("username", "??", "Alice").Decr("uid", 10).Where(uidEq100).Gen(w)

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtSelectAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var args []interface{}
		var err error
		var stuInterface interface{} = stuStruct
		var stuColsStr interface{} = []string{"uid", "username", "nickname", "email", "age", "enrolled", "gpa", "tokens", "comp", "create_time", "update_time"}
		var uidEq100 = sqlBuilderV3.Expr("uid = ??", 100)
		var usernameEqAlice = sqlBuilderV3.Expr("username = ??", "Alice")

		evalStruct := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student WHERE uid = ?", sql)
			assert.EqualValues(b, []interface{}{100}, args)
		}

		evalAnd := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student WHERE (uid = ?) AND (username = ?)", sql)
			assert.EqualValues(b, []interface{}{100, "Alice"}, args)
		}

		evalAny := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT * FROM student", sql)
			assert.EqualValues(b, []interface{}{}, args)
		}

		evalAS := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT * FROM student AS S", sql)
			assert.EqualValues(b, []interface{}{}, args)
		}

		evalJoin := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT * FROM student,student", sql)
			assert.EqualValues(b, []interface{}{}, args)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			stmt := sqlBuilderV3.SQL().Select(stuInterface, uidEq100)
			sql, args, err = stmt.Gen(w)
			evalStruct()
			stmt.Reset()

			stmt.Select(stuColsStr).From(stuInterface).Where(uidEq100)
			sql, args, err = stmt.Gen(w)
			evalStruct()
			stmt.Reset()

			stmt.Select(stuInterface).Where(uidEq100, usernameEqAlice).Gen(w)
			sql, args, err = stmt.Gen(w)
			evalAnd()
			stmt.Reset()

			stmt.Select().From("student")
			sql, args, err = stmt.Gen(w)
			evalAny()
			stmt.Reset()

			stmt.Select().From(stuInterface, "S")
			sql, args, err = stmt.Gen(w)
			evalAS()
			stmt.Reset()

			stmt.Select().From(stuInterface).From(stuInterface)
			sql, args, err = stmt.Gen(w)
			evalJoin()

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtSelect(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var uidEq100 = sqlBuilderV3.Expr("uid = ??", 100)
		var stuColsStr interface{} = []string{"uid", "username", "nickname", "email", "age", "enrolled", "gpa", "tokens", "comp", "create_time", "update_time"}
		var stuInterface interface{} = stuStruct
		var usernameEqAlice = sqlBuilderV3.Expr("username = ??", "Alice")

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			stmt := sqlBuilderV3.SQL().Select(stuInterface, uidEq100)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Select(stuColsStr).From(stuInterface).Where(uidEq100)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Select(stuInterface).Where(uidEq100, usernameEqAlice).Gen(w)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Select().From("student")
			stmt.Gen(w)
			stmt.Reset()

			stmt.Select().From(stuInterface, "S")
			stmt.Gen(w)
			stmt.Reset()

			stmt.Select().From(stuInterface).From(stuInterface)
			stmt.Gen(w)

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtOrderByAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var args []interface{}
		var err error

		eval := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT * FROM student ORDER BY uid ASC, username DESC", sql)
			assert.EqualValues(b, []interface{}{}, args)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			stmt := sqlBuilderV3.Select().From("student").Asc("uid").Desc("username")
			sql, args, err = stmt.Gen(w)
			eval()

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtOrderBy(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			stmt := sqlBuilderV3.Select().From("student").Asc("uid").Desc("username")
			stmt.Gen(w)

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtLimitAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var args []interface{}
		var err error

		evalLimitOffset := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT * FROM student LIMIT 10 OFFSET 5", sql)
			assert.EqualValues(b, []interface{}{}, args)
		}

		evalLimit := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT * FROM student LIMIT 10", sql)
			assert.EqualValues(b, []interface{}{}, args)
		}

		evalLimit0 := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT * FROM student LIMIT 10", sql)
			assert.EqualValues(b, []interface{}{}, args)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			stmt := sqlBuilderV3.Select().From("student").Limit(10, 5)
			sql, args, err = stmt.Gen(w)
			evalLimitOffset()
			stmt.Reset()

			stmt.Select().From("student").Limit(10)
			sql, args, err = stmt.Gen(w)
			evalLimit()
			stmt.Reset()

			stmt.Select().From("student").Limit(10, 0)
			sql, args, err = stmt.Gen(w)
			evalLimit0()

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtLimit(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			stmt := sqlBuilderV3.Select().From("student").Limit(10, 5)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Select().From("student").Limit(10)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Select().From("student").Limit(10, 0)
			stmt.Gen(w)

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtGroupByHavingAssert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var sql string
		var args []interface{}
		var err error
		var stuInterface interface{} = stuStruct
		var uidGe100 = sqlBuilderV3.Expr("COUNT(uid)>??", uid)

		eval := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student GROUP BY username HAVING COUNT(uid)>?", sql)
			assert.EqualValues(b, []interface{}{uid}, args)
		}

		evalAnd := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student GROUP BY username HAVING (COUNT(uid)>?) AND (COUNT(uid)>?)", sql)
			assert.EqualValues(b, []interface{}{uid, uid}, args)
		}

		evalOr := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student GROUP BY username HAVING (COUNT(uid)>?) OR (COUNT(uid)>?)", sql)
			assert.EqualValues(b, []interface{}{uid, uid}, args)
		}

		evalComplex := func() {
			assert.NoError(b, err)
			assert.EqualValues(b, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student GROUP BY username, nickname HAVING (COUNT(uid)>?) AND (COUNT(uid)>?)", sql)
			assert.EqualValues(b, []interface{}{uid, uid}, args)
		}

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			stmt := sqlBuilderV3.SQL().Select(stuInterface).GroupBy("username").Having(uidGe100)
			sql, args, err = stmt.Gen(w)
			eval()
			stmt.Reset()

			stmt.Select(stuInterface).GroupBy("username").Having(uidGe100, uidGe100)
			sql, args, err = stmt.Gen(w)
			evalAnd()
			stmt.Reset()

			stmt.Select(stuInterface).GroupBy("username").HavingOr(uidGe100, uidGe100)
			sql, args, err = stmt.Gen(w)
			evalOr()
			stmt.Reset()

			stmt.Select(stuInterface).GroupBy().GroupBy("username").GroupBy("nickname").Having(uidGe100).HavingAnd(uidGe100)
			sql, args, err = stmt.Gen(w)
			evalComplex()

			stmt.Destroy()
			w.Destroy()
		}
	})
}

func BenchmarkSQLStmtGroupByHaving(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		var stuInterface interface{} = stuStruct
		var uidGe100 = sqlBuilderV3.Expr("COUNT(uid)>??", uid)

		for pb.Next() {
			w := sqlBuilderV3.NewWriter()

			stmt := sqlBuilderV3.SQL().Select(stuInterface).GroupBy("username").Having(uidGe100)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Select(stuInterface).GroupBy("username").Having(uidGe100, uidGe100)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Select(stuInterface).GroupBy("username").HavingOr(uidGe100, uidGe100)
			stmt.Gen(w)
			stmt.Reset()

			stmt.Select(stuInterface).GroupBy().GroupBy("username").GroupBy("nickname").Having(uidGe100).HavingAnd(uidGe100)
			stmt.Gen(w)

			stmt.Destroy()
			w.Destroy()
		}
	})
}

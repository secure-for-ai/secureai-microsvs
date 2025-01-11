package sqlBuilderV3_test

import (
	"testing"
	"time"

	"github.com/secure-for-ai/secureai-microsvs/db"
	"github.com/secure-for-ai/secureai-microsvs/db/sqlBuilderV3"
	"github.com/stretchr/testify/assert"
)

type student struct {
	Uid        int64      `db:"uid"`
	Username   string     `db:"username"`
	Nickname   string     `db:"nickname"`
	Email      string     `db:"email"`
	Age        uint64     `db:"age"`
	Enrolled   bool       `db:"enrolled"`
	GPA        float64    `db:"gpa"`
	Tokens     []string   `db:"tokens"`
	Comp       complex128 `db:"comp"`
	CreateTime int64      `db:"create_time"`
	UpdateTime int64      `db:"update_time"`
}

// in no tag version, the member name need to match the column name exactly.
// This actually losses the flexibility. That's why we recommend to use the
// `db:"[column_name]"` Tag
type studentNoTag struct {
	uid         int64
	username    string
	nickname    string
	email       string
	age         uint64
	enrolled    bool
	gpa         float64
	tokens      []string
	comp        complex128
	create_time int64
	update_time int64
}

func (s student) GetTableName() string {
	return "student"
}

func (s student) Size() int {
	return len(s.Username) + len(s.Nickname) + len(s.Nickname) + 24
}

var (
	uid       int64 = 100
	ts, _           = time.Parse(time.UnixDate, "Sat Mar 7 11:06:39 PST 2015")
	stuStruct       = student{
		uid,
		"Alice",
		"Ali",
		"ali@gmail.com",
		20,
		true,
		3.5,
		[]string{"token1", "token2"},
		complex(10, 11),
		ts.Unix(),
		ts.Unix(),
	}
	stuStructArr = []interface{}{
		uid,
		"Alice",
		"Ali",
		"ali@gmail.com",
		uint64(20),
		true,
		float64(3.5),
		[]string{"token1", "token2"},
		complex(10, 11),
		ts.Unix(),
		ts.Unix(),
	}
	stuStructArrSorted = []interface{}{
		uint64(20),
		complex(10, 11),
		ts.Unix(),
		"ali@gmail.com",
		true,
		float64(3.5),
		"Ali",
		[]string{"token1", "token2"},
		uid,
		ts.Unix(),
		"Alice",
	}
	stuList = []student{
		stuStruct,
		stuStruct,
	}
	eqCond = sqlBuilderV3.Map{
		"uid":      uid,
		"username": "Alice",
	}
	stuMapUid = sqlBuilderV3.Map{
		"uid": uid,
	}
	stuMapUsername = sqlBuilderV3.Map{
		"username": "Alice",
	}
	stuVal = sqlBuilderV3.Map{
		"uid":         uid,
		"username":    "Alice",
		"nickname":    "Ali",
		"email":       "ali@gmail.com",
		"age":         uint64(20),
		"enrolled":    true,
		"gpa":         3.5,
		"tokens":      []string{"token1", "token2"},
		"comp":        complex(10, 11),
		"create_time": ts.Unix(),
		"update_time": ts.Unix(),
	}
	stuValExpr = sqlBuilderV3.Map{
		"uid":         sqlBuilderV3.Expr(db.Para, uid),
		"username":    sqlBuilderV3.Expr(db.Para, "Alice"),
		"nickname":    sqlBuilderV3.Expr(db.Para, "Ali"),
		"email":       sqlBuilderV3.Expr(db.Para, "ali@gmail.com"),
		"age":         sqlBuilderV3.Expr(db.Para, uint64(20)),
		"enrolled":    sqlBuilderV3.Expr(db.Para, true),
		"gpa":         sqlBuilderV3.Expr(db.Para, 3.5),
		"tokens":      sqlBuilderV3.Expr(db.Para, []string{"token1", "token2"}),
		"comp":        sqlBuilderV3.Expr(db.Para, complex(10, 11)),
		"create_time": sqlBuilderV3.Expr(db.Para, ts.Unix()),
		"update_time": sqlBuilderV3.Expr(db.Para, ts.Unix()),
	}
	stuStructNoTag = studentNoTag{
		uid,
		"Alice",
		"Ali",
		"ali@gmail.com",
		20,
		true,
		3.5,
		[]string{"token1", "token2"},
		complex(10, 11),
		ts.Unix(),
		ts.Unix(),
	}
	stuStructArrExpr = []interface{}{
		sqlBuilderV3.Expr(db.Para, uid),
		sqlBuilderV3.Expr(db.Para, "Alice"),
		sqlBuilderV3.Expr(db.Para, "Ali"),
		sqlBuilderV3.Expr(db.Para, "ali@gmail.com"),
		sqlBuilderV3.Expr(db.Para, uint64(20)),
		sqlBuilderV3.Expr(db.Para, true),
		sqlBuilderV3.Expr(db.Para, 3.5),
		sqlBuilderV3.Expr(db.Para, []string{"token1", "token2"}),
		sqlBuilderV3.Expr(db.Para, complex(10, 11)),
		sqlBuilderV3.Expr(db.Para, ts.Unix()),
		sqlBuilderV3.Expr(db.Para, ts.Unix()),
	}
	stuMapExpr = sqlBuilderV3.Map{
		"uid": sqlBuilderV3.Expr(db.Para, uid),
	}
)

func TestSQLStmt_Insert(t *testing.T) {
	var sql string
	var args []interface{}
	var err error
	w := sqlBuilderV3.NewWriter()

	evalSingle := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time) VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
		assert.EqualValues(t, stuStructArr, args)
	}

	evalSingleMap := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (age,comp,create_time,email,enrolled,gpa,nickname,tokens,uid,update_time,username) VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
		assert.EqualValues(t, stuStructArrSorted, args)
	}

	evalSingleValueOnly := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
		assert.EqualValues(t, stuStructArr, args)
	}

	// Test insert struct
	sql, args, err = sqlBuilderV3.SQL().Insert(&stuStruct).Gen(w)
	evalSingle()
	sql, args, err = sqlBuilderV3.Insert().IntoTable(&stuStruct).Values(&stuStruct).Gen(w)
	evalSingle()
	sql, args, err = sqlBuilderV3.Insert().IntoTable(&stuStruct).Values(&stuStruct).Values().Gen(w)
	evalSingle()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").Values(&stuStructNoTag).Values().Gen(w)
	evalSingle()
	// Test insert []interface{}
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(&stuStruct).Values(stuStructArr).Gen(w)
	evalSingle()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(&stuStruct).Values(stuStructArrExpr).Gen(w)
	evalSingle()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").
		IntoColumns([]string{"uid", "username", "nickname", "email", "age", "enrolled", "gpa", "tokens", "comp", "create_time", "update_time"}).
		Values(stuStructArr).Gen(w)
	evalSingle()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").
		IntoColumns(sqlBuilderV3.Columns{"uid", "username", "nickname", "email", "age", "enrolled", "gpa", "tokens", "comp", "create_time", "update_time"}).
		Values(stuStructArr).Gen(w)
	evalSingle()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").
		IntoColumns("uid", "username", "nickname", "email", "age", "enrolled", "gpa", "tokens", "comp", "create_time", "update_time").
		Values(stuStructArr).Gen(w)
	evalSingle()
	// Test insert Map
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").Values(stuVal).Gen(w)
	evalSingleMap()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").Values(stuValExpr).Gen(w)
	evalSingleMap()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(stuVal).Values(stuVal).Gen(w)
	evalSingleMap()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(stuVal).Values(stuValExpr).Gen(w)
	evalSingleMap()
	// Test insert []interface{} value only
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").
		Values(stuStructArr).Gen(w)
	evalSingleValueOnly()

	// sql, args, err = sqlBuilderV3.Insert().IntoTable("student").ValuesBulk([]interface{}{stuVal}).Gen(w)
	// evalSingleMap()
	// sql, args, err = sqlBuilderV3.Insert().IntoTable("student").ValuesBulk([]interface{}{stuStruct}).Gen(w)
	// evalSingle()

	evalBulk := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time) VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
		assert.EqualValues(t, []*[]interface{}{&stuStructArr, &stuStructArr}, w.BulkArgs())
	}

	evalBulkMap := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (age,comp,create_time,email,enrolled,gpa,nickname,tokens,uid,update_time,username) VALUES (?,?,?,?,?,?,?,?,?,?,?)", sql)
		assert.EqualValues(t, []*[]interface{}{&stuStructArrSorted, &stuStructArrSorted}, w.BulkArgs())
	}

	// test insert bulk corner cases: ValuesBulk() nil input and slice having size 0.
	sql, args, err = sqlBuilderV3.Insert(&stuStruct, &stuStruct).ValuesBulk(nil).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.Insert(&stuStruct, &stuStruct).ValuesBulk([]interface{}{}).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.InsertBulk(stuList).InsertBulk(nil).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.InsertBulk(stuList).InsertBulk([]interface{}{}).Gen(w)
	evalBulk()
	// test insert bulk struct
	sql, args, err = sqlBuilderV3.Insert(&stuStruct, &stuStruct).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.InsertBulk(stuList).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.InsertBulk(&stuList).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.InsertBulk([]interface{}{stuStruct, stuStruct}).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.InsertBulk(&[]interface{}{stuStruct, stuStruct}).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.InsertBulk([]interface{}{&stuStruct, &stuStruct}).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.InsertBulk(&[]interface{}{&stuStruct, &stuStruct}).Gen(w)
	evalBulk()
	// Test insert bulk []interface{}
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(&stuStruct).Values(stuStructArr, stuStructArr).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(&stuStruct).ValuesBulk([]interface{}{stuStructArr, stuStructArr}).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(&stuStruct).ValuesBulk(&[]interface{}{stuStructArr, stuStructArr}).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(&stuStruct).Values(stuStructArrExpr, stuStructArrExpr).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(&stuStruct).ValuesBulk([]interface{}{stuStructArrExpr, stuStructArrExpr}).Gen(w)
	evalBulk()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(&stuStruct).ValuesBulk(&[]interface{}{stuStructArrExpr, stuStructArrExpr}).Gen(w)
	evalBulk()
	// Test insert bulk Map
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").Insert(stuVal, stuVal).Gen(w)
	evalBulkMap()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").InsertBulk([]interface{}{stuVal, stuVal}).Gen(w)
	evalBulkMap()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").InsertBulk(&[]interface{}{stuVal, stuVal}).Gen(w)
	evalBulkMap()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").Insert(stuValExpr, stuValExpr).Gen(w)
	evalBulkMap()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").InsertBulk([]interface{}{stuValExpr, stuValExpr}).Gen(w)
	evalBulkMap()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").InsertBulk(&[]interface{}{stuValExpr, stuValExpr}).Gen(w)
	evalBulkMap()

	// Test insert select
	evalSelect1 := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}

	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").Select(stuStruct).Where(stuMapUid).Gen(w)
	evalSelect1()
	sql, args, err = sqlBuilderV3.Insert().IntoTable(&stuStruct).Select(&stuStruct).Where(stuMapUid).Gen(w)
	evalSelect1()

	evalSelect2 := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time) SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilderV3.Insert(&stuStruct).Select(&stuStruct).Where(stuMapUid).Gen(w)
	evalSelect2()
}

func TestSQLStmt_Delete(t *testing.T) {
	var sql string
	var args []interface{}
	var err error
	var uidEq100 = sqlBuilderV3.CondExpr("uid = ??", uid)
	var usernameEqAlice = sqlBuilderV3.CondExpr("username = ??", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilderV3.SQL().Delete(&stuStruct, uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct, uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct, stuMapUid).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Delete().From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Delete("student", uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Delete("student").Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Delete().From("student").Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Delete().From("student").Where("uid = ??", uid).Gen(w)
	evalStruct()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student WHERE (uid = ?) AND (username = ?)", sql)
		assert.EqualValues(t, []interface{}{uid, "Alice"}, args)
	}
	sql, args, err = sqlBuilderV3.Delete(&stuStruct, uidEq100, usernameEqAlice).Gen(w)
	evalAnd()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct, uidEq100).Where(usernameEqAlice).Gen(w)
	evalAnd()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct).Where(uidEq100, usernameEqAlice).Gen(w)
	evalAnd()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct).Where(uidEq100).And(usernameEqAlice).Gen(w)
	evalAnd()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct).Where(stuMapUid).And(stuMapUsername).Gen(w)
	evalAnd()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct, eqCond).Gen(w)
	evalAnd()

	evalOr := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student WHERE (uid = ?) OR (username = ?)", sql)
		assert.EqualValues(t, []interface{}{uid, "Alice"}, args)
	}
	sql, args, err = sqlBuilderV3.Delete(&stuStruct).Where(uidEq100).Or(usernameEqAlice).Gen(w)
	evalOr()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct).Where(nil).Or(uidEq100, usernameEqAlice).Gen(w)
	evalOr()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct).Where(stuMapUid).Or(stuMapUsername).Gen(w)
	evalOr()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct).Or(eqCond).Gen(w)
	evalOr()

	evalEmpty := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilderV3.Delete(&stuStruct).Gen(w)
	evalEmpty()
}

func TestSQLStmt_Update(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	var uidEq100 = sqlBuilderV3.CondExpr("uid = ??", uid)
	var usernameEqAlice = sqlBuilderV3.CondExpr("username = ??", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET "+
			"uid = ?,username = ?,nickname = ?,email = ?,age = ?,enrolled = ?,gpa = ?,tokens = ?,comp = ?,create_time = ?,update_time = ? "+
			"WHERE uid = ?", sql)
		assert.EqualValues(t, append(stuStructArr, uid), args)
	}

	sql, args, err = sqlBuilderV3.SQL().Update(&stuStruct, uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Update(&stuStruct, uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Update(&stuStruct, sqlBuilderV3.Map{"uid": uid}).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Update(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Update("student").Set(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Update().From(&stuStruct).Set(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Update().From("student").Set(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Update("student").Set(&stuStructNoTag).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Update().From(&stuStruct).
		Set("uid", sqlBuilderV3.Expr("??", uid)).
		Set("username", "??", "Alice").
		Set("nickname", "??", "Ali").
		Set("email", "ali@gmail.com").
		Set("age", "??", uint64(20)).
		Set("enrolled", "??", true).
		Set("gpa", "??", float64(3.5)).
		Set("tokens", "??", []string{"token1", "token2"}).
		Set("comp", "??", complex(10, 11)).
		Set("create_time", "??", ts.Unix()).
		Set("update_time", "??", ts.Unix()).
		Where("uid = ??", uid).Gen(w)
	evalStruct()

	evalMap := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilderV3.Update(stuMapUid).From(&stuStruct).Gen(w)
	evalMap()
	sql, args, err = sqlBuilderV3.Update(stuMapUid).From("student").Gen(w)
	evalMap()
	sql, args, err = sqlBuilderV3.Update().From(&stuStruct).Set(stuMapUid).Gen(w)
	evalMap()
	sql, args, err = sqlBuilderV3.Update().From("student").Set(stuMapUid).Gen(w)
	evalMap()
	sql, args, err = sqlBuilderV3.Update("student").Set(stuMapUid).Gen(w)
	evalMap()
	sql, args, err = sqlBuilderV3.Update("student").Set(stuMapExpr).Gen(w)
	evalMap()

	evalWhereAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET "+
			"uid = ?,username = ?,nickname = ?,email = ?,age = ?,enrolled = ?,gpa = ?,tokens = ?,comp = ?,create_time = ?,update_time = ? "+
			"WHERE (uid = ?) AND (username = ?)", sql)
		assert.EqualValues(t, append(stuStructArr, uid, "Alice"), args)
	}

	sql, args, err = sqlBuilderV3.Update(&stuStruct, eqCond).Gen(w)
	evalWhereAnd()
	sql, args, err = sqlBuilderV3.Update(&stuStruct).Where(eqCond).Gen(w)
	evalWhereAnd()
	sql, args, err = sqlBuilderV3.Update(&stuStruct, uidEq100, usernameEqAlice).Gen(w)
	evalWhereAnd()
	sql, args, err = sqlBuilderV3.Update(&stuStruct, uidEq100).Where(usernameEqAlice).Gen(w)
	evalWhereAnd()
	sql, args, err = sqlBuilderV3.Update(&stuStruct).Where(uidEq100, usernameEqAlice).Gen(w)
	evalWhereAnd()
	sql, args, err = sqlBuilderV3.Update(&stuStruct).Where(uidEq100).And(usernameEqAlice).Gen(w)
	evalWhereAnd()

	evalWhereOr := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET "+
			"uid = ?,username = ?,nickname = ?,email = ?,age = ?,enrolled = ?,gpa = ?,tokens = ?,comp = ?,create_time = ?,update_time = ? "+
			"WHERE (uid = ?) OR (username = ?)", sql)
		assert.EqualValues(t, append(stuStructArr, uid, "Alice"), args)
	}

	sql, args, err = sqlBuilderV3.Update(&stuStruct).Or(uidEq100, usernameEqAlice).Gen(w)
	evalWhereOr()

	evalIncr := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET uid = uid + ? WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{10, uid}, args)
	}
	sql, args, err = sqlBuilderV3.Update("student").Incr("uid", 10).Where(uidEq100).Gen(w)
	evalIncr()

	evalDecr := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET uid = uid - ? WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{10, uid}, args)
	}
	sql, args, err = sqlBuilderV3.Update("student").Decr("uid", 10).Where(uidEq100).Gen(w)
	evalDecr()
}

func TestSQLStmt_Select(t *testing.T) {
	var sql string
	var args []interface{}
	var err error
	var uidEq100 = sqlBuilderV3.Expr("uid = ??", 100)
	var usernameEqAlice = sqlBuilderV3.Expr("username = ??", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{100}, args)
	}

	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct, uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select(&stuStruct, uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select([]string{"uid", "username", "nickname", "email", "age", "enrolled", "gpa", "tokens", "comp", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select(sqlBuilderV3.Columns{"uid", "username", "nickname", "email", "age", "enrolled", "gpa", "tokens", "comp", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select().SelectColumns(&stuStruct).From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select().SelectColumns([]string{"uid", "username", "nickname", "email", "age", "enrolled", "gpa", "tokens", "comp", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select().SelectColumns(sqlBuilderV3.Columns{"uid", "username", "nickname", "email", "age", "enrolled", "gpa", "tokens", "comp", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select().SelectColumns("uid", "username", "nickname", "email", "age", "enrolled", "gpa", "tokens", "comp", "create_time", "update_time").From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student WHERE (uid = ?) AND (username = ?)", sql)
		assert.EqualValues(t, []interface{}{100, "Alice"}, args)
	}

	sql, args, err = sqlBuilderV3.Select(&stuStruct, uidEq100, usernameEqAlice).Gen(w)
	evalAnd()
	sql, args, err = sqlBuilderV3.Select(&stuStruct, uidEq100).Where(usernameEqAlice).Gen(w)
	evalAnd()
	sql, args, err = sqlBuilderV3.Select(&stuStruct).Where(uidEq100, usernameEqAlice).Gen(w)
	evalAnd()
	sql, args, err = sqlBuilderV3.Select(&stuStruct).Where(uidEq100).And(usernameEqAlice).Gen(w)
	evalAnd()

	evalAny := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Gen(w)
	evalAny()
	sql, args, err = sqlBuilderV3.Select().From("student").Gen(w)
	evalAny()
	sql, args, err = sqlBuilderV3.Select("student").Gen(w)
	evalAny()

	evalAS := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student AS S", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct, "S").Gen(w)
	evalAS()

	evalJoin := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student,student", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).From(&stuStruct).Gen(w)
	evalJoin()
}

// TestSQLStmt_OrderBy test sqlStmt.OrderBy, sqlStmt.Asc, and sqlStmt.Desc
func TestSQLStmt_OrderBy(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	eval := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student ORDER BY uid ASC, username DESC", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Asc("uid").Desc("username").Gen(w)
	eval()
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).OrderBy().Asc("uid").Desc("username").Gen(w)
	eval()
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).OrderBy("uid ASC", "username DESC").Gen(w)
	eval()
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).OrderBy("uid ASC").OrderBy("username DESC").Gen(w)
	eval()
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).OrderBy([]string{"uid ASC", "username DESC"}...).Gen(w)
	eval()
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).OrderBy(sqlBuilderV3.Columns{"uid ASC", "username DESC"}...).Gen(w)
	eval()

	evalASC := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student ORDER BY uid ASC, username ASC", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Asc("uid").Asc("username").Gen(w)
	evalASC()
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Asc().Asc("uid").Asc("username").Gen(w)
	evalASC()
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Asc("uid", "username").Gen(w)
	evalASC()

	evalDESC := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student ORDER BY uid DESC, username DESC", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Desc("uid").Desc("username").Gen(w)
	evalDESC()
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Desc().Desc("uid").Desc("username").Gen(w)
	evalDESC()
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Desc("uid", "username").Gen(w)
	evalDESC()
}

// TestSQLStmt_Limit test sqlStmt.Limit
func TestSQLStmt_Limit(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Limit(10, 5).Gen(w)
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM student LIMIT 10 OFFSET 5", sql)
	assert.EqualValues(t, []interface{}{}, args)

	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Limit(10).Gen(w)
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM student LIMIT 10", sql)
	assert.EqualValues(t, []interface{}{}, args)

	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Limit(10, 0).Gen(w)
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM student LIMIT 10", sql)
	assert.EqualValues(t, []interface{}{}, args)
}

// TestSQLStmt_GroupBy_Having test sqlStmt.GroupBy and sqlStmt.Havng
func TestSQLStmt_GroupBy_Having(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	var uidGe100 = sqlBuilderV3.Expr("COUNT(uid)>??", uid)

	eval := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student GROUP BY username HAVING COUNT(uid)>?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).Gen(w)
	eval()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student GROUP BY username HAVING (COUNT(uid)>?) AND (COUNT(uid)>?)", sql)
		assert.EqualValues(t, []interface{}{uid, uid}, args)
	}

	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100, uidGe100).Gen(w)
	evalAnd()
	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).HavingAnd(uidGe100).Gen(w)
	evalAnd()
	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct).GroupBy("username").HavingAnd(uidGe100, uidGe100).Gen(w)
	evalAnd()

	evalOr := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student GROUP BY username HAVING (COUNT(uid)>?) OR (COUNT(uid)>?)", sql)
		assert.EqualValues(t, []interface{}{uid, uid}, args)
	}

	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).HavingOr(uidGe100).Gen(w)
	evalOr()
	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct).GroupBy("username").HavingOr(uidGe100, uidGe100).Gen(w)
	evalOr()

	evalComplex := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student GROUP BY username, nickname HAVING (COUNT(uid)>?) AND (COUNT(uid)>?)", sql)
		assert.EqualValues(t, []interface{}{uid, uid}, args)
	}

	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct).GroupBy().GroupBy("username").GroupBy("nickname").Having(uidGe100).HavingAnd(uidGe100).Gen(w)
	evalComplex()
}

func TestSQLStmt_Error(t *testing.T) {
	var err error

	// missing the table name
	_, _, err = sqlBuilderV3.Insert().Gen(w)
	assert.EqualError(t, err, sqlBuilderV3.ErrNoTableName.Error())
	_, _, err = sqlBuilderV3.Delete().Gen(w)
	assert.EqualError(t, err, sqlBuilderV3.ErrNoTableName.Error())
	_, _, err = sqlBuilderV3.Update().Gen(w)
	assert.EqualError(t, err, sqlBuilderV3.ErrNoTableName.Error())
	_, _, err = sqlBuilderV3.Select().Gen(w)
	assert.EqualError(t, err, sqlBuilderV3.ErrNoTableName.Error())

	// insert errors
	_, _, err = sqlBuilderV3.Insert().IntoTable("student").Gen(w)
	assert.EqualError(t, err, sqlBuilderV3.ErrNoValueToInsert.Error())
	_, _, err = sqlBuilderV3.Insert().IntoTable("student").IntoColumns(&stuStruct).Gen(w)
	assert.EqualError(t, err, sqlBuilderV3.ErrNoValueToInsert.Error())

	// limit error
	_, _, err = sqlBuilderV3.Select().From(&stuStruct).Limit(-1, 0).Gen(w)
	assert.EqualError(t, err, sqlBuilderV3.ErrInvalidLimitation.Error())
	_, _, err = sqlBuilderV3.Select().From(&stuStruct).Limit(10, -1).Gen(w)
	assert.EqualError(t, err, sqlBuilderV3.ErrInvalidLimitation.Error())
}

func TestSQLStmt_Nest(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	// test insert select
	evalInsertSelect := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time) (SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student WHERE uid = ?)", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	selectStmt := sqlBuilderV3.Select(&stuStruct).Where(stuMapUid)
	sql, args, err = sqlBuilderV3.Insert(&stuStruct).Values(selectStmt).Gen(w)
	evalInsertSelect()

	evalJoin := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM (SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student WHERE uid = ?) AS S1,(SELECT uid,username,nickname,email,age,enrolled,gpa,tokens,comp,create_time,update_time FROM student WHERE uid = ?) AS S2", sql)
		assert.EqualValues(t, []interface{}{uid, uid}, args)
	}
	selectStmt2 := sqlBuilderV3.Select(&stuStruct).Where(stuMapUid)
	sql, args, err = sqlBuilderV3.Select().From(selectStmt, "S1").From(selectStmt2, "S2").Gen(w)
	evalJoin()
}

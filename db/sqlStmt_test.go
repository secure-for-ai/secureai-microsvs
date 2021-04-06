package db_test

import (
	"github.com/secure-for-ai/secureai-microsvs/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

type student struct {
	Uid      int    `db:"uid"`
	Username string `db:"username"`
}

func (s student) GetTableName() string {
	return "student"
}

var (
	stuStruct = student{
		100,
		"Alice",
	}
	stuList = []student{
		stuStruct,
		stuStruct,
	}
	eqCond = db.Map{
		"uid":      100,
		"username": "Alice",
	}
	stuMap = db.Map{
		"uid": 100,
	}
	stuVal = db.Map{
		"uid":      100,
		"username": "Alice",
	}
)

func TestSQLStmt_Insert(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	evalSingle := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username) VALUES (?,?)", sql)
		assert.EqualValues(t, []interface{}{100, "Alice"}, args)
	}

	sql, args, err = db.SQL().Insert(&stuStruct).Gen()
	evalSingle()
	sql, args, err = db.Insert().IntoTable(&stuStruct).Values(&stuStruct).Gen()
	evalSingle()
	sql, args, err = db.Insert().IntoTable("student").Values(stuVal).Gen()
	evalSingle()
	sql, args, err = db.Insert().IntoTable("student").ValuesBulk([]interface{}{stuVal}).Gen()
	evalSingle()
	sql, args, err = db.Insert().IntoTable("student").ValuesBulk([]interface{}{stuStruct}).Gen()
	evalSingle()

	evalBulk := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username) VALUES (?,?)", sql)
		assert.EqualValues(t, []interface{}{[]interface{}{100, "Alice"}, []interface{}{100, "Alice"}}, args)
	}
	sql, args, err = db.Insert(&stuStruct, &stuStruct).Gen()
	evalBulk()
	sql, args, err = db.InsertBulk(stuList).Gen()
	evalBulk()
	sql, args, err = db.InsertBulk(&stuList).Gen()
	evalBulk()
	sql, args, err = db.InsertBulk([]interface{}{stuStruct, stuStruct}).Gen()
	evalBulk()
	sql, args, err = db.InsertBulk(&[]interface{}{stuStruct, stuStruct}).Gen()
	evalBulk()
	sql, args, err = db.InsertBulk([]interface{}{&stuStruct, &stuStruct}).Gen()
	evalBulk()
	sql, args, err = db.InsertBulk(&[]interface{}{&stuStruct, &stuStruct}).Gen()
	evalBulk()
}

func TestSQLStmt_Delete(t *testing.T) {
	var sql string
	var args []interface{}
	var err error
	var uidEq100 = db.Expr("uid=?", 100)
	var usernameEqAlice = db.Expr("username=?", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student WHERE uid=?", sql)
		assert.EqualValues(t, []interface{}{100}, args)
	}
	sql, args, err = db.SQL().Delete(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Delete(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Delete(&stuStruct, stuMap).Gen()
	evalStruct()
	sql, args, err = db.Delete(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Delete().From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Delete("student", uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Delete("student").Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Delete().From("student").Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Delete().From("student").Where("uid=?", 100).Gen()
	evalStruct()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student WHERE (uid=?) AND (username=?)", sql)
		assert.EqualValues(t, []interface{}{100, "Alice"}, args)
	}
	sql, args, err = db.Delete(&stuStruct, uidEq100, usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = db.Delete(&stuStruct, uidEq100).Where(usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = db.Delete(&stuStruct).Where(uidEq100, usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = db.Delete(&stuStruct).Where(uidEq100).And(usernameEqAlice).Gen()
	evalAnd()

	evalEmpty := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = db.Delete(&stuStruct).Gen()
	evalEmpty()

}

func TestSQLStmt_Update(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	var uidEq100 = db.Expr("uid=?", 100)
	var usernameEqAlice = db.Expr("username=?", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET uid=?,username=? WHERE uid=?", sql)
		assert.EqualValues(t, []interface{}{100, "Alice", 100}, args)
	}

	sql, args, err = db.SQL().Update(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Update(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Update(&stuStruct, db.Map{"uid": 100}).Gen()
	evalStruct()
	sql, args, err = db.Update(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Update("student").Set(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Update().From(&stuStruct).Set(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Update().From("student").Set(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.
		Update().From(&stuStruct).
		Set("uid", db.Expr("?", 100)).
		Set("username", "?", "Alice").
		Where("uid=?", 100).Gen()
	evalStruct()

	evalMap := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET uid=?", sql)
		assert.EqualValues(t, []interface{}{100}, args)
	}
	sql, args, err = db.Update(stuMap).From(&stuStruct).Gen()
	evalMap()
	sql, args, err = db.Update(stuMap).From("student").Gen()
	evalMap()
	sql, args, err = db.Update().From(&stuStruct).Set(stuMap).Gen()
	evalMap()
	sql, args, err = db.Update().From("student").Set(stuMap).Gen()
	evalMap()
	sql, args, err = db.Update("student").Set(stuMap).Gen()
	evalMap()

	evalWhere := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET uid=?,username=? WHERE (uid=?) AND (username=?)", sql)
		assert.EqualValues(t, []interface{}{100, "Alice", 100, "Alice"}, args)
	}

	sql, args, err = db.Update(&stuStruct, eqCond).Gen()
	evalWhere()
	sql, args, err = db.Update(&stuStruct).Where(eqCond).Gen()
	evalWhere()
	sql, args, err = db.Update(&stuStruct, uidEq100, usernameEqAlice).Gen()
	evalWhere()
	sql, args, err = db.Update(&stuStruct, uidEq100).Where(usernameEqAlice).Gen()
	evalWhere()
	sql, args, err = db.Update(&stuStruct).Where(uidEq100, usernameEqAlice).Gen()
	evalWhere()
	sql, args, err = db.Update(&stuStruct).Where(uidEq100).And(usernameEqAlice).Gen()
	evalWhere()
}

func TestSQLStmt_Select(t *testing.T) {
	var sql string
	var args []interface{}
	var err error
	var uidEq100 = db.Expr("uid=?", 100)
	var usernameEqAlice = db.Expr("username=?", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username FROM student WHERE uid=?", sql)
		assert.EqualValues(t, []interface{}{100}, args)
	}

	sql, args, err = db.SQL().Select(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Select(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Select(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Select([]string{"uid", "username"}).From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Select(db.Columns{"uid", "username"}).From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Select().SelectColumns(&stuStruct).From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Select().SelectColumns([]string{"uid", "username"}).From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = db.Select().SelectColumns(db.Columns{"uid", "username"}).From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username FROM student WHERE (uid=?) AND (username=?)", sql)
		assert.EqualValues(t, []interface{}{100, "Alice"}, args)
	}

	sql, args, err = db.Select(&stuStruct, uidEq100, usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = db.Select(&stuStruct, uidEq100).Where(usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = db.Select(&stuStruct).Where(uidEq100, usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = db.Select(&stuStruct).Where(uidEq100).And(usernameEqAlice).Gen()
	evalAnd()

	evalAny := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = db.Select().From(&stuStruct).Gen()
	evalAny()
	sql, args, err = db.Select().From("student").Gen()
	evalAny()
	sql, args, err = db.Select("student").Gen()
	evalAny()
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
	sql, args, err = db.Select().From(&stuStruct).Asc("uid").Desc("username").Gen()
	eval()
	sql, args, err = db.Select().From(&stuStruct).OrderBy().Asc("uid").Desc("username").Gen()
	eval()
	sql, args, err = db.Select().From(&stuStruct).OrderBy("uid ASC", "username DESC").Gen()
	eval()
	sql, args, err = db.Select().From(&stuStruct).OrderBy([]string{"uid ASC", "username DESC"}...).Gen()
	eval()
	sql, args, err = db.Select().From(&stuStruct).OrderBy(db.Columns{"uid ASC", "username DESC"}...).Gen()
	eval()

	evalASC := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student ORDER BY uid ASC, username ASC", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = db.Select().From(&stuStruct).Asc("uid").Asc("username").Gen()
	evalASC()
	sql, args, err = db.Select().From(&stuStruct).Asc("uid", "username").Gen()
	evalASC()

	evalDESC := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student ORDER BY uid DESC, username DESC", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = db.Select().From(&stuStruct).Desc("uid").Desc("username").Gen()
	evalDESC()
	sql, args, err = db.Select().From(&stuStruct).Desc("uid", "username").Gen()
	evalDESC()
}

// TestSQLStmt_Limit test sqlStmt.Limit
func TestSQLStmt_Limit(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	sql, args, err = db.Select().From(&stuStruct).Limit(10, 5).Gen()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM student LIMIT 10 OFFSET 5", sql)
	assert.EqualValues(t, []interface{}{}, args)

	sql, args, err = db.Select().From(&stuStruct).Limit(10).Gen()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM student LIMIT 10", sql)
	assert.EqualValues(t, []interface{}{}, args)

	sql, args, err = db.Select().From(&stuStruct).Limit(10, 0).Gen()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM student LIMIT 10", sql)
	assert.EqualValues(t, []interface{}{}, args)

	_, _, err = db.Select().From(&stuStruct).Limit(-1, 0).Gen()
	assert.EqualError(t, err, db.ErrInvalidLimitation.Error())
	_, _, err = db.Select().From(&stuStruct).Limit(10, -1).Gen()
	assert.EqualError(t, err, db.ErrInvalidLimitation.Error())
}

// TestSQLStmt_GroupBy_Having test sqlStmt.GroupBy and sqlStmt.Havng
func TestSQLStmt_GroupBy_Having(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	var uidGe100 = db.Expr("COUNT(uid)>?", 100)
	//var usernameEqAlice = db.Expr("username=?", "Alice")

	eval := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username FROM student GROUP BY username HAVING COUNT(uid)>?", sql)
		assert.EqualValues(t, []interface{}{100}, args)
	}
	sql, args, err = db.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).Gen()
	eval()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username FROM student GROUP BY username HAVING (COUNT(uid)>?) AND (COUNT(uid)>?)", sql)
		assert.EqualValues(t, []interface{}{100, 100}, args)
	}

	sql, args, err = db.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100, uidGe100).Gen()
	evalAnd()
	sql, args, err = db.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).HavingAnd(uidGe100).Gen()
	evalAnd()
	sql, args, err = db.SQL().Select(&stuStruct).GroupBy("username").HavingAnd(uidGe100, uidGe100).Gen()
	evalAnd()

	evalOr := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username FROM student GROUP BY username HAVING (COUNT(uid)>?) OR (COUNT(uid)>?)", sql)
		assert.EqualValues(t, []interface{}{100, 100}, args)
	}

	sql, args, err = db.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).HavingOr(uidGe100).Gen()
	evalOr()
	sql, args, err = db.SQL().Select(&stuStruct).GroupBy("username").HavingOr(uidGe100, uidGe100).Gen()
	evalOr()
}

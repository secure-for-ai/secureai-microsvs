package sqlBuilder_test

import (
	"github.com/secure-for-ai/secureai-microsvs/db/sqlBuilder"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type student struct {
	Uid        int64  `db:"uid"`
	Username   string `db:"username"`
	Nickname   string `db:"nickname"`
	Email      string `db:"email"`
	CreateTime int64  `db:"create_time"`
	UpdateTime int64  `db:"update_time"`
}

func (s student) GetTableName() string {
	return "student"
}

func (s student) Size() int {
	return len(s.Username) + len(s.Nickname) + len(s.Nickname) + 24
}

var (
	uid       int64 = 100
	ts, _           = time.Parse(time.UnixDate, "Sat Mar  7 11:06:39 PST 2015")
	stuStruct       = student{
		uid,
		"Alice",
		"Ali",
		"ali@gmail.com",
		ts.Unix(),
		ts.Unix(),
	}
	stuStructArr       = []interface{}{uid, "Alice", "Ali", "ali@gmail.com", ts.Unix(), ts.Unix()}
	stuStructArrSorted = []interface{}{ts.Unix(), "ali@gmail.com", "Ali", uid, ts.Unix(), "Alice"}
	stuList            = []student{
		stuStruct,
		stuStruct,
	}
	eqCond = sqlBuilder.Map{
		"uid":      uid,
		"username": "Alice",
	}
	stuMap = sqlBuilder.Map{
		"uid": uid,
	}
	stuVal = sqlBuilder.Map{
		"uid":         uid,
		"username":    "Alice",
		"nickname":    "Ali",
		"email":       "ali@gmail.com",
		"create_time": ts.Unix(),
		"update_time": ts.Unix(),
	}
)

func TestSQLStmt_Insert(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	evalSingle := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username,nickname,email,create_time,update_time) VALUES (?,?,?,?,?,?)", sql)
		assert.EqualValues(t, stuStructArr, args)
	}

	evalSingleMap := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (create_time,email,nickname,uid,update_time,username) VALUES (?,?,?,?,?,?)", sql)
		assert.EqualValues(t, stuStructArrSorted, args)
	}

	sql, args, err = sqlBuilder.SQL().Insert(&stuStruct).Gen()
	evalSingle()
	sql, args, err = sqlBuilder.Insert().IntoTable(&stuStruct).Values(&stuStruct).Gen()
	evalSingle()
	sql, args, err = sqlBuilder.Insert().IntoTable("student").Values(stuVal).Gen()
	evalSingleMap()
	sql, args, err = sqlBuilder.Insert().IntoTable("student").ValuesBulk([]interface{}{stuVal}).Gen()
	evalSingleMap()
	sql, args, err = sqlBuilder.Insert().IntoTable("student").ValuesBulk([]interface{}{stuStruct}).Gen()
	evalSingle()

	evalBulk := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username,nickname,email,create_time,update_time) VALUES (?,?,?,?,?,?)", sql)
		assert.EqualValues(t, []interface{}{stuStructArr, stuStructArr}, args)
	}
	sql, args, err = sqlBuilder.Insert(&stuStruct, &stuStruct).Gen()
	evalBulk()
	sql, args, err = sqlBuilder.InsertBulk(stuList).Gen()
	evalBulk()
	sql, args, err = sqlBuilder.InsertBulk(&stuList).Gen()
	evalBulk()
	sql, args, err = sqlBuilder.InsertBulk([]interface{}{stuStruct, stuStruct}).Gen()
	evalBulk()
	sql, args, err = sqlBuilder.InsertBulk(&[]interface{}{stuStruct, stuStruct}).Gen()
	evalBulk()
	sql, args, err = sqlBuilder.InsertBulk([]interface{}{&stuStruct, &stuStruct}).Gen()
	evalBulk()
	sql, args, err = sqlBuilder.InsertBulk(&[]interface{}{&stuStruct, &stuStruct}).Gen()
	evalBulk()

	evalSelect1 := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student SELECT uid,username,nickname,email,create_time,update_time FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}

	sql, args, err = sqlBuilder.Insert().IntoTable("student").Select(stuStruct).Where(stuMap).Gen()
	evalSelect1()
	sql, args, err = sqlBuilder.Insert().IntoTable(&stuStruct).Select(&stuStruct).Where(stuMap).Gen()
	evalSelect1()

	evalSelect2 := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username,nickname,email,create_time,update_time) SELECT uid,username,nickname,email,create_time,update_time FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilder.Insert(&stuStruct).Select(&stuStruct).Where(stuMap).Gen()
	evalSelect2()
}

func TestSQLStmt_Delete(t *testing.T) {
	var sql string
	var args []interface{}
	var err error
	var uidEq100 = sqlBuilder.Expr("uid = ??", uid)
	var usernameEqAlice = sqlBuilder.Expr("username = ??", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilder.SQL().Delete(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Delete(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Delete(&stuStruct, stuMap).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Delete(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Delete().From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Delete("student", uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Delete("student").Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Delete().From("student").Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Delete().From("student").Where("uid = ??", uid).Gen()
	evalStruct()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student WHERE (uid = ?) AND (username = ?)", sql)
		assert.EqualValues(t, []interface{}{uid, "Alice"}, args)
	}
	sql, args, err = sqlBuilder.Delete(&stuStruct, uidEq100, usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = sqlBuilder.Delete(&stuStruct, uidEq100).Where(usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = sqlBuilder.Delete(&stuStruct).Where(uidEq100, usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = sqlBuilder.Delete(&stuStruct).Where(uidEq100).And(usernameEqAlice).Gen()
	evalAnd()

	evalEmpty := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilder.Delete(&stuStruct).Gen()
	evalEmpty()

}

func TestSQLStmt_Update(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	var uidEq100 = sqlBuilder.Expr("uid = ??", uid)
	var usernameEqAlice = sqlBuilder.Expr("username = ??", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET "+
			"uid = ?,username = ?,nickname = ?,email = ?,create_time = ?,update_time = ? "+
			"WHERE uid = ?", sql)
		assert.EqualValues(t, append(stuStructArr, uid), args)
	}

	sql, args, err = sqlBuilder.SQL().Update(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Update(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Update(&stuStruct, sqlBuilder.Map{"uid": uid}).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Update(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Update("student").Set(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Update().From(&stuStruct).Set(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Update().From("student").Set(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Update().From(&stuStruct).
		Set("uid", sqlBuilder.Expr("??", uid)).
		Set("username", "??", "Alice").
		Set("nickname", "??", "Ali").
		Set("email", "??", "ali@gmail.com").
		Set("create_time", "??", ts.Unix()).
		Set("update_time", "??", ts.Unix()).
		Where("uid = ??", uid).Gen()
	evalStruct()

	evalMap := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilder.Update(stuMap).From(&stuStruct).Gen()
	evalMap()
	sql, args, err = sqlBuilder.Update(stuMap).From("student").Gen()
	evalMap()
	sql, args, err = sqlBuilder.Update().From(&stuStruct).Set(stuMap).Gen()
	evalMap()
	sql, args, err = sqlBuilder.Update().From("student").Set(stuMap).Gen()
	evalMap()
	sql, args, err = sqlBuilder.Update("student").Set(stuMap).Gen()
	evalMap()

	evalWhere := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET "+
			"uid = ?,username = ?,nickname = ?,email = ?,create_time = ?,update_time = ? "+
			"WHERE (uid = ?) AND (username = ?)", sql)
		assert.EqualValues(t, append(stuStructArr, uid, "Alice"), args)
	}

	sql, args, err = sqlBuilder.Update(&stuStruct, eqCond).Gen()
	evalWhere()
	sql, args, err = sqlBuilder.Update(&stuStruct).Where(eqCond).Gen()
	evalWhere()
	sql, args, err = sqlBuilder.Update(&stuStruct, uidEq100, usernameEqAlice).Gen()
	evalWhere()
	sql, args, err = sqlBuilder.Update(&stuStruct, uidEq100).Where(usernameEqAlice).Gen()
	evalWhere()
	sql, args, err = sqlBuilder.Update(&stuStruct).Where(uidEq100, usernameEqAlice).Gen()
	evalWhere()
	sql, args, err = sqlBuilder.Update(&stuStruct).Where(uidEq100).And(usernameEqAlice).Gen()
	evalWhere()

	evalIncr := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET uid = uid + ? WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{10, uid}, args)
	}
	sql, args, err = sqlBuilder.Update("student").Incr("uid", 10).Where(uidEq100).Gen()
	evalIncr()

	evalDecr := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET uid = uid - ? WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{10, uid}, args)
	}
	sql, args, err = sqlBuilder.Update("student").Decr("uid", 10).Where(uidEq100).Gen()
	evalDecr()
}

func TestSQLStmt_Select(t *testing.T) {
	var sql string
	var args []interface{}
	var err error
	var uidEq100 = sqlBuilder.Expr("uid = ??", 100)
	var usernameEqAlice = sqlBuilder.Expr("username = ??", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,create_time,update_time FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{100}, args)
	}

	sql, args, err = sqlBuilder.SQL().Select(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Select(&stuStruct, uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Select(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Select([]string{"uid", "username", "nickname", "email", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Select(sqlBuilder.Columns{"uid", "username", "nickname", "email", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Select().SelectColumns(&stuStruct).From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Select().SelectColumns([]string{"uid", "username", "nickname", "email", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()
	sql, args, err = sqlBuilder.Select().SelectColumns(sqlBuilder.Columns{"uid", "username", "nickname", "email", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen()
	evalStruct()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,create_time,update_time FROM student WHERE (uid = ?) AND (username = ?)", sql)
		assert.EqualValues(t, []interface{}{100, "Alice"}, args)
	}

	sql, args, err = sqlBuilder.Select(&stuStruct, uidEq100, usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = sqlBuilder.Select(&stuStruct, uidEq100).Where(usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = sqlBuilder.Select(&stuStruct).Where(uidEq100, usernameEqAlice).Gen()
	evalAnd()
	sql, args, err = sqlBuilder.Select(&stuStruct).Where(uidEq100).And(usernameEqAlice).Gen()
	evalAnd()

	evalAny := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilder.Select().From(&stuStruct).Gen()
	evalAny()
	sql, args, err = sqlBuilder.Select().From("student").Gen()
	evalAny()
	sql, args, err = sqlBuilder.Select("student").Gen()
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
	sql, args, err = sqlBuilder.Select().From(&stuStruct).Asc("uid").Desc("username").Gen()
	eval()
	sql, args, err = sqlBuilder.Select().From(&stuStruct).OrderBy().Asc("uid").Desc("username").Gen()
	eval()
	sql, args, err = sqlBuilder.Select().From(&stuStruct).OrderBy("uid ASC", "username DESC").Gen()
	eval()
	sql, args, err = sqlBuilder.Select().From(&stuStruct).OrderBy([]string{"uid ASC", "username DESC"}...).Gen()
	eval()
	sql, args, err = sqlBuilder.Select().From(&stuStruct).OrderBy(sqlBuilder.Columns{"uid ASC", "username DESC"}...).Gen()
	eval()

	evalASC := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student ORDER BY uid ASC, username ASC", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilder.Select().From(&stuStruct).Asc("uid").Asc("username").Gen()
	evalASC()
	sql, args, err = sqlBuilder.Select().From(&stuStruct).Asc("uid", "username").Gen()
	evalASC()

	evalDESC := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student ORDER BY uid DESC, username DESC", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilder.Select().From(&stuStruct).Desc("uid").Desc("username").Gen()
	evalDESC()
	sql, args, err = sqlBuilder.Select().From(&stuStruct).Desc("uid", "username").Gen()
	evalDESC()
}

// TestSQLStmt_Limit test sqlStmt.Limit
func TestSQLStmt_Limit(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	sql, args, err = sqlBuilder.Select().From(&stuStruct).Limit(10, 5).Gen()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM student LIMIT 10 OFFSET 5", sql)
	assert.EqualValues(t, []interface{}{}, args)

	sql, args, err = sqlBuilder.Select().From(&stuStruct).Limit(10).Gen()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM student LIMIT 10", sql)
	assert.EqualValues(t, []interface{}{}, args)

	sql, args, err = sqlBuilder.Select().From(&stuStruct).Limit(10, 0).Gen()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM student LIMIT 10", sql)
	assert.EqualValues(t, []interface{}{}, args)

	_, _, err = sqlBuilder.Select().From(&stuStruct).Limit(-1, 0).Gen()
	assert.EqualError(t, err, sqlBuilder.ErrInvalidLimitation.Error())
	_, _, err = sqlBuilder.Select().From(&stuStruct).Limit(10, -1).Gen()
	assert.EqualError(t, err, sqlBuilder.ErrInvalidLimitation.Error())
}

// TestSQLStmt_GroupBy_Having test sqlStmt.GroupBy and sqlStmt.Havng
func TestSQLStmt_GroupBy_Having(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	var uidGe100 = sqlBuilder.Expr("COUNT(uid)>??", uid)
	//var usernameEqAlice = db.Expr("username = ?", "Alice")

	eval := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,create_time,update_time FROM student GROUP BY username HAVING COUNT(uid)>?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilder.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).Gen()
	eval()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,create_time,update_time FROM student GROUP BY username HAVING (COUNT(uid)>?) AND (COUNT(uid)>?)", sql)
		assert.EqualValues(t, []interface{}{uid, uid}, args)
	}

	sql, args, err = sqlBuilder.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100, uidGe100).Gen()
	evalAnd()
	sql, args, err = sqlBuilder.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).HavingAnd(uidGe100).Gen()
	evalAnd()
	sql, args, err = sqlBuilder.SQL().Select(&stuStruct).GroupBy("username").HavingAnd(uidGe100, uidGe100).Gen()
	evalAnd()

	evalOr := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,create_time,update_time FROM student GROUP BY username HAVING (COUNT(uid)>?) OR (COUNT(uid)>?)", sql)
		assert.EqualValues(t, []interface{}{uid, uid}, args)
	}

	sql, args, err = sqlBuilder.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).HavingOr(uidGe100).Gen()
	evalOr()
	sql, args, err = sqlBuilder.SQL().Select(&stuStruct).GroupBy("username").HavingOr(uidGe100, uidGe100).Gen()
	evalOr()
}

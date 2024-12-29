package sqlBuilderV3_test

import (
	"testing"
	"time"

	"github.com/secure-for-ai/secureai-microsvs/db/sqlBuilderV3"
	"github.com/stretchr/testify/assert"
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
	eqCond = sqlBuilderV3.Map{
		"uid":      uid,
		"username": "Alice",
	}
	stuMap = sqlBuilderV3.Map{
		"uid": uid,
	}
	stuVal = sqlBuilderV3.Map{
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
	w := sqlBuilderV3.NewWriter()

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

	sql, args, err = sqlBuilderV3.SQL().Insert(&stuStruct).Gen(w)
	evalSingle()
	sql, args, err = sqlBuilderV3.Insert().IntoTable(&stuStruct).Values(&stuStruct).Gen(w)
	evalSingle()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").Values(stuVal).Gen(w)
	evalSingleMap()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").ValuesBulk([]interface{}{stuVal}).Gen(w)
	evalSingleMap()
	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").ValuesBulk([]interface{}{stuStruct}).Gen(w)
	evalSingle()

	evalBulk := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username,nickname,email,create_time,update_time) VALUES (?,?,?,?,?,?)", sql)
		assert.EqualValues(t, []interface{}{stuStructArr, stuStructArr}, args)
	}
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

	evalSelect1 := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student SELECT uid,username,nickname,email,create_time,update_time FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}

	sql, args, err = sqlBuilderV3.Insert().IntoTable("student").Select(stuStruct).Where(stuMap).Gen(w)
	evalSelect1()
	sql, args, err = sqlBuilderV3.Insert().IntoTable(&stuStruct).Select(&stuStruct).Where(stuMap).Gen(w)
	evalSelect1()

	evalSelect2 := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "INSERT INTO student (uid,username,nickname,email,create_time,update_time) SELECT uid,username,nickname,email,create_time,update_time FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilderV3.Insert(&stuStruct).Select(&stuStruct).Where(stuMap).Gen(w)
	evalSelect2()
}

func TestSQLStmt_Delete(t *testing.T) {
	var sql string
	var args []interface{}
	var err error
	var uidEq100 = sqlBuilderV3.Expr("uid = ??", uid)
	var usernameEqAlice = sqlBuilderV3.Expr("username = ??", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "DELETE FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilderV3.SQL().Delete(&stuStruct, uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct, uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Delete(&stuStruct, stuMap).Gen(w)
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

	var uidEq100 = sqlBuilderV3.Expr("uid = ??", uid)
	var usernameEqAlice = sqlBuilderV3.Expr("username = ??", "Alice")

	evalStruct := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET "+
			"uid = ?,username = ?,nickname = ?,email = ?,create_time = ?,update_time = ? "+
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
	sql, args, err = sqlBuilderV3.Update().From(&stuStruct).
		Set("uid", sqlBuilderV3.Expr("??", uid)).
		Set("username", "??", "Alice").
		Set("nickname", "??", "Ali").
		Set("email", "??", "ali@gmail.com").
		Set("create_time", "??", ts.Unix()).
		Set("update_time", "??", ts.Unix()).
		Where("uid = ??", uid).Gen(w)
	evalStruct()

	evalMap := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET uid = ?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilderV3.Update(stuMap).From(&stuStruct).Gen(w)
	evalMap()
	sql, args, err = sqlBuilderV3.Update(stuMap).From("student").Gen(w)
	evalMap()
	sql, args, err = sqlBuilderV3.Update().From(&stuStruct).Set(stuMap).Gen(w)
	evalMap()
	sql, args, err = sqlBuilderV3.Update().From("student").Set(stuMap).Gen(w)
	evalMap()
	sql, args, err = sqlBuilderV3.Update("student").Set(stuMap).Gen(w)
	evalMap()

	evalWhere := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "UPDATE student SET "+
			"uid = ?,username = ?,nickname = ?,email = ?,create_time = ?,update_time = ? "+
			"WHERE (uid = ?) AND (username = ?)", sql)
		assert.EqualValues(t, append(stuStructArr, uid, "Alice"), args)
	}

	sql, args, err = sqlBuilderV3.Update(&stuStruct, eqCond).Gen(w)
	evalWhere()
	sql, args, err = sqlBuilderV3.Update(&stuStruct).Where(eqCond).Gen(w)
	evalWhere()
	sql, args, err = sqlBuilderV3.Update(&stuStruct, uidEq100, usernameEqAlice).Gen(w)
	evalWhere()
	sql, args, err = sqlBuilderV3.Update(&stuStruct, uidEq100).Where(usernameEqAlice).Gen(w)
	evalWhere()
	sql, args, err = sqlBuilderV3.Update(&stuStruct).Where(uidEq100, usernameEqAlice).Gen(w)
	evalWhere()
	sql, args, err = sqlBuilderV3.Update(&stuStruct).Where(uidEq100).And(usernameEqAlice).Gen(w)
	evalWhere()

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
		assert.EqualValues(t, "SELECT uid,username,nickname,email,create_time,update_time FROM student WHERE uid = ?", sql)
		assert.EqualValues(t, []interface{}{100}, args)
	}

	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct, uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select(&stuStruct, uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select([]string{"uid", "username", "nickname", "email", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select(sqlBuilderV3.Columns{"uid", "username", "nickname", "email", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select().SelectColumns(&stuStruct).From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select().SelectColumns([]string{"uid", "username", "nickname", "email", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()
	sql, args, err = sqlBuilderV3.Select().SelectColumns(sqlBuilderV3.Columns{"uid", "username", "nickname", "email", "create_time", "update_time"}).From(&stuStruct).Where(uidEq100).Gen(w)
	evalStruct()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,create_time,update_time FROM student WHERE (uid = ?) AND (username = ?)", sql)
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
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Asc("uid", "username").Gen(w)
	evalASC()

	evalDESC := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM student ORDER BY uid DESC, username DESC", sql)
		assert.EqualValues(t, []interface{}{}, args)
	}
	sql, args, err = sqlBuilderV3.Select().From(&stuStruct).Desc("uid").Desc("username").Gen(w)
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

	_, _, err = sqlBuilderV3.Select().From(&stuStruct).Limit(-1, 0).Gen(w)
	assert.EqualError(t, err, sqlBuilderV3.ErrInvalidLimitation.Error())
	_, _, err = sqlBuilderV3.Select().From(&stuStruct).Limit(10, -1).Gen(w)
	assert.EqualError(t, err, sqlBuilderV3.ErrInvalidLimitation.Error())
}

// TestSQLStmt_GroupBy_Having test sqlStmt.GroupBy and sqlStmt.Havng
func TestSQLStmt_GroupBy_Having(t *testing.T) {
	var sql string
	var args []interface{}
	var err error

	var uidGe100 = sqlBuilderV3.Expr("COUNT(uid)>??", uid)
	//var usernameEqAlice = db.Expr("username = ?", "Alice")

	eval := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,create_time,update_time FROM student GROUP BY username HAVING COUNT(uid)>?", sql)
		assert.EqualValues(t, []interface{}{uid}, args)
	}
	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).Gen(w)
	eval()

	evalAnd := func() {
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT uid,username,nickname,email,create_time,update_time FROM student GROUP BY username HAVING (COUNT(uid)>?) AND (COUNT(uid)>?)", sql)
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
		assert.EqualValues(t, "SELECT uid,username,nickname,email,create_time,update_time FROM student GROUP BY username HAVING (COUNT(uid)>?) OR (COUNT(uid)>?)", sql)
		assert.EqualValues(t, []interface{}{uid, uid}, args)
	}

	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct).GroupBy("username").Having(uidGe100).HavingOr(uidGe100).Gen(w)
	evalOr()
	sql, args, err = sqlBuilderV3.SQL().Select(&stuStruct).GroupBy("username").HavingOr(uidGe100, uidGe100).Gen(w)
	evalOr()
}

func BenchmarkSQLStmt_Insert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		w := sqlBuilderV3.NewWriter()
		for pb.Next() {
			stmt := sqlBuilderV3.SQL().Insert(&stuStruct)
			stmt.Gen(w)
			stmt.Destroy()
		}
	})
}

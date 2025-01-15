package pgdb_test

import (
	"context"
	"fmt"
	"github.com/secure-for-ai/secureai-microsvs/db/pgdb"
	"github.com/secure-for-ai/secureai-microsvs/db/sqlBuilderV3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var client *pgdb.PGClient
var ts, _ = time.Parse(time.UnixDate, "Sat Mar  7 11:06:39 PST 2015")

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

func initPG() {
	var err error
	conf := pgdb.PGPoolConf{
		Host:   "postgres",
		Port:   "5432",
		DBName: "test",
		User:   "test",
		PW:     "password",
		Verbose: true,
	}

	client, err = pgdb.NewPGClient(conf)

	if err != nil {
		// fmt.Println(err)
		fmt.Println("cannot connect to postgres")
		os.Exit(1)
	}

	if conf.Verbose {
		fmt.Println("connect to postgres")
	}
}

func TestPGConn(t *testing.T) {
	initPG()
	defer client.Close()

	exUid := int64(10000)
	exStu := student{
		exUid,
		"Alice",
		"Ali",
		"ali@gmail.com",
		ts.Unix(),
		ts.Unix(),
	}
	exStuMap := map[string]interface{}{
		"uid":         exStu.Uid,
		"username":    exStu.Username,
		"nickname":    exStu.Nickname,
		"email":       exStu.Email,
		"create_time": exStu.CreateTime,
		"update_time": exStu.UpdateTime,
	}
	exStuArr := [][]interface{}{{
		exStu.Uid,
		exStu.Username,
		exStu.Nickname,
		exStu.Email,
		exStu.CreateTime,
		exStu.UpdateTime,
	}}
	reStu := student{}
	var reStuSlice []student
	var resMaps = make([]map[string]interface{}, 0, 10)
	var resArr [][]interface{}
	ctx := context.Background()
	conn, err := client.GetConn(ctx)

	if err != nil {
		fmt.Println(err)
		panic("cannot acquire pg connection")
	}
	defer conn.Release()

	affectedRow, err := sqlBuilderV3.Insert(&exStu).ExecPG(conn, ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)

	affectedRow, err = sqlBuilderV3.Select(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).ExecPG(conn, ctx, &reStu)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStu, reStu)

	affectedRow, err = sqlBuilderV3.Select(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).Limit(10).ExecPG(conn, ctx, &reStuSlice)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, 1, len(reStuSlice))
	assert.EqualValues(t, exStu, reStuSlice[0])

	affectedRow, err = sqlBuilderV3.Select(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).Limit(10).ExecPG(conn, ctx, &resMaps)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, 1, len(resMaps))
	assert.EqualValues(t, exStuMap, resMaps[0])

	affectedRow, err = sqlBuilderV3.Select(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).ExecPG(conn, ctx, &resArr)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStuArr, resArr)

	now := ts.Unix()
	exStu.Username = "Bob"
	exStu.UpdateTime = now
	affectedRow, err = sqlBuilderV3.Update(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).ExecPG(conn, ctx, &resArr)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)

	affectedRow, err = sqlBuilderV3.Select(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).ExecPG(conn, ctx, &reStu)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStu, reStu)

	affectedRow, err = sqlBuilderV3.Delete(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).ExecPG(conn, ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
}

func TestPGTx(t *testing.T) {
	initPG()
	defer client.Close()

	exUid := int64(10000)
	exStu := student{
		exUid,
		"Alice",
		"Ali",
		"ali@gmail.com",
		ts.Unix(),
		ts.Unix(),
	}
	exStuMap := map[string]interface{}{
		"uid":         exStu.Uid,
		"username":    exStu.Username,
		"nickname":    exStu.Nickname,
		"email":       exStu.Email,
		"create_time": exStu.CreateTime,
		"update_time": exStu.UpdateTime,
	}
	exStuArr := [][]interface{}{{
		exStu.Uid,
		exStu.Username,
		exStu.Nickname,
		exStu.Email,
		exStu.CreateTime,
		exStu.UpdateTime,
	}}
	reStu := student{}
	var reStuSlice []student
	var resMaps = make([]map[string]interface{}, 0, 10)
	var resArr [][]interface{}
	ctx := context.Background()
	tx, err := client.Begin(ctx)

	if err != nil {
		panic("cannot start a transaction")
	}
	defer tx.RollBackDefer(ctx)

	affectedRow, err := sqlBuilderV3.Insert(&exStu).ExecPG(tx, ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)

	affectedRow, err = sqlBuilderV3.Select(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).ExecPG(tx, ctx, &reStu)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStu, reStu)

	affectedRow, err = sqlBuilderV3.Select(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).Limit(10).ExecPG(tx, ctx, &reStuSlice)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, 1, len(reStuSlice))
	assert.EqualValues(t, exStu, reStuSlice[0])

	affectedRow, err = sqlBuilderV3.Select(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).Limit(10).ExecPG(tx, ctx, &resMaps)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, 1, len(resMaps))
	assert.EqualValues(t, exStuMap, resMaps[0])

	affectedRow, err = sqlBuilderV3.Select(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).ExecPG(tx, ctx, &resArr)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStuArr, resArr)

	now := ts.Unix()
	exStu.Username = "Bob"
	exStu.UpdateTime = now
	affectedRow, err = sqlBuilderV3.Update(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).ExecPG(tx, ctx, &resArr)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)

	affectedRow, err = sqlBuilderV3.Select(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).ExecPG(tx, ctx, &reStu)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStu, reStu)

	affectedRow, err = sqlBuilderV3.Delete(&exStu).Where(sqlBuilderV3.Map{"uid": exUid}).ExecPG(tx, ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)

	err = tx.Commit(ctx)
	assert.NoError(t, err)
}

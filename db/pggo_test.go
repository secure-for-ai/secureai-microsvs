package db_test

import (
	"context"
	"fmt"
	"github.com/secure-for-ai/secureai-microsvs/db"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var client *db.PGClient

func initPG() {
	var err error
	conf := db.PGPoolConf{
		Host:   "localhost",
		Port:   "7000",
		DBName: "test",
		User:   "test",
		PW:     "password",
	}

	client, err = db.NewPGClient(conf)

	if err != nil {
		fmt.Println("cannot connect to postgres")
		os.Exit(1)
	}

	fmt.Println("connect to postgres")
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
		panic("cannot acquire pg connection")
	}
	defer conn.Release()

	affectedRow, err := db.Insert(&exStu).ExecPG(conn, ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)

	affectedRow, err = db.Select(&exStu).Where(db.Map{"uid": exUid}).ExecPG(conn, ctx, &reStu)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStu, reStu)

	affectedRow, err = db.Select(&exStu).Where(db.Map{"uid": exUid}).Limit(10).ExecPG(conn, ctx, &reStuSlice)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, 1, len(reStuSlice))
	assert.EqualValues(t, exStu, reStuSlice[0])

	affectedRow, err = db.Select(&exStu).Where(db.Map{"uid": exUid}).Limit(10).ExecPG(conn, ctx, &resMaps)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, 1, len(resMaps))
	assert.EqualValues(t, exStuMap, resMaps[0])

	affectedRow, err = db.Select(&exStu).Where(db.Map{"uid": exUid}).ExecPG(conn, ctx, &resArr)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStuArr, resArr)

	now := ts.Unix()
	exStu.Username = "Bob"
	exStu.UpdateTime = now
	affectedRow, err = db.Update(&exStu).Where(db.Map{"uid": exUid}).ExecPG(conn, ctx, &resArr)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)

	affectedRow, err = db.Select(&exStu).Where(db.Map{"uid": exUid}).ExecPG(conn, ctx, &reStu)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStu, reStu)

	affectedRow, err = db.Delete(&exStu).Where(db.Map{"uid": exUid}).ExecPG(conn, ctx)
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

	affectedRow, err := db.Insert(&exStu).ExecPG(tx, ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)

	affectedRow, err = db.Select(&exStu).Where(db.Map{"uid": exUid}).ExecPG(tx, ctx, &reStu)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStu, reStu)

	affectedRow, err = db.Select(&exStu).Where(db.Map{"uid": exUid}).Limit(10).ExecPG(tx, ctx, &reStuSlice)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, 1, len(reStuSlice))
	assert.EqualValues(t, exStu, reStuSlice[0])

	affectedRow, err = db.Select(&exStu).Where(db.Map{"uid": exUid}).Limit(10).ExecPG(tx, ctx, &resMaps)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, 1, len(resMaps))
	assert.EqualValues(t, exStuMap, resMaps[0])

	affectedRow, err = db.Select(&exStu).Where(db.Map{"uid": exUid}).ExecPG(tx, ctx, &resArr)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStuArr, resArr)

	now := ts.Unix()
	exStu.Username = "Bob"
	exStu.UpdateTime = now
	affectedRow, err = db.Update(&exStu).Where(db.Map{"uid": exUid}).ExecPG(tx, ctx, &resArr)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)

	affectedRow, err = db.Select(&exStu).Where(db.Map{"uid": exUid}).ExecPG(tx, ctx, &reStu)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)
	assert.EqualValues(t, exStu, reStu)

	affectedRow, err = db.Delete(&exStu).Where(db.Map{"uid": exUid}).ExecPG(tx, ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affectedRow)

	err = tx.Commit(ctx)
	assert.NoError(t, err)
}

package db_test

import (
	"context"
	"github.com/secure-for-ai/secureai-microsvs/db"
	"github.com/secure-for-ai/secureai-microsvs/snowflake"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func randStuSize() int {
	username, _ := util.GenerateRandomKey(15)
	nickname, _ := util.GenerateRandomKey(15)
	email, _ := util.GenerateRandomKey(15)
	exStu := student{
		0,
		util.Base64Encode(username),
		util.Base64Encode(nickname),
		util.Base64Encode(email) + "ali@gmail.com",
		ts.Unix(),
		ts.Unix(),
	}

	return exStu.Size()
}

func benchmarkInsert(b *testing.B, pageLen int) {
	conf := snowflake.NewNodeConf(1288834974657, 10, 12)
	var mu sync.Mutex
	nodeID := int64(0)
	initPG()
	defer client.Close()

	b.ReportAllocs()
	b.SetBytes(int64(randStuSize() * pageLen))
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		mu.Lock()
		nodeID++
		node, _ := snowflake.NewNode(nodeID, &conf)
		mu.Unlock()

		ctx := context.Background()
		conn, err := client.GetConn(ctx)

		if err != nil {
			panic("cannot acquire pg connection")
		}
		defer conn.Release()

		for pb.Next() {
			for i := 0; i < pageLen; i++ {
				username, _ := util.GenerateRandomKey(15)
				nickname, _ := util.GenerateRandomKey(15)
				email, _ := util.GenerateRandomKey(15)
				exStu := student{
					node.Generate().Int64(),
					util.Base64Encode(username),
					util.Base64Encode(nickname),
					util.Base64Encode(email) + "@gmail.com",
					ts.Unix(),
					ts.Unix(),
				}
				affectedRow, err := db.Insert(&exStu).ExecPG(conn, ctx)
				assert.NoError(b, err)
				assert.EqualValues(b, 1, affectedRow)
			}
		}
	})
}

func BenchmarkInsert1(b *testing.B)   { benchmarkInsert(b, 1) }
func BenchmarkInsert5(b *testing.B)   { benchmarkInsert(b, 5) }
func BenchmarkInsert10(b *testing.B)  { benchmarkInsert(b, 10) }
func BenchmarkInsert20(b *testing.B)  { benchmarkInsert(b, 20) }
func BenchmarkInsert50(b *testing.B)  { benchmarkInsert(b, 50) }
func BenchmarkInsert100(b *testing.B) { benchmarkInsert(b, 100) }

func benchmarkInsertBulk(b *testing.B, pageLen int) {
	conf := snowflake.NewNodeConf(1288834974657, 10, 12)
	var mu sync.Mutex
	nodeID := int64(0)
	initPG()
	defer client.Close()

	b.ReportAllocs()
	b.SetBytes(int64(randStuSize() * pageLen))
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		mu.Lock()
		nodeID++
		node, _ := snowflake.NewNode(nodeID, &conf)
		mu.Unlock()

		ctx := context.Background()
		conn, err := client.GetConn(ctx)

		if err != nil {
			panic("cannot acquire pg connection")
		}
		defer conn.Release()

		stus := make([]student, pageLen)

		for pb.Next() {
			for i := 0; i < pageLen; i++ {
				username, _ := util.GenerateRandomKey(15)
				nickname, _ := util.GenerateRandomKey(15)
				email, _ := util.GenerateRandomKey(15)
				stus[i] = student{
					node.Generate().Int64(),
					util.Base64Encode(username),
					util.Base64Encode(nickname),
					util.Base64Encode(email) + "@gmail.com",
					ts.Unix(),
					ts.Unix(),
				}
			}

			affectedRow, err := db.InsertBulk(&stus).ExecPG(conn, ctx)
			assert.NoError(b, err)
			assert.EqualValues(b, pageLen, affectedRow)
		}
	})
}

func BenchmarkInsertBulk1(b *testing.B)   { benchmarkInsertBulk(b, 1) }
func BenchmarkInsertBulk5(b *testing.B)   { benchmarkInsertBulk(b, 5) }
func BenchmarkInsertBulk10(b *testing.B)  { benchmarkInsertBulk(b, 10) }
func BenchmarkInsertBulk20(b *testing.B)  { benchmarkInsertBulk(b, 20) }
func BenchmarkInsertBulk50(b *testing.B)  { benchmarkInsertBulk(b, 50) }
func BenchmarkInsertBulk100(b *testing.B) { benchmarkInsertBulk(b, 100) }

func BenchmarkSelectOne(b *testing.B) {
	initPG()
	defer client.Close()

	b.ReportAllocs()
	b.SetBytes(int64(randStuSize()))
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()
		conn, err := client.GetConn(ctx)
		if err != nil {
			panic("cannot acquire pg connection")
		}
		defer conn.Release()

		stu := student{}
		uid := int64(10000)
		for pb.Next() {
			_, _ = db.Select(&stu).Where("uid > ??", uid).
				OrderBy("uid").Limit(1, 0).ExecPG(conn, ctx, &stu)
			uid = stu.Uid
		}
	})
}

func benchmarkSelectPage(b *testing.B, pageLen int) {
	initPG()
	defer client.Close()

	b.ReportAllocs()
	b.SetBytes(int64(randStuSize() * pageLen))
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()
		conn, err := client.GetConn(ctx)
		if err != nil {
			panic("cannot acquire pg connection")
		}
		defer conn.Release()

		stu := student{}
		var stus []student
		uid := int64(10000)
		for pb.Next() {
			_, _ = db.Select(&stu).Where("uid > ??", uid).
				OrderBy("uid").Limit(pageLen, 0).ExecPG(conn, ctx, &stus)
			uid = stus[len(stus)-1].Uid
			stus = nil
		}
	})
}

func BenchmarkSelectPage10(b *testing.B)  { benchmarkSelectPage(b, 10) }
func BenchmarkSelectPage20(b *testing.B)  { benchmarkSelectPage(b, 20) }
func BenchmarkSelectPage50(b *testing.B)  { benchmarkSelectPage(b, 50) }
func BenchmarkSelectPage100(b *testing.B) { benchmarkSelectPage(b, 100) }

func benchmarkUpdate(b *testing.B, pageLen int) {
	initPG()
	defer client.Close()

	b.ReportAllocs()
	b.SetBytes(int64(randStuSize() * pageLen))
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()
		conn, err := client.GetConn(ctx)
		if err != nil {
			panic("cannot acquire pg connection")
		}
		defer conn.Release()

		stu := student{}
		var stus []student
		uid := int64(10000)
		for pb.Next() {
			_, _ = db.Select(&stu).Where("uid > ??", uid).
				OrderBy("uid").Limit(pageLen, rand.Int()%32).ExecPG(conn, ctx, &stus)

			for i := 0; i < pageLen; i++ {
				uname, _ := util.GenerateRandomKey(15)
				stus[i].Username = util.Base64Encode(uname)
				stus[i].UpdateTime = time.Now().Unix()

				affectedRow, err := db.Update(&stus[i]).Where(db.Map{"uid": stus[i].Uid}).ExecPG(conn, ctx)
				assert.NoError(b, err)
				assert.EqualValues(b, 1, affectedRow)
			}

			uid = stus[len(stus)-1].Uid
			stus = nil
		}
	})
}

func BenchmarkUpdate1(b *testing.B)   { benchmarkUpdate(b, 1) }
func BenchmarkUpdate5(b *testing.B)   { benchmarkUpdate(b, 5) }
func BenchmarkUpdate10(b *testing.B)  { benchmarkUpdate(b, 10) }
func BenchmarkUpdate20(b *testing.B)  { benchmarkUpdate(b, 20) }
func BenchmarkUpdate50(b *testing.B)  { benchmarkUpdate(b, 50) }
func BenchmarkUpdate100(b *testing.B) { benchmarkUpdate(b, 100) }

func BenchmarkDelete(b *testing.B) {
	initPG()
	defer client.Close()

	b.ReportAllocs()
	b.SetBytes(int64(randStuSize()))
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()
		conn, err := client.GetConn(ctx)
		if err != nil {
			panic("cannot acquire pg connection")
		}
		defer conn.Release()

		stu := student{}
		uid := int64(10000)
		for pb.Next() {
			_, _ = db.Select(&stu).Where("uid > ??", uid).
				OrderBy("uid").Limit(1, rand.Int()%32).ExecPG(conn, ctx, &stu)
			_, _ = db.Delete(&stu).Where(db.Map{"uid": uid}).ExecPG(conn, ctx)
			uid = stu.Uid
		}
	})
}

func benchmarkTxInsert(b *testing.B, pageLen int) {
	conf := snowflake.NewNodeConf(1288834974657, 10, 12)
	var mu sync.Mutex
	nodeID := int64(0)
	initPG()
	defer client.Close()

	b.ReportAllocs()
	b.SetBytes(int64(randStuSize() * pageLen))
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		mu.Lock()
		nodeID++
		node, _ := snowflake.NewNode(nodeID, &conf)
		mu.Unlock()

		ctx := context.Background()

		for pb.Next() {
			tx, err := client.Begin(ctx)

			if err != nil {
				panic("cannot start tx")
			}

			for i := 0; i < pageLen; i++ {
				username, _ := util.GenerateRandomKey(15)
				nickname, _ := util.GenerateRandomKey(15)
				email, _ := util.GenerateRandomKey(15)
				exStu := student{
					node.Generate().Int64(),
					util.Base64Encode(username),
					util.Base64Encode(nickname),
					util.Base64Encode(email) + "@gmail.com",
					ts.Unix(),
					ts.Unix(),
				}
				affectedRow, err := db.Insert(&exStu).ExecPG(tx, ctx)
				assert.NoError(b, err)
				assert.EqualValues(b, 1, affectedRow)
			}

			_ = tx.Commit(ctx)
		}

	})
}

func BenchmarkTxInsert1(b *testing.B)   { benchmarkTxInsert(b, 1) }
func BenchmarkTxInsert5(b *testing.B)   { benchmarkTxInsert(b, 5) }
func BenchmarkTxInsert10(b *testing.B)  { benchmarkTxInsert(b, 10) }
func BenchmarkTxInsert20(b *testing.B)  { benchmarkTxInsert(b, 20) }
func BenchmarkTxInsert50(b *testing.B)  { benchmarkTxInsert(b, 50) }
func BenchmarkTxInsert100(b *testing.B) { benchmarkTxInsert(b, 100) }

func benchmarkTxUpdate(b *testing.B, pageLen int) {
	initPG()
	defer client.Close()

	b.ReportAllocs()
	b.SetBytes(int64(randStuSize() * pageLen))
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()

		stu := student{}
		var stus []student
		uid := int64(10000)
		for pb.Next() {
			tx, err := client.Begin(ctx)
			if err != nil {
				panic("cannot start tx")
			}
			defer tx.RollBackDefer(ctx)

			_, _ = db.Select(&stu).Where("uid > ??", uid).
				OrderBy("uid").Limit(pageLen, rand.Int()%32).ExecPG(tx, ctx, &stus)

			for i := 0; i < pageLen; i++ {
				uname, _ := util.GenerateRandomKey(15)
				stus[i].Username = util.Base64Encode(uname)
				stus[i].UpdateTime = time.Now().Unix()

				affectedRow, err := db.Update(&stus[i]).Where(db.Map{"uid": stus[i].Uid}).ExecPG(tx, ctx)
				assert.NoError(b, err)
				assert.EqualValues(b, 1, affectedRow)
			}

			uid = stus[len(stus)-1].Uid
			stus = nil
			_ = tx.Commit(ctx)
		}
	})
}

func BenchmarkTxUpdate1(b *testing.B)   { benchmarkTxUpdate(b, 1) }
func BenchmarkTxUpdate5(b *testing.B)   { benchmarkTxUpdate(b, 5) }
func BenchmarkTxUpdate10(b *testing.B)  { benchmarkTxUpdate(b, 10) }
func BenchmarkTxUpdate20(b *testing.B)  { benchmarkTxUpdate(b, 20) }
func BenchmarkTxUpdate50(b *testing.B)  { benchmarkTxUpdate(b, 50) }
func BenchmarkTxUpdate100(b *testing.B) { benchmarkTxUpdate(b, 100) }

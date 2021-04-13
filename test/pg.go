package main

import (
	"context"
	"fmt"
	"github.com/secure-for-ai/secureai-microsvs/db"
	"os"
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

type User struct {
	Uid      int64  `json:"uid" db:"uid"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func main() {
	initPG()
	defer client.Close()

	testQuery()
	testTransaction()
}

func testQuery() {
	fmt.Println("=============================")
	fmt.Println("======== Test Querys ========")
	fmt.Println("=============================")

	var affectedRows int64
	ctx := context.Background()
	conn, err := client.GetConn(ctx)

	if err != nil {
		panic("cannot acquire pg connection")
	}
	defer conn.Release()

	//fmt.Sprintf("SELECT uid, username, password FROM test.%s", pq.QuoteIdentifier("testUser"))
	var resultArray [][]interface{}
	affectedRows, err = conn.FindAllAsArray(ctx, "SELECT uid, username, password FROM test.test_user", &resultArray)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Array Find records:", affectedRows)
	fmt.Println("Array Scan:", resultArray)

	var resultMap []map[string]interface{}
	affectedRows, err = conn.FindAllAsMap(ctx, "SELECT uid, username, password FROM test.test_user", &resultMap)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Map Find records:", affectedRows)
	fmt.Println("Map Scan:", resultMap)

	var users []User
	affectedRows, err = conn.FindAll(ctx, "SELECT uid, username, password FROM test.test_user", &users)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("FindAll records:", affectedRows)
	fmt.Println("FindAll:", users)

	var user User
	err = conn.FindOne(ctx, "SELECT uid, username, password FROM test.test_user WHERE uid=$1", &user, 1)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("FindOne:", user)

	affectRows, err := conn.Insert(
		ctx, "INSERT INTO test.test_user (uid, username, password) VALUES ($1, $2, $3)",
		3, "hello world", "p@ssword")

	fmt.Println("Insert affected: ", affectRows)

	err = conn.FindOne(ctx, "SELECT uid, username, password FROM test.test_user WHERE uid=$1", &user, 3)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Check Insert with FindOne:", user)

	affectRows, err = conn.Update(
		ctx, "UPDATE test.test_user SET password=$1 WHERE uid=$2",
		"new_p@ssword", 3)

	fmt.Println("Update affected: ", affectRows)

	err = conn.FindOne(ctx, "SELECT uid, username, password FROM test.test_user WHERE uid=$1", &user, 3)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Check Update with FindOne:", user)

	affectRows, err = conn.Delete(
		ctx, "DELETE FROM test.test_user WHERE uid=$1",
		3)

	fmt.Println("Delete affected: ", affectRows)

	users = []User{}
	affectedRows, err = conn.FindAll(ctx, "SELECT uid, username, password FROM test.test_user WHERE uid=$1", &users, 3)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Check Deletion with FindAll:", users)

	user = User{}
	err = conn.FindOne(ctx, "SELECT uid, username, password FROM test.test_user WHERE uid=$1", &user, 3)

	if err != nil {
		fmt.Println("User uid=3 not found. Error", err.Error())
	}
}

func testTransaction() {
	fmt.Println("=============================")
	fmt.Println("====== Test Transaction =====")
	fmt.Println("=============================")

	ctx := context.Background()
	tx, err := client.Begin(ctx)

	if err != nil {
		panic("cannot start a transaction")
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	var resultArray [][]interface{}
	affectRows, err := tx.FindAllAsArray(ctx, "SELECT uid, username, password FROM test.test_user", &resultArray)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Array Find records:", affectRows)
	fmt.Println("Array Scan:", resultArray)

	var resultMap []map[string]interface{}
	affectRows, err = tx.FindAllAsMap(ctx, "SELECT uid, username, password FROM test.test_user", &resultMap)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Map Find records:", affectRows)
	fmt.Println("Map Scan:", resultMap)

	var users []User
	affectRows, err = tx.FindAll(ctx, "SELECT uid, username, password FROM test.test_user", &users)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("FindAll records:", affectRows)
	fmt.Println("FindAll:", users)

	var user User
	err = tx.FindOne(ctx, "SELECT uid, username, password FROM test.test_user WHERE uid=$1", &user, 1)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("FindOne:", user)

	affectRows, err = tx.Insert(
		ctx, "INSERT INTO test.test_user (uid, username, password) VALUES ($1, $2, $3)",
		3, "hello world", "p@ssword")

	fmt.Println("Insert affected: ", affectRows)

	err = tx.FindOne(ctx, "SELECT uid, username, password FROM test.test_user WHERE uid=$1", &user, 3)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Check Insert with FindOne:", user)

	affectRows, err = tx.Update(
		ctx, "UPDATE test.test_user SET password=$1 WHERE uid=$2",
		"new_p@ssword", 3)

	fmt.Println("Update affected: ", affectRows)

	err = tx.FindOne(ctx, "SELECT uid, username, password FROM test.test_user WHERE uid=$1", &user, 3)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Check Update with FindOne:", user)

	affectRows, err = tx.Delete(
		ctx, "DELETE FROM test.test_user WHERE uid=$1",
		3)

	fmt.Println("Delete affected: ", affectRows)

	users = []User{}
	affectRows, err = tx.FindAll(ctx, "SELECT uid, username, password FROM test.test_user WHERE uid=$1", &users, 3)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Check Deletion with FindAll:", users)

	user = User{}
	err = tx.FindOne(ctx, "SELECT uid, username, password FROM test.test_user WHERE uid=$1", &user, 3)

	if err != nil {
		fmt.Println("User uid=3 not found. Error", err.Error())
	}

	err = tx.Commit(ctx)

	if err != nil {
		fmt.Println("Commit Failed")
		panic(err.Error())
	}

	fmt.Println("Commit")
}

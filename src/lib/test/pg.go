package main

import (
	"context"
	"fmt"
	"os"
	"template2/lib/db"
)

var client *db.PGClient

func initPG() {
	var err error
	conf := db.PGPoolConf{
		Host:   "localhost",
		Port:   "7000",
		DBName: "postgres",
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
	Uid      int64  `json:"uid" pg:"uid"`
	Username string `json:"username" pg:"username"`
	Password string `json:"password" pg:"password"`
}

func main() {
	initPG()

	ctx := context.Background()
	conn, err := client.GetConn(ctx)

	if err != nil {
		panic("cannot acquire pg connection")
	}

	resultArray, err := conn.FindAllAsArray(ctx, "SELECT uid, username, password FROM public.test_user")

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Array Scan:", resultArray)

	resultMap, err := conn.FindAllAsMap(ctx, "SELECT uid, username, password FROM public.test_user")

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Map Scan:", resultMap)

	var users []User
	err = conn.FindAll(ctx, "SELECT uid, username, password FROM public.test_user", &users)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("FindAll:", users)

	var user User
	err = conn.FindOne(ctx, "SELECT uid, username, password FROM public.test_user WHERE uid=$1", &user, 1)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("FindOne:", user)

	affectRows, err := conn.InsertOne(
		ctx, "INSERT INTO public.test_user (uid, username, password) VALUES ($1, $2, $3)",
		3, "hello world", "p@ssword")

	fmt.Println("InsertOne affected: ", affectRows)

	err = conn.FindOne(ctx, "SELECT uid, username, password FROM public.test_user WHERE uid=$1", &user, 3)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Check Insert with FindOne:", user)

	affectRows, err = conn.UpdateOne(
		ctx, "UPdate public.test_user SET password=$1 WHERE uid=$2",
		"new_p@ssword", 3)

	fmt.Println("UpdateOne affected: ", affectRows)

	err = conn.FindOne(ctx, "SELECT uid, username, password FROM public.test_user WHERE uid=$1", &user, 3)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Check Update with FindOne:", user)

	affectRows, err = conn.DeleteOne(
		ctx, "DELETE FROM public.test_user WHERE uid=$1",
		3)

	fmt.Println("DeleteOne affected: ", affectRows)

	users = []User{}
	err = conn.FindAll(ctx, "SELECT uid, username, password FROM public.test_user WHERE uid=$1", &users, 3)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Check Deletion with FindAll:", users)

	// what happened if findone does not find anything
	err = conn.FindOne(ctx, "SELECT uid, username, password FROM public.test_user WHERE uid=$1", &user, 3)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Check Deletion with FindOne:", user)
	conn.Release()
}

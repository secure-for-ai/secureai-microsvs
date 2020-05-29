package config

import (
	"context"
	//"context"
	//"context"
	//"fmt"
	"log"
	"template2/lib/db"
)

type Config struct {
	db string
}

var Conf *Config
var MongoDBClient *db.MongoDBClient

func init() {
	log.Println("begin init all configs")
	Conf = &Config{}
	MongoDBClient = &db.MongoDBClient{}
	MongoDBClient1, err := db.NewMongoDB()
	MongoDBClient.Copy(MongoDBClient1)
	log.Println(MongoDBClient1)
	if err != nil {
		log.Fatal(err)
		//panic(nil)
	}

	err = MongoDBClient.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	MongoDBClient.UseDatabase("gtest")
	log.Println("Connected to MongoDB!")
	log.Println(MongoDBClient)

	log.Println("over init all configs")
}

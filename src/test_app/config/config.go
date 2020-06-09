package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"template2/lib/cache"
	"template2/lib/db"
)

type Config struct {
	MongoDB db.MongoDBConf
	Redis   cache.RedisConf
}

var Conf *Config
var MongoDBClient *db.MongoDBClient

var confPathPrefix = defaultConfPath("test_app/config")

func init() {
	log.Println("Begin init")

	initConf()
	log.Println(Conf)

	intiMongoDB()

	log.Println("Over init")
}

func defaultConfPath(dir string) string {
	wdPath, err := os.Getwd()

	if err != nil {
		log.Panic(err)
		return ""
	}

	s := path.Join(wdPath, dir)
	return s
}

func initConf() {
	log.Println("Begin init default config")

	Conf = &Config{}
	fileName := "default.json"

	if v, ok := os.LookupEnv("CONFIG_PATH"); ok {
		confPathPrefix = v
	}

	// read default config
	defaultConfFilePath := path.Join(confPathPrefix, fileName)
	data, err := ioutil.ReadFile(defaultConfFilePath)

	if err != nil {
		log.Println("config-initConf: read default.json error")
		log.Panic(err)
		return
	}

	err = json.Unmarshal(data, Conf)
	if err != nil {
		log.Println("config-initConf: unmarshal default.json error")
		log.Panic(err)
		return
	}

	// read env and config path
	if v, ok := os.LookupEnv("ENV"); ok {
		fileName = v + ".json"
	}

	if fileName != "default.json" {
		// read env config
		data, err = ioutil.ReadFile(path.Join(confPathPrefix, fileName))
		if err != nil {
			log.Println("config-initConf: read [env].json error")
			log.Panic(err)
			return
		}

		err = json.Unmarshal(data, Conf)
		if err != nil {
			log.Println("config-initConf: unmarshal [env].json error")
			log.Panic(err)
			return
		}
	}

	log.Println("Over init default config")
}

func intiMongoDB() {
	log.Println("Begin init mongoDB")

	client, err := db.NewMongoDB(Conf.MongoDB)

	if err != nil {
		log.Println("unable to init mongoDB")
		log.Panic(err)
		return
	}

	MongoDBClient = client

	log.Println("Over init mongoDB")
}

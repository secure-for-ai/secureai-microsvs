package config

import (
	"log"
)

type Config struct {
	db string
}

var Conf *Config

func init() {
	log.Println("begin init all configs")
	Conf = &Config{}
	log.Println("over init all configs")
}

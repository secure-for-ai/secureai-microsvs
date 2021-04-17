package db

type Schema int

const (
	Tag          = "db"
	Para         = "??"
	SchPG Schema = iota
	SchMYSQL
)

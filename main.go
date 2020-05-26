package main

import (
	"fmt"
	"net/http"
	"template2/test_app/gq"
	//"template2/test_app/config"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Worldkkk")
}

func main() {
	fmt.Println("hello world")
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/graphql", gq.Graphql)
	http.ListenAndServe(":6551", nil)
}

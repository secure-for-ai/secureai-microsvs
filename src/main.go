package main

import (
	"fmt"
	"net/http"
	"template2/test_app/graphql_api"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Worldkkk")
}

func main() {
	fmt.Println("hello world")
	//defer config.DB.Disconnect()
	//http.HandleFunc("/api/hello", helloWorld)
	http.HandleFunc("/api/graphql", graphql_api.Graphql)
	http.ListenAndServe(":6551", nil)
}

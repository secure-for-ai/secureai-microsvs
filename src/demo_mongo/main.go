package main

import (
	"fmt"
	"log"
	"net/http"
	"template2/demo_mongo/graphql"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Worldkkk")
}

func main() {
	log.Println("Start service")
	http.HandleFunc("/api/hello", helloWorld)
	http.HandleFunc("/api/graphql", graphql.Graphql)
	http.ListenAndServe(":6551", nil)
}

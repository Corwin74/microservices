package main

import (
	"fmt"
	"log"
	"net/http"
)


func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Chat messaging service say: Hello!")
}


func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}



func main() {
	http.HandleFunc("/", hello) 
	http.HandleFunc("/ping", pong)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

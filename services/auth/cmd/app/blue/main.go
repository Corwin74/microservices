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

	fmt.Fprintf(w, `		
		<!DOCTYPE html>
		<html lang="ru">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Приветствие</title>
			<style>
				.hello-text {
					color: blue;
				}
			</style>
		</head>
		<body>
			<h1 class="hello-text">Auth service say: Hello!</h1>
		</body>
		</html>`,
	)
}


func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}



func main() {
	http.HandleFunc("/", hello) 
	http.HandleFunc("/ping", pong)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

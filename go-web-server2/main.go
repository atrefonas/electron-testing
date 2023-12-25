package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested from web server 2: %s\n", r.URL.Path)
	})

	fmt.Println("Starting server at port 8091")
	if err := http.ListenAndServe(":8091", nil); err != nil {
		log.Fatal(err)
	}
}

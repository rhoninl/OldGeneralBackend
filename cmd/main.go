package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("111")
		fmt.Fprintf(w, "hello")
	})
	http.ListenAndServe(":8080", mux)
}

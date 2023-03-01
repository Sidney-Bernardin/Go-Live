package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello, World!")
		w.Write([]byte("Hello, World!"))
	})

	log.Fatal(http.ListenAndServe(":8020", nil))
}

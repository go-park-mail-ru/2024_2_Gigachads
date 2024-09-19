package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "hello, world!\n")
	}

	http.HandleFunc("/", helloHandler)
	log.Println("Server is running on port 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
	
}
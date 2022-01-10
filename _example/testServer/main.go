package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func info(w http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)

	log.Printf("body: [%s]", body)
	log.Printf("method: [%s]", req.Method)
}

func main() {
	http.HandleFunc("/", info)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println(err)
	}
}

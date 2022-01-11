package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func info(w http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)

	printIt := "\n"

	printIt += fmt.Sprintf("body: [%s]\n", body)
	printIt += fmt.Sprintf("method: [%s]\n", req.Method)
	printIt += fmt.Sprintf("headers: [%v]", req.Header)

	log.Print(printIt)

	_, _ = w.Write([]byte(printIt))
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

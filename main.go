package main

import (
	"log"
	"net/http"
	"temperature-backend/server"
)

func main() {
	mux := server.Setup()

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

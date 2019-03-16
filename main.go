package main

import (
	"log"
	"net/http"
	"temperature-backend/server"

)

func main() {
	server.Setup()

	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}



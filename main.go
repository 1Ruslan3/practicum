package main

import (
	"log"
	"net/http"
	"practicum/handlers"
)

func main() {

	handlers.InitDb()
	http.HandleFunc("/subjects", handlers.SubjectsHandler)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

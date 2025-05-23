package main

import (
	"log"
	"net/http"
	"practicum/DataBaseConnect"
	"practicum/handlers"
)

func main() {

	DataBaseConnect.InitDb()
	http.HandleFunc("/subjects", handlers.SubjectsHandler)
	http.HandleFunc("/programs", handlers.ProgramsHandler)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

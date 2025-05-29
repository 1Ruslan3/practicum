package main

import (
	"log"
	"net/http"
	"practicum/DataBaseConnect"
	"practicum/handlefunc"
)

func main() {

	DataBaseConnect.InitDb()
	handlefunc.SetupRoutes()
	DataBaseConnect.MigrateDb()
	DataBaseConnect.DataFiling()

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

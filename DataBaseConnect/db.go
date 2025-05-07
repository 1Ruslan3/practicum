package DataBaseConnect

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDb() {
	var err error
	Db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=practicum sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

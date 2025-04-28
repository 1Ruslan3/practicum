package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Subject struct {
	Subjects map[string]string `json:"subjects"`
}

var Db *sql.DB

func InitDb() {
	var err error
	Db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=practicum sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func SubjectsHandler(w http.ResponseWriter, r *http.Request) {

	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "ru"
	}

	var columnName string
	switch locale {
	case "en":
		columnName = "name_en"
	case "ru":
		columnName = "name_ru"
	default:
		http.Error(w, "Invalid locale parameter. Use 'ru' or 'en'", http.StatusBadRequest)
		return
	}

	query := "SELECT id_item, " + columnName + " FROM items"
	rows, err := Db.Query(query)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		log.Printf("Database query error: %v", err)
		return
	}
	defer rows.Close()

	subjectsMap := make(map[string]string)
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			log.Printf("Row scan error: %v", err)
			return
		}
		subjectsMap[id] = name
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating rows", http.StatusInternalServerError)
		log.Printf("Rows iteration error: %v", err)
		return
	}

	response := Subject{Subjects: subjectsMap}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding JSON: %v", err)
		return
	}

}

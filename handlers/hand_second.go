package handlers

import (
	"encoding/json"
	"net/http"
	"practicum/DataBaseConnect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Program struct {
	ID              int    `json:"id"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	FormatEducation string `json:"format_education"`
	Year            int    `json:"year"`
	ScoreBudget     int    `json:"score_budget"`
	ScorePaid       int    `json:"score_paid"`
	AvailablePlaces int    `json:"available_places"`
	EducationPrice  int    `json:"education_price"`
}
type SubjectScore struct {
	SubjectID string
	Score     int
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ProgramsHandler(w http.ResponseWriter, r *http.Request) {

	studentScores := map[string]int{}
	for key, vals := range r.URL.Query() {
		if len(vals) > 0 {
			score, err := strconv.Atoi(vals[0])
			if err == nil {
				studentScores[strings.ToLower(key)] = score
			}
		}
	}

	if len(studentScores) == 0 {
		http.Error(w, "No subjects or scores provided", http.StatusBadRequest)
		return
	}

	// Получаем список всех программ
	query := `
        SELECT 
            p.id, p.code, p.name, p.description, p.format_education,
            ps.year, ps.score_budget, ps.score_paid, ps.available_places, ps.education_price
        FROM programs p
        JOIN passing_scores ps ON ps.program_id = p.id
    `
	rows, err := DataBaseConnect.Db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Program

	for rows.Next() {
		var prog Program
		err := rows.Scan(
			&prog.ID, &prog.Code, &prog.Name, &prog.Description, &prog.FormatEducation,
			&prog.Year, &prog.ScoreBudget, &prog.ScorePaid, &prog.AvailablePlaces, &prog.EducationPrice,
		)
		if err != nil {
			continue
		}

		// Получаем обязательные предметы для программы
		subjQuery := `
            SELECT s.id
			FROM program_subject ps
			JOIN subjects s ON ps.subject_id = s.id
			WHERE ps.program_id = $1

        `
		subjRows, err := DataBaseConnect.Db.Query(subjQuery, prog.ID)
		if err != nil {
			continue
		}

		totalScore := 0
		meetsAll := true
		for subjRows.Next() {
			var subject string
			if err := subjRows.Scan(&subject); err != nil {
				continue
			}

			score, ok := studentScores[subject]
			if !ok {
				meetsAll = false
				break
			}
			totalScore += score
		}
		subjRows.Close()

		if meetsAll && totalScore >= prog.ScoreBudget {
			results = append(results, prog)
		}
	}

	// Ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

package handlers

import (
	"encoding/json"
	"log"
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

	query := `
		SELECT
		p.id, p.code, p.name, p.description, p.format_education,
		ps.score_budget, ps.score_paid, ps.available_places, ps.year, ps.education_price
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
			&prog.ScoreBudget, &prog.ScorePaid, &prog.AvailablePlaces, &prog.Year, &prog.EducationPrice,
		)
		if err != nil {
			log.Println("Error scanning program:", err)
			continue
		}

		groupQuery := `
			SELECT sg.id, sg.group_type
			FROM program_subject_groups psg
			JOIN subject_groups sg ON sg.id = psg.group_id
			WHERE psg.program_id = $1
		`
		groupRows, err := DataBaseConnect.Db.Query(groupQuery, prog.ID)
		if err != nil {
			log.Println("Error querying subject groups:", err)
			continue
		}

		totalScore := 0
		meetsAll := true

		for groupRows.Next() {
			var groupID int
			var groupType string
			if err := groupRows.Scan(&groupID, &groupType); err != nil {
				continue
			}

			subjectsQuery := `
				SELECT s.id
				FROM subject_group_items sgi
				JOIN subjects s ON sgi.subject_id = s.id
				WHERE sgi.group_id = $1
			`
			subjectRows, err := DataBaseConnect.Db.Query(subjectsQuery, groupID)
			if err != nil {
				continue
			}

			found := false
			maxScore := 0

			for subjectRows.Next() {
				var subj string
				if err := subjectRows.Scan(&subj); err != nil {
					continue
				}
				if score, ok := studentScores[subj]; ok {
					found = true
					if groupType == "required" {
						totalScore += score
					}
					if score > maxScore {
						maxScore = score
					}
				}
			}
			subjectRows.Close()

			if groupType == "optional" && found {
				totalScore += maxScore
			}
			if groupType == "required" && !found {
				meetsAll = false
				break
			}
		}
		groupRows.Close()

		if meetsAll && totalScore >= prog.ScoreBudget {
			results = append(results, prog)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

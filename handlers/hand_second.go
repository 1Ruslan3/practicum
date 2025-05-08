package handlers

import (
	"encoding/json"
	"fmt"
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

	if r.Method != http.MethodGet {
		log.Printf("Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Only GET method is allowed"})
		return
	}

	var subjectScores []SubjectScore
	for key, values := range r.URL.Query() {
		if len(values) == 0 {
			continue
		}
		score, err := strconv.Atoi(values[0])
		if err != nil || score < 0 || score > 100 {
			log.Printf("Invalid score for subject %s: %v", key, err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Invalid score for subject %s", key)})
			return
		}
		subjectScores = append(subjectScores, SubjectScore{SubjectID: key, Score: score})
	}

	if len(subjectScores) == 0 {
		log.Printf("No subjects provided")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No subjects provided"})
		return
	}

	var invalidSubjects []string
	for _, ss := range subjectScores {
		var count int
		err := DataBaseConnect.Db.QueryRow("SELECT COUNT(*) FROM subjects WHERE id = $1", ss.SubjectID).Scan(&count)
		if err != nil {
			log.Printf("Failed to check subject %s: %v", ss.SubjectID, err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to validate subjects"})
			return
		}
		if count == 0 {
			invalidSubjects = append(invalidSubjects, ss.SubjectID)
		}
	}
	if len(invalidSubjects) > 0 {
		log.Printf("Invalid subjects: %v", invalidSubjects)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Invalid subjects: %v", invalidSubjects)})
		return
	}

	var subjectIDs []string
	var params []interface{}
	var scoreConditions []string
	var totalScore int
	paramIndex := 1
	const year = 2024

	for _, ss := range subjectScores {
		subjectIDs = append(subjectIDs, fmt.Sprintf("$%d", paramIndex))
		params = append(params, ss.SubjectID)
		paramIndex++
		scoreConditions = append(scoreConditions, fmt.Sprintf("(psu.subject_id = $%d AND ss.score <= $%d)", paramIndex-1, paramIndex))
		params = append(params, ss.Score)
		paramIndex++
		totalScore += ss.Score
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT
			p.id,
			p.code,
			p.name,
			p.description,
			p.format_education,
			ps.year,
			ps.score_budget,
			ps.score_paid,
			ps.available_places,
			ps.education_price
		FROM
			programs p
			JOIN program_subject psu ON p.id = psu.program_id
			JOIN passing_scores ps ON p.id = ps.program_id
			JOIN subject_scores ss ON psu.subject_id = ss.subject_id
		WHERE
			ps.year = %d
			AND ss.year = %d
			AND (
				p.id IN (
					SELECT psu2.program_id
					FROM program_subject psu2
					WHERE psu2.subject_id IN (%s)
					GROUP BY psu2.program_id
					HAVING COUNT(DISTINCT psu2.subject_id) = %d
				)
			)
			AND (
				%s
			)
			AND (
				(SELECT SUM(ss2.score)
				 FROM subject_scores ss2
				 JOIN program_subject psu3 ON ss2.subject_id = psu3.subject_id
				 WHERE psu3.program_id = p.id AND ss2.year = %d) <= $%d
			)
		ORDER BY
			p.id;
	`, year, year, strings.Join(subjectIDs, ","), len(subjectScores), strings.Join(scoreConditions, " OR "), year, paramIndex)

	params = append(params, totalScore)

	log.Printf("Executing query: %s", query)
	log.Printf("Parameters: %v", params)

	rows, err := DataBaseConnect.Db.Query(query, params...)
	if err != nil {
		log.Printf("Database query failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Database query failed: %v", err)})
		return
	}
	defer rows.Close()

	var programs []Program
	for rows.Next() {
		var p Program
		if err := rows.Scan(&p.ID, &p.Code, &p.Name, &p.Description, &p.FormatEducation, &p.Year, &p.ScoreBudget, &p.ScorePaid,
			&p.AvailablePlaces, &p.EducationPrice); err != nil {
			log.Printf("Failed to scan results: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to scan results"})
			return
		}
		programs = append(programs, p)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error processing results: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error processing results"})
		return
	}

	log.Printf("Query executed successfully, found %d programs", len(programs))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(programs)
}

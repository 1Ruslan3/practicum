package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
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
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Invalid score for subject %s", key)})
			return
		}
		subjectScores = append(subjectScores, SubjectScore{SubjectID: key, Score: score})
	}

	if len(subjectScores) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No subjects provided"})
		return
	}
}

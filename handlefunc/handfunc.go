package handlefunc

import (
	"net/http"
	"practicum/handlers"
)

func SetupRoutes() {
	http.HandleFunc("/subjects", handlers.SubjectsHandler)
	http.HandleFunc("/programs", handlers.ProgramsHandler)
}

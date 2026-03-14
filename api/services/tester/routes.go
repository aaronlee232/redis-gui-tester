package tester

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("POST /run-scenario", RunAllScenarios)
	r.HandleFunc("POST /run-scenario/{id}", RunScenario)
	r.HandleFunc("GET /run-scenario/{id}", GetScenarioOutput)

	return r
}

package scenario

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /get-all", GetAllScenarios)
	r.HandleFunc("POST /create", CreateScenario)

	r.HandleFunc("GET /get/{id}", GetScenario)
	r.HandleFunc("UPDATE /update/{id}", UpdateScenario)
	r.HandleFunc("DELETE /delete/{id}", DeleteScenario)

	return r
}

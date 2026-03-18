package scenario

import (
	"database/sql"
	"net/http"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /get-all", h.GetAllScenarios)
	r.HandleFunc("POST /create", h.CreateScenario)
	r.HandleFunc("PUT /update/{id}", h.UpdateScenario)
	r.HandleFunc("DELETE /delete/{id}", h.DeleteScenario)

	r.HandleFunc("GET /form", h.Form)
	r.HandleFunc("GET /form/{id}", h.FormByID)
	r.HandleFunc("GET /get/{id}", h.GetScenario)

	return r
}

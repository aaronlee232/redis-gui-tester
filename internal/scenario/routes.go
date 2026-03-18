package scenario

import (
	"net/http"

	"github.com/aaronlee232/redis-gui-tester/internal/database"
)

type Handler struct {
	repo *database.Registry
}

func NewHandler(repo *database.Registry) *Handler {
	return &Handler{repo: repo}
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

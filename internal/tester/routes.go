package tester

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
	r.HandleFunc("POST /run-scenario", h.RunAllScenarios)
	r.HandleFunc("POST /run-scenario/{id}", h.RunScenario)

	return r
}

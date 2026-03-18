package scenario

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aaronlee232/redis-gui-tester/internal/models"
	"github.com/aaronlee232/redis-gui-tester/internal/utils"
)

func (h *Handler) GetScenario(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "missing scenario id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid scenario id", http.StatusBadRequest)
		return
	}

	s, err := h.repo.Scenarios.GetById(r.Context(), id)
	if err != nil {
		http.Error(w, "scenario not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *Handler) GetAllScenarios(w http.ResponseWriter, r *http.Request) {
	scenarios, err := h.repo.Scenarios.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get scenarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scenarios)
}

func (h *Handler) CreateScenario(w http.ResponseWriter, r *http.Request) {
	var payload models.Scenario
	if err := utils.DecodeRequestJSON(w, r, &payload); err != nil {
		return
	}
	if payload.Status == "" {
		payload.Status = models.StatusUntested
	}

	if err := h.repo.Scenarios.Create(r.Context(), &payload); err != nil {
		http.Error(w, "failed to create scenario", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateScenario(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "missing scenario id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid scenario id", http.StatusBadRequest)
		return
	}

	// Load existing scenario so we can decide whether to clear expected responses.
	existing, err := h.repo.Scenarios.GetById(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to load existing scenario", http.StatusInternalServerError)
		return
	}

	var payload models.Scenario
	if err := utils.DecodeRequestJSON(w, r, &payload); err != nil {
		return
	}
	if payload.Status == "" {
		payload.Status = models.StatusUntested
	}

	// If commands changed, clear expected responses; otherwise preserve them.
	commandsChanged := len(existing.Commands) != len(payload.Commands)
	if !commandsChanged {
		for i := range existing.Commands {
			if existing.Commands[i] != payload.Commands[i] {
				commandsChanged = true
				break
			}
		}
	}
	if commandsChanged {
		// Non-nil empty slice => explicit "clear expected responses".
		payload.ExpectedResponses = []string{}
	} else {
		// Nil => repository leaves expected responses untouched.
		payload.ExpectedResponses = nil
	}

	if err := h.repo.Scenarios.Update(r.Context(), id, &payload); err != nil {
		http.Error(w, "failed to update scenario", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteScenario(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "missing scenario id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid scenario id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Scenarios.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete scenario", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

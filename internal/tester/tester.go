package tester

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type runResponse struct {
	Outputs []string `json:"outputs"`
}

func run(c context.Context, commands []string) ([]string, error) {
	var outputs []string
	for _, cmdArgs := range commands {
		args := strings.Fields(cmdArgs)
		// Use CommandContext so it dies if the request times out
		cmd := exec.CommandContext(c, "redis-cli", args...)

		// CombinedOutput captures both Stdout and Stderr
		out, err := cmd.CombinedOutput()
		if err != nil {
			return outputs, fmt.Errorf("cmd redis-cli [%s] failed: %w: %s", cmdArgs, err, string(out))
		}
		outputs = append(outputs, string(out))
	}
	return outputs, nil
}

func (h *Handler) RunScenario(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "failed to retrieve scenario with id", http.StatusInternalServerError)
	}

	// Bound execution time for redis-cli commands.
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	outs, err := run(ctx, s.Commands)
	if err != nil {
		log.Printf("RunScenario: scenario_id=%d failed: %v", s.ID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Persist expected responses for this scenario.
	s.ExpectedResponses = outs
	if err := h.repo.Scenarios.Update(ctx, s.ID, &s); err != nil {
		log.Printf("RunScenario: failed to persist expected responses for scenario_id=%d: %v", s.ID, err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(runResponse{Outputs: outs})
}

func (h *Handler) RunAllScenarios(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	scenarios, err := h.repo.Scenarios.GetAll(r.Context())
	if err != nil {
		log.Printf("RunAllScenarios: failed to retrieve all scenarios: %v", err)
		http.Error(w, "failed to retrieve all scenarios", http.StatusInternalServerError)
		return
	}

	type scenarioRun struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Output  string `json:"output"`
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}

	results := make([]scenarioRun, 0, len(scenarios))
	for _, s := range scenarios {
		outs, err := run(ctx, s.Commands)
		if err != nil {
			log.Printf("RunAllScenarios: scenario_id=%d failed: %v", s.ID, err)
			results = append(results, scenarioRun{
				ID:      s.ID,
				Title:   s.Title,
				Output:  strings.Join(outs, "\n"),
				Success: false,
				Error:   err.Error(),
			})
			continue
		}
		results = append(results, scenarioRun{
			ID:      s.ID,
			Title:   s.Title,
			Output:  strings.Join(outs, "\n"),
			Success: true,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

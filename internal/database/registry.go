package database

import (
	"database/sql"

	"github.com/aaronlee232/redis-gui-tester/internal/models"
)

type Registry struct {
	Scenarios models.ScenarioStore
	// Testers models.TesterStore
}

func NewRegistry(db *sql.DB) *Registry {
	return &Registry{
		Scenarios: NewScenarioRepository(db),
	}
}

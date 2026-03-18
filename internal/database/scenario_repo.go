package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aaronlee232/redis-gui-tester/internal/models"
)

type ScenarioRepository struct {
	db *sql.DB
}

func NewScenarioRepository(db *sql.DB) *ScenarioRepository {
	return &ScenarioRepository{db}
}

func (r *ScenarioRepository) Create(c context.Context, s *models.Scenario) error {
	// 1. Use a Transaction (atomic operations)
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return fmt.Errorf("failed to create scenario. could not initiate transaction: %w", err)
	}
	// Safely roll back if something goes wrong
	defer tx.Rollback()

	// Insert scenario
	res, err := tx.ExecContext(c,
		`INSERT INTO scenarios (title, description, status) VALUES (?, ?, ?)`,
		s.Title, s.Description, s.Status)
	if err != nil {
		return fmt.Errorf("failed to create scenario. could not insert row into `scenarios`: %w", err)
	}

	// Record Scenario ID
	scenarioID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to create scenario. Failed to fetch ID of inserted scenario: %w", err)
	}

	// Insert commands
	cmdStmt, err := tx.Prepare(`INSERT INTO commands (scenario_id, step_order, command_text) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to create scenario. error when preparing statment for insert into `commands`: %w", err)
	}
	defer cmdStmt.Close()

	for stepOrder, commandText := range s.Commands {
		if _, err := cmdStmt.Exec(scenarioID, stepOrder, commandText); err != nil {
			return fmt.Errorf("failed to create scenario. Failed to insert row into `commands`: %w", err)
		}
	}

	// Insert expected responses
	rspStmt, err := tx.Prepare(`INSERT INTO expected_responses (scenario_id, step_order, response_text) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to create scenario. error when preparing statment for insert into `expected_responses`: %w", err)
	}
	defer rspStmt.Close()

	for stepOrder, responseText := range s.ExpectedResponses {
		if _, err := rspStmt.Exec(scenarioID, stepOrder, responseText); err != nil {
			return fmt.Errorf("failed to create scenario. Failed to insert row into `expected_responses`: %w", err)
		}
	}

	// Commit whole batch
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to create scenario. Failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ScenarioRepository) GetById(c context.Context, id int) (models.Scenario, error) {
	var s models.Scenario
	var cmdsJSON, respsJSON string
	err := r.db.QueryRowContext(c, `
		SELECT
			s.scenario_id, s.title, s.description, s.status,
			(SELECT json_group_array(command_text) FROM (SELECT command_text FROM commands WHERE scenario_id = s.scenario_id ORDER BY step_order)) AS commands,
			(SELECT json_group_array(response_text) FROM (SELECT response_text FROM expected_responses WHERE scenario_id = s.scenario_id ORDER BY step_order)) AS responses
		FROM scenarios s WHERE s.scenario_id = ?
	`, id).Scan(&s.ID, &s.Title, &s.Description, &s.Status, &cmdsJSON, &respsJSON)
	if err != nil {
		return models.Scenario{}, fmt.Errorf("failed to get scenario: %w", err)
	}

	json.Unmarshal([]byte(cmdsJSON), &s.Commands)
	json.Unmarshal([]byte(respsJSON), &s.ExpectedResponses)

	return s, nil
}

func (r *ScenarioRepository) GetAll(c context.Context) ([]models.Scenario, error) {
	// Get scenario with commands and expected responses
	rows, err := r.db.Query(`
	SELECT 
	s.scenario_id, s.title, s.description, s.status,
	(SELECT json_group_array(command_text) 
	FROM (SELECT command_text FROM commands 
	WHERE scenario_id = s.scenario_id ORDER BY step_order)) as commands,
	(SELECT json_group_array(response_text) 
	FROM (SELECT response_text FROM expected_responses 
	WHERE scenario_id = s.scenario_id ORDER BY step_order)) as responses
	FROM scenarios s
	ORDER BY s.created_at DESC;
	`)
	if err != nil {
		return []models.Scenario{}, fmt.Errorf("failed to get all scenarios : %w", err)
	}
	defer rows.Close()

	scanErrors := []error{}
	scenarios := []models.Scenario{}
	for rows.Next() {
		var s models.Scenario
		var cmdsJSON, respsJSON string

		// Scan column values into variables
		err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.Status, &cmdsJSON, &respsJSON)
		if err != nil {
			scanErrors = append(scanErrors, err)
			continue
		}

		// Unmarshal the JSON strings into your Go slices
		json.Unmarshal([]byte(cmdsJSON), &s.Commands)
		json.Unmarshal([]byte(respsJSON), &s.ExpectedResponses)

		scenarios = append(scenarios, s)
	}

	// Check for errors that happened during iteration
	if err = rows.Err(); err != nil {
		combinedErrors := errors.Join(scanErrors...)
		return []models.Scenario{}, fmt.Errorf("failed to get all scenarios. Error(s) during iteration: %w", combinedErrors)
	}

	return scenarios, nil
}

func (r *ScenarioRepository) Update(c context.Context, id int, s *models.Scenario) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return fmt.Errorf("failed to update scenario. could not initiate transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(c, `UPDATE scenarios SET title = ?, description = ?, status = ? WHERE scenario_id = ?`,
		s.Title, s.Description, s.Status, id)
	if err != nil {
		return fmt.Errorf("failed to update scenario. could not update row in `scenarios`: %w", err)
	}

	_, err = tx.ExecContext(c, `DELETE FROM commands WHERE scenario_id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to update scenario. could not delete existing rows in `commands`: %w", err)
	}
	cmdStmt, err := tx.Prepare(`INSERT INTO commands (scenario_id, step_order, command_text) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to update scenario. could not prepare statment to insert new rows in `commands`: %w", err)
	}
	defer cmdStmt.Close()
	for stepOrder, commandText := range s.Commands {
		if _, err := cmdStmt.Exec(id, stepOrder, commandText); err != nil {
			return fmt.Errorf("failed to update scenario. could insert new rows in `commands`: %w", err)
		}
	}

	// Only touch expected_responses if the slice is non-nil.
	// nil means "leave existing expected responses as-is".
	if s.ExpectedResponses != nil {
		_, err = tx.ExecContext(c, `DELETE FROM expected_responses WHERE scenario_id = ?`, id)
		if err != nil {
			return fmt.Errorf("failed to update scenario. could not delete existing rows in `expected_responses`: %w", err)
		}
		rspStmt, err := tx.Prepare(`INSERT INTO expected_responses (scenario_id, step_order, response_text) VALUES (?, ?, ?)`)
		if err != nil {
			return fmt.Errorf("failed to update scenario. could not prepare statment to insert new rows in `expected_responses`: %w", err)
		}
		defer rspStmt.Close()
		for stepOrder, responseText := range s.ExpectedResponses {
			if _, err := rspStmt.Exec(id, stepOrder, responseText); err != nil {
				return fmt.Errorf("failed to update scenario. could not insert new rows in `expected_responses`: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to update scenario. Failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ScenarioRepository) Delete(c context.Context, id int) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return fmt.Errorf("failed to delete scenario. could not initiate transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(c, `DELETE FROM commands WHERE scenario_id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to update scenario. could not delete exisiting rows in `commands`: %w", err)
	}
	_, err = tx.ExecContext(c, `DELETE FROM expected_responses WHERE scenario_id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to update scenario. could not delete exisiting rows in `expected_responses`: %w", err)
	}
	res, err := tx.ExecContext(c, `DELETE FROM scenarios WHERE scenario_id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to update scenario. could not delete exisiting row in `scenarios`: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("failed to update scenario. no rows were affected `scenarios`: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to delete scenario. Failed to commit transaction: %w", err)
	}

	return nil
}

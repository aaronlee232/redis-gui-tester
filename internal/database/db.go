package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	// Get current working directory as the project root
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Set up app.db file
	dbPath := filepath.Join(root, "data", "app.db")

	// Open DB Connection
	dsn := dbPath + "?_foreign_keys=on"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatal(err)
	}

	createScenariosTable(db)
	createCommandsTable(db)
	createExpectedResponsesTable(db)

	return db
}

func createScenariosTable(db *sql.DB) {
	createTableQuery := `CREATE TABLE IF NOT EXISTS scenarios (
		scenario_id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func createCommandsTable(db *sql.DB) {
	createTableQuery := `CREATE TABLE IF NOT EXISTS commands (
		command_id INTEGER PRIMARY KEY,
		scenario_id INTEGER NOT NULL,
		step_order INTEGER NOT NULL,
		command_text TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(scenario_id) REFERENCES scenarios(scenario_id)
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func createExpectedResponsesTable(db *sql.DB) {
	createTableQuery := `CREATE TABLE IF NOT EXISTS expected_responses (
		expected_response_id INTEGER PRIMARY KEY,
		scenario_id INTEGER NOT NULL,
		step_order INTEGER NOT NULL,
		response_text TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(scenario_id) REFERENCES scenarios(scenario_id)
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

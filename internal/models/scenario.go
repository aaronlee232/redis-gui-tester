package models

import (
	"context"
)

type ScenarioStatus string

const (
	StatusUntested ScenarioStatus = "untested"
	StatusPassed   ScenarioStatus = "passed"
	StatusFailed   ScenarioStatus = "failed"
)

type Scenario struct {
	ID                int            `json:"ID"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	Commands          []string       `json:"commands"`
	ExpectedResponses []string       `json:"expected_responses"`
	Status            ScenarioStatus `json:"status"`
}

func NewScenario(title string, description string, commands []string, expectedResponses []string, status ScenarioStatus) *Scenario {
	return &Scenario{
		Title:             title,
		Description:       description,
		Commands:          commands,
		ExpectedResponses: expectedResponses,
		Status:            status,
	}
}

type ScenarioStore interface {
	Create(c context.Context, s *Scenario) error
	GetById(c context.Context, id int) (Scenario, error)
	GetAll(c context.Context) ([]Scenario, error)
	Update(c context.Context, id int, s *Scenario) error
	Delete(c context.Context, id int) error
	// Run(c context.Context, id int) (string, error)
	// RunAll(c context.Context) ([]string, error)
}

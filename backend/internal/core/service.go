package core

import (
	"context"
	"fmt"

	"github.com/rohaaaaaan/devair-backend/internal/db"
	"github.com/rohaaaaaan/devair-backend/internal/gateway"
	"github.com/rohaaaaaan/devair-backend/internal/models"
)

// Service handles core business logic
type Service struct {
	// db *pgxpool.Pool
}

func NewService() *Service {
	return &Service{}
}

// GetProjects returns the list of projects from DB
func (s *Service) GetProjects() []models.Project {
	if db.Pool == nil {
		// Fallback to mock if DB not connected
		return []models.Project{
			{ID: "mock-1", Name: "Mock-Project", Status: "DB Not Connected", State: "error"},
		}
	}

	rows, err := db.Pool.Query(context.Background(), "SELECT id, name, repo_url, status, state FROM projects")
	if err != nil {
		fmt.Printf("Error querying projects: %v\n", err)
		return []models.Project{}
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.RepoURL, &p.Status, &p.State); err != nil {
			continue
		}
		projects = append(projects, p)
	}

	return projects
}

// TriggerBuild creates a new job and dispatches it to the agent
func (s *Service) TriggerBuild(projectID string) (models.Job, error) {
	return s.TriggerJob(projectID, "BUILD")
}

// TriggerJob allows executing arbitrary commands (e.g. OPEN_IDE, OPEN_APP, AI_INSTRUCTION, UI_ACTION)
func (s *Service) TriggerJob(projectID string, jobType string) (models.Job, error) {
	return s.TriggerJobWithParams(projectID, jobType, "", "", "", "", "")
}

// TriggerJobWithParams allows passing app name, prompts, or UI action details
func (s *Service) TriggerJobWithParams(projectID string, jobType string, appName string, prompt string, action string, target string, value string) (models.Job, error) {
	// 1. Create Job in DB
	var jobID string
	err := db.Pool.QueryRow(context.Background(),
		"INSERT INTO jobs (project_id, type, status, input_params) VALUES ($1, $2, $3, $4) RETURNING id",
		projectID, jobType, "QUEUED", "{}").Scan(&jobID)

	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return models.Job{}, err
	}

	// 2. Dispatch to Agent
	// Default command for known types
	cmdStr := ""
	switch jobType {
	case models.CommandTypeBuild:
		cmdStr = "npm run build"
	case models.CommandTypeOpenIDE:
		cmdStr = "code ."
	}

	cmd := models.CommandPayload{
		JobID:   jobID,
		Type:    jobType,
		Command: cmdStr,
		App:     appName,
		Prompt:  prompt,
		Action:  action,
		Target:  target,
		Value:   value,
	}

	msg := models.WSMessage{
		Type:    models.EventTypeCommand,
		Payload: cmd,
	}

	if sent := gateway.GlobalManager.SendToAgent(projectID, msg); sent {
		fmt.Printf("Command dispatched to Agent for Project %s\n", projectID)
	} else {
		fmt.Printf("No agent connected for Project %s. Job queued.\n", projectID)
	}

	return models.Job{
		ID:        jobID,
		ProjectID: projectID,
		Type:      jobType,
		Status:    "QUEUED",
	}, nil
}

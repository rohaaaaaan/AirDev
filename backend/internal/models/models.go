package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Project struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	RepoURL       string    `json:"repo_url"`
	Status        string    `json:"status"` // "last run: 2 hours ago"
	State         string    `json:"state"`  // "success", "error", "idle"
	LastBuildID   string    `json:"last_build_id"`
}

type Job struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Type      string    `json:"type"`   // BUILD, TEST, DEPLOY
	Status    string    `json:"status"` // CREATED, RUNNING, COMPLETED, FAILED
	Result    string    `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}

type Agent struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"project_id"`
	Status       string    `json:"status"` // ONLINE, OFFLINE
	LastSeenAt   time.Time `json:"last_seen_at"`
}

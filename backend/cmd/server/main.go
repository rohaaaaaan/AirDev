package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rohaaaaaan/devair-backend/internal/core"
	"github.com/rohaaaaaan/devair-backend/internal/db"
	"github.com/rohaaaaaan/devair-backend/internal/gateway"
)

func main() {
	// Load .env
	_ = godotenv.Load()

	// Init DB (Skip for now if no creds, but structure is there)
	if os.Getenv("DATABASE_URL") != "" {
		if err := db.InitDB(); err != nil {
			log.Printf("Warning: Database init failed: %v", err)
		} else {
			defer db.CloseDB()
		}
	}

	// Init Service
	svc := core.NewService()

	// Init Gin
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow all for dev
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	api := r.Group("/api")
	{
		api.GET("/projects", func(c *gin.Context) {
			projects := svc.GetProjects()
			c.JSON(http.StatusOK, projects)
		})

		api.POST("/projects/:id/build", func(c *gin.Context) {
			projectID := c.Param("id")
			job, _ := svc.TriggerBuild(projectID)
			c.JSON(http.StatusOK, job)
		})

		api.POST("/projects/:id/command", func(c *gin.Context) {
			projectID := c.Param("id")
			var req struct {
				Type string `json:"type"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}
			job, err := svc.TriggerJob(projectID, req.Type)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, job)
		})

		// AI Analysis Endpoint
		aiSvc := core.NewAIService()
		api.POST("/ai/analyze", func(c *gin.Context) {
			var req struct {
				Logs string `json:"logs"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			analysis := aiSvc.AnalyzeLogs(req.Logs)
			c.JSON(http.StatusOK, analysis)
		})

		// AI Command Endpoint (OPEN_APP, AI_INSTRUCTION, UI_ACTION)
		api.POST("/projects/:id/ai-command", func(c *gin.Context) {
			projectID := c.Param("id")
			var req struct {
				Type   string `json:"type"`   // OPEN_APP, AI_INSTRUCTION, UI_ACTION
				App    string `json:"app"`    // For OPEN_APP
				Prompt string `json:"prompt"` // For AI_INSTRUCTION
				Action string `json:"action"` // For UI_ACTION
				Target string `json:"target"` // For UI_ACTION
				Value  string `json:"value"`  // For UI_ACTION
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}
			job, err := svc.TriggerJobWithParams(projectID, req.Type, req.App, req.Prompt, req.Action, req.Target, req.Value)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, job)
		})
	}

	// WebSocket
	r.GET("/ws", gateway.HandleWebSocket)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

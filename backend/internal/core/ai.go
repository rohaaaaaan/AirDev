package core

import (
	"strings"
	"time"
)

type AIService struct{}

func NewAIService() *AIService {
	return &AIService{}
}

type AIAnalysisResponse struct {
	Analysis   string `json:"analysis"`
	Suggestion string `json:"suggestion"`
	Confidence int    `json:"confidence"`
}

// AnalyzeLogs mocks an LLM by looking for keywords
func (s *AIService) AnalyzeLogs(logs string) AIAnalysisResponse {
	// Simulate "Thinking" time
	time.Sleep(1 * time.Second)

	logsLower := strings.ToLower(logs)

	if strings.Contains(logsLower, "error") || strings.Contains(logsLower, "fail") {
		return AIAnalysisResponse{
			Analysis:   "I detected a build failure in the logs.",
			Suggestion: "It looks like a syntax error or missing dependency. Try running 'npm install' to ensure all packages are available.",
			Confidence: 85,
		}
	}

	if strings.Contains(logsLower, "success") || strings.Contains(logsLower, "done") {
		return AIAnalysisResponse{
			Analysis:   "The build appears to be successful.",
			Suggestion: "You can proceed to deployment.",
			Confidence: 95,
		}
	}

	return AIAnalysisResponse{
		Analysis:   "I analyzed the logs but didn't find any obvious errors.",
		Suggestion: "Check the agent connection or try running the build again.",
		Confidence: 50,
	}
}

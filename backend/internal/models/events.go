package models

type EventType string

const (
	EventTypeIdentify      EventType = "IDENTIFY"
	EventTypeCommand       EventType = "COMMAND"
	EventTypeLogChunk      EventType = "LOG_CHUNK"
	EventTypeJobUpdate     EventType = "JOB_UPDATE"
	EventTypeAIStageUpdate EventType = "AI_STAGE_UPDATE" // New: For streaming AI progress
)

const (
	RoleAgent  = "AGENT"
	RoleClient = "CLIENT"
)

// Command Types for Job Dispatching
const (
	CommandTypeBuild         = "BUILD"
	CommandTypeOpenIDE       = "OPEN_IDE"
	CommandTypeOpenApp       = "OPEN_APP"       // New: Open arbitrary apps
	CommandTypeAIInstruction = "AI_INSTRUCTION" // New: Natural language instruction
	CommandTypeUIAction      = "UI_ACTION"      // New: Low-level UI control
)

// Base WebSocket Message
type WSMessage struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

// Payload for "IDENTIFY" (Agent -> Server)
type IdentifyPayload struct {
	ProjectID string `json:"project_id"`
	Secret    string `json:"secret"` // Simple auth for now
	Role      string `json:"role"`   // "AGENT" or "CLIENT"
}

// Payload for "COMMAND" (Server -> Agent)
type CommandPayload struct {
	JobID   string            `json:"job_id"`
	Type    string            `json:"type"`              // BUILD, OPEN_APP, AI_INSTRUCTION, UI_ACTION
	Command string            `json:"command,omitempty"` // e.g., "npm run build"
	App     string            `json:"app,omitempty"`     // New: For OPEN_APP
	Prompt  string            `json:"prompt,omitempty"`  // New: For AI_INSTRUCTION
	Action  string            `json:"action,omitempty"`  // New: For UI_ACTION (find, click, type)
	Target  string            `json:"target,omitempty"`  // New: For UI_ACTION (window name, element name)
	Value   string            `json:"value,omitempty"`   // New: For UI_ACTION (text to type)
	Params  map[string]string `json:"params"`
}

// Payload for "JOB_UPDATE" (Agent -> Server)
type JobUpdatePayload struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
}

// Payload for "AI_STAGE_UPDATE" (Agent -> Server -> Clients)
type AIStagePayload struct {
	JobID   string `json:"job_id"`
	Stage   string `json:"stage"`   // e.g., "Analyzing", "Creating files", "Done"
	Message string `json:"message"` // Optional details
}

package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rohaaaaaan/devair-backend/internal/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

// Manager tracks connections
type Manager struct {
	agents  map[string]*websocket.Conn
	clients map[string][]*websocket.Conn // ProjectID -> List of Clients
	lock    sync.RWMutex
}

var GlobalManager = &Manager{
	agents:  make(map[string]*websocket.Conn),
	clients: make(map[string][]*websocket.Conn),
}

func (m *Manager) Register(projectID string, conn *websocket.Conn, role string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if role == models.RoleClient {
		m.clients[projectID] = append(m.clients[projectID], conn)
		log.Printf("Client connected to Project: %s", projectID)
	} else {
		m.agents[projectID] = conn
		log.Printf("Agent registered for Project: %s", projectID)
	}
}

func (m *Manager) Unregister(projectID string, conn *websocket.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Check Agents
	if current, ok := m.agents[projectID]; ok && current == conn {
		conn.Close()
		delete(m.agents, projectID)
		log.Printf("Agent disconnected from Project: %s", projectID)
		return
	}

	// Check Clients
	if clients, ok := m.clients[projectID]; ok {
		for i, c := range clients {
			if c == conn {
				conn.Close()
				// Remove from slice
				m.clients[projectID] = append(clients[:i], clients[i+1:]...)
				log.Printf("Client disconnected from Project: %s", projectID)
				return
			}
		}
	}
}

func (m *Manager) SendToAgent(projectID string, msg models.WSMessage) bool {
	m.lock.RLock()
	conn, ok := m.agents[projectID]
	m.lock.RUnlock()

	if !ok {
		return false
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending to agent: %v", err)
		m.Unregister(projectID, conn)
		return false
	}
	return true
}

func (m *Manager) BroadcastToClients(projectID string, msg models.WSMessage) {
	m.lock.RLock()
	clients := m.clients[projectID]
	m.lock.RUnlock()

	for _, conn := range clients {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error sending to client: %v", err)
		}
	}
}

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return
	}
	// Don't close immediately, let connection live

	// Wait for IDENTIFY message
	var msg models.WSMessage
	if err := conn.ReadJSON(&msg); err != nil {
		log.Println("Failed to read initial message:", err)
		conn.Close()
		return
	}

	if msg.Type == models.EventTypeIdentify {
		// Parse payload to get Project ID
		payloadBytes, _ := json.Marshal(msg.Payload)
		var identify models.IdentifyPayload
		if err := json.Unmarshal(payloadBytes, &identify); err == nil && identify.ProjectID != "" {

			role := identify.Role
			if role == "" {
				role = models.RoleAgent // Default to Agent for backward compat
			}

			GlobalManager.Register(identify.ProjectID, conn, role)

			// Listen loop to keep connection open (and handle updates)
			for {
				var incomingMsg models.WSMessage
				if err := conn.ReadJSON(&incomingMsg); err != nil {
					GlobalManager.Unregister(identify.ProjectID, conn)
					break
				}

				// Broadcast if it's an Agent Log or Job Update or AI Stage
				if role == models.RoleAgent && (incomingMsg.Type == models.EventTypeLogChunk || incomingMsg.Type == models.EventTypeJobUpdate || incomingMsg.Type == models.EventTypeAIStageUpdate) {
					GlobalManager.BroadcastToClients(identify.ProjectID, incomingMsg)
				}
			}
		} else {
			log.Println("Invalid IDENTIFY payload")
			conn.Close()
		}
	} else {
		log.Println("First message must be IDENTIFY")
		conn.Close()
	}
}

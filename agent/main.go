package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Copied from backend for now (should share code later)
type EventType string

const (
	EventTypeIdentify      EventType = "IDENTIFY"
	EventTypeCommand       EventType = "COMMAND"
	EventTypeLogChunk      EventType = "LOG_CHUNK"
	EventTypeJobUpdate     EventType = "JOB_UPDATE"
	EventTypeAIStageUpdate EventType = "AI_STAGE_UPDATE"
)

type WSMessage struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

type IdentifyPayload struct {
	ProjectID string `json:"project_id"`
	Secret    string `json:"secret"`
	Role      string `json:"role"`
}

type CommandPayload struct {
	JobID   string            `json:"job_id"`
	Type    string            `json:"type"`
	Command string            `json:"command"`
	App     string            `json:"app"`    // For OPEN_APP
	Prompt  string            `json:"prompt"` // For AI_INSTRUCTION
	Action  string            `json:"action"` // For UI_ACTION
	Target  string            `json:"target"` // For UI_ACTION
	Value   string            `json:"value"`  // For UI_ACTION
	Params  map[string]string `json:"params"`
}

// App launcher mapping (friendly name -> executable)
var appLauncher = map[string]string{
	"cursor":             "cursor",
	"blender":            "blender",
	"code":               "code",
	"vs code":            "code",
	"visual studio code": "code",
	"notepad":            "notepad",
	"calc":               "calc",
	"explorer":           "explorer",
}

func main() {
	// Define Flags
	serverURLPtr := flag.String("server", "ws://localhost:8080/ws", "WebSocket server URL")
	apiURLPtr := flag.String("api", "http://localhost:8080/api/projects", "API URL for project auto-discovery")
	projectIDPtr := flag.String("project", "", "Project ID (optional, will auto-fetch if empty)")
	secretPtr := flag.String("secret", "my-secret-token", "Authentication secret")
	wdPtr := flag.String("wd", ".", "Working directory for executed commands")

	flag.Parse()

	serverURL := *serverURLPtr
	apiURL := *apiURLPtr
	projectID := *projectIDPtr
	secret := *secretPtr
	workDir := *wdPtr

	// Resolve Project ID
	if projectID == "" {
		// Auto-fetch from API
		log.Println("No Project ID provided, fetching from API...")
		resp, err := http.Get(apiURL)
		if err != nil {
			log.Fatalf("Failed to fetch projects from %s: %v", apiURL, err)
		}
		defer resp.Body.Close()

		var projects []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
			log.Fatalf("Failed to decode projects: %v", err)
		}

		if len(projects) > 0 {
			projectID = projects[0].ID
			log.Printf("Auto-selected Project: %s (%s)", projects[0].Name, projectID)
		} else {
			log.Fatal("No projects found in Backend. Please seed the DB or provide -project flag.")
		}
	}

	log.Printf("Connecting to %s as Agent for Project %s...", serverURL, projectID)
	log.Printf("Working Directory: %s", workDir)

	c, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// IDENTIFY
	identifyMsg := WSMessage{
		Type: EventTypeIdentify,
		Payload: IdentifyPayload{
			ProjectID: projectID,
			Secret:    secret,
			Role:      "AGENT", // Explicitly set role
		},
	}
	if err := c.WriteJSON(identifyMsg); err != nil {
		log.Println("write identify:", err)
		return
	}
	log.Println("Identified with Backend.")

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			var msg WSMessage
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("Received Message Type: %s", msg.Type)

			if msg.Type == EventTypeCommand {
				// Parse Payload
				payloadBytes, _ := json.Marshal(msg.Payload)
				var cmdPayload CommandPayload
				if err := json.Unmarshal(payloadBytes, &cmdPayload); err != nil {
					log.Printf("Error processing command payload: %v", err)
					continue
				}

				log.Printf(">>> EXECUTING: %s", cmdPayload.Type)

				switch cmdPayload.Type {
				case "OPEN_APP":
					// Launch an app by friendly name
					appName := cmdPayload.App
					executable, ok := appLauncher[strings.ToLower(appName)]
					if !ok {
						log.Printf("Unknown app: %s", appName)
						c.WriteJSON(WSMessage{
							Type: EventTypeJobUpdate,
							Payload: map[string]string{
								"job_id": cmdPayload.JobID,
								"status": "FAILED",
							},
						})
						continue
					}
					log.Printf("Launching app: %s (%s)", appName, executable)
					cmd := exec.Command(executable)
					cmd.Dir = workDir
					if err := cmd.Start(); err != nil {
						log.Printf("Failed to launch app: %v", err)
					}
					c.WriteJSON(WSMessage{
						Type: EventTypeJobUpdate,
						Payload: map[string]string{
							"job_id": cmdPayload.JobID,
							"status": "COMPLETED",
						},
					})

				case "AI_INSTRUCTION":
					// Simulate AI processing stages
					prompt := cmdPayload.Prompt
					log.Printf("AI Instruction received: %s", prompt)

					// Simulate sending AI stage updates
					stages := []string{"Analyzing request...", "Planning execution...", "Generating code...", "Done!"}
					for _, stage := range stages {
						log.Printf("AI Stage: %s", stage)
						c.WriteJSON(WSMessage{
							Type: EventTypeAIStageUpdate,
							Payload: map[string]string{
								"job_id":  cmdPayload.JobID,
								"stage":   stage,
								"message": fmt.Sprintf("Processing: %s", prompt),
							},
						})
						time.Sleep(1 * time.Second) // Simulate work
					}

					c.WriteJSON(WSMessage{
						Type: EventTypeJobUpdate,
						Payload: map[string]string{
							"job_id": cmdPayload.JobID,
							"status": "COMPLETED",
						},
					})

				case "UI_ACTION":
					action := cmdPayload.Action
					target := cmdPayload.Target
					value := cmdPayload.Value

					log.Printf("UI Action: %s on %s with value '%s'", action, target, value)

					switch action {
					case "FIND":
						// FIND/FOCUS
						// Use PowerShell to focus window, return extensive error if fails
						// Also try partial match by iterating processes if direct AppActivate fails
						psScript := fmt.Sprintf(`
							$wshell = New-Object -ComObject wscript.shell
							if ($wshell.AppActivate('%s')) { exit 0 }
							# Try to find by process main window title
							$proc = Get-Process | Where-Object { $_.MainWindowTitle -match '%s' } | Select-Object -First 1
							if ($proc) {
								if ($wshell.AppActivate($proc.Id)) { exit 0 }
							}
							exit 1
						`, target, target)

						cmd := exec.Command("powershell", "-Command", psScript)
						if err := cmd.Run(); err != nil {
							log.Printf("Failed to focus window '%s': %v", target, err)

							// If we can't find it, launch it?
							// Only if it's a known app
							if executable, ok := appLauncher[strings.ToLower(target)]; ok {
								log.Printf("Launching %s...", target)
								exec.Command(executable).Start()

								// Wait longer for app to actually open
								time.Sleep(3 * time.Second)
							} else {
								// REPORT FAILURE and STOP
								c.WriteJSON(WSMessage{
									Type: EventTypeJobUpdate,
									Payload: map[string]string{
										"job_id": cmdPayload.JobID,
										"status": "FAILED",
										"error":  fmt.Sprintf("Window '%s' not found", target),
									},
								})
								continue
							}
						}

					case "TYPE":
						// TYPE text
						// Safety check: If a target is specified, ensure it's focused!
						if target != "" {
							// Reuse the focus logic (simplified here for brevity, or extract to helper)
							// Use robust FIND logic: AppActivate OR Process Title Match
							psFocus := fmt.Sprintf(`
								$wshell = New-Object -ComObject wscript.shell
								if ($wshell.AppActivate('%s')) { exit 0 }
								# Try to find by process main window title
								$proc = Get-Process | Where-Object { $_.MainWindowTitle -match '%s' } | Select-Object -First 1
								if ($proc) {
									if ($wshell.AppActivate($proc.Id)) { exit 0 }
								}
								exit 1
							`, target, target)
							if err := exec.Command("powershell", "-Command", psFocus).Run(); err != nil {
								// REPORT FAILURE and STOP
								c.WriteJSON(WSMessage{
									Type: EventTypeJobUpdate,
									Payload: map[string]string{
										"job_id": cmdPayload.JobID,
										"status": "FAILED",
										"error":  fmt.Sprintf("Target '%s' not focused. Aborting TYPE.", target),
									},
								})
								continue
							}
						}

						// Use PowerShell SendKeys with escaping
						// +^%~(){}[] need escaping with {}
						safeValue := value

						// Basic escaping for SendKeys:
						// + -> {+}
						// ( -> {(}
						// ! -> {!}
						// etc.
						replacer := strings.NewReplacer(
							"+", "{+}",
							"^", "{^}",
							"%", "{%}",
							"~", "{~}",
							"(", "{(}",
							")", "{)}",
							"[", "{[}",
							"]", "{]}",
							"{", "{{}",
							"}", "{}}",
							"!", "{!}",
						)
						safeValue = replacer.Replace(safeValue)

						// Also escape single quotes for PowerShell string
						psSafeValue := strings.ReplaceAll(safeValue, "'", "''")

						psScript := fmt.Sprintf(`
							$wshell = New-Object -ComObject wscript.shell
							$wshell.SendKeys('%s')
						`, psSafeValue)
						exec.Command("powershell", "-Command", psScript).Run()

					case "CLICK":
						// Simple click/shortcut simulation
						// e.g. "enter" -> SendKeys("{ENTER}")
						keys := ""
						switch value {
						case "enter":
							keys = "{ENTER}"
						case "tab":
							keys = "{TAB}"
						case "space":
							keys = " "
						default:
							keys = value
						}
						psScript := fmt.Sprintf(`
							$wshell = New-Object -ComObject wscript.shell
							$wshell.SendKeys('%s')
						`, keys)
						exec.Command("powershell", "-Command", psScript).Run()
					}

					c.WriteJSON(WSMessage{
						Type: EventTypeJobUpdate,
						Payload: map[string]string{
							"job_id": cmdPayload.JobID,
							"status": "COMPLETED",
						},
					})

				case "OPEN_IDE":
					cmd := exec.Command("code", ".")
					cmd.Dir = workDir
					if err := cmd.Start(); err != nil {
						log.Printf("Failed to open IDE: %v", err)
					}
					c.WriteJSON(WSMessage{
						Type: EventTypeJobUpdate,
						Payload: map[string]string{
							"job_id": cmdPayload.JobID,
							"status": "COMPLETED",
						},
					})

				default:
					// Generic command execution (BUILD, etc.)
					cmd := exec.Command("cmd", "/C", cmdPayload.Command)
					cmd.Dir = workDir

					stdout, _ := cmd.StdoutPipe()
					stderr, _ := cmd.StderrPipe()

					if err := cmd.Start(); err != nil {
						log.Printf("Failed to start command: %v", err)
						continue
					}

					// Reader routine
					go func() {
						reader := io.MultiReader(stdout, stderr)
						buf := make([]byte, 1024)
						for {
							n, err := reader.Read(buf)
							if n > 0 {
								chunk := string(buf[:n])
								fmt.Print(chunk)
								c.WriteJSON(WSMessage{
									Type: EventTypeLogChunk,
									Payload: map[string]string{
										"job_id": cmdPayload.JobID,
										"chunk":  chunk,
									},
								})
							}
							if err != nil {
								break
							}
						}
					}()

					cmd.Wait()
					log.Println(">>> COMMAND FINISHED")

					c.WriteJSON(WSMessage{
						Type: EventTypeJobUpdate,
						Payload: map[string]string{
							"job_id": cmdPayload.JobID,
							"status": "COMPLETED",
						},
					})
				}
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	log.Println("interrupt")

	// Cleanly close connection
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}
}

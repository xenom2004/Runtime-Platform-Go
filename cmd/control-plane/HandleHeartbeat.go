package main

import (
	"encoding/json"
	"net/http"

	"github.com/xenom2004/Runtime-Platform-Go/cmd/helpers"
)

type Heartbeat struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func handleheartbeat(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, Post only", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var heartbeat Heartbeat

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&heartbeat); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	newagent, exists := agentstore.Get(heartbeat.Id)
	if !exists {
		http.Error(w, "Agent not found", http.StatusBadRequest)
		return
	}
	agentstore.SetLast_Alive(newagent.Id)

	// fmt.Printf("%s\n", heartbeat.Message)

	helpers.Jsonhelper(w, map[string]string{"HeartbeatResponse": heartbeat.Message}, http.StatusOK)

}

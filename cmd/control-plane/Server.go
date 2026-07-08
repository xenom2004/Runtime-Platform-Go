package main

import (
	"fmt"
	"net/http"
)

func serverhealth(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"message": "Go Server is Running, Health normal",
	}
	jsonhelper(w, res, http.StatusOK)
}

func showagents(w http.ResponseWriter, r *http.Request) {
	jsonhelper(w, map[string]any{"Agents": agentstore.agents, "AgentIDs": agentstore.AgentIDs}, http.StatusOK)
}

func main() {
	http.HandleFunc("/HandleAgents", handleagent)
	http.HandleFunc("/health", serverhealth)
	http.HandleFunc("/status", serverstatus)
	http.HandleFunc("/load", serverload)
	http.HandleFunc("/agents", showagents)
	http.HandleFunc("/HandleHeartbeat", handleheartbeat)
	http.HandleFunc("/jobs", handleJob)
	port := 9000
	fmt.Printf("Starting Server on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("error while running server :", err)
	}

}

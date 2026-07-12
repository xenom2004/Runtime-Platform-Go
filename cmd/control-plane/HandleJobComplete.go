package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/xenom2004/Runtime-Platform-Go/cmd/helpers"
)

type CompletedJob struct {
	JobID   string `json:"job_id"`
	AgentID string `json:"agent_id"`
}

func handleJobComplete(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, Post only", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var completedjob CompletedJob

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&completedjob); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("Got completed job request")
	jobstore.Update_status_and_finishtime(completedjob.JobID, "Completed", time.Now().Format(time.RFC3339))
	agentstore.Status(completedjob.AgentID, "alive")
	agentstore.SetJob(completedjob.AgentID, nil)
	fmt.Printf("Job %s completed by Agent %s\n", completedjob.JobID, completedjob.AgentID)
	helpers.Jsonhelper(w, map[string]string{"message": "Job Completed"}, http.StatusOK)

}

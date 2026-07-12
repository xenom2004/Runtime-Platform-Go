package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xenom2004/Runtime-Platform-Go/cmd/helpers"
)

type AssignedJobs struct {
	Job     Job    `json:"job"`
	AgentId string `json:"agent_id"`
}

func executejob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, Post only", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var assignedjob AssignedJobs

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&assignedjob); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("Recieved request for %s", assignedjob.Job.JobID)
	Agents[assignedjob.AgentId].CurrentJob = &assignedjob.Job
	Agents[assignedjob.AgentId].pending_jobs <- assignedjob.Job

	fmt.Printf("Job %s assigned to Agent %s\n", assignedjob.Job.JobID, assignedjob.AgentId)

	helpers.Jsonhelper(w, map[string]string{"Assigned": "Success"}, http.StatusOK)

}

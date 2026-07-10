package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xenom2004/Runtime-Platform-Go/cmd/helpers"
)

type Job struct {
	JobID string `json:"job_id"`
}

type AssignedJobs struct {
	Job     Job    `json:"job"`
	AgentId string `json:"agent_id"`
}

func handleJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, Post only", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var job Job

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&job); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var AgentAssign = agentstore.Next()
	if AgentAssign == nil {
		fmt.Printf("Job %s failed to assign\n", job.JobID)
		helpers.Jsonhelper(w, map[string]string{"Failed": "No Agent Available"}, http.StatusOK)
		return
	}

	resp := helpers.JsonPost("http://localhost:9001/executejob", AssignedJobs{Job: job, AgentId: AgentAssign.(string)})
	if resp == nil {
		helpers.Jsonhelper(w, map[string]string{"Failed": "Failed to assign job"}, http.StatusOK)
		return
	}
	mes := fmt.Sprintf("Job %s assigned to Agent %s \n", job.JobID, AgentAssign)
	fmt.Println(mes)
	helpers.Jsonhelper(w, AssignedJobs{Job: job, AgentId: AgentAssign.(string)}, http.StatusOK)

}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Job struct {
	JobID string `json:"job_id"`
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
		jsonhelper(w, map[string]string{"Failed": "No Agent Available"}, http.StatusOK)
		return
	}
	mes := fmt.Sprintf("Job %s assigned to Agent %s \n", job.JobID, AgentAssign)
	fmt.Println(mes)
	jsonhelper(w, map[string]string{"Assigned": mes}, http.StatusOK)

}

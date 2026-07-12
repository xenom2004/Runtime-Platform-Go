package main

import (
	"fmt"
	"time"

	"github.com/xenom2004/Runtime-Platform-Go/cmd/helpers"
)

type CompletedJob struct {
	JobID   string `json:"job_id"`
	AgentID string `json:"agent_id"`
}

func process(agentId string, job Job) {

	fmt.Printf("Processing Job %s by Agent %s \n", job.JobID, agentId)
	if Agents[agentId].deadstate == true {

		return
	}
	time.Sleep(5 * time.Second)
	Agents[agentId].CurrentJob = nil
	helpers.JsonPost(controlplaneurl+"/JobComplete", CompletedJob{JobID: job.JobID, AgentID: agentId})

}

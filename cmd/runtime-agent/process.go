package main

import (
	"fmt"
	"time"
)

func process(agentId string, job Job) {

	fmt.Printf("Processing Job %s by Agent %s \n", job.JobID, agentId)
	time.Sleep(2 * time.Second)
	fmt.Printf("Job %s completed\n", job.JobID)

}

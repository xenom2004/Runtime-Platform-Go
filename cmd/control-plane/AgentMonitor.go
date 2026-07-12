package main

import (
	"fmt"
	"time"
)

func Monitor() {
	go func() {

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// fmt.Printf("Running for Agent Id: %s\n", agent.Id)
			for _, agent := range agentstore.agents {
				lastAliveTime, err := time.Parse(time.RFC3339, agent.Last_Alive)
				if err != nil {
					fmt.Printf("Error parsing Last Active time: %v \n", err)
					return
				}

				if time.Since(lastAliveTime) > 5*time.Second {
					fmt.Printf("Agent ID %s is dead\n", agent.Id)
					agentstore.Status(agent.Id, AgentDead)
					if agent.CurrentJob != nil {
						var job = *agent.CurrentJob
						agentstore.SetJob(agent.Id, nil)
						Requeue(job)
						fmt.Printf("Job Requeue process initiated for job %s \n", job.JobID)
					}
					return
				}
			}
		}
	}()
}

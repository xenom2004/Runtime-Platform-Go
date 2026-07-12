package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Job struct {
	JobID      string `json:"job_id"`
	Status     string `json:"status"`
	AgentID    string `json:"agent_id"`
	CreatedAt  string `json:"created_at"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

type Agent struct {
	Id           string `json:"id"`
	Time_started string `json:"time_started"`
	Last_Alive   string `json:"last_alive"`
	Status       string `json:"status"`
	pending_jobs chan Job
	CurrentJob   *Job
	deadstate    bool
}

func InitializeAgent(url string, heartbeatUrl string, payload Agent) {

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error json:%v\n", err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("Agent registered successfully")

	go func() {
		for job := range Agents[payload.Id].pending_jobs {
			process(payload.Id, job)
		}
	}()

	go func() {

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		timeout := time.After(50 * time.Second)
		if Agents[payload.Id].deadstate == true {
			timeout = time.After(6 * time.Second)
		}
		for {
			select {
			case <-ticker.C:
				sendHeartbeat(heartbeatUrl,
					Heartbeat{Id: payload.Id, Message: fmt.Sprintf("Agent %s is alive", payload.Id)})
			case <-timeout:
				fmt.Printf("agent ID = %s dead\n", payload.Id)
				return
			}
		}
	}()
}

var Agents = make(map[string]*Agent)
var controlplaneurl = "http://localhost:9000"

func main() {

	url := controlplaneurl + "/HandleAgents"
	heartbeatUrl := controlplaneurl + "/HandleHeartbeat"

	Agents["Agent01"] = &Agent{Id: "Agent01", Time_started: time.Now().Format(time.RFC3339), Status: "alive", pending_jobs: make(chan Job, 10), deadstate: false}
	Agents["Agent02"] = &Agent{Id: "Agent02", Time_started: time.Now().Format(time.RFC3339), Status: "alive", pending_jobs: make(chan Job, 10), deadstate: true}
	Agents["Agent03"] = &Agent{Id: "Agent03", Time_started: time.Now().Format(time.RFC3339), Status: "alive", pending_jobs: make(chan Job, 10), deadstate: false}
	InitializeAgent(url, heartbeatUrl, *Agents["Agent01"])
	InitializeAgent(url, heartbeatUrl, *Agents["Agent02"])
	InitializeAgent(url, heartbeatUrl, *Agents["Agent03"])

	http.HandleFunc("/executejob", executejob)
	port := 9001
	fmt.Printf("Starting Server on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("error while running server :", err)
	}

}

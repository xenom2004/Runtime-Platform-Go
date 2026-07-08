package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Agent struct {
	Id           string `json:"id"`
	Time_started string `json:"time_started"`
	Last_Alive   string `json:"last_alive"`
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

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		timeout := time.After(50 * time.Second)
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

func main() {

	url := "http://localhost:9000/HandleAgents"
	heartbeatUrl := "http://localhost:9000/HandleHeartbeat"
	InitializeAgent(url, heartbeatUrl, Agent{Id: "Agent01", Time_started: time.Now().Format(time.RFC3339)})
	time.Sleep((2000 * time.Millisecond))
	InitializeAgent(url, heartbeatUrl, Agent{Id: "Agent02", Time_started: time.Now().Format(time.RFC3339)})
	time.Sleep((2000 * time.Millisecond))
	InitializeAgent(url, heartbeatUrl, Agent{Id: "Agent03", Time_started: time.Now().Format(time.RFC3339)})
	// InitializeAgent(url, heartbeatUrl, Agent{Id: "Agent01", Time_started: time.Now().Format(time.RFC3339)})
	time.Sleep(20 * time.Second)

}

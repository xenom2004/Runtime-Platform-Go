package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Heartbeat struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func sendHeartbeat(url string, payload Heartbeat) {
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

	fmt.Printf("Agent Id %s sent heartbeat successfully \n", payload.Id)

}

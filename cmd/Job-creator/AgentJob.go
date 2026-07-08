package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Job struct {
	JobID string `json:"job_id"`
}

func SendJob(url string, payload Job) {

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

	fmt.Printf("Job %s sent successfully\n", payload.JobID)
}

func main() {

	url := "http://localhost:9000/jobs"

	SendJob(url, Job{JobID: "Job01"})
	time.Sleep((2000 * time.Millisecond))
	SendJob(url, Job{JobID: "Job02"})
	time.Sleep((2000 * time.Millisecond))
	SendJob(url, Job{JobID: "Job03"})
	time.Sleep((2000 * time.Millisecond))
	SendJob(url, Job{JobID: "Job04"})
	time.Sleep((2000 * time.Millisecond))
	SendJob(url, Job{JobID: "Job05"})

}

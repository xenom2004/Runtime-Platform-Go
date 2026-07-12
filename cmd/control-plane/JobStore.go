package main

import "sync"

type JobStore struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

func (store *JobStore) Set(job Job) {

	store.mu.Lock()
	defer store.mu.Unlock()

	store.jobs[job.JobID] = &job
}

func (store *JobStore) Get(jobID string) (*Job, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	job, exists := store.jobs[jobID]
	return job, exists
}

func (store *JobStore) Update_job_agent(jobID string, agentID string, status string, starttime string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.jobs[jobID].AgentID = agentID
	store.jobs[jobID].Status = status
	store.jobs[jobID].StartedAt = starttime
}

func (store *JobStore) Update_status_and_finishtime(jobID string, status string, finishtime string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.jobs[jobID].Status = status
	store.jobs[jobID].FinishedAt = finishtime
}

var jobstore = &JobStore{jobs: make(map[string]*Job)}

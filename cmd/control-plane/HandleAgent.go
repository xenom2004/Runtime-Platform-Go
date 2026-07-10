package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/xenom2004/Runtime-Platform-Go/cmd/helpers"
)

type Agent struct {
	Id           string `json:"id"`
	Time_started string `json:"time_started"`
	Last_Alive   string `json:"last_alive"`
	Status       string `json:"status"`
}

type AgentStore struct {
	mu             sync.RWMutex
	agents         map[string]*Agent
	last_usedIndex int
	AgentIDs       []string
}

var agentstore = &AgentStore{
	agents:         make(map[string]*Agent),
	AgentIDs:       make([]string, 0),
	last_usedIndex: -1,
}

func (store *AgentStore) Get(id string) (*Agent, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	agent, exists := store.agents[id]
	return agent, exists
}

func (store *AgentStore) Set(agent Agent) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.agents[agent.Id] = &agent
	store.AgentIDs = append(store.AgentIDs, agent.Id)
}

func (store *AgentStore) SetLast_Alive(id string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.agents[id].Last_Alive = time.Now().Format(time.RFC3339)
}

func (store *AgentStore) Dead(id string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.agents[id].Status = "dead"
}

func (store *AgentStore) Next() any {
	store.mu.Lock()
	defer store.mu.Unlock()
	Maxsteps := len(store.AgentIDs)
	for i := (store.last_usedIndex + 1) % len(store.AgentIDs); i < len(store.AgentIDs); i = (i + 1) % len(store.AgentIDs) {
		if Maxsteps == 0 {
			return nil
		}
		Maxsteps--
		if store.agents[store.AgentIDs[i]].Status == "dead" ||
			store.agents[store.AgentIDs[i]].Status == "busy" {
			continue
		}
		store.last_usedIndex = i
		// store.agents[store.AgentIDs[i]].Status = "busy"
		return store.AgentIDs[store.last_usedIndex]
	}
	return nil

}

func handleagent(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, Post only", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var agent Agent

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&agent); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if _, exists := agentstore.Get(agent.Id); exists {
		http.Error(w, "Agent already registered", http.StatusBadRequest)
		return
	}
	agent.Last_Alive = time.Now().Format(time.RFC3339)
	agentstore.Set(agent)

	fmt.Printf("Registered Agent ID = %s and start time =%s \n", agent.Id, agent.Time_started)
	msg := fmt.Sprintf("Agent %s registered successfully started at %s", agent.Id, agent.Time_started)

	go func() {

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// fmt.Printf("Running for Agent Id: %s\n", agent.Id)
			newagent, _ := agentstore.Get(agent.Id)
			lastAliveTime, err := time.Parse(time.RFC3339, newagent.Last_Alive)
			if err != nil {
				fmt.Printf("Error parsing Last Active time: %v \n", err)
				return
			}

			if time.Since(lastAliveTime) > 4*time.Second {
				fmt.Printf("Agent ID %s is dead\n", agent.Id)
				agentstore.Dead(agent.Id)
				return
			}
		}
	}()
	helpers.Jsonhelper(w, map[string]string{"Registered": msg}, http.StatusOK)

}

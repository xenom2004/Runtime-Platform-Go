package main

import (
	"net/http"
	"runtime"
)

type Serverload struct {
	Goroutines int     `json:"goroutines"`
	Memory     float64 `json:"memory_allocated_mb"`
}

func serverload(w http.ResponseWriter, r *http.Request) {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	res := Serverload{
		Goroutines: runtime.NumGoroutine(),
		Memory:     float64(m.Alloc) / (1024 * 1024),
	}

	jsonhelper(w, res, 200)

}

package main

import (
	"net/http"
	"time"

	"github.com/xenom2004/Runtime-Platform-Go/cmd/helpers"
	"github.com/xenom2004/Runtime-Platform-Go/config"
)

type Serverstatus struct {
	Status      string `json:"status"`
	Version     string `json:"version"`
	Appname     string `json:"appname"`
	CurrentTime string `json:"current_time"`
}

func serverstatus(w http.ResponseWriter, r *http.Request) {
	res := Serverstatus{
		Status:      "operational",
		Version:     config.Version,
		Appname:     config.Appname,
		CurrentTime: time.Now().Format(time.RFC3339),
	}
	helpers.Jsonhelper(w, res, 200)

}

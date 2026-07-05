package main

import (
	"fmt"
	"net/http"
)

func serverhealth(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"message": "Go Server is Running, Health normal",
	}
	jsonhelper(w, res, 200)
}

func main() {

	http.HandleFunc("/health", serverhealth)
	http.HandleFunc("/status", serverstatus)
	http.HandleFunc("/load", serverload)
	port := 9000
	fmt.Printf("Starting Server on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("error while running server :", err)
	}

}

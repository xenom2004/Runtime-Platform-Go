package main

import (
	"fmt"
	"net/http"
)

func ServerCondition(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Go Server is Running, condition nominal")
}

func main() {

	http.HandleFunc("/health", ServerCondition)
	port := 9000
	fmt.Printf("Starting Server on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("error while running server :", err)
	}

}

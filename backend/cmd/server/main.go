package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jekabso21/IoT-Dashboard/backend/internal/api"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the IoT Backend!")
	})
	// New API route
	http.HandleFunc("/api/device-data", api.DeviceDataHandler)

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

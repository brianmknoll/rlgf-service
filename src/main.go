package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func main() {
	s := NewRlgServer()
	s.runMigrations()

	mux := http.NewServeMux()

	mux.HandleFunc("/events", s.eventsPostHandler)

	err := http.ListenAndServe(":8888", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

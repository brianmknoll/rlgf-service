package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/brianmknoll/rlgf-service/internal/db"
)

type ApiEvent struct {
	GuildId string `json:"guildId"`
	Name    string `json:"name"`
}

func main() {
	db := db.NewFirestoreDatabase()

	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var e ApiEvent

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&e)
		if err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Write to the database.
		err = db.CreateEvent(e.GuildId, e.Name)
		if err != nil {
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Create the Discord event.
		// Commented out for now. We should do this off of a Firestore write event instead.
		// discord.CreateDiscordEvent(e.Name)

		// TODO:
		// 1. Consider a DB rollback if the Discord request fails.
		// 2. Make Discord event optional.

		fmt.Printf("Received: Name=%s\n", e.Name)
		w.WriteHeader(http.StatusCreated)
	})

	err := http.ListenAndServe(":8888", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/brianmknoll/rlgf-service/internal/db"
	"github.com/brianmknoll/rlgf-service/internal/discord"
)

type ApiEvent struct {
	Name string `json:"name"`
}

func main() {
	discord := discord.NewDiscordClient()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := db.NewDynamoDatabase(sess)
	db.RunMigrations()

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
		err = db.CreateEvent(e.Name)
		if err != nil {
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Create the Discord event.
		discord.CreateDiscordEvent(e.Name)

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

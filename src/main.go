package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brianmknoll/rlgf-service/internal/db"
	"github.com/brianmknoll/rlgf-service/internal/discord"
)

type ApiMessage struct {
	Channel string `json:"channel"`
	GuildId string `json:"guildId"`
	Message string `json:"message"`
	Epoch   int64  `json:"timestamp"`
}

type ApiEvent struct {
	GuildId string `json:"guildId"`
	Name    string `json:"name"`
}

func main() {
	d := discord.NewDiscordClient()
	database := db.NewFirestoreDatabase()

	mux := http.NewServeMux()

	mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var m ApiMessage

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&m)
		if err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusUnprocessableEntity)
		}
		defer r.Body.Close()

		fmt.Printf("Received %v\n", m)

		seconds := m.Epoch / 1000
		nanos := (m.Epoch % 1000) * int64(time.Millisecond)

		newMsg := db.DbMessage{
			Message:   m.Message,
			Timestamp: time.Unix(seconds, nanos),
		}

		recentMessages, err := database.ReadRecentMessages(m.GuildId, m.Channel)
		if err != nil {
			log.Printf("Failed to read recent messages: %v\n", err)
			http.Error(w, "Failed to read recent messages", http.StatusInternalServerError)
			return
		}

		err = database.CreateMessage(m.GuildId, m.Channel, newMsg)
		if err != nil {
			log.Printf("Failed to create new message: %v\n", err)
			http.Error(w, "Failed to create new message", http.StatusInternalServerError)
			return
		}

		allMessages := append(recentMessages, newMsg)
		log.Printf("Total messages to send: %d\n", len(allMessages))

		jsonResponse, err := json.Marshal(allMessages)
		if err != nil {
			log.Printf("Error marshaling messages to JSON: %v", err)
			http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(jsonResponse)
		if err != nil {
			// If writing the response fails (e.g., client disconnected), log it.
			// It's often too late to send an HTTP error code here as headers might have been sent.
			log.Printf("Error writing JSON response to ResponseWriter: %v", err)
		}

		log.Println("Successfully sent JSON response with messages.")
	})

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

		fmt.Printf("Received: %v\n", e)

		err = database.CreateEvent(e.GuildId, e.Name)
		if err != nil {
			fmt.Printf("Internal server error creating DB event: %v\n", err.Error())
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Create the Discord event.
		// TODO: We should do this off of a Firestore write event instead.
		err = d.CreateDiscordEvent(e.GuildId, e.Name)
		if err != nil {
			fmt.Printf("Internal server error creating Discord event: %v\n", err.Error())
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

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

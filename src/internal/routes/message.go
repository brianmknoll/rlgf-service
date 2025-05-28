package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/brianmknoll/rlgf-service/internal/db"
)

type ApiMessage struct {
	Author  string `json:"author"`
	GuildId string `json:"guildId"`
	Channel string `json:"channel"`
	Message string `json:"message"`
	Epoch   int64  `json:"timestamp"`
}

func (router *RlgfRouter) HandleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		router.createMessage(w, r)
	} else if r.Method == http.MethodGet {
		router.getMessages(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router *RlgfRouter) createMessage(w http.ResponseWriter, r *http.Request) {
	var m ApiMessage

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received message %v\n", m)

	seconds := m.Epoch / 1000
	nanos := (m.Epoch % 1000) * int64(time.Millisecond)

	newMsg := db.DbMessage{
		Author:    m.Author,
		Message:   m.Message,
		Timestamp: time.Unix(seconds, nanos),
	}

	err = router.database.CreateMessage(m.GuildId, m.Channel, newMsg)
	if err != nil {
		log.Printf("Failed to create new message: %v\n", err)
		if errors.Is(err, db.ErrAlreadyExists) {
			http.Error(w, "Message already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create new message", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (router *RlgfRouter) getMessages(w http.ResponseWriter, r *http.Request) {
	guildId := r.URL.Query().Get("guild_id")
	channelId := r.URL.Query().Get("channel_id")
	msgs, err := router.database.ReadRecentMessages(guildId, channelId)
	if err != nil {
		http.Error(w, "Failed to read messages", http.StatusInternalServerError)
	}
	jsonResponse, err := json.Marshal(msgs)
	if err != nil {
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
}

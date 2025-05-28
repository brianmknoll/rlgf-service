package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/brianmknoll/rlgf-service/internal/db"
)

type ApiMemory struct {
	GuildId string `json:"guildId"`
	Memory  string `json:"memory"`
}

func (router *RlgfRouter) HandleMemory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		router.createMemory(w, r)
	} else if r.Method == http.MethodGet {
		router.getMemory(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (router *RlgfRouter) createMemory(w http.ResponseWriter, r *http.Request) {
	var m ApiMemory

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusUnprocessableEntity)
	}
	defer r.Body.Close()

	fmt.Printf("Received memory %v\n", m)

	memory := db.DbMemory{
		Memory: m.Memory,
	}

	err = router.database.CreateMemory(m.GuildId, memory.Memory)
	if err != nil {
		log.Printf("Failed to create new message: %v\n", err)
		http.Error(w, "Failed to create new message", http.StatusInternalServerError)
		return
	}
}

func (router *RlgfRouter) getMemory(w http.ResponseWriter, r *http.Request) {
	guildId := r.URL.Query().Get("guild_id")
	memories, err := router.database.ReadMemories(guildId)
	if err != nil {
		log.Printf("Failed to read memories: %v\n", err)
		http.Error(w, "Failed to read memories", http.StatusInternalServerError)
		return
	}
	apiMemory := ApiMemory{
		GuildId: guildId,
		Memory:  memories,
	}
	jsonResponse, err := json.Marshal(apiMemory)
	if err != nil {
		log.Printf("Error marshaling memory to JSON: %v", err)
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		// If writing the response fails (e.g., client disconnected), log it.
		// It's often too late to send an HTTP error code here as headers might have been sent.
		log.Printf("Error writing JSON response to ResponseWriter: %v", err)
		return
	}
	log.Println("Successfully sent JSON response with messages.")
}

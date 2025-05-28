package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiEvent struct {
	GuildId string `json:"guildId"`
	Name    string `json:"name"`
}

func (router *RlgfRouter) HandleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		router.createEvent(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router *RlgfRouter) createEvent(w http.ResponseWriter, r *http.Request) {
	var e ApiEvent

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&e)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received: %v\n", e)

	err = router.database.CreateEvent(e.GuildId, e.Name)
	if err != nil {
		fmt.Printf("Internal server error creating DB event: %v\n", err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the Discord event.
	// TODO: We should do this off of a Firestore write event instead.
	err = (*router.discord).CreateDiscordEvent(e.GuildId, e.Name)
	if err != nil {
		fmt.Printf("Internal server error creating Discord event: %v\n", err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received: Name=%s\n", e.Name)

	w.WriteHeader(http.StatusCreated)
}

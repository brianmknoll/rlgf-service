package routes

import (
	"encoding/json"
	"net/http"
)

type ApiGuild struct {
	Users map[string]string
}

func (router *RlgfRouter) HandleGuild(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		router.getGuild(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router *RlgfRouter) getGuild(w http.ResponseWriter, r *http.Request) {
	guildId := r.URL.Query().Get("guild_id")
	guild, err := router.database.GetGuild(guildId)
	if err != nil {
		http.Error(w, "No guild found", http.StatusNotFound)
	}
	g := ApiGuild{
		Users: guild.Users,
	}
	jsonResponse, err := json.Marshal(g)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
	}
}

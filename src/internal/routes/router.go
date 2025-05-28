package routes

import (
	"net/http"

	"github.com/brianmknoll/rlgf-service/internal/db"
	"github.com/brianmknoll/rlgf-service/internal/discord"
)

type Router interface {
	HandleEvent(w http.ResponseWriter, r *http.Request)
	HandleMessage(w http.ResponseWriter, r *http.Request)
	HandleMemory(w http.ResponseWriter, r *http.Request)
	HandleGuild(w http.ResponseWriter, r *http.Request)
}

type RlgfRouter struct {
	database db.Database
	discord  *discord.DiscordClient
}

func NewRouter(database db.Database, discord *discord.DiscordClient) Router {
	return &RlgfRouter{
		database: database,
		discord:  discord,
	}
}

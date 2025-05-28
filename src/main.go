package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/brianmknoll/rlgf-service/internal/db"
	"github.com/brianmknoll/rlgf-service/internal/discord"
	"github.com/brianmknoll/rlgf-service/internal/routes"
)

func main() {
	router := routes.NewRouter(
		db.NewFirestoreDatabase(),
		discord.NewDiscordClient(),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/events", router.HandleEvent)
	mux.HandleFunc("/guilds", router.HandleGuild)
	mux.HandleFunc("/memory", router.HandleMemory)
	mux.HandleFunc("/message", router.HandleMessage)

	err := http.ListenAndServe(":8888", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

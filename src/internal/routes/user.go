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

type ApiUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Joined   int64  `json:"joined"`
}

func (router *RlgfRouter) HandleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		router.createUser(w, r)
	} else if r.Method == http.MethodGet {
		router.getUsers(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router *RlgfRouter) createUser(w http.ResponseWriter, r *http.Request) {
	var user ApiUser

	guildId := r.PathValue("guildId")
	fmt.Printf("Creating user for guild %s\n", guildId)
	if guildId == "" {
		http.Error(w, "Guild ID is required", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received user %v\n", user)

	seconds := user.Joined / 1000
	nanos := (user.Joined % 1000) * int64(time.Millisecond)

	newUser := db.DbUser{
		Id:       user.Id,
		Username: user.Username,
		Name:     user.Name,
		Joined:   time.Unix(seconds, nanos),
	}

	err = router.database.CreateUser(guildId, newUser)

	if err != nil {
		log.Printf("Failed to create new user: %v\n", err)
		if errors.Is(err, db.ErrAlreadyExists) {
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create new user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (router *RlgfRouter) getUsers(w http.ResponseWriter, r *http.Request) {
	guildId := r.PathValue("guildId")
	fmt.Printf("Creating user for guild %s\n", guildId)
	if guildId == "" {
		http.Error(w, "Guild ID is required", http.StatusBadRequest)
		return
	}

	dbUsers, err := router.database.ReadUsers(guildId)
	if err != nil {
		http.Error(w, "Failed to read users", http.StatusInternalServerError)
	}

	var apiUsers []ApiUser
	for _, dbUser := range dbUsers {
		apiUser := ApiUser{
			Id:       dbUser.Id,
			Username: dbUser.Username,
			Name:     dbUser.Name,
			Joined:   dbUser.Joined.UnixNano() / int64(time.Millisecond),
		}
		apiUsers = append(apiUsers, apiUser)
	}

	jsonResponse, err := json.Marshal(apiUsers)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResponse)
	if err != nil {
		// If writing the response fails (e.g., client disconnected), log it.
		// It's often too late to send an HTTP error code here as headers might have been sent.
		log.Printf("Error writing JSON response to ResponseWriter: %v", err)
	}
}

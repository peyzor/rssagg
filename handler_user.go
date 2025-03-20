package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/peyzor/rssagg/internal/auth"
	"github.com/peyzor/rssagg/internal/database"
	"net/http"
	"time"
)

func (sc *ServerConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error parsing json: %v", err))
		return
	}

	user, err := sc.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("couldn't create user: %v", err))
		return
	}

	responseWithJson(w, 201, databaseUserToUser(user))
}

func (sc *ServerConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		responseWithError(w, 403, fmt.Sprintf("auth error: %v", err))
		return
	}

	user, err := sc.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("user not found: %v", err))
		return
	}

	responseWithJson(w, 200, databaseUserToUser(user))
}

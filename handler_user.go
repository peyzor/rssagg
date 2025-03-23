package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/peyzor/rssagg/internal/database"
	"net/http"
	"time"
)

func (sc *serverConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
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

func (sc *serverConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	responseWithJson(w, 200, databaseUserToUser(user))
}

func (sc *serverConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := sc.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprint("couldn't retrieve posts for user: ", err))
		return
	}

	responseWithJson(w, 200, databasePostsToPosts(posts))
}

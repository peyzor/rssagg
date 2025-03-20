package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/peyzor/rssagg/internal/database"
	"net/http"
	"time"
)

func (sc *serverConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error parsing json: %v", err))
		return
	}

	feed, err := sc.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("couldn't create feed: %v", err))
		return
	}

	responseWithJson(w, 201, databaseFeedToFeed(feed))
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/peyzor/rssagg/internal/database"
	"net/http"
	"time"
)

func (sc *serverConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error parsing json: %v", err))
		return
	}

	feedFollow, err := sc.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error creating feed follow: %v", err))
		return
	}

	responseWithJson(w, 200, databaseFeedFollowToFeedFollow(feedFollow))
}

func (sc *serverConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := sc.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("couldn't get feed follows: %v", err))
		return
	}

	responseWithJson(w, 200, databaseFeedFollowsToFeedFollows(feedFollows))
}

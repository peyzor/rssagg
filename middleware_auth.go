package main

import (
	"fmt"
	"github.com/peyzor/rssagg/internal/auth"
	"github.com/peyzor/rssagg/internal/database"
	"net/http"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (sc *serverConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		handler(w, r, user)
	}
}

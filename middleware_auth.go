package main

import (
	"database/sql"
	"net/http"
	"rss_server/internal/database"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		auth := strings.Split(authHeader, " ")
		if len(auth) != 2 {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if auth[0] != "ApiKey" {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		apiKey := auth[1]
		if apiKey == "" {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			respondWithError(w, http.StatusInternalServerError, "error getting user")
			return
		}

		handler(w, r, user)
	})
}

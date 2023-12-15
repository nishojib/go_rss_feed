package main

import (
	"encoding/json"
	"net/http"
	"rss_server/internal/database"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type request struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		FeedID    uuid.UUID `json:"feed_id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	follow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		UserID:    user.ID,
		FeedID:    req.FeedID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:        follow.ID,
		FeedID:    follow.FeedID,
		UserID:    follow.UserID,
		CreatedAt: follow.CreatedAt,
		UpdatedAt: follow.UpdatedAt,
	})
}

func (cfg *apiConfig) deleteFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID := chi.URLParam(r, "feedFollowID")

	params := database.DeleteFeedFollowParams{
		ID:     uuid.MustParse(feedFollowID),
		UserID: user.ID,
	}

	err := cfg.DB.DeleteFeedFollow(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func (cfg *apiConfig) getFeedFollowsHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type response struct {
		ID        uuid.UUID `json:"id"`
		FeedID    uuid.UUID `json:"feed_id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	feedFollows, err := cfg.DB.GetFeedFollowsByUserId(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting feed follows")
		return
	}

	var responses []response
	for _, feedFollow := range feedFollows {
		responses = append(responses, response{
			ID:        feedFollow.ID,
			FeedID:    feedFollow.FeedID,
			UserID:    feedFollow.UserID,
			CreatedAt: feedFollow.CreatedAt,
			UpdatedAt: feedFollow.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, responses)
}

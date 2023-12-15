package main

import (
	"net/http"
	"rss_server/internal/database"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getPostsByUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type response struct {
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Title       string    `json:"title"`
		Url         string    `json:"url"`
		Description string    `json:"description"`
		PublishedAt time.Time `json:"published_at"`
		FeedID      uuid.UUID `json:"feed_id"`
	}

	posts, err := cfg.DB.GetPostsByUserId(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var responses []response
	for _, post := range posts {
		responses = append(responses, response{
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			Url:         post.Url,
			Description: post.Description,
			PublishedAt: post.PublishedAt,
			FeedID:      post.FeedID,
		})
	}

	respondWithJSON(w, http.StatusOK, responses)
}

-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: CreateFeedFollow :one
INSERT INTO feeds_users (created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feeds_users WHERE id = $1 AND user_id = $2;

-- name: GetFeedFollowsByUserId :many
SELECT * FROM feeds_users WHERE user_id = $1;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;
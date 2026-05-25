-- name: CreateFeedFollow :one
INSERT INTO feed_follows (user_id, feed_id)
VALUES (
    sqlc.arg(user_id),
    sqlc.arg(feed_id)
)
RETURNING *;

-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows
WHERE user_id = sqlc.arg(user_id);

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);
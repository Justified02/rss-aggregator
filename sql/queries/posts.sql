-- name: CreatePost :one
INSERT INTO posts (
    feed_id,
    title,
    url,
    description,
    published_at
) VALUES (
    sqlc.arg(feed_id),
    sqlc.arg(title),
    sqlc.arg(url),
    sqlc.arg(description),
    sqlc.arg(published_at)
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed.id
WHERE feed_follows.user_id = sqlc.arg(user_id)
ORDER BY posts.published_at DESC
LIMIT sqlc.arg(max_posts);
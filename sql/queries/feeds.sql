-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedList :many
SELECT feeds.name as feed_name, feeds.url as url, users.name as creator_name
FROM feeds, users
WHERE feeds.user_id = users.id;
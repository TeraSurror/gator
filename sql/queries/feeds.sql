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

-- name: CreateFeedFollow :one
WITH inserted_follow_feed AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT iff.id, iff.created_at, iff.updated_at, iff.user_id, iff.feed_id, u.name as user_name, f.name as feed_name
FROM inserted_follow_feed iff, users u, feeds f
WHERE iff.user_id = u.id AND iff.feed_id = f.id;

-- name: GetFeedByURL :one
SELECT id, created_at, updated_at, name, url, user_id FROM feeds WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT ff.id, ff.user_id, ff.feed_id, u.name as user_name, f.name as feed_name
FROM feed_follows ff, users u, feeds f
WHERE ff.user_id = u.id AND ff.feed_id = f.id AND u.name = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds SET updated_at = NOW(), last_fetched_at = NOW() WHERE id = $1 RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;
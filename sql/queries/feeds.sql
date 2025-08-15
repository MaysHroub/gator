-- name: CreateFeed :one
INSERT INTO
    feeds (id, created_at, updated_at, name, url, user_id)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
*;

-- name: GetAllFeeds :many
SELECT
    f.name AS feedName,
    f.url,
    u.name AS userName
FROM
    feeds f
    JOIN users u ON u.id=f.user_id;

-- name: GetFeedByURL :one
SELECT
    *
FROM
    feeds
WHERE
    feeds.url=$1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET
    last_fetched_at=NOW(),
    updated_at=NOW()
WHERE
    id=$1;

-- name: GetNextFeedToFetch :one
SELECT
    id,
    name,
    url,
    last_fetched_at
FROM
    feeds
ORDER BY
    last_fetched_at ASC NULLS FIRST
LIMIT
    1;
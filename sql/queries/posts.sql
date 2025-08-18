-- name: CreatePost :one
INSERT INTO
    posts (
        id,
        created_at,
        updated_at,
        title,
        url,
        description,
        published_at,
        feed_id
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
*;

-- name: GetPostsForUser :many
SELECT
    p.*
FROM
    posts p
    JOIN feed_follows ff ON ff.feed_id=p.feed_id
    JOIN users u ON u.user_id=ff.user_id
WHERE
    u.name=$1
LIMIT
    $2;
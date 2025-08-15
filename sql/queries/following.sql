-- name: CreateFeedFollow :many
WITH
    inserted_feed_follow AS (
        INSERT INTO
            feed_follows (id, created_at, updated_at, user_id, feed_id)
        VALUES
            ($1, $2, $3, $4, $5)
        RETURNING *
    )
SELECT
    iff.*,
    f.name AS feed_name,
    u.name AS user_name
FROM
    inserted_feed_follow iff
    JOIN users u ON u.id=iff.user_if
    JOIN feeds f ON f.id=iff.feed_id
GROUP BY
    iff.id;
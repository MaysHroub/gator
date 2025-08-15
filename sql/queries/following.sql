-- name: CreateFeedFollow :many
WITH
    inserted_feed_follow AS (
        INSERT INTO
            feed_follows (id, created_at, updated_at, user_id, feed_id)
        VALUES
            ($1, $2, $3, $4, $5)
        RETURNING
*
    )
SELECT
    iff.*,
    f.name AS feed_name,
    u.name AS user_name
FROM
    inserted_feed_follow iff
    JOIN users u ON u.id=iff.user_id
    JOIN feeds f ON f.id=iff.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.*,
    f.name AS feed_name
FROM
    feed_follows ff
    JOIN users u ON u.id=ff.user_id
    JOIN feeds f ON f.id=ff.feed_id
WHERE
    u.name=$1;

-- name: DeleteFeedFollowByUserAndURL :exec
DELETE FROM feed_follows
WHERE
    user_id=(
        SELECT
            id
        FROM
            users
        WHERE
            users.name=$1
    )
    AND feed_id=(
        SELECT
            id
        FROM
            feeds
        WHERE
            feeds.url=$2
    );

-- Follow queries
-- name: CreateFollow :one
INSERT INTO follows (follower_id, following_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteFollow :exec
DELETE FROM follows WHERE follower_id = $1 AND following_id = $2;

-- name: GetFollow :one
SELECT * FROM follows WHERE follower_id = $1 AND following_id = $2;

-- name: ListFollowers :many
SELECT u.* FROM users u
JOIN follows f ON u.id = f.follower_id
WHERE f.following_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListFollowing :many
SELECT u.* FROM users u
JOIN follows f ON u.id = f.following_id
WHERE f.follower_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;

-- name: IsFollowing :one
SELECT EXISTS(
    SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2
);

-- name: CountFollowers :one
SELECT COUNT(*) FROM follows WHERE following_id = $1;

-- name: CountFollowing :one
SELECT COUNT(*) FROM follows WHERE follower_id = $1;

-- Like queries
-- name: CreateLike :one
INSERT INTO likes (user_id, log_id, review_id) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteLike :exec
DELETE FROM likes WHERE user_id = $1 AND ((log_id = $2) OR (review_id = $3));

-- name: GetLikeByLog :one
SELECT * FROM likes WHERE user_id = $1 AND log_id = $2;

-- name: GetLikeByReview :one
SELECT * FROM likes WHERE user_id = $1 AND review_id = $2;

-- name: CountLikesForLog :one
SELECT COUNT(*) FROM likes WHERE log_id = $1;

-- name: CountLikesForReview :one
SELECT COUNT(*) FROM likes WHERE review_id = $1;

-- name: HasUserLikedLog :one
SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND log_id = $2);

-- name: HasUserLikedReview :one
SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND review_id = $2);

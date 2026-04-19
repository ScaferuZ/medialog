-- name: CreateReview :one
INSERT INTO reviews (user_id, media_id, log_id, title, content, rating, contains_spoilers)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetReviewByID :one
SELECT * FROM reviews WHERE id = $1;

-- name: ListReviewsByUser :many
SELECT * FROM reviews WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListReviewsByMedia :many
SELECT r.*, u.username, u.display_name 
FROM reviews r
JOIN users u ON r.user_id = u.id
WHERE r.media_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateReview :one
UPDATE reviews SET
    title = COALESCE($2, title),
    content = COALESCE($3, content),
    rating = COALESCE($4, rating),
    contains_spoilers = COALESCE($5, contains_spoilers),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteReview :exec
DELETE FROM reviews WHERE id = $1;

-- name: GetReviewByLog :one
SELECT * FROM reviews WHERE log_id = $1;

-- name: CreateLog :one
INSERT INTO logs (
    user_id, media_id, status, rating, started_at, completed_at,
    rewatch_count, progress, total, note, contains_spoilers
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetLogByID :one
SELECT * FROM logs WHERE id = $1;

-- name: GetLogByUserAndMedia :one
SELECT * FROM logs WHERE user_id = $1 AND media_id = $2;

-- name: ListLogsByUser :many
SELECT * FROM logs 
WHERE user_id = $1 
AND (NULLIF($2::varchar, '') IS NULL OR status = $2)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListLogsByMedia :many
SELECT l.*, u.username, u.display_name FROM logs l
JOIN users u ON l.user_id = u.id
WHERE l.media_id = $1 AND l.status = 'completed'
ORDER BY l.rating DESC NULLS LAST, l.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateLog :one
UPDATE logs SET
    status = COALESCE($2, status),
    rating = COALESCE($3, rating),
    started_at = COALESCE($4, started_at),
    completed_at = COALESCE($5, completed_at),
    rewatch_count = COALESCE($6, rewatch_count),
    progress = COALESCE($7, progress),
    total = COALESCE($8, total),
    note = COALESCE($9, note),
    contains_spoilers = COALESCE($10, contains_spoilers),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteLog :exec
DELETE FROM logs WHERE id = $1;

-- name: GetUserStats :one
SELECT 
    COUNT(*) FILTER (WHERE status = 'completed') as completed_count,
    COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress_count,
    COUNT(*) FILTER (WHERE status = 'planned') as planned_count,
    COUNT(*) FILTER (WHERE status = 'dropped') as dropped_count,
    COALESCE(AVG(rating) FILTER (WHERE rating IS NOT NULL), 0) as average_rating,
    COUNT(DISTINCT media_id) as total_media
FROM logs 
WHERE user_id = $1;

-- name: GetTimeline :many
-- Get logs from users that $1 follows
SELECT l.*, u.username, u.display_name FROM logs l
JOIN users u ON l.user_id = u.id
WHERE l.user_id IN (
    SELECT following_id FROM follows WHERE follower_id = $1
)
OR l.user_id = $1
ORDER BY l.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLatestPublicActivity :many
SELECT 
    l.id,
    l.media_id,
    l.status,
    l.rating,
    l.note,
    l.created_at,
    u.username,
    u.display_name,
    m.title,
    m.cover_image,
    m.type
FROM logs l
JOIN users u ON l.user_id = u.id
JOIN media m ON l.media_id = m.id
WHERE l.status = 'completed'
  AND u.is_public = true
ORDER BY l.created_at DESC
LIMIT $1 OFFSET $2;

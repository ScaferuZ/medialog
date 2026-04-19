-- name: CreateMedia :one
INSERT INTO media (
    type, title, original_title, description, cover_image,
    release_date, metadata, tmdb_id, mal_id, google_books_id, igdb_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetMediaByID :one
SELECT * FROM media WHERE id = $1;

-- name: GetMediaByTMDBID :one
SELECT * FROM media WHERE tmdb_id = $1;

-- name: SearchMedia :many
SELECT * FROM media 
WHERE (
    to_tsvector('english', title) @@ plainto_tsvector('english', $1)
    OR title ILIKE '%' || $1 || '%'
)
AND ($2::varchar IS NULL OR type = $2)
ORDER BY 
    CASE WHEN title ILIKE $1 || '%' THEN 0 ELSE 1 END,
    created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListMedia :many
SELECT * FROM media 
WHERE ($1::varchar IS NULL OR type = $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateMedia :one
UPDATE media SET
    title = COALESCE($2, title),
    original_title = COALESCE($3, original_title),
    description = COALESCE($4, description),
    cover_image = COALESCE($5, cover_image),
    release_date = COALESCE($6, release_date),
    metadata = COALESCE($7, metadata),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteMedia :exec
DELETE FROM media WHERE id = $1;

-- name: CreateGenre :one
INSERT INTO genres (name) VALUES ($1) RETURNING *;

-- name: GetGenreByName :one
SELECT * FROM genres WHERE name = $1;

-- name: ListGenres :many
SELECT * FROM genres ORDER BY name;

-- name: AddMediaGenre :exec
INSERT INTO media_genres (media_id, genre_id) VALUES ($1, $2);

-- name: GetMediaGenres :many
SELECT g.* FROM genres g
JOIN media_genres mg ON g.id = mg.genre_id
WHERE mg.media_id = $1;

-- name: RemoveMediaGenre :exec
DELETE FROM media_genres WHERE media_id = $1 AND genre_id = $2;

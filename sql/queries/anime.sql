-- name: CreateAnime :one
INSERT INTO anime(
	created_at,
	updated_at,
	romaji_name,
	japanese_name,
	english_name,
	type,
	release_date
)
VALUES(NOW(), NOW(), $1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAnime :many
SELECT * FROM anime;

-- name: GetAnimeById :one
SELECT * FROM anime WHERE id = $1;

-- name: DeleteAnimeById :exec
DELETE from anime WHERE id = $1;

-- name: UpdateAnimeById :exec
UPDATE anime 
SET updated_at = NOW(), romaji_name = $1, japanese_name = $2, english_name = $3, type = $4, release_date = $5;

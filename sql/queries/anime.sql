-- name: InsertAnime :one
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

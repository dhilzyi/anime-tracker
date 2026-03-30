-- +goose Up
CREATE TABLE anime(
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	romaji_name TEXT NOT NULL,
	japanese_name TEXT,
	english_name TEXT,
	type TEXT,
	release_date TIMESTAMP
);
-- +goose Down
DROP TABLE anime;

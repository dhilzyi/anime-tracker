package api

import "time"

type CreateAnimeRequest struct {
	RomajiName   string    `json:"romaji_name"`
	JapaneseName *string   `json:"japanese_name"`
	EnglishName  *string   `json:"english_name"`
	Type         *string   `json:"type"`
	ReleaseDate  time.Time `json:"release_date"`
}

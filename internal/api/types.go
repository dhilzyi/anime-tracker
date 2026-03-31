package api

type Anime struct {
	Id           int32   `json:"id"`
	RomajiName   string  `json:"romaji_name"`
	JapaneseName *string `json:"japanese_name"`
	EnglishName  *string `json:"english_name"`
	Type         *string `json:"type"`
	ReleaseDate  *string `json:"release_date"`
}

type UpdateAnimeRequest struct {
	RomajiName   *string `json:"romaji_name"`
	JapaneseName *string `json:"japanese_name"`
	EnglishName  *string `json:"english_name"`
	Type         *string `json:"type"`
	ReleaseDate  *string `json:"release_date"`
}

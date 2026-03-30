package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dhilzyi/anime-tracker/internal/database"
)

type Handler struct {
	db *database.Queries
}

func NewHandler(db *database.Queries) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetAnime(w http.ResponseWriter, r *http.Request) {
	dataAnime, err := h.db.GetAnime(context.Background())
	if err != nil {
		respondWithError(w, 500, "", err)
		return
	}

	var remapDataAnime []CreateAnimeRequest
	for i := range dataAnime {
		instance := CreateAnimeRequest{
			RomajiName:   dataAnime[i].RomajiName,
			EnglishName:  &dataAnime[i].EnglishName.String,
			JapaneseName: &dataAnime[i].JapaneseName.String,
			Type:         &dataAnime[i].Type.String,
			ReleaseDate:  dataAnime[i].ReleaseDate.Time,
		}
		remapDataAnime = append(remapDataAnime, instance)
	}

	respondWithJSON(w, 200, remapDataAnime)
}

func (h *Handler) PostAnime(w http.ResponseWriter, r *http.Request) {
	var req CreateAnimeRequest
	if err := decodeJson(r.Body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "", err)
		return
	}

	params := database.InsertAnimeParams{
		RomajiName:   req.RomajiName,
		JapaneseName: toNullString(req.JapaneseName),
		EnglishName:  toNullString(req.EnglishName),
		Type:         toNullString(req.Type),
		ReleaseDate:  toNullTime(req.ReleaseDate),
	}

	if _, err := h.db.InsertAnime(context.Background(), params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
	}
	w.WriteHeader(204)
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func decodeJson(raw io.Reader, placeholder any) error {
	decoder := json.NewDecoder(raw)
	err := decoder.Decode(placeholder)
	if err != nil {
		return err
	}

	return nil
}

func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func toNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Time: t, Valid: true}
}

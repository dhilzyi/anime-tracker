package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

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

	var remapDataAnime []Anime
	for i := range dataAnime {
		instance := Anime{
			RomajiName:   dataAnime[i].RomajiName,
			EnglishName:  &dataAnime[i].EnglishName.String,
			JapaneseName: &dataAnime[i].JapaneseName.String,
			Type:         &dataAnime[i].Type.String,
			ReleaseDate:  validateTime(dataAnime[i].ReleaseDate),
		}

		remapDataAnime = append(remapDataAnime, instance)
	}

	respondWithJSON(w, 200, remapDataAnime)
}

func (h *Handler) PostAnime(w http.ResponseWriter, r *http.Request) {
	var req Anime
	if err := decodeJson(r.Body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "", err)
		return
	}

	params := database.CreateAnimeParams{
		RomajiName:   req.RomajiName,
		JapaneseName: toNullString(req.JapaneseName),
		EnglishName:  toNullString(req.EnglishName),
		Type:         toNullString(req.Type),
		ReleaseDate:  toNullTime(req.ReleaseDate),
	}

	if _, err := h.db.CreateAnime(context.Background(), params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
	}
	w.WriteHeader(204)
}

func (h *Handler) DeleteAnimeById(w http.ResponseWriter, r *http.Request) {
	rawAnimeID := r.PathValue("animeID")
	animeID, err := strconv.Atoi(rawAnimeID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err := h.db.DeleteAnimeById(context.Background(), int32(animeID)); err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Anime by following id is not exist", err)
			return
		}
		respondWithError(w, 500, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAnimeById(w http.ResponseWriter, r *http.Request) {
	rawAnimeID := r.PathValue("animeID")
	animeID, err := strconv.Atoi(rawAnimeID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	animeData, err := h.db.GetAnimeById(context.Background(), int32(animeID))
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Anime by following id is not exist", err)
			return
		}
		respondWithError(w, 500, err.Error(), err)
		return
	}
	remap := Anime{
		RomajiName:   animeData.RomajiName,
		EnglishName:  &animeData.EnglishName.String,
		JapaneseName: &animeData.JapaneseName.String,
		Type:         &animeData.Type.String,
		ReleaseDate:  validateTime(animeData.ReleaseDate),
	}

	respondWithJSON(w, http.StatusOK, remap)

}

func (h *Handler) UpdateAnimeById(w http.ResponseWriter, r *http.Request) {
	rawAnimeID := r.PathValue("animeID")
	animeID, err := strconv.Atoi(rawAnimeID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	inputParam := database.UpdateAnimeByIdParams{}
}

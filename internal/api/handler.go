package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
			Id:           dataAnime[i].ID,
			RomajiName:   dataAnime[i].RomajiName,
			EnglishName:  &dataAnime[i].EnglishName.String,
			JapaneseName: &dataAnime[i].JapaneseName.String,
			Type:         &dataAnime[i].Type.String,
			ReleaseDate:  validateTime(dataAnime[i].ReleaseDate),
		}

		remapDataAnime = append(remapDataAnime, instance)
	}

	respondWithJSON(w, http.StatusOK, remapDataAnime)
}

func (h *Handler) PostAnime(w http.ResponseWriter, r *http.Request) {
	var req Anime
	if err := decodeJson(r.Body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	if strings.TrimSpace(req.RomajiName) == "" {
		respondWithError(w, http.StatusBadRequest, "romaji_name can't be empty", fmt.Errorf("romaji_name can't be empty"))
		return
	}

	params := database.CreateAnimeParams{
		RomajiName:   req.RomajiName,
		JapaneseName: toNullString(req.JapaneseName),
		EnglishName:  toNullString(req.EnglishName),
		Type:         toNullString(req.Type),
		ReleaseDate:  toNullTime(req.ReleaseDate),
	}

	animeData, err := h.db.CreateAnime(context.Background(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	remap := Anime{
		Id:           animeData.ID,
		RomajiName:   animeData.RomajiName,
		EnglishName:  &animeData.EnglishName.String,
		JapaneseName: &animeData.JapaneseName.String,
		Type:         &animeData.Type.String,
		ReleaseDate:  validateTime(animeData.ReleaseDate),
	}

	respondWithJSON(w, http.StatusOK, remap)
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
		Id:           animeData.ID,
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

	ctx := context.Background()
	animeData, err := h.db.GetAnimeById(ctx, int32(animeID))
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Anime by following id is not exist", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	var req UpdateAnimeRequest
	if err := decodeJson(r.Body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid json", err)
		return
	}

	if req.RomajiName != nil {
		animeData.RomajiName = *req.RomajiName
	}
	if req.JapaneseName != nil {
		animeData.JapaneseName = toNullString(req.JapaneseName)
	}
	if req.EnglishName != nil {
		animeData.EnglishName = toNullString(req.EnglishName)
	}
	if req.ReleaseDate != nil {
		animeData.ReleaseDate = toNullTime(req.ReleaseDate)
	}
	if req.Type != nil {
		animeData.Type = toNullString(req.Type)
	}

	if err := h.db.UpdateAnimeById(ctx, database.UpdateAnimeByIdParams{
		RomajiName:   animeData.RomajiName,
		JapaneseName: animeData.JapaneseName,
		EnglishName:  animeData.EnglishName,
		ReleaseDate:  animeData.ReleaseDate,
		Type:         animeData.Type,
		ID:           int32(animeID),
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PostAnimeBulk(w http.ResponseWriter, r *http.Request) {
	var reqBulk []Anime
	if err := decodeJson(r.Body, &reqBulk); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	type bulkResponse struct {
		Inserted int `json:"inserted"`
		Failed   int `json:"failed"`
	}
	var res bulkResponse
	for _, req := range reqBulk {
		if strings.TrimSpace(req.RomajiName) == "" {
			res.Failed += 1
			continue
		}

		params := database.CreateAnimeParams{
			RomajiName:   req.RomajiName,
			JapaneseName: toNullString(req.JapaneseName),
			EnglishName:  toNullString(req.EnglishName),
			Type:         toNullString(req.Type),
			ReleaseDate:  toNullTime(req.ReleaseDate),
		}

		if _, err := h.db.CreateAnime(context.Background(), params); err != nil {
			res.Failed += 1
		} else {
			res.Inserted += 1
		}
	}
	if res.Inserted == 0 {
		respondWithError(w, http.StatusBadRequest, "failed to insert", fmt.Errorf("failed to insert"))
		return
	}
	respondWithJSON(w, http.StatusMultiStatus, res)
}

package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

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

func toNullTime(s *string) sql.NullTime {
	if s == nil {
		return sql.NullTime{Valid: false}
	}

	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func validateTime(t sql.NullTime) *string {
	var date *string
	if t.Valid {
		val := t.Time.Format("2006-01-02")
		date = &val
	} else {
		date = nil
	}

	return date
}

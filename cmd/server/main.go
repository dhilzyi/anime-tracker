package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dhilzyi/anime-tracker/internal/api"
	"github.com/dhilzyi/anime-tracker/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	dbUrl := os.Getenv("db_url")
	port := os.Getenv("port")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	handler := api.NewHandler(dbQueries)

	mux := http.NewServeMux()
	srv := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	mux.HandleFunc("GET /api/anime", handler.GetAnime)
	mux.HandleFunc("GET /api/anime/{animeID}", handler.GetAnimeById)
	mux.HandleFunc("POST /api/anime", handler.PostAnime)
	mux.HandleFunc("DELETE /api/anime/{animeID}", handler.DeleteAnimeById)

	fmt.Printf("Serving at http://localhost:%s/\n", port)
	log.Fatal(srv.ListenAndServe())

}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dhilzyi/anime-tracker/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
}

func main() {
	godotenv.Load(".env")

	dbUrl := os.Getenv("db_url")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	state := state{
		db: dbQueries,
	}

	input := database.InsertAnimeParams{
		RomajiName: "shingeki no kyojin",
	}
	input2 := database.InsertAnimeParams{
		RomajiName: "kuroko no basket",
	}

	if data, err := state.db.InsertAnime(context.Background(), input); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(data)
	}
	if data, err := state.db.InsertAnime(context.Background(), input2); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(data)
	}

}

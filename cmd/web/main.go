package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/dylan0804/Llamarama/cmd/internal/app"
	db "github.com/dylan0804/Llamarama/cmd/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	workDir, err := os.Getwd()
    if err != nil {
        log.Fatal("Error getting working directory:", err)
    }
    err = godotenv.Load(filepath.Join(workDir, ".env"))
    if err != nil {
        log.Fatal("Error loading .env file:", err)
    }

	ctx := context.Background()

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to parse database URL:", err)
		return
	}
	
	config.ConnConfig.DefaultQueryExecMode = 5

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return
	}
	defer conn.Close()

	queries := db.New(conn)

	app := app.New(queries)

	err = app.Run()
	
	if err != nil {
		log.Fatal("Error running application:", err)
		return
	}
}
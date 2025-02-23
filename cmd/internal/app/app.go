package app

import (
	"time"

	db "github.com/dylan0804/Llamarama/cmd/internal/db/sqlc"
	"github.com/dylan0804/Llamarama/cmd/internal/handlers"
	"github.com/dylan0804/Llamarama/cmd/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type application struct {
	router *gin.Engine
	handler *handlers.Handler
}

func New(queries *db.Queries) *application {
	router := gin.Default()

	router.Use(middleware.RequestLogger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler := handlers.NewHandler(queries)

	app := &application{
		router: router,
		handler: handler,
	}

	return app
}

func (app *application) Run() error {
	if err := app.routes(); err != nil {
		return err
	}

	return nil
}

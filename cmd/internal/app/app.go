package app

import (
	db "github.com/dylan0804/Llamarama/cmd/internal/db/sqlc"
	"github.com/dylan0804/Llamarama/cmd/internal/handlers"
	"github.com/dylan0804/Llamarama/cmd/internal/utils"
	"github.com/gin-gonic/gin"
)

type application struct {
	router *gin.Engine
	handler *handlers.Handler
	sessionStore *utils.SessionStore
}

func New(queries *db.Queries) *application {
	router := gin.Default()

	sessionStore := utils.NewSessionStore()

	handler := handlers.NewHandler(queries, sessionStore)
	
	app := &application{
		router: router,
		handler: handler,
		sessionStore: sessionStore,
	}

	return app
}

func (app *application) Run() error {
	if err := app.routes(); err != nil {
		return err
	}

	return nil
}

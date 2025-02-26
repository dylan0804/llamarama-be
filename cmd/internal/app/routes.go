package app

import (
	"time"

	"github.com/dylan0804/Llamarama/cmd/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *application) routes() error {
	app.router.Use(middleware.RequestLogger())
	app.router.Use(gin.Recovery())
	app.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// websocket
	ws := app.router.Group("/ws")
	ws.Use(middleware.AuthMiddleware(app.sessionStore))
	ws.GET("/rooms/:id", app.handler.WsHandler)

	// other routes
	v1 := app.router.Group("/api/v1")
	{
			
		auth := v1.Group("/auth")
		{
			auth.POST("/register", app.handler.Register)
			auth.POST("/login", app.handler.Login)

			auth.Use(middleware.AuthMiddleware(app.sessionStore))
			auth.POST("/logout", app.handler.Logout)
		}

		rooms := v1.Group("/rooms")
		rooms.Use(middleware.AuthMiddleware(app.sessionStore))
		{
			rooms.POST("", app.handler.CreateRoom)
			rooms.GET("", app.handler.ListRooms)
			rooms.GET("/:id", app.handler.GetRoom)
		}
	}

	return app.router.Run(":8080")
}

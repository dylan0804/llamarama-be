package app

func (app *application) routes() error {
	// websocket
	app.router.GET("/ws/rooms/:id", app.handler.WsHandler)

	// other routes
	v1 := app.router.Group("/api/v1")
	{
		rooms := v1.Group("/rooms")
		{
			rooms.POST("", app.handler.CreateRoom)
			rooms.GET("", app.handler.ListRooms)
			rooms.GET("/:id", app.handler.GetRoom)
		}
	}

	return app.router.Run(":8080")
}

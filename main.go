package main

import (
	"./handlers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	// Routing
	e.GET("/api/v1/systemd/:unit", handlers.Get())
	e.POST("/api/v1/systemd/:unit", handlers.Post())
	e.DELETE("/api/v1/systemd/:unit", handlers.Delete())
	e.PUT("/api/v1/systemd/:unit", handlers.Put())

	e.Logger.Fatal(e.Start(":1323"))
}

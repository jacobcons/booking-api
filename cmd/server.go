package main

import (
	"booking-api/routes"
	"booking-api/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status}\nError: ${error}\n\n",
	}))
	e.Use(middleware.Recover())
	// turn on debugging if in development env
	env := os.Getenv("APP_ENV")
	if env == "development" {
		e.Debug = true
	}
	// validator
	e.Validator = utils.Validator

	// routes
	routes.RegisterBookingsRoutes(e)

	e.Logger.Fatal(e.Start(":3000"))
}

package main

import (
	"booking-api/routes"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "Request: ${method} ${uri} ${status}, Error: ${error}\n\n",
	}))
	e.Use(middleware.Recover())
	// turn on debugging if in development env
	env := os.Getenv("APP_ENV")
	if env == "development" {
		e.Debug = true
	}
	// validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// routes
	routes.RegisterBookingsRoutes(e)

	e.Logger.Fatal(e.Start(":3000"))
}

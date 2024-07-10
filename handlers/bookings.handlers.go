package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetBookings(c echo.Context) error {
	return c.String(http.StatusOK, "get")
}

func CreateBooking(c echo.Context) error {
	return c.String(http.StatusOK, "create")
}

func DeleteBooking(c echo.Context) error {
	return c.String(http.StatusOK, "del")
}

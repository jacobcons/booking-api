package routes

import (
	"booking-api/handlers"
	"booking-api/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterBookingsRoutes(e *echo.Echo) {
	g := e.Group("/bookings")
	g.GET("", handlers.GetBookings)
	g.POST("", handlers.CreateBooking, middlewares.VerifyJwt)
	g.DELETE("/:id", handlers.DeleteBooking, middlewares.VerifyJwt)
}

package routes

import (
	"booking-api/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterBookingsRoutes(e *echo.Echo) {
	g := e.Group("/bookings")
	g.GET("", handlers.GetBookings)
	g.POST("", handlers.CreateBooking)
	g.DELETE("/:id", handlers.DeleteBooking)
}

package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func GetBookings(c echo.Context) error {
	return c.String(http.StatusOK, "get")
}

func isValidBookingDate(datetimeStr string) bool {
	const DATETIME_FORMAT = "2006-01-02 15:04:05"
	const TIME_FORMAT = "15:04:05"
	datetime, _ := time.Parse(DATETIME_FORMAT, datetimeStr)
	weekday := datetime.Weekday()

	// must be monday-friday
	if weekday > time.Friday {
		return false
	}

	// must be between 9-5
	extractedTime, _ := time.Parse(TIME_FORMAT, datetime.Format("15:04:05"))
	openingTime, _ := time.Parse(TIME_FORMAT, "09:00:00")
	closingTime, _ := time.Parse(TIME_FORMAT, "17:00:00")
	if extractedTime.Before(openingTime) || extractedTime.After(closingTime) {
		return false
	}

	// must be on the hour
	return true
}

func CreateBooking(c echo.Context) (err error) {
	type Body struct {
		StartDate string `validate:"required,datetime=2006-01-02 15:04:05"`
		EndDate   string `validate:"required,datetime=2006-01-02 15:04:05"`
	}

	body := new(Body)
	if err = c.Bind(body); err != nil {
		return err
	}
	if err = c.Validate(body); err != nil {
		return err
	}

	if !isValidBookingDate(body.StartDate) || !isValidBookingDate(body.EndDate) {
		return echo.NewHTTPError(http.StatusBadRequest, "startDate and endDate must be on a weekday between 9-5, on the hour and in the future")
	}

	return c.JSON(http.StatusOK, body)
}

func DeleteBooking(c echo.Context) error {
	return c.String(http.StatusOK, "del")
}

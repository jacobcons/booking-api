package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
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

	// must be between 9am-5pm
	extractedTime, _ := time.Parse(TIME_FORMAT, datetime.Format("15:04:05"))
	openingTime, _ := time.Parse(TIME_FORMAT, "09:00:00")
	closingTime, _ := time.Parse(TIME_FORMAT, "17:00:00")
	if extractedTime.Before(openingTime) || extractedTime.After(closingTime) {
		return false
	}

	// must be on the hour
	timeToNearestHour := extractedTime.Truncate(time.Hour)
	if !extractedTime.Equal(timeToNearestHour) {
		return false
	}

	// must be in the future
	ukLocation, _ := time.LoadLocation("Europe/London")
	currentDatetimeStr := time.Now().In(ukLocation).Format(DATETIME_FORMAT)
	if datetimeStr < currentDatetimeStr {
		return false
	}

	return true
}

func CreateBooking(c echo.Context) (err error) {
	type Body struct {
		StartDateTime string `validate:"required,datetime=2006-01-02 15:04:05"`
		EndDateTime   string `validate:"required,datetime=2006-01-02 15:04:05"`
	}

	// load body into struct
	body := new(Body)
	if err = c.Bind(body); err != nil {
		return err
	}

	// validation
	if err = c.Validate(body); err != nil {
		return err
	}
	if !isValidBookingDate(body.StartDateTime) || !isValidBookingDate(body.EndDateTime) {
		return echo.NewHTTPError(http.StatusBadRequest, "The datetimes must be on a weekday, between 9am-5pm, on the hour and in the future")
	}
	if body.EndDateTime <= body.StartDateTime {
		return echo.NewHTTPError(http.StatusBadRequest, "endDateTime must be greater than startDateTime")
	}
	startDate := strings.Split(body.StartDateTime, " ")[0]
	endDate := strings.Split(body.EndDateTime, " ")[0]
	if startDate != endDate {
		return echo.NewHTTPError(http.StatusBadRequest, "The dates must be the same")
	}

	return c.JSON(http.StatusOK, body)
}

func DeleteBooking(c echo.Context) error {
	return c.String(http.StatusOK, "del")
}

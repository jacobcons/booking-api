package handlers

import (
	. "booking-api/db"
	. "booking-api/types"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

type Booking struct {
	Id            string `json:"id"`
	UserId        string `json:"userId"`
	StartDatetime string `json:"startDateTime"`
	EndDatetime   string `json:"endDateTime"`
}

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
		StartDatetime string `validate:"required,datetime=2006-01-02 15:04:05"`
		EndDatetime   string `validate:"required,datetime=2006-01-02 15:04:05"`
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
	startDatetime, endDatetime := body.StartDatetime, body.EndDatetime
	if !isValidBookingDate(startDatetime) || !isValidBookingDate(endDatetime) {
		return echo.NewHTTPError(http.StatusBadRequest, "The datetimes must be on a weekday, between 9am-5pm, on the hour and in the future")
	}
	if endDatetime <= startDatetime {
		return echo.NewHTTPError(http.StatusBadRequest, "endDateTime must be greater than startDateTime")
	}
	startDate := strings.Split(startDatetime, " ")[0]
	endDate := strings.Split(endDatetime, " ")[0]
	if startDate != endDate {
		return echo.NewHTTPError(http.StatusBadRequest, "The dates must be the same")
	}

	// extract userId from jwt
	userId := c.Get("user").(*UserJwtClaims).Id
	// run queries in transaction to prevent race conditions
	tx := DB.MustBegin()
	// ensure booking doesn't overlap with existing bookings
	var isOverlap bool
	_ = DB.Get(
		&isOverlap,
		`
			SELECT EXISTS (
				SELECT 1
				FROM booking
				WHERE (start_datetime, end_datetime) OVERLAPS ($1, $2) 
			)
		`,
		startDatetime, endDatetime,
	)
	if isOverlap {
		return echo.NewHTTPError(http.StatusConflict, "booking must not overlap with existing bookings")
	}
	// insert new booking
	newBooking := Booking{}
	_ = tx.QueryRowx(
		`INSERT INTO booking(user_id, start_datetime, end_datetime) VALUES ($1, $2, $3) RETURNING *`,
		userId, startDatetime, endDatetime,
	).StructScan(&newBooking)
	_ = tx.Commit()

	return c.JSON(http.StatusOK, newBooking)
}

func DeleteBooking(c echo.Context) error {
	return c.String(http.StatusOK, "del")
}

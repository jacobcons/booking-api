package handlers

import (
	. "booking-api/db"
	. "booking-api/types"
	"fmt"
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

func GetBookings(c echo.Context) (err error) {
	type QueryParams struct {
		StartDatetime string `query:"startDatetime" validate:"required,datetime=2006-01-02T15:04:05Z"`
		EndDatetime   string `query:"endDatetime" validate:"required,datetime=2006-01-02T15:04:05Z"`
	}

	// load query params into struct
	queryParams := new(QueryParams)
	if err = c.Bind(queryParams); err != nil {
		return err
	}

	// validation
	if err = c.Validate(queryParams); err != nil {
		return err
	}

	// get all bookings between the start and end datetimes
	startDatetime, endDatetime := queryParams.StartDatetime, queryParams.EndDatetime

	bookings := []Booking{}
	DB.Select(
		&bookings,
		`
			SELECT *
			FROM booking b
			WHERE start_datetime >= $1 AND end_datetime <= $2
			ORDER BY start_datetime
		`,
		startDatetime, endDatetime,
	)

	return c.JSON(http.StatusOK, bookings)
}

func isValidBookingDate(datetimeStr string) bool {
	const DATETIME_FORMAT = "2006-01-02T15:04:05Z"
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
		StartDatetime string `validate:"required,datetime=2006-01-02T15:04:05Z"`
		EndDatetime   string `validate:"required,datetime=2006-01-02T15:04:05Z"`
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

	startDate := strings.Split(startDatetime, "T")[0]
	endDate := strings.Split(endDatetime, "T")[0]
	if startDate != endDate {
		return echo.NewHTTPError(http.StatusBadRequest, "The dates must be the same")
	}

	// create booking as long as it doesn't overlap with existing bookings
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
	// insert booking
	newBooking := Booking{}
	_ = tx.QueryRowx(
		`INSERT INTO booking(user_id, start_datetime, end_datetime) VALUES ($1, $2, $3) RETURNING *`,
		userId, startDatetime, endDatetime,
	).StructScan(&newBooking)
	_ = tx.Commit()

	return c.JSON(http.StatusCreated, newBooking)
}

func DeleteBooking(c echo.Context) (err error) {
	type Params struct {
		Id string `param:"id" validate:"required,uuid"`
	}

	// load params into struct
	params := new(Params)
	if err = c.Bind(params); err != nil {
		return err
	}

	// validation
	if err = c.Validate(params); err != nil {
		return err
	}

	// delete booking with id route param belonging to the logged-in user
	id := params.Id
	userId := c.Get("user").(*UserJwtClaims).Id

	rowsDeleted, _ := DB.MustExec(`DELETE FROM booking WHERE id = $1 AND user_id = $2`, id, userId).RowsAffected()
	if rowsDeleted == 0 {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No booking<%v> belonging to user<%v> found", id, userId))
	}

	return c.NoContent(http.StatusNoContent)
}

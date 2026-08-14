package validation

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

func validateNotBeforeToday(fl validator.FieldLevel) bool {
	inputDate, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return !inputDate.Before(today)
}

func New() *CustomValidator {
	v := validator.New()
	_ = v.RegisterValidation("not_before_today", validateNotBeforeToday)

	return &CustomValidator{validator: v}
}

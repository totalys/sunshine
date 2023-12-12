package controler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// TemperatureReqParams parameters for the temperature weather api
type TemperatureReqParams struct {
	City    string `json:"city" validate:"required"`
	Country string `json:"country"`
}

// CustomValidator hooks any validator framework that implements Validate interface
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator creates a new validator for the api controler
func NewCustomValidator(v *validator.Validate) *CustomValidator {
	return &CustomValidator{
		validator: v,
	}
}

// Validate validates the given parameters
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

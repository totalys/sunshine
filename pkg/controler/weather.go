package controler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/swaggo/swag/example/celler/httputil"
	"github.com/totalys/sunshine/pkg/weather-service"
	"go.uber.org/zap"
)

const (
	errInternalMsg = "internal server error"

	pCity    = "city"
	pCountry = "country"
)

// WeatherControler weather check controler
type WeatherControler interface {
	GetTemperatureForCity(c echo.Context) error
}

type weatherControler struct {
	l              *zap.Logger
	weatherService weather.WeatherService
}

var _ WeatherControler = (*weatherControler)(nil)

// NewWeatherControler creates a controler for weathercheck
func NewWeatherControler(l *zap.Logger, service weather.WeatherService) *weatherControler {
	return &weatherControler{
		l:              l,
		weatherService: service,
	}
}

// GetTemperatureForCity returns wheater .
// @Summary Gets the temperature in Kelvin, Celsius and Fahrenheit for a given city
// @Description Gets the application status. sunshine means it is working.
// @Tags temperature
// @Accept json
// @Produce json
// @Param city query string true "Sao Paulo"
// @Param country query string false "Brazil"
// @Success 200 {string} message "sunshine"
// @Failure 500 {object}  httputil.HTTPError
// @Router /api/temperature [get]
func (ctl *weatherControler) GetTemperatureForCity(c echo.Context) error {
	ctl.l.Debug("getTemperatureForCity called")

	params := TemperatureReqParams{
		City:    c.QueryParam(pCity),
		Country: c.QueryParam(pCountry),
	}

	if err := c.Echo().Validator.Validate(params); err != nil {
		ctl.l.Error("validation failed", zap.Error(err))
		return err
	}

	info, err := ctl.weatherService.GetTemperatureForCity(params.City, params.Country)
	if err != nil {
		ctl.l.Debug("getTemperatureForCity returned error",
			zap.Error(err),
			zap.String("service", "weatherService"))
		return echo.NewHTTPError(http.StatusInternalServerError, userError(err))
	}

	return c.JSON(http.StatusOK, info)
}

// userError hides internal errors that could expose internal
// structures and returns an user friendly message
func userError(err error) string {
	log.Println("error calling api", err)
	return errInternalMsg
}

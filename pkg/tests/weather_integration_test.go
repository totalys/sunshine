package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/totalys/sunshine/configs"
	"github.com/totalys/sunshine/pkg/controler"
	"github.com/totalys/sunshine/pkg/logger"
	"github.com/totalys/sunshine/pkg/mocks/geofinder"
	openweathermock "github.com/totalys/sunshine/pkg/mocks/open-weather"
	openweather "github.com/totalys/sunshine/pkg/open-weather"
	"github.com/totalys/sunshine/pkg/weather-service"
)

func TestGetTemperatureForCity(t *testing.T) {

	// Arrange
	city := "Cajobi"
	country := "Brazil"
	expectLat := float64(-20.8806453)
	expectLng := float64(-48.8103486)
	expectedOpWeatherInfo := &openweather.WeatherInfo{
		Lat:            expectLat,
		Lon:            expectLng,
		Timezone:       "",
		TimezoneOffset: 0,
		Current: openweather.Current{
			Temp: 297.8,
		},
	}
	expectedInfo := weather.Info{
		Lat:     -20.8806453,
		Lng:     -48.8103486,
		City:    "Cajobi",
		Country: "Brazil",
		Temperature: weather.Temperature{
			Kelvin:    "297.80",
			Celsius:   "24.65",
			Farenheit: "76.37",
		},
	}

	e := echo.New()
	e.Validator = controler.NewCustomValidator(validator.New())

	req := httptest.NewRequest(http.MethodGet, configs.APIHealth, nil)
	q := req.URL.Query()
	q.Add("city", city)
	q.Add("country", country)
	req.URL.RawQuery = q.Encode()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	l, err := logger.New("info")
	if err != nil {
		t.Fatal("error returned is not HTTPError")
	}

	geomock := new(geofinder.GeoFinder)
	geomock.On("FindLatLngByCity", mock.Anything, mock.Anything).
		Return(expectLat, expectLng, nil)

	wthmock := new(openweathermock.WeatherOfficer)
	wthmock.On("GetWeather", mock.Anything, mock.Anything).
		Return(expectedOpWeatherInfo, nil)

	service, err := weather.NewWeatherService(geomock, wthmock)
	if err != nil {
		t.Fatal(err)
	}

	c := controler.NewWeatherControler(l, service)

	// Act
	err = c.GetTemperatureForCity(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var result weather.Info

	err = json.Unmarshal(rec.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, expectedInfo, result, "weather should be as expected")
}

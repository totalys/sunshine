package controler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/totalys/sunshine/configs"
	"github.com/totalys/sunshine/pkg/controler"
	"github.com/totalys/sunshine/pkg/logger"
	weathermock "github.com/totalys/sunshine/pkg/mocks/weather-service"
	"github.com/totalys/sunshine/pkg/weather-service"
)

const errInternalMsg = "internal server error"

func TestNewWeatherControler(t *testing.T) {

	// Arrange & Act
	c := controler.NewWeatherControler(nil, nil)

	// assert
	assert.NotNil(t, c)
}

func TestGetTemperatureForCity(t *testing.T) {

	// Arrange
	city := "Cajobi"
	country := "Brazil"
	e := echo.New()
	e.Validator = controler.NewCustomValidator(validator.New())
	req := httptest.NewRequest(http.MethodGet, configs.APIHealth, nil)
	q := req.URL.Query()
	q.Add("city", city)
	q.Add("country", country)
	req.URL.RawQuery = q.Encode()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
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

	serviceMock := new(weathermock.WeatherService)
	serviceMock.On("GetTemperatureForCity",
		testifymock.Anything, testifymock.Anything).
		Return(&expectedInfo, nil)

	l, err := logger.New("info")
	if err != nil {
		t.Fatal("error returned is not HTTPError")
	}

	c := controler.NewWeatherControler(l, serviceMock)

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

func TestGetTemperatureForCityWhenInvalidRequestShouldHandleError(t *testing.T) {

	// Arrange
	e := echo.New()
	e.Validator = controler.NewCustomValidator(validator.New())
	req := httptest.NewRequest(http.MethodGet, configs.APITemperature, nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	serviceErr := errors.New("Key: 'TemperatureReqParams.City' Error:Field validation for 'City' failed on the 'required' tag")
	serviceMock := new(weathermock.WeatherService)

	l, err := logger.New("info")
	if err != nil {
		t.Fatal("error returned is not HTTPError")
	}

	c := controler.NewWeatherControler(l, serviceMock)

	// Act
	err = c.GetTemperatureForCity(ctx)

	assert.IsType(t, &echo.HTTPError{}, err)

	httperr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatal("error returned is not HTTPError")
	}

	// Assert
	assert.Equal(t, http.StatusBadRequest, httperr.Code)
	assert.Equal(t, serviceErr.Error(), httperr.Message)
}

func TestGetTemperatureForCityWhenServiceErrorShouldHandleError(t *testing.T) {

	// Arrange
	city := "Cajobi"
	country := "Brazil"
	e := echo.New()
	e.Validator = controler.NewCustomValidator(validator.New())
	req := httptest.NewRequest(http.MethodGet, configs.APITemperature, nil)
	q := req.URL.Query()
	q.Add("city", city)
	q.Add("country", country)
	req.URL.RawQuery = q.Encode()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	serviceErr := errors.New("some internal error")

	serviceMock := new(weathermock.WeatherService)
	serviceMock.On("GetTemperatureForCity",
		testifymock.Anything, testifymock.Anything).
		Return(nil, serviceErr)

	l, err := logger.New("info")
	if err != nil {
		t.Fatal("error returned is not HTTPError")
	}

	c := controler.NewWeatherControler(l, serviceMock)

	// Act
	err = c.GetTemperatureForCity(ctx)

	assert.IsType(t, &echo.HTTPError{}, err)

	httperr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatal("error returned is not HTTPError")
	}

	// Assert
	assert.Equal(t, http.StatusInternalServerError, httperr.Code)
	assert.Equal(t, errInternalMsg, httperr.Message)
	assert.NotEqual(t, serviceErr, err, "server error should be ommited")
}

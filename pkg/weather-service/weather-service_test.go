package weather_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/totalys/sunshine/pkg/mocks/geofinder"
	openweathermock "github.com/totalys/sunshine/pkg/mocks/open-weather"
	openweather "github.com/totalys/sunshine/pkg/open-weather"
	"github.com/totalys/sunshine/pkg/weather-service"
)

func TestNewWeatherService(t *testing.T) {

	// Arrange
	geofindermock := new(geofinder.GeoFinder)
	openweathermock := new(openweathermock.WeatherOfficer)

	// Act
	w, err := weather.NewWeatherService(geofindermock, openweathermock)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, w)
}

func TestGetTemperatureForCityShouldReturnTemperature(t *testing.T) {

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

	expectInfo := &weather.Info{
		Lat:     expectLat,
		Lng:     expectLng,
		City:    city,
		Country: country,
		Temperature: weather.Temperature{
			Kelvin:    "297.80",
			Celsius:   "24.65",
			Farenheit: "76.37",
		},
	}

	geofindermock := new(geofinder.GeoFinder)
	geofindermock.On("FindLatLngByCity", mock.Anything, mock.Anything).
		Return(expectLat, expectLng, nil)

	openweathermock := new(openweathermock.WeatherOfficer)
	openweathermock.On("GetWeather", mock.Anything, mock.Anything).
		Return(expectedOpWeatherInfo, nil)

	w, err := weather.NewWeatherService(geofindermock, openweathermock)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	info, err := w.GetTemperatureForCity(city, country)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, expectInfo, info)
}

func TestGetTemperatureForCityShouldReturnErrorWhenOpenWeatherFails(t *testing.T) {

	// Arrange
	city := "Cajobi"
	country := "Brazil"
	expectLat := float64(-20.8806453)
	expectLng := float64(-48.8103486)

	geofindermock := new(geofinder.GeoFinder)
	geofindermock.On("FindLatLngByCity", mock.Anything, mock.Anything).
		Return(expectLat, expectLng, nil)

	openweathermock := new(openweathermock.WeatherOfficer)
	openweathermock.On("GetWeather", mock.Anything, mock.Anything).
		Return(nil, errors.New("some open-weather error"))

	w, err := weather.NewWeatherService(geofindermock, openweathermock)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	info, err := w.GetTemperatureForCity(city, country)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, info)
}

func TestGetTemperatureForCityShouldReturnErrorWhenGeofinderFails(t *testing.T) {

	// Arrange
	city := "Cajobi"
	country := "Brazil"

	geofindermock := new(geofinder.GeoFinder)
	geofindermock.On("FindLatLngByCity", mock.Anything, mock.Anything).
		Return(float64(0), float64(0), errors.New("some geofinder error"))

	openweathermock := new(openweathermock.WeatherOfficer)

	w, err := weather.NewWeatherService(geofindermock, openweathermock)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	info, err := w.GetTemperatureForCity(city, country)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, info)
}

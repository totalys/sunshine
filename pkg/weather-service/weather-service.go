package weather

import (
	"fmt"
	"math"

	"github.com/totalys/sunshine/pkg/geolocation"
	openweather "github.com/totalys/sunshine/pkg/open-weather"
)

var _ WeatherService = (*service)(nil)

// Info weather data of a given location
type Info struct {
	Lat         float64     `json:"latitude"`
	Lng         float64     `json:"longitude"`
	City        string      `json:"city"`
	Country     string      `json:"country"`
	Temperature Temperature `json:"temperature"`
}

type Temperature struct {
	Kelvin    string `json:"kelvin"`
	Celsius   string `json:"celsius"`
	Farenheit string `json:"farenheit"`
}

type service struct {
	geofinder      geolocation.GeoFinder
	weatherofficer openweather.WeatherOfficer
}

// WeatherService application contract for the weather services offered
type WeatherService interface {

	// GetTemperatureForCity gets the weather info for a given city from a given country
	GetTemperatureForCity(city, country string) (*Info, error)
}

// NewWeatherService creates a weather service
func NewWeatherService(g geolocation.GeoFinder, w openweather.WeatherOfficer) (*service, error) {
	return &service{
		geofinder:      g,
		weatherofficer: w,
	}, nil
}

func (s *service) GetTemperatureForCity(city, country string) (*Info, error) {

	lat, lng, err := s.geofinder.FindLatLngByCity(city, country)
	if err != nil {
		return nil, fmt.Errorf("failed to get latitude and longitude from geo service: %w", err)
	}

	weatherInfo, err := s.weatherofficer.GetWeather(lat, lng)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather info from open weather service: %w", err)
	}

	k := weatherInfo.Current.Temp
	c := kelvinToCelsius(k)
	f := celsiusToFarenheit(c)

	return &Info{
		Lat:     lat,
		Lng:     lng,
		City:    city,
		Country: country,
		Temperature: Temperature{
			Kelvin:    fmt.Sprintf("%.2f", k),
			Celsius:   fmt.Sprintf("%.2f", c),
			Farenheit: fmt.Sprintf("%.2f", f),
		},
	}, nil
}

func kelvinToCelsius(kelvin float64) float64 {

	celsius := kelvin - float64(273.15)

	return math.Round(celsius*100) / 100
}

func celsiusToFarenheit(c float64) float64 {

	faren := (c * 1.8) + 32

	return math.Round(faren*100) / 100
}

package openweather

import (
	"fmt"
	"os"

	"gopkg.in/h2non/gentleman.v2"
)

const (
	openWeatherApiKey  = "SUNNY_EXTERNAL_WEATHER_APIKEY"
	openWeatherApiPath = "/data/3.0/onecall"

	pAppIDKey     = "appid"
	pLatitudeKey  = "lat"
	pLongitudeKey = "lon"
	pExcludeKey   = "exclude"
	pexcludeVals  = "minutely,hourly,daily,alert"
)

var _ WeatherOfficer = (*OpenWeather)(nil)

// WeatherInfo holds structured data from the open weather api
type WeatherInfo struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Current        Current `json:"current"`
}

type Current struct {
	Dt         float64 `json:"dt"`
	Sunrise    float64 `json:"sunrise"`
	Sunset     float64 `json:"sunset"`
	Temp       float64 `json:"temp"`
	FeelsLike  float64 `json:"feels_like"`
	Pressure   float64 `json:"pressure"`
	Humidity   float64 `json:"humidity"`
	DewPoint   float64 `json:"dew_point"`
	Uvi        float64 `json:"uvi"`
	Clouds     float64 `json:"clouds"`
	Visibility float64 `json:"visibility"`
	WindSpeed  float64 `json:"wind_speed"`
	WindDeg    float64 `json:"wind_deg"`
	Weather    []struct {
		ID          float64 `json:"id"`
		Main        string  `json:"main"`
		Description string  `json:"description"`
		Icon        string  `json:"icon"`
	}
}

// WeatherOfficer deliver weather information
type WeatherOfficer interface {

	// GetWeather retrives weather info from a given latitude and longitude
	GetWeather(lat, lng float64) (*WeatherInfo, error)
}

// OpenWeather provides a WeatherOfficer implementation that uses http protocol
// to gather weather information from openweathermap.org
type OpenWeather struct {
	client *gentleman.Client
	apikey string
}

// NewOpenWeather creates a new OpenWeather
func NewOpenWeather(client *gentleman.Client) (*OpenWeather, error) {

	return &OpenWeather{
		client: client,
		apikey: os.Getenv(openWeatherApiKey),
	}, nil
}

// GetWeather impl
func (o OpenWeather) GetWeather(lat, lng float64) (*WeatherInfo, error) {

	req := o.client.Request().
		Path(openWeatherApiPath).
		SetQuery(pLatitudeKey, fmt.Sprintf("%f", lat)).
		SetQuery(pLongitudeKey, fmt.Sprintf("%f", lng)).
		SetQuery(pExcludeKey, pexcludeVals).
		SetQuery(pAppIDKey, o.apikey)

	res, err := req.Send()
	if err != nil {
		return nil, fmt.Errorf("error retrieving weather from open wather api: $%w ", err)
	}

	if !res.Ok {
		return nil, fmt.Errorf("invalid response from open wather api: %s ", res.String())
	}

	var weatherInfo = &WeatherInfo{}

	err = res.JSON(weatherInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling open weather api response: $%w ", err)
	}

	return weatherInfo, nil
}

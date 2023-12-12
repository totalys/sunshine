package geolocation

import (
	"fmt"
	"os"

	"gopkg.in/h2non/gentleman.v2"
)

const (
	openWeatherApiKey  = "SUNNY_EXTERNAL_WEATHER_APIKEY"
	openWeatherApiPath = "/geo/1.0/direct"

	pAppIDKey = "appid"
	pCity     = "q"
	pLimit    = "limit"

	limit = "1"
)

var _ GeoFinder = (*OpenWeatherFinder)(nil)

// Geodata contains a city's latitude and longitude coordinates
type GeoData []struct {
	Name       string `json:"name"`
	LocalNames struct {
		FeatureName string `json:"feature_name"`
	} `json:"local_names,omitempty"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

// OpenWeatherFinder finds geo coordinates
type OpenWeatherFinder struct {
	client *gentleman.Client
	apikey string
}

// NewOpenWeatherFinder creates a OpenWeatherFinder
func NewOpenWeatherFinder(client *gentleman.Client) (*OpenWeatherFinder, error) {
	return &OpenWeatherFinder{
		client: client,
		apikey: os.Getenv(openWeatherApiKey),
	}, nil
}

// FindLatLngByCity invokes the open weather API to get latitude and longitude values for a given city.
func (o OpenWeatherFinder) FindLatLngByCity(city, country string) (float64, float64, error) {

	req := o.client.Request().
		Path(openWeatherApiPath).
		SetQuery(pCity, city).
		SetQuery(pLimit, limit).
		SetQuery(pAppIDKey, o.apikey)

	res, err := req.Send()
	if err != nil {
		return float64(0), float64(0), fmt.Errorf("error retrieving geo location api: $%w ", err)
	}

	if !res.Ok {
		return float64(0), float64(0), fmt.Errorf("invalid response geo location api: %s ", res.String())
	}

	var geoData = &GeoData{}

	err = res.JSON(geoData)
	if err != nil {
		return float64(0), float64(0), fmt.Errorf("error unmarshaling open weather api response: $%w ", err)
	}

	if len(*geoData) == 0 {
		return float64(0), float64(0), fmt.Errorf("no coordinates found for city [%s]: %s ", city, res.String())
	}

	return (*geoData)[0].Lat, (*geoData)[0].Lon, nil
}

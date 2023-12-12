package geolocation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
	"googlemaps.github.io/maps"
)

// GeoFinder provides geografical data about locations on earth
type GeoFinder interface {

	// FindLatLngByCity finds latitude and longitude values for a given city of a given country
	FindLatLngByCity(city, country string) (float64, float64, error)
}

// InternalCitiesFinder provides a GeoFinder implementation that uses an internal source to find geo data.
type InternalCitiesFinder struct {
	cities map[string]maps.LatLng
}

// NewInternalCitiesFinder creates a InternalCitiesFinder
func NewInternalCitiesFinder(citiesPath string, fs afero.Fs) (GeoFinder, error) {

	b, err := afero.ReadFile(fs, citiesPath)
	if err != nil {
		return nil, fmt.Errorf("error reading cities file: [%s] error: %w", citiesPath, err)
	}

	var cities map[string]maps.LatLng

	err = json.Unmarshal(b, &cities)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling cities: %w", err)
	}

	return &InternalCitiesFinder{
		cities: cities,
	}, nil
}

// FindLatLngByCity impl
func (i *InternalCitiesFinder) FindLatLngByCity(city, country string) (float64, float64, error) {

	cityWithCountry := fmt.Sprintf("%s, %s", city, country)

	latlng, ok := i.cities[cityWithCountry]
	if !ok {
		return 0, 0, fmt.Errorf("could not find geolocation for the city [%s] and country [%s] provided",
			city, country)
	}

	return latlng.Lat, latlng.Lng, nil
}

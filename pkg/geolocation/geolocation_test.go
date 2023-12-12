package geolocation_test

import (
	"math"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/totalys/sunshine/pkg/geolocation"
)

func TestNewCitiesFinderReturnErrorWhenNoFileIsFound(t *testing.T) {

	// Arrange
	scenarios := []struct {
		name     string
		citypath string
		fs       afero.Fs
	}{
		{
			name:     "the file provided is not found or the file system fails to retrieve it",
			citypath: ".",
			fs:       afero.NewOsFs(),
		},
		{
			name:     "there is some error in the json provided that causes an unmarshaling error",
			citypath: ".",
			fs:       afero.NewMemMapFs(),
		},
	}

	for _, s := range scenarios {

		// Act
		got, err := geolocation.NewInternalCitiesFinder(s.citypath, s.fs)

		// Assert
		assert.Nil(t, got, s.name)
		assert.Error(t, err, s.name)
	}

}

func TestFindByCityShouldReturnExpectedLatLng(t *testing.T) {

	// Arrange
	cfg_file_name := "local_city_source.json"
	fs := afero.NewMemMapFs()

	cfg_file := `{"Tokyo, Japan":  {"Lat":356897,"Lng":1396922},
		"Jakarta, Indonesia":  {"Lat":-61750,"Lng":1068275},
		"Delhi, India":  {"Lat":286100,"Lng":772300},
		"Guangzhou, China":  {"Lat":231300,"Lng":1132600},
		"Sao Paulo, Brazil":  {"Lat":-235500,"Lng":-466333}}`

	f, err := fs.Create(cfg_file_name)
	if err != nil {
		t.Fatal(err)
	}

	f.WriteString(cfg_file)
	f.Close()

	geofinder, err := geolocation.NewInternalCitiesFinder(cfg_file_name, fs)
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		city    string
		country string
	}{
		{
			city:    "Sao Paulo",
			country: "Brazil",
		},
		{
			city:    "Delhi",
			country: "India",
		},
	}

	for _, s := range scenarios {

		// Act
		lat, lng, err := geofinder.FindLatLngByCity(s.city, s.country)

		// Assert
		assert.Nil(t, err)
		assert.Greater(t, math.Abs(lat), float64(0))
		assert.Greater(t, math.Abs(lng), float64(0))
	}
}

func TestFindByCityShouldReturnErrorWhenCantFindCity(t *testing.T) {

	// Arrange
	cfg_file_name := "local_city_source.json"
	fs := afero.NewMemMapFs()

	cfg_file := `{"Tokyo, Japan":  {"Lat":356897,"Lng":1396922},
		"Jakarta, Indonesia":  {"Lat":-61750,"Lng":1068275},
		"Delhi, India":  {"Lat":286100,"Lng":772300},
		"Guangzhou, China":  {"Lat":231300,"Lng":1132600},
		"Sao Paulo, Brazil":  {"Lat":-235500,"Lng":-466333}}`

	f, err := fs.Create(cfg_file_name)
	if err != nil {
		t.Fatal(err)
	}

	f.WriteString(cfg_file)
	f.Close()

	geofinder, err := geolocation.NewInternalCitiesFinder(cfg_file_name, fs)
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		city    string
		country string
	}{
		{
			city:    "Nowhere",
			country: "NoCountry",
		},
	}

	for _, s := range scenarios {

		// Act
		lat, lng, err := geofinder.FindLatLngByCity(s.city, s.country)

		// Assert
		assert.NotNil(t, err)
		assert.Equal(t, float64(0), lat)
		assert.Equal(t, float64(0), lng)
	}
}

package openweather_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mock "gopkg.in/h2non/gentleman-mock.v2"
	"gopkg.in/h2non/gentleman.v2"

	openweather "github.com/totalys/sunshine/pkg/open-weather"
)

func TestNewOpenWeather(t *testing.T) {

	// Arrange
	const host = "http://localhost0"
	defer mock.Disable()
	mock.New(host).
		Get("/*").
		Reply(200).BodyString("")

	gentclient := gentleman.New()
	gentclient.Use(mock.Plugin)

	// Act
	w, err := openweather.NewOpenWeather(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.NotNil(t, w)
}

func TestGetWeatherShouldReturnWeatherInfo(t *testing.T) {

	// Arrange
	const host = "http://localhost1"
	defer mock.Disable()
	mock.New(host).
		Get("/*").
		Reply(200).BodyString(`{
			"lat": -23.5374,
			"lon": -46.6613,
			"timezone": "America/Sao_Paulo",
			"timezone_offset": -10800,
			"current": {
				"dt": 1696052522,
				"sunrise": 1696063697,
				"sunset": 1696107903,
				"temp": 291.22,
				"feels_like": 291.38,
				"pressure": 1020,
				"humidity": 88,
				"dew_point": 289.2,
				"uvi": 0,
				"clouds": 100,
				"visibility": 10000,
				"wind_speed": 2.06,
				"wind_deg": 90,
				"weather": [
					{
						"id": 804,
						"main": "Clouds",
						"description": "overcast clouds",
						"icon": "04n"
					}
				]
			}
		}`)

	gentclient := gentleman.New()
	gentclient.URL(host)
	gentclient.Use(mock.Plugin)

	w, err := openweather.NewOpenWeather(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	info, err := w.GetWeather(-23.5374, -46.6613)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.NotNil(t, w)
	assert.Equal(t, 291.22, info.Current.Temp)
}

func TestGetWeatherShouldReturnErrorWhenConnectionFails(t *testing.T) {

	// Arrange
	const host = "http://localhost2"
	defer mock.Disable()
	mock.New(host).
		Get("/*").EnableNetworking()

	gentclient := gentleman.New()
	gentclient.URL(host)
	gentclient.Use(mock.Plugin)

	w, err := openweather.NewOpenWeather(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	info, err := w.GetWeather(-23.5374, -46.6613)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, info)
}

func TestGetWeatherShouldReturnErrorWhenExternalServiceNotOk(t *testing.T) {

	// Arrange
	const host = "http://localhost3"
	defer mock.Disable()
	mock.New(host).
		Get("/*").
		Reply(401)

	gentclient := gentleman.New()
	gentclient.URL(host)
	gentclient.Use(mock.Plugin)

	w, err := openweather.NewOpenWeather(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	info, err := w.GetWeather(-23.5374, -46.6613)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, info)
}

func TestGetWeatherShouldReturnErrorWhenContractChange(t *testing.T) {

	// Arrange
	const host = "http://localhost4"
	defer mock.Disable()
	mock.New(host).
		Get("/*").
		Reply(200).BodyString(`<?xml version="1.0" encoding="UTF-8"?>
		<root>
		   <current>
			  <clouds>100</clouds>
			  <dew_point>289.2</dew_point>
			  <dt>1696052522</dt>
			  <feels_like>291.38</feels_like>
			  <humidity>88</humidity>
			  <pressure>1020</pressure>
			  <sunrise>1696063697</sunrise>
			  <sunset>1696107903</sunset>
			  <temp>291.22</temp>
			  <uvi>0</uvi>
			  <visibility>10000</visibility>
			  <weather>
				 <element>
					<description>overcast clouds</description>
					<icon>04n</icon>
					<id>804</id>
					<main>Clouds</main>
				 </element>
			  </weather>
			  <wind_deg>90</wind_deg>
			  <wind_speed>2.06</wind_speed>
		   </current>
		   <lat>-23.5374</lat>
		   <lon>-46.6613</lon>
		   <timezone>America/Sao_Paulo</timezone>
		   <timezone_offset>-10800</timezone_offset>
		</root>`)

	gentclient := gentleman.New()
	gentclient.URL(host)
	gentclient.Use(mock.Plugin)

	w, err := openweather.NewOpenWeather(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	info, err := w.GetWeather(-23.5374, -46.6613)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, info)
}

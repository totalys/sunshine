package geolocation_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "gopkg.in/h2non/gentleman-mock.v2"
	"gopkg.in/h2non/gentleman.v2"

	"github.com/totalys/sunshine/pkg/geolocation"
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
	w, err := geolocation.NewOpenWeatherFinder(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.NotNil(t, w)
}

func TestFindLatLngByCityShouldReturnLatLng(t *testing.T) {

	// Arrange
	const host = "http://localhost1"
	const city = "London"
	defer mock.Disable()
	mock.New(host).
		Get("/*").
		Reply(200).BodyString(`[
			{
				"name": "London",
				"local_names": {
					"ff": "London",
					"ms": "London",
					"av": "Лондон",
					"mk": "Лондон",
					"fo": "London",
					"el": "Λονδίνο",
					"te": "లండన్",
					"an": "Londres",
					"th": "ลอนดอน",
					"fj": "Lodoni",
					"kn": "ಲಂಡನ್",
					"nn": "London",
					"ka": "ლონდონი",
					"sc": "Londra",
					"ln": "Lóndɛlɛ",
					"kl": "London",
					"mt": "Londra",
					"si": "ලන්ඩන්",
					"nl": "Londen",
					"ce": "Лондон",
					"de": "London",
					"mn": "Лондон",
					"tg": "Лондон",
					"hr": "London",
					"ps": "لندن",
					"bh": "लंदन",
					"gd": "Lunnainn",
					"my": "လန်ဒန်မြို့",
					"lv": "Londona",
					"su": "London",
					"qu": "London",
					"io": "London",
					"is": "London",
					"pl": "Londyn",
					"fa": "لندن",
					"ia": "London",
					"fi": "Lontoo",
					"be": "Лондан",
					"es": "Londres",
					"bg": "Лондон",
					"se": "London",
					"hu": "London",
					"or": "ଲଣ୍ଡନ",
					"fr": "Londres",
					"eu": "Londres",
					"tk": "London",
					"lb": "London",
					"sr": "Лондон",
					"cs": "Londýn",
					"sv": "London",
					"tw": "London",
					"br": "Londrez",
					"da": "London",
					"tt": "Лондон",
					"gl": "Londres",
					"vi": "Luân Đôn",
					"gv": "Lunnin",
					"az": "London",
					"fy": "Londen",
					"kk": "Лондон",
					"kw": "Loundres",
					"ab": "Лондон",
					"rm": "Londra",
					"gn": "Lóndyre",
					"ru": "Лондон",
					"ta": "இலண்டன்",
					"sn": "London",
					"sd": "لنڊن",
					"et": "London",
					"mi": "Rānana",
					"en": "London",
					"mg": "Lôndôna",
					"pa": "ਲੰਡਨ",
					"vo": "London",
					"km": "ឡុងដ៍",
					"cv": "Лондон",
					"co": "Londra",
					"nv": "Tooh Dineʼé Bikin Haalʼá",
					"wo": "Londar",
					"tl": "Londres",
					"om": "Landan",
					"zh": "伦敦",
					"he": "לונדון",
					"ug": "لوندۇن",
					"ja": "ロンドン",
					"ku": "London",
					"id": "London",
					"it": "Londra",
					"uk": "Лондон",
					"sk": "Londýn",
					"li": "Londe",
					"ca": "Londres",
					"cy": "Llundain",
					"ar": "لندن",
					"jv": "London",
					"zu": "ILondon",
					"pt": "Londres",
					"am": "ለንደን",
					"gu": "લંડન",
					"mr": "लंडन",
					"ga": "Londain",
					"ig": "London",
					"so": "London",
					"ky": "Лондон",
					"hy": "Լոնդոն",
					"eo": "Londono",
					"sm": "Lonetona",
					"ny": "London",
					"no": "London",
					"ml": "ലണ്ടൻ",
					"sa": "लन्डन्",
					"ie": "London",
					"ro": "Londra",
					"ur": "علاقہ لندن",
					"ba": "Лондон",
					"lt": "Londonas",
					"bs": "London",
					"to": "Lonitoni",
					"tr": "Londra",
					"cu": "Лондонъ",
					"af": "Londen",
					"feature_name": "London",
					"bo": "ལོན་ཊོན།",
					"ascii": "London",
					"ko": "런던",
					"sh": "London",
					"hi": "लंदन",
					"st": "London",
					"os": "Лондон",
					"yo": "Lọndọnu",
					"sl": "London",
					"ee": "London",
					"ht": "Lonn",
					"bi": "London",
					"na": "London",
					"uz": "London",
					"ay": "London",
					"yi": "לאנדאן",
					"kv": "Лондон",
					"bn": "লন্ডন",
					"sw": "London",
					"wa": "Londe",
					"ha": "Landan",
					"sq": "Londra",
					"oc": "Londres",
					"lo": "ລອນດອນ",
					"ne": "लन्डन",
					"bm": "London"
				},
				"lat": 51.5073219,
				"lon": -0.1276474,
				"country": "GB",
				"state": "England"
			}
		]`)
	gentclient := gentleman.New()
	gentclient.URL(host)
	gentclient.Use(mock.Plugin)

	w, err := geolocation.NewOpenWeatherFinder(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	lat, lon, err := w.FindLatLngByCity(city, "")
	if err != nil {
		t.Fatal(err)
	}

	// Asser
	assert.Greater(t, lat, float64(0), "London should be at the north hemisphere")
	assert.Less(t, lon, float64(0), "London should be west from the prime meridian")
}

func TestFindLatLngByCityWhenNoCityFoundShouldReturnError(t *testing.T) {

	// Arrange
	const host = "http://localhost1"
	const city = "imaginary-city"
	defer mock.Disable()
	mock.New(host).
		Get("/*").
		Reply(200).BodyString(`[]`)
	gentclient := gentleman.New()
	gentclient.URL(host)
	gentclient.Use(mock.Plugin)

	w, err := geolocation.NewOpenWeatherFinder(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	lat, lon, err := w.FindLatLngByCity(city, "")

	// Asser
	assert.NotNil(t, err)
	assert.Equal(t, lat, float64(0))
	assert.Equal(t, lon, float64(0))
}

func TestFindLatLngByCityWhenInvalidDataReceivedShouldReturnError(t *testing.T) {

	// Arrange
	const host = "http://localhost1"
	const city = "London"
	defer mock.Disable()
	mock.New(host).
		Get("/*").
		Reply(200).BodyString(`} this is a broken json`)
	gentclient := gentleman.New()
	gentclient.URL(host)
	gentclient.Use(mock.Plugin)

	w, err := geolocation.NewOpenWeatherFinder(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	lat, lon, err := w.FindLatLngByCity(city, "")

	// Asser
	assert.NotNil(t, err)
	assert.Equal(t, lat, float64(0))
	assert.Equal(t, lon, float64(0))
}

func TestFindLatLngByCityWhenUnauthorizedShouldReturnError(t *testing.T) {

	// Arrange
	const host = "http://localhost1"
	const city = "London"
	defer mock.Disable()
	mock.New(host).
		Get("/*").
		Reply(401).BodyString(`unauthorized text`)
	gentclient := gentleman.New()
	gentclient.URL(host)
	gentclient.Use(mock.Plugin)

	w, err := geolocation.NewOpenWeatherFinder(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	lat, lon, err := w.FindLatLngByCity(city, "")

	// Asser
	assert.NotNil(t, err)
	assert.Equal(t, lat, float64(0))
	assert.Equal(t, lon, float64(0))
}

func TestFindLatLngByCityWhenExternalErrorShouldReturnError(t *testing.T) {

	// Arrange
	const host = "http://localhost1"
	const city = "London"
	defer mock.Disable()
	mock.New(host).
		Get("/*").
		ReplyError(errors.New(`internal server error`))
	gentclient := gentleman.New()
	gentclient.URL(host)
	gentclient.Use(mock.Plugin)

	w, err := geolocation.NewOpenWeatherFinder(gentclient)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	lat, lon, err := w.FindLatLngByCity(city, "")

	// Asser
	assert.NotNil(t, err)
	assert.Equal(t, lat, float64(0))
	assert.Equal(t, lon, float64(0))
}

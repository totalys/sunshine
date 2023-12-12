# :partly_sunny: Sunshine Weather API :partly_sunny:

|:book: Table of contents|
|-|
|[:space_invader: tech stack](#space_invader-tech-stack)|
|[:dart: features ](#dart-features)|
|[:key: environment variables ](#key-environment-variables)|
|[:eyes: live documentation ](#eyes-live-docs)|
|[:gear: usage ](#gear-usage)|
|[:whale2: docker ](#whale2-docker)|
|[:test_tube: tests ](#test_tube-tests)|
|[:books: architecture ](#book-architecture)|

<!-- TechStack-->
### :space_invader: Tech Stack

<details>
    <summary>Language</summary>
    <ul>
        <li><a href="https://go.dev/">go-lang</a></li>
    </ul>
</details>
<details>
    <summary>Third party software</summary>
    <ul>
        <li><a href="https://openweathermap.org/">Open Weather</a></li>
        <li><a href="https://geocode.maps.co/">Free geocoding API for off line cities data</a></li>
    </ul>
</details>
<details>
    <summary>Main libraries</summary>
    <ul>
        <li><a href="https://github.com/labstack/echo">Echo</a> for http hosting</li>
        <li><a href="gopkg.in/h2non/gentleman.v2">Gentleman</a> for http requests</li>
        <li><a href="gopkg.in/h2non/gentleman-retry.v2">Retry</a> for resilience when communicating with external services</li>
        <li><a href="github.com/go-playground/validator">Validator</a> for simplifying validation at the controler layer</li>
    </ul>
</details>
<details>
    <summary>CI/CD</summary>
    <ul>
        <li><a href="https://github.com/features/actions">Github actions</a></li>
        <li><a href="https://www.docker.com/">Docker</a></li>
    </ul>
</details>
<details>
    <summary>Acknowledges</summary>
    <ul>
        <li><a href="https://mholt.github.io/json-to-go/">converts json into go structs</a> for saving time!</li>
        <li><a href="https://github.com/ikatyang/emoji-cheat-sheet">emojis-cheat-sheet</a> for improving this doc</li>
    </ul>
</details>

<!-- Features-->
### :dart: Features
 - Get current temperature for a given city in Celsius, Fahrenheit and Kelvin units 

<!-- EnvironmentVariables-->
### :key: Environment Variables

It is possible to run the application with a config.json file or with environment variables.
The default container provided embbeds the default values below but all variables set overrides the ones in the config.json file.

|name|default value|description|
|-|-|-|
|SUNNY_CONFIGFILE|configs/config.json| leave it empty to only use environment variables|
|SUNNY_SERVER_PORT|8080|port where the app will listen for requests|
|SUNNY_SWAGGER_ENABLED|true| enable swagger to expose a live documentation for users to try the api|
|SUNNY_SWAGGER_PATH|/swagger/*|path from where the swagger doc is presented|
|SUNNY_GEOLOCATION_INTERNAL|true|indicates if should use a built-in city maps go get geo coordinates|
|SUNNY_GEOLOCATION_PATH|configs/local_city_source.json|the built-in city map with citie's geo coordinates|
|SUNNY_EXTERNAL_GEOLOCATION_HOST|https://geocode.maps.co/search|the google maps api for geo coordinates|
|SUNNY_EXTERNAL_GEOLOCATION_APIKEY| - |this must be contracted|
|SUNNY_EXTERNAL_WEATHER_HOST|https://api.openweathermap.org|the open weather api host that provides climate data|
|SUNNY_EXTERNAL_WEATHER_APIKEY| - |this must be contracted|

<!-- Live Docs-->
### :eyes: Live Docs

- Start with Swagger enabled: true
- Access [swagger index](http://localhost:8080/swagger/index.html) in your browser.

#### :v: Updating swagger docs

* First time doing it?
    * Install the swaggo

``` bash
go install github.com/swaggo/swag/cmd/swag@latest
```


* After updating the api routes, run:
``` bash
swag init --parseDependency --parseInternal --parseDepth 1 -g pkg/controler/health.go -g pkg/controler/weather.go
```

<!-- Usage-->
### :gear: Usage

|before you go...|
|-|
|:warning: It's is necessary [subscribe](https://home.openweathermap.org/subscriptions/billing_info/onecall_30/base?key=base&service=onecall_30) to open-weather org to get an api-key and set it as environment variable before running this application. <br>ex: `$ export SUNNY_EXTERNAL_WEATHER_APIKEY=[your api key]`


* build

``` bash
go build -o build/sunshine.exe cmd/sunshine/main.go
```


* quick run

obs: It is necessary to add an api key for the openWeather api in the configs/config.jon to make it work properly. Alternative: `$ export SUNNY_EXTERNAL_WEATHER_APIKEY=[your api key]`

``` bash
./build/sunshine.exe --config-file configs/config.json
```

#### How to call this application

* A quick go-lang program do be used as a client

``` go

package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "http://localhost:8080/api/temperature?city=Sao%20Paulo&country=Brazil"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", "insomnia/2023.5.8")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
```

* a curl command to test the application

``` shell
curl --request GET \
  --url 'http://localhost:8080/api/temperature?city=Sao%20Paulo&country=Brazil'
```

<!-- Docker -->
### :whale2: Docker

* building an image

``` bash
docker build --no-cache --build-arg CI_VERSION=1.0.0 --build-arg CI_COMMIT_SHA=devepment-stage -t totalys/go/sunshine .
```

* running it

``` bash
docker run -d --rm -p 8080:8080 -e SUNNY_CONFIGFILE=configs/config.json -e SUNNY_EXTERNAL_WEATHER_APIKEY=[your open-weather api key]  --name sunshine totalys/go/sunshine
```

<!-- Tests-->
### :test_tube: Tests

``` bash
go test --cover ./...
```

#### Mocking

* mockery library

installation ref: [mockery documentation](https://vektra.github.io/mockery/latest/installation/)

``` bash
go install github.com/vektra/mockery/v2@v2.34.2
```

creating a new mock for some type

``` bash
mockery --name=[interface-name] --dir=[diretory-name] --output=[destination-path] --outpkg=[package-name]
```
example:

``` bash
mockery --name=WeatherService --dir=pkg/weather-service --output=pkg/mocks/weather-service --outpkg=weatherservice --filename=weather-service-mock.go
```

<!-- Architecture discussions -->
### :books: Architecture

A layered architecture in go lang may not look exactly like in others O.O systems. Go has a package oriented way of organizing code. Each package should be able to handle there own responsabilities. I try to stick with the standard proposed by [go.dev](https://go.dev/doc/modules/layout)
It looks like code is spread around but it is conceptually organized in layers of responsability

|Layer| package|
|-|-|
presentation|controler
business logic|weather-service
data providers|geolocation<br>open-weather


The diagram below show how the components of the system interact with each other over time. It shows how the system behaves and what the user might expect.

#### Sequence diagram

``` mermaid
---
title: Sunshine sequence diagram for retrieving temperature data
---
sequenceDiagram

actor U as user
box lightblue *office*
participant W as weather API
participant WS as weather service
participant G as geodata finder
end
participant eG as external geodata finder
participant eW as external open weather api

U->>W: sunshine/api/temperature?<br>city=a&country=b
W->>WS: getTemperatureForCity()
WS->>G: getGeodataForCity()
G-->>WS: latitude and longitude
alt using external geodata finder
WS->>eG: /geocode.maps.co/search?q=a,b
eG-->>WS: latitude and longitude
end
WS->>eW: api.openweathermap.org/data/3.0/onecall<br>querystring:lat,lon,api-secret
eW-->>WS: weather data
WS-->>W: temperature <br> Kelvin, Celcius and Farenheit
W-->>U: temperature <br> Kelvin, Celcius and Farenheit
```




<!-- Architecture discussions -->
### :lotus_position: Architecture discussions


- If the project grows, it would be interesting to consider some dependency injection library such as [wire](https://github.com/google/wire) or [dig](https://github.com/uber-go/dig). I'm not using any because it's a quite simple project.

- It is prudent to use encryption in communication. [let's encrypt](https://echo.labstack.com/docs/cookbook/auto-tls#server) helps to implement it in a relatively simple way. It would be a future implementation after an analysis of the API's usage conditions and where it will be exposed.

- If everbody starts to love the API and gets crazy about it making a lot of requests we may scale this app running more instances and add a simple reverse proxy with NginX to aliviate the processing. Reference: [load-balancing](https://echo.labstack.com/docs/cookbook/load-balancing)

- [The touble with rounding floating point numbers](https://www.theregister.com/2006/08/12/floating_point_approximation/) address the situation of dealing with float numbers. An even detailed discussion is exposed by [Goldberg D, 1991](https://docs.oracle.com/cd/E19957-01/806-3568/ncg_goldberg.html). For the sake of simplicity and because one decimal issue will not make anyone really sick, I'll just perform a simple string formatting in the calculated number.

- I like the [AAA](https://robertmarshall.dev/blog/arrange-act-and-assert-pattern-the-three-as-of-unit-testing/) pattern for unit testing mainly because it make it easier for new developers to understand which method is being tested, what belongs to the test preparation and what is being asserted.
    - [google's api test](https://github.com/googlemaps/google-maps-services-go/blob/c722fc00e8d79b4399832ea9204265d21e0483bc/geocoding_test.go#L24C25-L24C25) helped me to mock and implement the unit testing for the google maps service providing data and guidance.
- We could increase the performance a little bit by using [ffjson](https://github.com/pquerna/ffjson) which uses static marshaling functions to reduce the runtime reflection time.
- We could increase user experience adding some caching-layers. [go-cache](https://pkg.go.dev/github.com/patrickmn/go-cache) migth be useful for that.
    - caching geodata from the geofinder would save time since cities do not use to change their coordinates, only one call is needed to get the location for a given city
        - this cache not only makes the api fasters but also avoids possible instability issues when user is getting data from a cached city.
    - caching the temperature for some seconds, even for a minute will also increase experience since the temperature does not usualy change dramatically inside a minute.
- We could add a fast and pretty logging api like [zap](https://github.com/uber-go/zap)
    - Implementation started for the controler level
- We could add key-values correlation ID data in the context of the api request and log the correlation ID during the application flow to trace what happens to any particular request for fast troubleshooting. This feature is a must when the application scales and there is issues with some requests in particular

### :sound: Other discussions

* Some cities has the same name, specially in China. <br>Ex: [Suzhou, China:  {Lat:33.6333,Lng:116.9683}](https://www.google.com/maps/place/33%C2%B037'59.9%22N+116%C2%B058'05.9%22E/@33.6044436,116.9857392,12.82z/data=!4m4!3m3!8m2!3d33.6333!4d116.9683?entry=ttu), and [Suzhou, China:  {Lat:31.3,Lng:120.6194}](https://www.google.com/maps/place/Suzhou,+Jiangsu,+China/@31.2847188,120.6017754,11.75z/data=!4m15!1m8!3m7!1s0x35b3a0d19bd25e07:0x21e57f85bd766004!2sSuzhou,+Jiangsu,+China!3b1!8m2!3d31.2983399!4d120.58319!16zL20vMDFxcTgw!3m5!1s0x35b3a0d19bd25e07:0x21e57f85bd766004!8m2!3d31.2983399!4d120.58319!16zL20vMDFxcTgw?entry=ttu).

Even though the API accepts the country as an optional parameter, it may not return accurate results for some cases. For the sake of simplicity and in favor of most scenarios, the temperature returned will be considered for the first city returned by the consulted sources.

The option to use the internal map (configs/local_city_source.json) does not contain duplicate cities within the same country since the source was manually curated.

An alternative would be to return the temperature as an array of all cities with the given name but I'll let this option as a possible v2 of this API. We could even add a second endpoint for this option and let the user choose.

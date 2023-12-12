package configs

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"go.uber.org/zap"
)

const (
	//API the base path for the services
	API = "/api"

	// APIHealth health check path
	APIHealth = API + "/health"

	// APITemperature weather api path
	APITemperature = API + "/temperature"
)

const (
	secret = "secret"
	pass   = "pass"
	apikey = "apikey"
)

type Config struct {
	Log struct {
		Level string `env:"SUNNY_LOG_LEVEL"`
	}
	Server struct {
		Port    int `env:"SUNNY_SERVER_PORT"`
		Swagger struct {
			Enabled bool   `env:"SUNNY_SWAGGER_ENABLED"`
			Path    string `env:"SUNNY_SWAGGER_PATH"`
		}
		Geolocation struct {
			Internal bool   `env:"SUNNY_GEOLOCATION_INTERNAL"`
			Path     string `env:"SUNNY_GEOLOCATION_PATH"`
		}
	}
	External struct {
		Geolocation struct {
			Host   string `env:"SUNNY_EXTERNAL_GEOLOCATION_HOST"`
			ApiKey string `env:"SUNNY_EXTERNAL_GEOLOCATION_APIKEY"`
		}
		Weather struct {
			Host   string `env:"SUNNY_EXTERNAL_WEATHER_HOST"`
			ApiKey string `env:"SUNNY_EXTERNAL_WEATHER_APIKEY"`
		}
	}
}

func (c *Config) PrintConfigs(l *zap.Logger, version string) {

	l.Info("starting Sunshine monitor with the following configurations",
		zap.String("version", version))

	printConfigs(l, *c, "")
}

func printConfigs(l *zap.Logger, obj interface{}, prefix string) {

	var cfgs []reflect.StructField
	var cfgtype = reflect.TypeOf(obj)

	getPublicFields(&cfgs, &cfgtype)

	for i, c := range cfgs {
		if shouldOmmit(c.Name) {
			l.Info(fmt.Sprintf("%s: <ommited>", c.Name))
		} else {
			if c.Type.Kind() == reflect.Struct {
				printConfigs(l, reflect.ValueOf(obj).Field(i).Interface(), c.Name)
			} else {
				l.Info(fmt.Sprintf("%s.%s: %v",
					prefix, c.Name, reflect.ValueOf(obj).Field(i).Interface()))
			}
		}
	}
}

func getPublicFields(arr *[]reflect.StructField, typ *reflect.Type) {
	for i := 0; i < (*typ).NumField(); i++ {
		isPublic := unicode.IsUpper([]rune((*typ).Field(i).Name)[0])
		if isPublic {
			*arr = append(*arr, (*typ).Field(i))
		}
	}
}

// ommit passwords and secrets in log
func shouldOmmit(cfg string) bool {
	if strings.Contains(strings.ToLower(cfg), secret) {
		return true
	}

	if strings.Contains(strings.ToLower(cfg), pass) {
		return true
	}

	if strings.Contains(strings.ToLower(cfg), apikey) {
		return true
	}

	return false
}

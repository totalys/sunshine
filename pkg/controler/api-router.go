package controler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	"github.com/totalys/sunshine/configs"
	"github.com/totalys/sunshine/pkg/weather-service"
)

// Start starts the weather service
func Start(cfg *configs.Config, e *echo.Echo, l *zap.Logger, service weather.WeatherService, startChan chan<- struct{}) error {

	setup(cfg, e, l, service)
	e.Server.Addr = fmt.Sprintf(":%d", cfg.Server.Port)

	log.Printf("starting server at port [%d]", cfg.Server.Port)

	startChan <- struct{}{}

	if err := e.Server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop stops the weather service
func Stop(ctx context.Context, e *echo.Echo) error {

	return e.Shutdown(ctx)
}

func setup(cfg *configs.Config, e *echo.Echo, l *zap.Logger, service weather.WeatherService) {

	setSwaggerRouts(cfg, e)

	setHeathCheckRouts(cfg, e)

	setWeatherRouts(cfg, e, l, service)

}

func setSwaggerRouts(cfg *configs.Config, e *echo.Echo) {
	if cfg.Server.Swagger.Enabled {
		e.GET(cfg.Server.Swagger.Path, echoSwagger.WrapHandler)
	}
}

func setHeathCheckRouts(cfg *configs.Config, e *echo.Echo) {
	health := NewHealthControler()
	e.GET(configs.APIHealth, func(c echo.Context) error {
		return health.GetHealthCheck(c)
	})
}

func setWeatherRouts(cfg *configs.Config, e *echo.Echo, l *zap.Logger, service weather.WeatherService) {
	weather := NewWeatherControler(l, service)

	e.GET(configs.APITemperature, func(c echo.Context) error {
		return weather.GetTemperatureForCity(c)
	})
}

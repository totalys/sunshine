package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/iamolegga/enviper"
	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gopkg.in/eapache/go-resiliency.v1/retrier"
	retry "gopkg.in/h2non/gentleman-retry.v2"
	"gopkg.in/h2non/gentleman.v2"

	"github.com/totalys/sunshine/configs"
	_ "github.com/totalys/sunshine/docs"
	"github.com/totalys/sunshine/pkg/controler"
	"github.com/totalys/sunshine/pkg/geolocation"
	"github.com/totalys/sunshine/pkg/logger"
	openweather "github.com/totalys/sunshine/pkg/open-weather"
	"github.com/totalys/sunshine/pkg/weather-service"
)

const (
	appName = "SUNNY"
)

var (
	version  string = "v0.0.0"
	commitID string = "debug environment"

	configfile *string = flag.String("config-file", os.Getenv("SUNNY_CONFIGFILE"), "configuration file")
)

func main() {
	flag.Parse()

	cfg, err := configure()
	if err != nil {
		log.Fatal("error starting server. could not parse configuration ", err)
	}

	ctx, stopctx := context.WithCancel(context.Background())

	l, err := logger.New(cfg.Log.Level)
	if err != nil {
		log.Fatal("error starting server. could create logger", err)
	}

	defer func() {
		_ = l.Sync()
	}()

	cfg.PrintConfigs(l, fmt.Sprintf("%s commit: %s", version, commitID))

	go handleSignalInterruption(l, stopctx)

	errg, ctx := errgroup.WithContext(ctx)

	var fs = func() afero.Fs { return afero.NewOsFs() }

	var g geolocation.GeoFinder
	var w openweather.WeatherOfficer

	gentclient := gentleman.New()
	gentclient.URL(cfg.External.Weather.Host)
	gentclient.Use(retry.New(retrier.New(retrier.ExponentialBackoff(3, 100*time.Millisecond), nil)))

	if cfg.Server.Geolocation.Internal {
		g, err = geolocation.NewInternalCitiesFinder(cfg.Server.Geolocation.Path, fs())
		if err != nil {
			l.Fatal("error creating internal geo service", zap.Error(err))
		}
	} else {
		g, err = geolocation.NewOpenWeatherFinder(gentclient)
		if err != nil {
			l.Fatal("error creating google maps geo service", zap.Error(err))
		}
	}

	w, err = openweather.NewOpenWeather(gentclient)
	if err != nil {
		l.Fatal("error creating open weather client service", zap.Error(err))
	}

	service, err := weather.NewWeatherService(g, w)
	if err != nil {
		l.Fatal("error starting weather service", zap.Error(err))
	}

	e := echo.New()
	e.Validator = controler.NewCustomValidator(validator.New())

	startChan := make(chan struct{})

	errg.Go(func() error {
		return controler.Start(cfg, e, l, service, startChan)
	})

	<-startChan
	l.Info("weather service started")

	errg.Go(func() error {

		<-ctx.Done()

		c, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer handleShutdown(l, cancel)

		return controler.Stop(c, e)
	})

	if err := errg.Wait(); err != nil {
		l.Info("shutting down", zap.Error(err))
		stopctx()
	}
}

func configure() (*configs.Config, error) {

	v := enviper.New(viper.New())
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	if cfgfile := *configfile; cfgfile != "" {
		log.Printf("using config file: %s", cfgfile)
		v.SetConfigFile(cfgfile)
	}

	var cfg configs.Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// handleSignalInterruption handles signals from the os.
func handleSignalInterruption(l *zap.Logger, stopctx context.CancelFunc) {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)

	s := <-sig

	l.Info("received signal [%s]", zap.String("signal", s.String()))
	stopctx()
}

// handleShutdown performes a gracefull shutdown
func handleShutdown(l *zap.Logger, cancel context.CancelFunc) {

	l.Info("performing a gracefull shutting down")
	defer cancel()
}

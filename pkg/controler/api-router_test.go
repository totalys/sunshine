package controler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/totalys/sunshine/configs"
	"github.com/totalys/sunshine/pkg/controler"
	"github.com/totalys/sunshine/pkg/logger"
	weathermock "github.com/totalys/sunshine/pkg/mocks/weather-service"
)

func TestStartStop(t *testing.T) {

	// Arrange
	ctx := context.Background()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, configs.API, nil)
	rec := httptest.NewRecorder()
	e.NewContext(req, rec)

	cfg := &configs.Config{}
	startChan := make(chan struct{})
	serviceMock := new(weathermock.WeatherService)

	l, err := logger.New("info")
	if err != nil {
		t.Fatal("error returned is not HTTPError")
	}

	go func() {
		err = controler.Start(cfg, e, l, serviceMock, startChan)
	}()

	<-startChan
	if err != nil {
		t.Fatal(err)
	}

	// Act
	err = controler.Stop(ctx, e)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Nil(t, err)

}

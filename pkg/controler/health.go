package controler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const stateHealthy = "sunshine"

// HealthControler health check controler
type HealthControler interface {
	GetHealthCheck(c echo.Context) error
}

type healthControler struct {
}

var _ HealthControler = (*healthControler)(nil)

// NewHealthControler creates a controler for healthcheck
func NewHealthControler() *healthControler {
	return &healthControler{}
}

// GetHealthCheck returns wheater this application is alive or not.
// @Summary Gets the application status
// @Description Gets the application status. sunshine means it is working.
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {string} message "sunshine"
// @Router /api/health [get]
func (ctl *healthControler) GetHealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, stateHealthy)
}

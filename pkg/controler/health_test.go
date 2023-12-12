package controler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/totalys/sunshine/configs"
	"github.com/totalys/sunshine/pkg/controler"
)

func TestNewHealthControler(t *testing.T) {

	// Arrange & Act
	c := controler.NewHealthControler()

	// assert
	assert.NotNil(t, c)
}

func TestGetHealthCheck(t *testing.T) {

	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, configs.APIHealth, nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	c := controler.NewHealthControler()

	// Act
	err := c.GetHealthCheck(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, "\"sunshine\"\n", rec.Body.String(), "weather should be sunshine")

}

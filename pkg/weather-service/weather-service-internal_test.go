package weather

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKelvinToCelsius(t *testing.T) {

	// Arrange
	scenarios := []struct {
		input  float64
		expect float64
	}{
		{input: float64(0), expect: -273.15},
		{input: float64(100), expect: -173.15},
		{input: float64(273.15), expect: 0},
		{input: float64(373.15), expect: 100},
		{input: float64(500), expect: 226.85},
		{input: float64(1000), expect: 726.85},
		{input: float64(2000), expect: 1726.85},
		{input: float64(0.0), expect: -273.15},
		{input: float64(100.0023), expect: -173.15},
		{input: float64(273.1523), expect: 0},
		{input: float64(373.1534), expect: 100},
		{input: float64(500.0023), expect: 226.85},
		{input: float64(1000.00134), expect: 726.85},
		{input: float64(2000.0001), expect: 1726.85},
	}

	for _, s := range scenarios {
		// Act
		got := kelvinToCelsius(s.input)

		// Assert
		assert.Equal(t, s.expect, got)
	}

}

func TestCelsiusToFarenheit(t *testing.T) {

	// Arrange
	scenarios := []struct {
		input  float64
		expect float64
	}{
		{input: float64(-273.15), expect: -459.67},
		{input: float64(-173.15), expect: -279.67},
		{input: float64(0), expect: 32.00},
		{input: float64(100), expect: 212.00},
		{input: float64(273.15), expect: 523.67},
		{input: float64(373.15), expect: 703.67},
		{input: float64(500), expect: 932},
		{input: float64(1000), expect: 1832},
		{input: float64(-273.15012), expect: -459.67},
		{input: float64(-173.150123), expect: -279.67},
		{input: float64(0.0023), expect: 32.00},
		{input: float64(100.0013), expect: 212.00},
		{input: float64(273.1524), expect: 523.67},
		{input: float64(373.15021), expect: 703.67},
		{input: float64(500.00123), expect: 932},
		{input: float64(1000.000321), expect: 1832},
	}

	for _, s := range scenarios {
		// Act
		got := celsiusToFarenheit(s.input)

		// Assert
		assert.Equal(t, s.expect, got)
	}

}

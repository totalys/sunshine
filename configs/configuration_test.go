package configs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

const sample_cfg = `{
    "Server": {
        "Port" : 8080,
        "Swagger": {
            "Enabled": true,
            "Path": "/swagger/*"
        },
        "GeoLocation": {
            "Internal": false,
            "Path" : "configs/local_city_source.json"
        }
    },
    "External": {
       "GeoLocation": {
        "Host": "https://api.openweathermap.org",
        "ApiKey": ""
       },
       "Weather": {
        "Host": "https://api.openweathermap.org",
        "ApiKey": ""
       } 
    }
}`

type MemorySink struct {
	*bytes.Buffer
}

func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

func TestPrintConfigsShouldPrint(t *testing.T) {

	// Arrange
	var cfg Config

	if err := json.Unmarshal([]byte(sample_cfg), &cfg); err != nil {
		t.Fatal(err)
	}

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	l := zap.New(observedZapCore)

	// Act & Assert
	cfg.PrintConfigs(l, fmt.Sprintf("%s commit: %s", "1.0.0", "some_commit_id"))

	// Assert
	assert.Equal(t, 11, observedLogs.Len())
}

func TestShouldOmit(t *testing.T) {

	//Arrange
	scenarios := []struct {
		input    string
		expected bool
	}{
		{input: "secret", expected: true},
		{input: "Secret", expected: true},
		{input: "Secrets", expected: true},
		{input: "SecretsConfig", expected: true},
		{input: "pass", expected: true},
		{input: "password", expected: true},
		{input: "apikey", expected: true},
		{input: "ApiKey", expected: true},
		{input: "APIKEY", expected: true},
		{input: "MY_APIKEY", expected: true},
		{input: "MYAPIKEY", expected: true},
		{input: "any other word", expected: false},
		{input: "", expected: false},
	}

	for _, s := range scenarios {
		// Act

		got := shouldOmmit(s.input)

		// Assert
		assert.Equal(t, s.expected, got)
	}
}

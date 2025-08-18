package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAppConfig_AppAddress(t *testing.T) {
	tests := []struct {
		name     string
		cfg      AppConfig
		expected string
	}{
		{"http address", AppConfig{Address: "localhost", Port: 8080, SSLEnable: false}, "http://localhost:8080"},
		{"https address", AppConfig{Address: "localhost", Port: 443, SSLEnable: true}, "https://localhost:443"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.cfg.AppAddress())
		})
	}
}

func TestAppConfig_AppPort(t *testing.T) {
	cfg := AppConfig{Port: 3000}
	assert.Equal(t, ":3000", cfg.AppPort())
}

func TestDatabaseConfig_DSN_DSNSchema(t *testing.T) {
	db := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		Name:     "dbname",
		SSLMode:  "disable",
	}

	expectedDSN := "host=localhost port=5432 user=user password=pass dbname=dbname sslmode=disable"
	assert.Equal(t, expectedDSN, db.DSN())

	expectedURL := "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
	assert.Equal(t, expectedURL, db.DSNSchema())
}

func TestKafkaConfig_GetBrokersList(t *testing.T) {
	kc := KafkaConfig{Brokers: "broker1:9092;broker2:9092;;"}
	expected := []string{"broker1:9092", "broker2:9092"}
	assert.Equal(t, expected, kc.GetBrokersList())
}

func TestCORS_DefaultValues(t *testing.T) {
	c := CORSConfig{}
	if len(c.AllowMethods) == 0 {
		c.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if len(c.AllowHeaders) == 0 {
		c.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}
	if c.MaxAge == 0 {
		c.MaxAge = 5 * time.Minute
	}

	assert.Equal(t, 5*time.Minute, c.MaxAge)
	assert.Contains(t, c.AllowMethods, "GET")
	assert.Contains(t, c.AllowHeaders, "Content-Type")
}

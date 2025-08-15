package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App            AppConfig      `yaml:"app" env-prefix:"L0_"`
	Log            LogConfig      `yaml:"log" env-prefix:"L0_LOG_"`
	CORS           CORSConfig     `yaml:"cors" env-prefix:"L0_CORS_"`
	DatabaseConfig DatabaseConfig `yaml:"database" env-prefix:"L0_DB_"`
	KafkaConfig    KafkaConfig    `yaml:"kafka" env-prefix:"L0_KAFKA_"`
	CacheConfig    CacheConfig    `yaml:"cache" env-prefix:"L0_CACHE_"`
}

type AppConfig struct {
	Name      string `yaml:"name" env:"NAME" env-default:"gp"`
	Version   string `yaml:"version" env:"VERSION" env-default:"dev"`
	Env       string `yaml:"env" env:"ENV" env-default:"development"`
	Address   string `yaml:"address" env:"ADDRESS" env-default:"localhost"`
	Port      int    `yaml:"port" env:"PORT" env-default:"80"`
	SSL       bool   `yaml:"ssl" env:"SSL" env-default:"false"`
	SSLEnable bool   `yaml:"ssl_enable" env:"SSL_ENABLE"`
	SSLConfig struct {
		CertPath string `yaml:"cert_path" env:"CERT_PATH"`
		KeyPath  string `yaml:"key_path" env:"KEY_PATH"`
	} `yaml:"ssl_config"`
}

type CacheConfig struct {
	Capacity        int           `yaml:"capacity" env:"CAPACITY"`
	CleanupInterval time.Duration `yaml:"cleanup_interval" env:"CLEANUP_INTERVAL"`
}

type DatabaseConfig struct {
	Host            string        `yaml:"host" env:"HOST"`
	Port            int           `yaml:"port" env:"PORT"`
	User            string        `yaml:"user" env:"USER"`
	Password        string        `yaml:"password" env:"PASSWORD"`
	Name            string        `yaml:"name" env:"NAME"`
	SSLMode         string        `yaml:"sslmode" env:"SSLMODE"`
	MaxOpenConns    int           `yaml:"max_open_conns" env:"MAX_OPEN_CONNS"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"CONN_MAX_LIFETIME"`
}

func (d *DatabaseConfig) DSN() string {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
	if d.SSLMode == "disable" {
		return connString
	}
	return connString
}

func (d *DatabaseConfig) DSNSchema() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.SSLMode,
	)
}

type KafkaConfig struct {
	Brokers string `yaml:"brokers" env:"BROKERS"`
	Topic   string `yaml:"topic" env:"TOPIC"`
	GroupID string `yaml:"group_id" env:"GROUP_ID"`
}

func (kc *KafkaConfig) GetBrokersList() (br []string) {
	for _, broker := range strings.Split(kc.Brokers, ";") {
		if broker != "" {
			br = append(br, broker)
		}

	}
	return br
}

type LogConfig struct {
	Level        string `yaml:"level" env:"LEVEL"`
	Format       string `yaml:"format" env:"FORMAT"`
	ServiceName  string `yaml:"service_name" env:"SERVICE_NAME"`
	Environment  string `yaml:"environment" env:"ENVIRONMENT"`
	EnableCaller bool   `yaml:"enable_caller" env:"ENABLE_CALLER"`
}

type CORSConfig struct {
	AllowOrigins     []string      `yaml:"allow_origins" env:"ALLOW_ORIGINS"`
	AllowMethods     []string      `yaml:"allow_methods" env:"ALLOW_METHODS"`
	AllowHeaders     []string      `yaml:"allow_headers" env:"ALLOW_HEADERS"`
	AllowCredentials bool          `yaml:"allow_credentials" env:"ALLOW_CREDENTIALS"`
	MaxAge           time.Duration `yaml:"max_age" env:"MAX_AGE"`
}

func Load() (*Config, error) {
	var cfg Config

	// Читаем конфиг из файла и переменных окружения
	err := cleanenv.ReadConfig("./config/config.yaml", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	// Устанавливаем значения по умолчанию для CORS
	if len(cfg.CORS.AllowMethods) == 0 {
		cfg.CORS.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if len(cfg.CORS.AllowHeaders) == 0 {
		cfg.CORS.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}
	if cfg.CORS.MaxAge == 0 {
		cfg.CORS.MaxAge = 5 * time.Minute
	}
	return &cfg, err
}

func (c *AppConfig) AppAddress() string {
	prefix := "http://"
	address := "localhost"
	if c.SSLEnable {
		prefix = "https://"
		address = c.Address
	}
	return fmt.Sprintf("%s%s:%d", prefix, address, c.Port)
}

func (c *AppConfig) AppPort() string {
	return fmt.Sprintf(":%d", c.Port)
}

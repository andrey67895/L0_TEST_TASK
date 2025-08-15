package migrations

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/lib/pq"
)

//go:embed sql/*.sql
var migrations embed.FS

// Register registers the migrations filesystem
func Register() (source.Driver, error) {
	httpFS := http.FS(migrations)
	driver, err := httpfs.New(httpFS, "sql")
	if err != nil {
		return nil, fmt.Errorf("не удается создать драйвер: %v", err)
	}
	return driver, nil
}

// createMigrate creates a new migrate instance
func createMigrate(dsn string) (*migrate.Migrate, error) {
	driver, err := Register()
	if err != nil {
		return nil, fmt.Errorf("ошибка во время реегистрации миграции: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("httpfs", driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания инстанса миграции: %v", err)
	}

	return m, nil
}

// Validate checks if all migrations are applied
func Validate(dsn string) (bool, error) {
	m, err := createMigrate(dsn)
	if err != nil {
		return false, err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err == migrate.ErrNilVersion {
		return false, nil // no migrations applied yet
	}
	if err != nil {
		return false, fmt.Errorf("ошибка получения версии миграции: %v", err)
	}
	if dirty {
		return false, fmt.Errorf("проблемы с миграцией: dirty, %v, версия: %d", err, version)
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("ошибка проверки миграции: %v", err)
	}
	return false, nil
}

// Run executes database migrations
func Run(dsn string) error {
	m, err := createMigrate(dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка во время запуска миграции: %v", err)
	}
	return nil
}

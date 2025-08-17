package main

import (
	"context"
	bLog "log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/andrey67895/L0_TEST_TASK/internal/cache/in_memory"
	"github.com/andrey67895/L0_TEST_TASK/internal/config"
	"github.com/andrey67895/L0_TEST_TASK/internal/kafka"
	"github.com/andrey67895/L0_TEST_TASK/internal/logger"
	"github.com/andrey67895/L0_TEST_TASK/internal/migrations"
	"github.com/andrey67895/L0_TEST_TASK/internal/repository/postgres"
	"github.com/andrey67895/L0_TEST_TASK/internal/service"
	"github.com/andrey67895/L0_TEST_TASK/internal/transport/http/html"
	"github.com/andrey67895/L0_TEST_TASK/internal/transport/http/order"
	"github.com/andrey67895/L0_TEST_TASK/internal/transport/http/server"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		bLog.Fatalf("Ошибка при загрузке конфигурационного файла")
	}

	// Initialize logger
	log, err := logger.New(logger.Config{
		Level:        cfg.Log.Level,
		Format:       cfg.Log.Format,
		ServiceName:  cfg.Log.ServiceName,
		Environment:  cfg.Log.Environment,
		EnableCaller: cfg.Log.EnableCaller,
	})
	if err != nil {
		bLog.Fatalf("Ошибка при инициализации логгера")
	}
	defer log.Sync()
	valid, err := migrations.Validate(cfg.DatabaseConfig.DSNSchema())
	if err != nil {
		log.Fatal("Ошибка проверки статуса миграции", zap.Error(err))
	}

	if valid {
		log.Info("Миграция не требуется")
	} else {
		log.Info("Запускаем миграцию на Базе данных...")
		if err := migrations.Run(cfg.DatabaseConfig.DSNSchema()); err != nil {
			log.Fatal("Ошибка во время выполнения миграции", zap.Error(err))
		}
		log.Info("Все миграции произведены успешно")
	}

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Initialize database connection
	db, err := postgres.New(&cfg.DatabaseConfig)
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных", zap.Error(err))
	}
	defer db.Close()

	//repository
	orderRepo := postgres.NewOrderRepository(db)

	//cache
	orderCache := in_memory.NewInMemoryCache(cfg.CacheConfig.Capacity, cfg.CacheConfig.CleanupInterval)

	//service
	orderService := service.NewOrderService(orderRepo, orderCache)

	//prepare work
	go func() {
		if err = orderService.AddLastOrderInCache(ctx, cfg.CacheConfig.Capacity); err != nil {
			log.Error("Ошибка заполнения кеша актуальными данными из базы данных", zap.Error(err))
		}
	}()

	//handler
	orderHandler := order.NewHandler(log, orderService)
	htmlHandler := html.NewHandler(log)

	// Start server
	go func() {
		log.Info("Starting server",
			zap.String("service", cfg.App.Name),
			zap.String("version", cfg.App.Version),
			zap.String("environment", cfg.App.Env),
			zap.String("address", cfg.App.AppAddress()),
		)

		handlers := server.NewAPIHandlers(orderHandler, htmlHandler)
		middlewares := server.CreateMiddlewares(cfg)
		// Create server instance
		apiServer := server.New(cfg, log, handlers, middlewares)

		log.Info("LO - Test Task HTTP SERVER is UP")
		if err := apiServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("Ошибка в работе сервиса", zap.Error(err))
		}

	}()
	ctx, cancel := context.WithCancel(ctx)
	go kafka.NewKafkaService(ctx, cfg.KafkaConfig.GetBrokersList(), cfg.KafkaConfig.Topic, cfg.KafkaConfig.GroupID, log, orderService)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	cancel()

	// Shutdown gracefully
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("Shutting down HTTP SERVER")
	if err := e.Shutdown(ctx); err != nil {
		log.Error("Error during server shutdown", zap.Error(err))
	}

}

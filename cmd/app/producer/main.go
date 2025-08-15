package main

import (
	"context"
	"encoding/json"
	"fmt"
	bLog "log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/andrey67895/L0_TEST_TASK/internal/config"
	"github.com/andrey67895/L0_TEST_TASK/internal/logger"
)

func randInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func randString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func randomDelivery() map[string]string {
	return map[string]string{
		"name":    "Test Testov",
		"phone":   fmt.Sprintf("+972%07d", randInt(0, 9999999)),
		"zip":     fmt.Sprintf("%07d", randInt(1000000, 9999999)),
		"city":    "Kiryat Mozkin",
		"address": "Ploshad Mira 15",
		"region":  "Kraiot",
		"email":   fmt.Sprintf("test%d@gmail.com", randInt(1, 1000)),
	}
}

func randomPayment() map[string]interface{} {
	return map[string]interface{}{
		"transaction":   randString(16),
		"request_id":    "",
		"currency":      "USD",
		"provider":      "wbpay",
		"amount":        randInt(100, 2000),
		"payment_dt":    time.Now().Unix(),
		"bank":          "alpha",
		"delivery_cost": randInt(100, 500),
		"goods_total":   randInt(50, 1500),
		"custom_fee":    0,
	}
}

func randomItem() map[string]interface{} {
	return map[string]interface{}{
		"chrt_id":      randInt(1000000, 9999999),
		"track_number": "WBIL" + randString(8),
		"price":        randInt(50, 500),
		"rid":          randString(16),
		"name":         "Item " + randString(5),
		"sale":         randInt(0, 50),
		"size":         "0",
		"total_price":  randInt(50, 500),
		"nm_id":        randInt(100000, 999999),
		"brand":        "Brand " + randString(4),
		"status":       202,
	}
}

func randomOrder() map[string]interface{} {
	return map[string]interface{}{
		"order_uid":          randString(16),
		"track_number":       "WBIL" + randString(8),
		"entry":              "WBIL",
		"delivery":           randomDelivery(),
		"payment":            randomPayment(),
		"items":              []map[string]interface{}{randomItem()},
		"locale":             "en",
		"internal_signature": "",
		"customer_id":        "test" + randString(4),
		"delivery_service":   "meest",
		"shardkey":           "9",
		"sm_id":              99,
		"date_created":       time.Now().Format(time.RFC3339),
		"oof_shard":          "1",
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		bLog.Fatalf("Ошибка при загрузке конфигурационного файла")
	}
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.KafkaConfig.GetBrokersList()...),
		Topic:        cfg.KafkaConfig.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 100 * time.Millisecond,
	}
	defer writer.Close()

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
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		order := randomOrder()

		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Error("Ошибка сериализации JSON:", zap.Error(err))
			continue
		}

		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(order["order_uid"].(string)),
				Value: orderJSON,
			},
		)
		if err != nil {
			log.Error("Ошибка отправки в Kafka:", zap.Error(err))
			continue
		}
		log.Info("Отправлено сообщение в Kafka:", zap.Any("order_uid", order["order_uid"]))
	}
}

package model_factory

import (
	"fmt"
	"math/rand"
	"time"
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

func RandomOrder() map[string]interface{} {
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

package test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"

	"test/integration/clients"
	"test/integration/kafka_factory"
	"test/integration/model_factory"
)

func TestOrderUIDNotExist(t *testing.T) {
	ctx := t.Context()
	cs := setup(ctx)
	defer cleanup(ctx, cs)
	client := clients.GetClient(t)
	response, err := client.ApiV1GetOrderByOrderUid(ctx, uuid.New().String())
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestOrderUIDExist(t *testing.T) {
	ctx := t.Context()
	cs := setup(ctx)
	defer cleanup(ctx, cs)

	order := model_factory.RandomOrder()
	writer := kafka_factory.GetKafkaWriter()

	orderJSON, err := json.Marshal(order)
	require.NoError(t, err, "Ошибка сериализации JSON")
	msg := kafka.Message{
		Key:   []byte(order["order_uid"].(string)),
		Value: orderJSON,
	}
	t.Log(writer.Addr.String())
	t.Log(writer.Topic)

	err = writer.WriteMessages(ctx, msg)
	require.NoError(t, err, "Ошибка отправки в Kafka")

	timeout := time.After(60 * time.Second)
	tick := time.Tick(100 * time.Millisecond)

	client := clients.GetClient(t)

	var response *http.Response
	for {
		select {
		case <-timeout:
			require.Fail(t, "Сообщение не появилось в API за отведенное время")
		case <-tick:
			response, err = client.ApiV1GetOrderByOrderUid(ctx, string(msg.Key))
			if err == nil && response.StatusCode == http.StatusOK {
				goto done
			}
		}
	}
done:
	require.Equal(t, http.StatusOK, response.StatusCode)
}

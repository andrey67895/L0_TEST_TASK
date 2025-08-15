package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/andrey67895/L0_TEST_TASK/internal/domain"
	"github.com/andrey67895/L0_TEST_TASK/internal/logger"
	"github.com/andrey67895/L0_TEST_TASK/internal/service"
)

type KafkaService struct {
	KafkaReader  *kafka.Reader
	log          *logger.Logger
	orderService *service.OrderService
}

func NewKafkaReader(ctx context.Context, brokers []string, topic, groupID string, log *logger.Logger, orderService *service.OrderService) {
	ctx, cansel := context.WithCancel(ctx)
	ks := &KafkaService{
		KafkaReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupID,
			MinBytes:    1,
			MaxBytes:    10e6,
			MaxWait:     time.Second * 2,
			StartOffset: kafka.FirstOffset,
		}),
		log:          log,
		orderService: orderService,
	}
	messages := make(chan kafka.Message, 100)
	var wg sync.WaitGroup

	ks.log.Info("Потребитель Kafka начал работу...")
	ks.startWorkers(ctx, 5, messages, &wg)
	go ks.readMessages(ctx, messages)
	wg.Wait()
	ks.KafkaReader.Close()
	cansel()
}

func (k *KafkaService) processMessage(ctx context.Context, msg kafka.Message) error {
	k.log.Info(fmt.Sprintf("Обработка сообщения со смещением %d", msg.Offset))
	var order domain.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		k.log.Error(fmt.Sprintf("Некорректный JSON: %v", err))
		return nil // возвращаем nil, чтобы всё равно коммитить offset
	} else {
		if order.IsValid() {
			err = k.orderService.CreateOrder(ctx, order)
			if err != nil {
				return err
			}
		} else {
			//Пустой ORDER_ID
			return nil
		}
	}
	return nil
}

func (k *KafkaService) startWorkers(ctx context.Context, workerPoolSize int, messages <-chan kafka.Message, wg *sync.WaitGroup) {
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for msg := range messages {
				k.log.Info(fmt.Sprintf("Исполнитель %d обрабатывает сообщение...", id))
				err := k.processMessage(ctx, msg)
				if err != nil {
					k.log.Error(fmt.Sprintf("Исполнитель %d: сообщение об ошибке обработки: %v", id, err))
					continue
				}
				if err = k.KafkaReader.CommitMessages(ctx, msg); err != nil {
					k.log.Error(fmt.Sprintf("Исполнитель %d: не удалось зафиксировать сообщение: %v", id, err))
				}
			}
		}(i + 1)
	}
}

func (k *KafkaService) readMessages(ctx context.Context, messages chan<- kafka.Message) {
	for {
		msg, err := k.KafkaReader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			k.log.Error(fmt.Sprintf("ошибка при получении сообщения: %v", err))
			continue
		}
		k.log.Info("Вычитывание из Kafka произошло успешно", zap.String("Key", string(msg.Key)))
		messages <- msg
	}
	close(messages)
}

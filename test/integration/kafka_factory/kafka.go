package kafka_factory

import (
	"time"

	"github.com/segmentio/kafka-go"

	"test/integration/config"
)

func GetKafkaWriter() *kafka.Writer {
	conf := config.CreateKafkaConfigTest()
	writer := &kafka.Writer{
		Addr:         kafka.TCP(conf.GetBrokersList()...),
		Topic:        conf.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 100 * time.Millisecond,
	}
	return writer
}

package config

import "strings"

func CreateKafkaConfigTest() KafkaConfig {
	return KafkaConfig{
		Brokers: "localhost:29092",
		Topic:   "orders.events",
	}

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

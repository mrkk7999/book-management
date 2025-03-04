package kafka

import (
	"time"
)

type KafkaConfig struct {
	Brokers       []string
	Topic         string
	MaxRetries    int
	RetryInterval time.Duration
}

func NewKafkaConfig(brokers []string, topic string, maxRetries int, retryInterval time.Duration) *KafkaConfig {
	return &KafkaConfig{
		Brokers:       brokers,
		Topic:         topic,
		MaxRetries:    maxRetries,
		RetryInterval: retryInterval,
	}
}

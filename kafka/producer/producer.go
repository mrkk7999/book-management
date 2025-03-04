package producer

import (
	"book-management/kafka"
	"log"

	"github.com/IBM/sarama"
)

func NewAsyncProducer(config *kafka.KafkaConfig) sarama.AsyncProducer {
	producerConfig := sarama.NewConfig()

	producerConfig.Producer.Retry.Max = config.MaxRetries
	producerConfig.Producer.Retry.Backoff = config.RetryInterval
	producerConfig.Producer.Return.Errors = true
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(config.Brokers, producerConfig)
	if err != nil {
		log.Println("Error creating Kafka async producer: ", err)
		return nil
	}

	// Goroutine to log successes
	go func() {
		for msg := range producer.Successes() {
			log.Printf("Message delivered to topic %s, partition %d, offset %d",
				msg.Topic, msg.Partition, msg.Offset)
		}
	}()

	// Goroutine to log errors
	go func() {
		for err := range producer.Errors() {
			log.Printf("%v", err)
		}
	}()

	return producer
}

// PublishMessageAsynchronous
func PublishMessageAsynchronous(producer sarama.AsyncProducer, topic, message string) {

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	producer.Input() <- msg

	// Goroutine to check for errors or success
	go func() {
		select {
		case success := <-producer.Successes():
			log.Println(
				"topic", success.Topic,
				"partition", success.Partition,
				"offset", success.Offset,
			)
		case err := <-producer.Errors():
			log.Println("Failed to send message", err)
		}
	}()

	return
}

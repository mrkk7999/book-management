package consumer

import (
	"book-management/iface"
	models "book-management/request_response/books"
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

// StartConsumer
// initializes and starts a Kafka consumer group
func StartConsumer(brokers []string, topic string, groupID string, consumerSvc iface.ConsumerService) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// new consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Printf("Failed to start Sarama consumer group: %v", err)
		return
	}

	consumer := Consumer{
		ready:           make(chan bool),
		consumerService: consumerSvc,
	}

	ctx := context.Background()

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.Printf("Error from consumer: %v", err)
				return
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	log.Printf("Sarama consumer up and running!")
}

// Consumer
type Consumer struct {
	consumerService iface.ConsumerService
	ready           chan bool
}

// Setup
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Signal that the consumer is ready
	close(consumer.ready)
	log.Println("Consumer session setup complete")
	return nil
}

// Cleanup
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer session cleanup complete")
	return nil
}

// ConsumeClaim
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		go func(msg *sarama.ConsumerMessage) {
			var event models.BookEvent
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Printf("Error decoding message: %v", err)
				return
			}

			switch event.EventType {
			case "create_book":
				consumer.consumerService.CreateBook(event)
			case "update_book":
				consumer.consumerService.UpdateBook(event)
			case "delete_book":
				consumer.consumerService.DeleteBook(event)
			default:
				log.Printf("Unknown event type: %s", event.EventType)
			}

			session.MarkMessage(msg, "")
		}(message)
	}
	return nil
}

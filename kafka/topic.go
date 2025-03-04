package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// EnsureTopicExists
func EnsureTopicExists(brokers []string, topic string) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Producer.Return.Successes = true

	// Create a new Kafka admin client
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Printf("Error creating Kafka admin client: %v", err)
		return err
	}
	defer admin.Close()

	topics, err := admin.ListTopics()
	if err != nil {
		log.Printf("Error listing Kafka topics: %v", err)
		return err
	}

	if _, exists := topics[topic]; exists {
		log.Printf("Topic '%s' already exists, skipping creation.", topic)
		return nil
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1, // Should match KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR
	}

	// Creates topic
	err = admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		log.Printf("Failed to create topic '%s': %v", topic, err)
		return err
	}

	log.Printf("Topic '%s' created successfully.", topic)
	return nil
}

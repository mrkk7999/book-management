package main

import (
	"book-management/caches/redis"
	"book-management/controller"
	"book-management/implementation"
	"book-management/kafka"
	"book-management/kafka/consumer"
	"book-management/kafka/producer"
	"book-management/repository"
	httpTransport "book-management/transport/http"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.Info("Starting application...")

	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.WithError(err).Fatal("Error loading .env file")
	}

	var (
		httpAddr = os.Getenv("HTTP_ADDR")
		// Kafka
		brokers = []string{os.Getenv("BROKERS")}
		// brokers = []string{"localhost:9092"} // Hardcoded for debugging

		topic         = os.Getenv("TOPIC")
		groupID       = os.Getenv("GROUP_ID")
		maxRetries    = os.Getenv("MAX_RETRIES")
		retryInterval = os.Getenv("RETRY_INTERVAL")
		// Redis
		redisUrl = os.Getenv("REDIS_URL")
		// DB
		dbHost     = os.Getenv("DB_HOST")
		dbUser     = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbName     = os.Getenv("DB_NAME")
		dbPort     = os.Getenv("DB_PORT")
		dbSSLMode  = os.Getenv("DB_SSLMODE")
		dbTimeZone = os.Getenv("DB_TIMEZONE")
	)

	fmt.Println("Kafka Brokers:", os.Getenv("BROKERS"))
	fmt.Println("kafka borkers variable value", brokers)

	cache, err := redis.NewRedis(redisUrl)
	if err != nil {
		log.Println("Failed to connect to Redis")
	}
	log.Println(cache)

	err = kafka.EnsureTopicExists(brokers, topic)
	if err != nil {
		log.Println("Failed to ensure Kafka topic exists")
	}

	// Convert Kafka configurations
	maxRetriesInt, err := strconv.Atoi(maxRetries)
	if err != nil {
		log.Println("Failed to convert MAX_RETRIES to int")
	}

	retryIntervalInt, err := strconv.Atoi(retryInterval)
	if err != nil {
		log.Println("Failed to convert RETRY_INTERVAL to int")
	}

	// Producer
	kafkaConfig := kafka.NewKafkaConfig(brokers, topic, maxRetriesInt, time.Duration(retryIntervalInt)*time.Second)
	asyncProducer := producer.NewAsyncProducer(kafkaConfig)

	log.Println("Kafka producer initialized successfully.", asyncProducer)

	// PostgreSQL DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimeZone)

	log.Info("Connecting to database...")

	// Connect to Database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}

	log.Info("Successfully connected to the database.")

	repo := repository.New(db)
	svc := implementation.New(repo, topic, asyncProducer, cache, log)
	consumerSvc := implementation.NewConsumerService(repo, log)

	// Consumer
	consumer.StartConsumer(brokers, topic, groupID, consumerSvc)

	ctrl := controller.New(svc)
	handler := httpTransport.SetUpRouter(ctrl,log)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.WithField("address", httpAddr).Info("Starting HTTP server...")

	go func() {
		server := &http.Server{
			Addr:    httpAddr,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
	}()

	log.Error("exit", <-errs)
}

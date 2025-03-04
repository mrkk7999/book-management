package implementation

import (
	"book-management/caches"
	"book-management/iface"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type service struct {
	repo          iface.Respository
	asyncProducer sarama.AsyncProducer
	topic         string

	log   *logrus.Logger
	cache caches.Cache
}

func New(repo iface.Respository, topic string, asyncProducer sarama.AsyncProducer, cache caches.Cache, log *logrus.Logger) iface.Service {
	return &service{
		repo:          repo,
		asyncProducer: asyncProducer,
		topic:         topic,
		log:           log,
		cache:         cache,
	}
}

type consumerService struct {
	repository iface.Respository
	log        *logrus.Logger
}

func NewConsumerService(repository iface.Respository, log *logrus.Logger) iface.ConsumerService {
	return &consumerService{
		repository: repository,
		log:        log,
	}
}

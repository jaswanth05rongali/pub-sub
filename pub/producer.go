package pub

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func Producer(kafkaBrokerURL string) (p *kafka.Producer, e error) {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBrokerURL})

	p = producer
	e = err
	return p, e
}

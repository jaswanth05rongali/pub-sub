package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

var writer *kafka.Writer

func Configure(kafkaBrokerURLs []string, clientId string, topic string) (w *kafka.Writer, err error) {
	//Connects to kafka cluster
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}

	config := kafka.WriterConfig{
		Brokers:          kafkaBrokerURLs,
		Topic:            topic,
		Dialer:           dialer,
		Balancer:         &kafka.LeastBytes{},
		ReadTimeout:      10 * time.Second,
		WriteTimeout:     10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}

	w = kafka.NewWriter(config)
	writer = w
	return w, nil
}

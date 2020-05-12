package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaBrokerURL     string
	kafkaVerbose       bool
	kafkaTopic         string
	kafkaConsumerGroup string
	kafkaClientId      string
)

func main() {
	flag.StringVar(&kafkaBrokerURL, "kafka-brokers", "localhost:19092,localhost:29092,localhost:39092", "kafka-brokers URLs, All of them seperated with comma")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "kafka-verbose logging")
	flag.StringVar(&kafkaTopic, "kafka-topic", "foo", "Kafka Topic")
	flag.StringVar(&kafkaConsumerGroup, "kafka-consumers", "consumer-group", "Kafka consumer group")
	flag.StringVar(&kafkaClientId, "kafka-clientId", "my-client-id", "Kafka client Id")

	flag.Parse()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT)

	brokers := strings.Split(kafkaBrokerURL, ",")

	config := kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         kafkaClientId,
		Topic:           kafkaTopic,
		MinBytes:        10e3,
		MaxBytes:        10e6,
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Msgf("Error while receiving message: %s", err.Error())
			continue
		}

		value := m.Value
		fmt.Printf("Message in the partition:%v of topic:%v with offset:%v is %s", m.Partition, m.Topic, m.Offset, string(value))
	}
}

package producer

import (
	"fmt"
	"os"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// P stores the created producer instance
var P *kafka.Producer
var err error

//Init will initialize the producer function
func Init(kafkaBrokerURL string) (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBrokerURL})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	P = p
	fmt.Printf("Created Producer %v\n", P)
	return p, err
}

package producer

import (
	"fmt"
	"os"

	"github.com/jaswanth05rongali/pub-sub/logger"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// P stores the created producer instance
var P *kafka.Producer
var err error

//Init will initialize the producer function
func Init(kafkaBrokerURL string) {
	producerLogger := logger.Getlogger()
	P, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBrokerURL})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		producerLogger.Errorf("Failed to create producer:%s", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", P)
	producerLogger.Infof("Created Producer %v", P)
}

//GetProducer returns the producer variable
func GetProducer() *kafka.Producer {
	return P
}

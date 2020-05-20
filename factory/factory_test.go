package factory

import (
	"testing"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type TestDataItem struct {
	inputs   []Input
	result   kafka.Message
	hasError bool
}
type Input struct {
	kafkaPubMessageType string
	key                 string
	pubPartition        string
	kafkaTopic          string
	value               string
}

func TestFactoryGetMessage(t *testing.T) {

}

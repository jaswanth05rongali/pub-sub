package producer

import (
	"log"
	"testing"

	"github.com/jaswanth05rongali/pub-sub/logger"
)

func TestProducerInit(t *testing.T) {
	FileLocation := "./testlogs/producer_test.log"
	err := logger.NewLogger(FileLocation)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}

	broker := "localhost:19092"
	Init(broker)
	switch GetProducer().String() {
	case "rdkafka#producer-1":
		if broker == "localhost:19092" {
			t.Logf("Init Working Fine")
		} else {
			t.Errorf("Expected Init to return a nil value but got a producer instance")
		}
	default:
		if broker == "localhost:19092" {
			t.Errorf("Init should create a legitimate Producer. But not happened...")
		} else {
			t.Logf("Init Working Fine")
		}
	}
}

func TestGetProducer(t *testing.T) {
	Init("localhost:19092")
	result := GetProducer()
	if result != nil {
		t.Logf("GetProducer() PASSED,expected producer got %v", result)
	} else {
		t.Errorf("GetMessage() FAILED,expected <nil> got %v", result)
	}
}

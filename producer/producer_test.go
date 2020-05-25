package producer

import (
	"log"
	"testing"

	"github.com/jaswanth05rongali/pub-sub/logger"
)

func TestGetProducer(t *testing.T) {
	FileLocation := "./testlogs/producer_test.log"
	err := logger.NewLogger(FileLocation)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}
	Init("localhost:19092")
	result := GetProducer()
	if result != nil {
		t.Logf("GetProducer() PASSED,expected producer got %v", result)
	} else {
		t.Errorf("GetMessage() FAILED,expected <nil> got %v", result)
	}
}

func TestProducerInit(t *testing.T) {

}

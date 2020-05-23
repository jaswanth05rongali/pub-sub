package producer

import (
	"testing"
)

func TestGetProducer(t *testing.T) {
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

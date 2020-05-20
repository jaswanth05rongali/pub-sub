package worker

import (
	"testing"
)

func TestInit(t *testing.T) {
	broker := "localhost:2120"
	group := "Email"
	Init(broker, group)
	switch GetConsumer().String() {
	case "rdkafka#consumer-1":
		if broker == "localhost:19092" {
			t.Logf("Init Working Fine")
		} else {
			t.Errorf("Expected Init to return a nil value but got a consumer instance")
		}
	default:
		if broker == "localhost:19092" {
			t.Errorf("Init should create a legitimate Consumer. But not happened...")
		} else {
			t.Logf("Init Working Fine")
		}
	}
}

func TestConsume(t *testing.T) {

}

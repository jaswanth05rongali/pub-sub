package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/jaswanth05rongali/pub-sub/logger"
	"github.com/jaswanth05rongali/pub-sub/producer"
)

func TestMainServer(t *testing.T) {
	FileLocation := "./testlogs/main_test.log"
	err := logger.NewLogger(FileLocation)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}

	serverAddr := "0.0.0.0:9000"
	Broker := "192.168.99.100:19092"
	producer.Init(Broker)
	defer producer.GetProducer().Close()
	err = server(serverAddr)
	fmt.Println(err)
}

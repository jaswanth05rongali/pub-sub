package main

import (
	"fmt"
	"log"

	"github.com/jaswanth05rongali/pub-sub/client"
	"github.com/jaswanth05rongali/pub-sub/config"
	"github.com/jaswanth05rongali/pub-sub/logger"
	"github.com/jaswanth05rongali/pub-sub/worker"
	"github.com/namsral/flag"

	"github.com/spf13/viper"
)

var consumer *worker.ConsumerObject

var topics string

func main() {

	flag.StringVar(&topics, "topic", "Email", "Gets the topic from command line")
	flag.Parse()

	config.Init(false)

	FileLocation := "./consumer.log"
	err := logger.NewLogger(FileLocation)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}
	consumerLogger := logger.Getlogger()
	consumerLogger.Infof("Starting %v Consumer......", topics)

	broker := viper.GetString("broker")
	group := topics + "Group"
	// topics := viper.GetString("topic")

	client := client.Object{}
	consumer = &worker.ConsumerObject{ClientInterface: client}

	consumer.Init(broker, group)

	err = consumer.GetConsumer().Subscribe(topics, nil)
	if err != nil {
		fmt.Printf("Error:%v while subscribing to topic:%v\n", err, topics)
		consumerLogger.Errorf("Error:%v while subscribing to topic:%v", err, topics)
	}

	consumerChan := make(chan string)

	go func() {
		consumerChan <- consumer.Consume(false)
	}()

	outputString := <-consumerChan
	fmt.Println(outputString)
	consumerLogger.Infof(outputString)
	consumer.GetConsumer().Close()
}

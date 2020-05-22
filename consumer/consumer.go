package main

import (
	"fmt"

	"github.com/jaswanth05rongali/pub-sub/client"
	"github.com/jaswanth05rongali/pub-sub/config"
	"github.com/jaswanth05rongali/pub-sub/worker"

	"github.com/spf13/viper"
)

var consumer *worker.ConsumerObject

func main() {

	config.Init(false)

	broker := viper.GetString("broker")
	group := viper.GetString("group")
	topics := viper.GetString("topic")

	client := client.Object{}
	consumer = &worker.ConsumerObject{ClientInterface: client}

	consumer.Init(broker, group)

	err := consumer.GetConsumer().Subscribe(topics, nil)
	if err != nil {
		fmt.Printf("Error:%v while subscribing to topic:%v", err, topics)
	}

	consumer.Consume()
}

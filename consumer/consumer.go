package main

import (
	"fmt"

	"github.com/jaswanth05rongali/pub-sub/config"
	"github.com/jaswanth05rongali/pub-sub/worker"

	"github.com/spf13/viper"
)

func main() {

	config.Init(false)

	broker := viper.GetString("broker")
	group := viper.GetString("group")
	topics := viper.GetString("topic")

	worker.Init(broker, group)

	err := worker.C.Subscribe(topics, nil)
	if err != nil {
		fmt.Printf("Error:%v while subscribing to topic:%v", err, topics)
	}

	worker.Consume()
}

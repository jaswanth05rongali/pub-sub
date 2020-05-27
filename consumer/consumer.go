package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	if topics == "Phone" {
		go apiHandler("0.0.0.0:5001")
	} else {
		go apiHandler("0.0.0.0:5000")
	}

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

func apiHandler(listenAddr string) {
	gin.SetMode(gin.ReleaseMode)
	apiLogger := logger.Getlogger()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/health", healthController)
	err := r.Run(listenAddr)
	if err != nil {
		apiLogger.Errorf("Server not able to startup with error: ", err)
		os.Exit(1)
	}
}

func healthController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success",
	})

	ctx.Abort()
	return
}

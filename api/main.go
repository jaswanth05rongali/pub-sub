package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaswanth05rongali/pub-sub/config"
	"github.com/jaswanth05rongali/pub-sub/factory"
	"github.com/jaswanth05rongali/pub-sub/logger"
	"github.com/jaswanth05rongali/pub-sub/producer"

	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var producerLogger logger.Logger

var (
	listenAddrAPI  string
	kafkaBrokerURL string
	kafkaTopic     string
)

func main() {

	config.Init(true)

	FileLocation := "./api/producer.log"
	err := logger.NewLogger(FileLocation)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}
	producerLogger = logger.Getlogger()
	producerLogger.Infof("Starting API......")

	listenAddrAPI = viper.GetString("listenAddrAPI")
	kafkaBrokerURL = viper.GetString("kafkaBrokerURL")
	kafkaTopic = viper.GetString("kafkaTopic")

	producer.Init(kafkaBrokerURL)
	defer producer.GetProducer().Close()

	errChan := make(chan error, 1)

	go func() {
		producerLogger.Infof("starting server at %s", listenAddrAPI)
		errChan <- server(listenAddrAPI)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		producerLogger.Infof("got an interrupt, exiting...")
	case err := <-errChan:
		if err != nil {
			producerLogger.Errorf("error while running api, exiting...")
		}
	}
}

func server(listenAddr string) error {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.POST("/api/v1/data", postDataToKafka)

	for _, routeInfo := range router.Routes() {
		producerLogger.Debugf("Path: %v, handler: %v, method: %v, registered routes.", routeInfo.Path, routeInfo.Handler, routeInfo.Method)
	}

	return router.Run(listenAddr)
}

func postDataToKafka(ctx *gin.Context) {
	parent := context.Background()
	defer parent.Done()

	form := &struct {
		Requestid      string `json:"request_id"`
		Topicname      string `json:"topic_name"`
		Messagebody    string `json:"message_body"`
		Transactionid  string `json:"transaction_id"`
		Email          string `json:"email"`
		Phone          string `json:"phone"`
		Customerid     string `json:"customer_id"`
		Key            string `json:"key"`
		PubMessageType string `json:"pubMessageType"`
		PubPartition   string `json:"pubPartition"`
	}{}

	ctx.Bind(form)
	formInBytes, err := json.Marshal(form)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while marshalling json: %s", err.Error()),
			},
		})

		producerLogger.Errorf("error while marshalling json: %s", err.Error())
		ctx.Abort()
		return
	}

	deliveryChan := make(chan kafka.Event)

	value := string(formInBytes)
	kafkaTopic = form.Topicname
	message := factory.GetMessage(form.PubMessageType, form.Key, form.PubPartition, kafkaTopic, value)

	err = producer.GetProducer().Produce(&message, deliveryChan)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while push message into kafka: %s", err.Error()),
			},
		})

		producerLogger.Errorf("error while push message into kafka: %s", err.Error())
		ctx.Abort()
		return
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		producerLogger.Infof("Delivery failed: %v", m.TopicPartition.Error)
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Message push failed",
		})
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		producerLogger.Infof("Delivered message to topic %s [%d] at offset %v",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "success push data into kafka",
			"data":    form,
		})
	}

	close(deliveryChan)
}

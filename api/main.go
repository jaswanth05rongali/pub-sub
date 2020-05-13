package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jaswanth05rongali/second-app/pub"
	"github.com/rs/zerolog/log"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var logger = log.With().Str("pkg", "main").Logger()

var p *kafka.Producer

var (
	listenAddrAPI  string
	kafkaBrokerURL string
	kafkaTopic     string
)

func main() {
	flag.StringVar(&listenAddrAPI, "listen-address", "0.0.0.0:9000", "Listen address for api")
	flag.StringVar(&kafkaBrokerURL, "kafkaBroker", "localhost:19092", "URL of kafka broker")
	flag.StringVar(&kafkaTopic, "kafkaTopic", "foo", "kafka topic to push")

	flag.Parse()

	var err error
	p, err = pub.Producer(kafkaBrokerURL)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)
	defer p.Close()

	var errChan = make(chan error, 1)

	go func() {
		log.Info().Msgf("starting server at %s", listenAddrAPI)
		errChan <- server(listenAddrAPI)
	}()

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		logger.Info().Msg("got an interrupt, exiting...")
	case err := <-errChan:
		if err != nil {
			logger.Error().Err(err).Msg("error while running api, exiting...")
		}
	}
}

func server(listenAddr string) error {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.POST("/api/v1/data", postDataToKafka)

	for _, routeInfo := range router.Routes() {
		logger.Debug().
			Str("path", routeInfo.Path).
			Str("handler", routeInfo.Handler).
			Str("method", routeInfo.Method).
			Msg("registered routes")
	}

	return router.Run(listenAddr)
}

func postDataToKafka(ctx *gin.Context) {
	parent := context.Background()
	defer parent.Done()

	form := &struct {
		Text string `form:"text" json:"text"`
	}{}

	ctx.Bind(form)
	formInBytes, err := json.Marshal(form)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while marshalling json: %s", err.Error()),
			},
		})

		ctx.Abort()
		return
	}

	deliveryChan := make(chan kafka.Event)

	value := string(formInBytes)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, deliveryChan)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while push message into kafka: %s", err.Error()),
			},
		})

		ctx.Abort()
		return
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Message push failed",
		})
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "success push data into kafka",
			"data":    form,
		})
	}

	close(deliveryChan)
}

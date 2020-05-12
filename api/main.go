package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jaswanth05rongali/pub-sub/kafka"
	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
)

var logger = log.With().Str("pkg", "main").Logger()

var (
	listenAddrApi  string
	kafkaBrokerURL string
	kafkaVerbose   bool
	kafkaTopic     string
	kafkaClientId  string
)

func main() {
	flag.StringVar(&listenAddrApi, "listen-address", "0.0.0.0:9000", "Listen address for api")
	flag.StringVar(&kafkaBrokerURL, "kafka-brokers", "localhost:19092,localhost:29092,localhost:39092", "kafka-brokers URLs, All of them seperated with comma")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "kafka-verbose logging")
	flag.StringVar(&kafkaTopic, "kafka-topic", "foo", "Kafka Topic")
	flag.StringVar(&kafkaClientId, "kafka-clientId", "my-client-id", "Kafka client Id")

	flag.Parse()

	kafkaProducer, err := kafka.Configure(strings.Split(kafkaBrokerURL, ","), kafkaClientId, kafkaTopic)
	if err != nil {
		logger.Error().Str("error", err.Error()).Msg("unable to configure kafka")
		return
	}
	defer kafkaProducer.Close()

	var errChan = make(chan error, 1)
	go func() {
		log.Info().Msgf("starting server at %s", listenAddrApi)
		errChan <- server(listenAddrApi)
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
		logger.Debug().Str("path", routeInfo.Path).Str("handler", routeInfo.Handler).Str("method", routeInfo.Method).Msg("registered routes")
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

	err = kafka.Push(parent, nil, formInBytes)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while push message into kafka: %s", err.Error()),
			},
		})

		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "success push data into kafka",
		"data":    form,
	})
}

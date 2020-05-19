package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jaswanth05rongali/pub-sub/config"
	"github.com/jaswanth05rongali/pub-sub/worker"

	"github.com/spf13/viper"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var prevRetryTime int64

var (
	serverStatus   bool
	prevServerTime int64
	serverRunTime  int64
	serverDownTime int64
)

func main() {

	config.Init(false)

	broker := viper.GetString("broker")
	group := viper.GetString("group")
	topics := viper.GetString("topic")
	worker.Init(broker, group)
	os.Exit(1)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"broker.address.family": "v4",
		"group.id":              group,
		"session.timeout.ms":    6000,
		"enable.auto.commit":    false,
		"auto.offset.reset":     "latest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.Subscribe(topics, nil)

	run := true
	serverStatus = true
	prevServerTime = time.Now().Unix()
	serverRunTime = viper.GetInt64("serverRunTime")
	serverDownTime = viper.GetInt64("serverDownTime")
	// waitTime := viper.GetInt64("waitTime")
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)

			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				sentStatus := sendMessage(string(e.Value))
				if sentStatus == false {
					retrySendingMessage(string(e.Value))
					c.Commit()
				} else {
					c.Commit()
				}
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}

func sendMessage(value string) bool {
	currentTime := time.Now().Unix()
	if serverStatus == true {
		if currentTime > (prevServerTime + serverRunTime) {
			serverStatus = !serverStatus
			prevServerTime = currentTime
		}
	} else {
		if currentTime > (prevServerTime + serverDownTime) {
			serverStatus = !serverStatus
			prevServerTime = currentTime
		}
	}

	dataStrings := strings.Split(strings.Split(strings.Split(value, "{")[1], "}")[0], ",")
	requestString := strings.Split(dataStrings[0], ":")[1]
	requestBody := requestString[1 : len(requestString)-1]
	messageString := strings.Split(dataStrings[2], ":")[1]
	messageBody := messageString[1 : len(messageString)-1]
	emailString := strings.Split(dataStrings[4], ":")[1]
	emailID := emailString[1 : len(emailString)-1]
	// phoneNumber := strings.Split(strings.Split(dataStrings[5],":")[1],"'")[1]

	if serverStatus {
		fmt.Printf("Message: '%v' sent successfully to %v. Request ID: %v\n", messageBody, emailID, requestBody)
		return serverStatus
	} else {
		fmt.Printf("Message: '%v' delivery to %v failed. Server Down!! Request ID: %v\n", messageBody, emailID, requestBody)
		return serverStatus
	}
}

func retrySendingMessage(message string) {
	for {
		sent := sendMessage(message)
		if sent {
			break
		}
		time.Sleep(10000 * time.Millisecond)
	}
}

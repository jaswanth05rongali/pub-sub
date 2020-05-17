package main

import (
	"container/list"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jaswanth05rongali/first-app/config"

	"github.com/spf13/viper"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var prevRetryTime int64

var (
	serverStatus   bool
	prevServerTime int64
)

func main() {

	// if len(os.Args) < 4 {
	// 	fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <topics..>\n",
	// 		os.Args[0])
	// 	os.Exit(1)
	// }

	// broker := os.Args[1]
	// group := os.Args[2]
	// topics := os.Args[3:]

	config.Init(false)

	broker := viper.GetString("broker")
	group := viper.GetString("group")
	topics := viper.GetString("topic")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"broker.address.family": "v4",
		"group.id":              group,
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.Subscribe(topics, nil)

	run := true
	queue := list.New()
	serverStatus = true
	prevServerTime = time.Now().Unix()
	serverRunTime := viper.GetInt64("serverRunTime")
	serverDownTime := viper.GetInt64("serverDownTime")
	waitTime := viper.GetInt64("waitTime")
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			var continueFromHere string
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

			if queue.Len() > 0 {
				if currentTime > (prevRetryTime + waitTime) {
					if ev != nil {
						switch e := ev.(type) {
						case *kafka.Message:
							fmt.Printf("%% Message on %s:\n%s\n",
								e.TopicPartition, string(e.Value))
							queue.PushBack(string(e.Value))
							continueFromHere = retrySendingMessage(queue)
						default:
							continueFromHere = retrySendingMessage(queue)
						}
					} else {
						continueFromHere = retrySendingMessage(queue)
					}
					prevRetryTime = currentTime
				} else {
					if ev != nil {
						switch e := ev.(type) {
						case *kafka.Message:
							fmt.Printf("%% Message on %s:\n%s\n",
								e.TopicPartition, string(e.Value))
							queue.PushBack(string(e.Value))
						default:
						}
					}
					continueFromHere = "YES"
				}
			}
			if continueFromHere == "YES" {
				continue
			}
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				sentStatus := sendMessage(string(e.Value))
				if sentStatus == false {
					queue.PushBack(string(e.Value))
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

func retrySendingMessage(mQueue *list.List) string {
	for mQueue.Len() > 0 {
		element := mQueue.Front()
		sent := sendMessage(element.Value.(string))
		if sent {
			mQueue.Remove(element)
		} else {
			break
		}
	}

	if mQueue.Len() > 0 {
		return "YES"
	}

	return "NO"
}

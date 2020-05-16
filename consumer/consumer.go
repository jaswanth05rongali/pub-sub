package main

import (
	"container/list"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// var prevTimeRetry

func main() {

	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <topics..>\n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	group := os.Args[2]
	topics := os.Args[3:]
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

	err = c.SubscribeTopics(topics, nil)

	run := true
	queue := list.New()
	// waitTime := 300000
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			var continueFromHere string
			if queue.Len() > 0 {
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
	// messageBody := strings.Split(strings.Split(dataStrings[2],":")[1],"'")[1]
	emailString := strings.Split(dataStrings[4], ":")[1]
	emailID := emailString[1 : len(emailString)-1]
	// phoneNumber := strings.Split(strings.Split(dataStrings[5],":")[1],"'")[1]

	serverStatus := false
	if serverStatus {
		fmt.Printf("Message Sent Successfully to %v\n", emailID)
		return serverStatus
	} else {
		fmt.Printf("Message delivery to %v failed.\n", emailID)
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

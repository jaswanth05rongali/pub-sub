package worker

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaswanth05rongali/pub-sub/client"
	"github.com/rs/zerolog/log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// C stores the created producer instance
var C *kafka.Consumer

//ConsumerObject defines a struct for entire consumer along with few methodsC
type ConsumerObject struct {
	ClientInterface client.Interface
}

//Init will initialize the consumer function
func (cons *ConsumerObject) Init(broker string, group string) {
	var err error
	C, err = kafka.NewConsumer(&kafka.ConfigMap{
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

	fmt.Printf("Created Consumer %v\n", C)
}

//GetConsumer returns the consumer variable
func (cons *ConsumerObject) GetConsumer() *kafka.Consumer {
	return C
}

//Consume will help consuming messages from the cluster and also in sending them to the clients
func (cons *ConsumerObject) Consume(testCall bool) string {

	cons.ClientInterface.Init()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	iterations := 0
	run := true
	for run {
		if testCall && iterations == 10 {
			break
		}
		iterations++
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			var ev kafka.Event
			if testCall {
				var val kafka.Event

				if rand.Intn(30) < 15 {
					var message kafka.Message
					topic := "newTopic"
					message = kafka.Message{
						TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: int32(0)},
						Value:          []byte("{\"request_id\":\"1\",\"topic_name\":\"foo\",\"message_body\":\"Transaction Successful\",\"transaction_id\":\"987456321\",\"email\":\"kafka@gopostman.com\",\"phone\":\"9876543210\",\"customer_id\":\"1\",\"key\":\"1254\",\"pubMessageType\":\"0\",\"pubPartition\":\"3\"}"),
						Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
					}
					val = &message
				} else if rand.Intn(30) < 22 {
					var er kafka.Error
					val = er
				} else {
					val = nil
				}
				ev = val
			} else {
				ev = C.Poll(100)
				if ev == nil {
					continue
				}
			}

			switch e := ev.(type) {
			case *kafka.Message:
				message := string(e.Value)
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, message)
				sentStatus := cons.ClientInterface.SendMessage(message)
				if !sentStatus {
					checkRetry := cons.ClientInterface.RetrySendingMessage(message)
					if !checkRetry {
						err := cons.ClientInterface.SaveToFile(message)
						if err != nil {
							log.Error().Err(err).Msgf("Error while saving failed message to log file.")
						}
					}
				}
				C.Commit()
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

	return "Closing consumer"
}

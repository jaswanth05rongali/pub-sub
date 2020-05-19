package worker

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaswanth05rongali/pub-sub/client"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// C stores the created producer instance
var C *kafka.Consumer

//Init will initialize the consumer function
func Init(broker string, group string) {
	var err error
	C, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"broker.address.family": "v4",
		"group.id":              group,
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "latest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", C)
}

//GetConsumer returns the consumer variable
func GetConsumer() *kafka.Consumer {
	return C
}

//Consume will help consuming messages from the cluster and also in sending them to the clients
func Consume() {
	client.Init()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := C.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				message := string(e.Value)
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, message)
				sentStatus := client.SendMessage(message)
				if !sentStatus {
					client.RetrySendingMessage(message)
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

	fmt.Printf("Closing consumer\n")
	C.Close()
}

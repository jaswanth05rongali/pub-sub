package client

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	serverStatus   bool
	prevServerTime int64
	serverRunTime  int64
	serverDownTime int64
)

//Init will initialize the client variables
func Init() {
	serverStatus = true
	prevServerTime = time.Now().Unix()
	serverRunTime = viper.GetInt64("serverRunTime")
	serverDownTime = viper.GetInt64("serverDownTime")
	// waitTime := viper.GetInt64("waitTime")
}

//SendMessage will check for the client. If it is up and running then the func sends the message to the client and returns a true value. If client is down, it returns a false.
func SendMessage(value string) bool {
	currentTime := time.Now().Unix()
	if serverStatus {
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

//RetrySendingMessage will try resending the messages
func RetrySendingMessage(message string) {
	for {
		sent := SendMessage(message)
		if sent {
			break
		}
		time.Sleep(10000 * time.Millisecond)
	}
}

package client

import (
	"testing"
)

func TestSendMessage(t *testing.T) {
	if serverStatus == SendMessage("{\"request_id\":\"1\",\"topic_name\":\"foo\",\"message_body\":\"Transaction Successful\",\"transaction_id\":\"987456321\",\"email\":\"kafka@gopostman.com\",\"phone\":\"9876543210\",\"customer_id\":\"1\",\"key\":\"1254\",\"pubMessageType\":\"0\",\"pubPartition\":\"3\"}") {
		t.Logf("Send Message working fine.")
	} else {
		t.Errorf("Expected true from SendMessage but not successful")
	}
}

func TestRetrySendingMessage(t *testing.T) {

}

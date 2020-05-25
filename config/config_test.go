package config

import (
	"log"
	"testing"

	"github.com/jaswanth05rongali/pub-sub/logger"
	"github.com/spf13/viper"
)

func TestConfigInit(t *testing.T) {
	FileLocation := "./testlogs/config_test.log"
	err := logger.NewLogger(FileLocation)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}

	Init(false)
	broker := viper.GetString("broker")
	serverRunTime := viper.GetInt("serverRunTime")
	serverDownTime := viper.GetInt("serverDownTime")
	waitTime := viper.GetInt("waitTime")

	if broker == "localhost:19092" {
		t.Logf("Init(false) PASSED, expected \"localhost:19092\" got %v", broker)
	} else {
		t.Errorf("Init(false) PASSED, expected \"localhost:19092\" got %v", broker)
	}

	if serverRunTime == 120 {
		t.Logf("Init(false) PASSED, expected 120 got %v", serverRunTime)
	} else {
		t.Errorf("Init(false) FAILED, expected 120 got %v", serverRunTime)
	}

	if serverDownTime == 30 {
		t.Logf("Init(false) PASSED, expected 30 got %v", serverDownTime)
	} else {
		t.Errorf("Init(false) FAILED, expected 30 got %v", serverDownTime)
	}

	if waitTime == 10 {
		t.Logf("Init(false) PASSED, expected 10 got %v", waitTime)
	} else {
		t.Errorf("Init(false) FAILED, expected 10 got %v", waitTime)
	}

}

package logger

import (
	"reflect"
	"testing"
)

var testLogger Logger

func TestNewLogger(t *testing.T) {
	err := NewLogger("./testlogs/logger_test.log")
	if err == nil {
		t.Logf("NewLogger working fine...")
	} else {
		t.Errorf("Expected nil but got error:%v", err)
	}
	testLogger = Getlogger()
}

func TestDebugf(t *testing.T) {
	testLogger.Debugf("")
}

func TestInfof(t *testing.T) {
	testLogger.Infof("")
}

func TestWarnf(t *testing.T) {
	testLogger.Warnf("")
}

func TestErrorf(t *testing.T) {
	testLogger.Errorf("")
}

func TestGetLogger(t *testing.T) {
	if reflect.TypeOf(Getlogger()).String() == "logger.Logger" {
		t.Logf("GetLogger working fine!!")
	} else {
		t.Errorf("Expected Logger but unsuccessful")
	}
}

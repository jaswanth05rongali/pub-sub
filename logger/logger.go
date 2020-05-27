package logger

import (
	"io"
	"os"

	"github.com/jaswanth05rongali/pub-sub/logger/loggerconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

//Logger Object
type Logger struct {
	logger *logrus.Logger
}

// A global variable so that log functions can be directly accessed
var log Logger

func getFormatter(isJSON bool) logrus.Formatter {
	if isJSON {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}
}

//NewLogger returns an instance of logger
func NewLogger(FileLocation string) error {
	loggerconfig.Init()
	logLevel := viper.GetString("ConsoleLevel")
	if logLevel == "" {
		logLevel = viper.GetString("FileLevel")
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	stdOutHandler := os.Stdout
	fileHandler := &lumberjack.Logger{
		Filename:   FileLocation,
		MaxSize:    2,
		MaxBackups: 6,
		Compress:   true,
		MaxAge:     28,
	}
	lLogger := &logrus.Logger{
		Out:       stdOutHandler,
		Formatter: getFormatter(viper.GetBool("ConsoleJSONFormat")),
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}

	if viper.GetBool("EnableConsole") && viper.GetBool("EnableFile") {
		lLogger.SetOutput(io.MultiWriter(stdOutHandler, fileHandler))
	} else {
		if viper.GetBool("EnableFile") {
			lLogger.SetOutput(fileHandler)
			lLogger.SetFormatter(getFormatter(viper.GetBool("FileJSONFormat")))
		}
	}

	log.logger = lLogger

	return nil
}

//Debugf executes the Debugf call on logger
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

//Infof executes the Infof call on logger
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

//Warnf executes the Warnf call on logger
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

//Errorf executes the Errorf call on logger
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

//Fatalf executes the Fatalf call on logger
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

//Panicf executes the Panicf call on logger
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

//Getlogger returns the logger object
func Getlogger() Logger {
	return log
}

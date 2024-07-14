package constants

import (
	"log"
)

const (
	messageSeparator = " - "
	Warn             = "WARN"
	Info             = "INFO"
)

type LoggingUtils struct {
	loggerName string
	level      string
}

func NewLoggingUtils(loggerName string, level string) *LoggingUtils {
	return &LoggingUtils{
		loggerName: loggerName,
		level:      level,
	}
}

func (l *LoggingUtils) Info(message string) {
	log.Println(l.loggerName + messageSeparator + message)
}

func (l *LoggingUtils) Printf(format string, v ...interface{}) {
	log.Printf(l.loggerName+messageSeparator+format, v...)
}

func (l *LoggingUtils) Println(v ...interface{}) {
	log.Println(l.loggerName+messageSeparator, v)
}

func (l *LoggingUtils) WarnInfo(message string) {
	if l.level == Warn {
		log.Println(l.loggerName + messageSeparator + message)
	}
}

func (l *LoggingUtils) Error(message string, err error) {
	log.Println(l.loggerName+messageSeparator+message, err)
}

func (l *LoggingUtils) Warn(message string, err error) {
	if l.level == Warn {
		log.Println(l.loggerName+messageSeparator+message, err)
	}
}

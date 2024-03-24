package logger

import (
	"github.com/sirupsen/logrus"
)

type ILog interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

var _logger *Logger

func UseLogger() *Logger {
	return _logger
}

type Logger struct {
	logger *logrus.Logger
}

func InitLogger(logLevel string) error {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	logger := logrus.New()
	logger.SetLevel(level)

	_logger = &Logger{
		logger: logger,
	}

	return nil
}

func (l Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

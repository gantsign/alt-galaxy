package logging

import (
	"github.com/gantsign/alt-galaxy/internal/util"
)

type SerialLoggerFactory struct {
	loggers         chan SerialLogger
	completionLatch util.CompletionLatch
}

func (factory SerialLoggerFactory) NewLogger() SerialLogger {
	logger := SerialLogger{
		outputBuffer: make(chan message, 20),
	}
	factory.loggers <- logger
	return logger
}

func (factory SerialLoggerFactory) printOutput() {
	for logger := range factory.loggers {
		logger.printOutput()
	}
	factory.completionLatch.Success()
}

func (factory SerialLoggerFactory) StartOutput() {
	go factory.printOutput()
}

func (factory SerialLoggerFactory) Close() {
	close(factory.loggers)
}

func (factory SerialLoggerFactory) AwaitOutputComplete() {
	factory.completionLatch.Await()
}

func NewSerialLoggerFactory(queueSize int) SerialLoggerFactory {
	return SerialLoggerFactory{
		loggers:         make(chan SerialLogger, queueSize),
		completionLatch: util.NewCompletionLatch(1),
	}
}

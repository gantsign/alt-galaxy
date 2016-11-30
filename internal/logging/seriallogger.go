package logging

import (
	"fmt"

	"github.com/gantsign/alt-galaxy/internal/message"
)

type SerialLogger struct {
	outputBuffer chan message.Message
}

func (logger SerialLogger) Progressf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.outputBuffer <- message.Message{
		MessageType: message.OutMsg,
		Body:        msg,
	}
}

func (logger SerialLogger) Errorf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.outputBuffer <- message.Message{
		MessageType: message.ErrorMsg,
		Body:        msg,
	}
}

func (logger SerialLogger) Close() {
	close(logger.outputBuffer)
}

func NewSerialLogger(outputBuffer chan message.Message) SerialLogger {
	return SerialLogger{outputBuffer}
}

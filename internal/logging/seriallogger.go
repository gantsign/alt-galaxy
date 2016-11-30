package logging

import (
	"bufio"
	"fmt"
	"os"
	"time"

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

func (logger SerialLogger) PrintOutput() {
	stdOut := bufio.NewWriter(os.Stdout)
	stdErr := bufio.NewWriter(os.Stderr)
	for msg := range logger.outputBuffer {
		switch msg.MessageType {
		case message.OutMsg:
			fmt.Fprintln(stdOut, "- ", msg.Body)
			stdOut.Flush()
		case message.ErrorMsg:
			// A short sleep helps the stdout and stderr render in the correct order
			time.Sleep(time.Second)

			fmt.Fprintln(stdErr, "ERROR! ", msg.Body)
			stdErr.Flush()

			// A short sleep helps the stdout and stderr render in the correct order
			time.Sleep(time.Second)
		default:
			// A short sleep helps the stdout and stderr render in the correct order
			time.Sleep(time.Second)

			fmt.Fprintln(stdErr, fmt.Sprintf("ERROR! Unsupported MessageType: %d", msg.MessageType))
			stdErr.Flush()

			// A short sleep helps the stdout and stderr render in the correct order
			time.Sleep(time.Second)
		}
	}
}

func NewSerialLogger() SerialLogger {
	return SerialLogger{
		outputBuffer: make(chan message.Message, 20),
	}
}

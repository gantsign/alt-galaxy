package logging

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	outMsg   messageType = iota
	errorMsg messageType = iota
)

type messageType int

type message struct {
	messageType messageType
	body        string
}

type SerialLogger struct {
	outputBuffer chan message
}

func (logger SerialLogger) Progressf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.outputBuffer <- message{
		messageType: outMsg,
		body:        msg,
	}
}

func (logger SerialLogger) Errorf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.outputBuffer <- message{
		messageType: errorMsg,
		body:        msg,
	}
}

func (logger SerialLogger) Close() {
	close(logger.outputBuffer)
}

func (logger SerialLogger) PrintOutput() {
	stdOut := bufio.NewWriter(os.Stdout)
	stdErr := bufio.NewWriter(os.Stderr)
	for msg := range logger.outputBuffer {
		switch msg.messageType {
		case outMsg:
			fmt.Fprintln(stdOut, "- ", msg.body)
			stdOut.Flush()
		case errorMsg:
			// A short sleep helps the stdout and stderr render in the correct order
			time.Sleep(time.Second)

			fmt.Fprintln(stdErr, "ERROR! ", msg.body)
			stdErr.Flush()

			// A short sleep helps the stdout and stderr render in the correct order
			time.Sleep(time.Second)
		default:
			// A short sleep helps the stdout and stderr render in the correct order
			time.Sleep(time.Second)

			fmt.Fprintln(stdErr, fmt.Sprintf("ERROR! Unsupported MessageType: %d", msg.messageType))
			stdErr.Flush()

			// A short sleep helps the stdout and stderr render in the correct order
			time.Sleep(time.Second)
		}
	}
}

func NewSerialLogger() SerialLogger {
	return SerialLogger{
		outputBuffer: make(chan message, 20),
	}
}

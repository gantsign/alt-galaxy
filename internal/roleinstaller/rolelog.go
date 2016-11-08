package roleinstaller

import (
	"fmt"

	"github.com/gantsign/alt-galaxy/internal/message"
)

type roleLog struct {
	outputBuffer chan message.Message
}

func (roleLog roleLog) Progressf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	roleLog.outputBuffer <- message.Message{
		MessageType: message.OutMsg,
		Body:        msg,
	}
}

func (roleLog roleLog) Errorf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	roleLog.outputBuffer <- message.Message{
		MessageType: message.ErrorMsg,
		Body:        msg,
	}
}

func (roleLog roleLog) close() {
	close(roleLog.outputBuffer)
}

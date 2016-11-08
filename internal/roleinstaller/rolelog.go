package roleinstaller

import (
	"fmt"

	"github.com/gantsign/alt-galaxy/internal/logging"
	"github.com/gantsign/alt-galaxy/internal/message"
)

type roleLog struct {
	installer *roleInstaller
	roleIndex int
}

func (installer *roleInstaller) roleLog(role installerRole) logging.Log {
	return roleLog{installer, role.Index}
}

func (roleLog roleLog) Progressf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	output := roleLog.installer.roleOutputBuffers[roleLog.roleIndex]
	output <- message.Message{
		MessageType: message.OutMsg,
		Body:        msg,
	}
}

func (roleLog roleLog) Errorf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	output := roleLog.installer.roleOutputBuffers[roleLog.roleIndex]
	output <- message.Message{
		MessageType: message.ErrorMsg,
		Body:        msg,
	}
}

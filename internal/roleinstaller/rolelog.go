package roleinstaller

import (
	"fmt"

	"github.com/gantsign/alt-galaxy/internal/logging"
	"github.com/gantsign/alt-galaxy/internal/message"
	"github.com/gantsign/alt-galaxy/internal/rolesfile"
)

type roleLog struct {
	installer *roleInstaller
	roleIndex int
}

func (installer *roleInstaller) roleLog(role rolesfile.Role) logging.Log {
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

package pipeline

import (
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
)

type StepBase struct {
	ctx        model.Context
	pipeline   Pipeline
	RoleQueue  chan model.Role
	onComplete OnComplete
}

func (step *StepBase) init(ctx model.Context, pipeline Pipeline, queueSize int, onComplete OnComplete) {
	step.ctx = ctx
	step.pipeline = pipeline
	step.RoleQueue = make(chan model.Role, queueSize)
	step.onComplete = onComplete
}

func (step *StepBase) Context() model.Context {
	return step.ctx
}

func (step *StepBase) Queue(role model.Role) {
	step.RoleQueue <- role
}

func (step *StepBase) Fail(role model.Role) {
	step.pipeline.fail(role)
}

func (step *StepBase) Success(role model.Role) {
	step.onComplete(role)
}

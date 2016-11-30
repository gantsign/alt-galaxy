package pipeline

import (
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
	"github.com/gantsign/alt-galaxy/internal/util"
)

type ProcessRole func(ctx model.Context, role model.Role) (model.Role, error)

type StepBase struct {
	ctx        model.Context
	pipeline   Pipeline
	RoleQueue  chan model.Role
	Semaphore  util.Semaphore
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

func (step *StepBase) ConcurrentlyProcessRole(role model.Role, processor ProcessRole) {
	ctx := step.Context()
	step.Semaphore.Acquire()

	go func() {
		role, err := processor(ctx, role)
		step.Semaphore.Release()
		if err != nil {
			role.Errorf(err.Error())
			step.Fail(role)
		} else {
			step.Success(role)
		}
	}()
}

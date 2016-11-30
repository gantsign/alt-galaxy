package pipeline

import (
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
	"github.com/gantsign/alt-galaxy/internal/util"
)

type ProcessRole func(ctx model.Context, step Step, role model.Role)

type factoryStep struct {
	StepBase
	processRole ProcessRole
	semaphore   util.Semaphore
}

func (step *factoryStep) Fail(role model.Role) {
	step.semaphore.Release()
	step.pipeline.fail(role)
}

func (step *factoryStep) Success(role model.Role) {
	step.semaphore.Release()
	step.onComplete(role)
}

func (step *factoryStep) processRoles() {
	ctx := step.Context()
	for role := range step.RoleQueue {
		step.semaphore.Acquire()

		go step.processRole(ctx, step, role)
	}
}

func (step *factoryStep) Start() {
	go step.processRoles()
}

func NewStep(processRole ProcessRole, maxConcurrent int) Step {
	return &factoryStep{
		processRole: processRole,
		semaphore:   util.NewSemaphore(maxConcurrent),
	}
}

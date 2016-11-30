package pipeline

import (
	"github.com/gantsign/alt-galaxy/internal/logging"
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
	"github.com/gantsign/alt-galaxy/internal/rolesfile"
	"github.com/gantsign/alt-galaxy/internal/util"
)

type Pipeline interface {
	InstallRole(fileRole rolesfile.Role)

	Start()

	fail(role model.Role)

	Await() bool
}

type pipeline struct {
	ctx           model.Context
	loggerFactory logging.SerialLoggerFactory
	queueSize     int
	Steps         []Step
	roleLatch     util.CompletionLatch
}

func (pipeline *pipeline) InstallRole(fileRole rolesfile.Role) {
	logger := pipeline.loggerFactory.NewLogger()

	role := model.NewRole(fileRole, logger)

	pipeline.roleLatch.TaskAdded()
	pipeline.Steps[0].Queue(role)
}

func (pipeline *pipeline) Start() {
	pipeline.loggerFactory.StartOutput()

	ctx := pipeline.ctx
	queueSize := pipeline.queueSize

	onComplete := func(role model.Role) {
		pipeline.success(role)
	}

	for i := len(pipeline.Steps) - 1; i >= 0; i-- {
		pipeline.Steps[i].init(ctx, pipeline, queueSize, onComplete)
		pipeline.Steps[i].Start()

		onComplete = func(nextStep Step) func(role model.Role) {
			return func(role model.Role) {
				nextStep.Queue(role)
			}
		}(pipeline.Steps[i])
	}
}

func (pipeline *pipeline) fail(role model.Role) {
	role.Progressf("%s install failed", role.Name)
	role.Close()
	pipeline.roleLatch.Failure()
}

func (pipeline *pipeline) success(role model.Role) {
	role.Progressf("%s was installed successfully", role.Name)
	role.Close()
	pipeline.roleLatch.Success()
}

func (pipeline *pipeline) Await() bool {
	success := pipeline.roleLatch.Await()

	pipeline.loggerFactory.Close()
	pipeline.loggerFactory.AwaitOutputComplete()

	return success
}

func NewInstallPipeline(ctx model.Context, queueSize int, steps []Step) Pipeline {
	return &pipeline{
		ctx:           ctx,
		loggerFactory: logging.NewSerialLoggerFactory(queueSize),
		queueSize:     queueSize,
		Steps:         steps,
		roleLatch:     util.NewCompletionLatch(0),
	}
}

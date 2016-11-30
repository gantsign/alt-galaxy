package pipeline

import (
	"github.com/gantsign/alt-galaxy/internal/roleinstaller/internal/model"
)

type OnComplete func(role model.Role)

type Step interface {
	init(ctx model.Context, pipeline Pipeline, queueSize int, onComplete OnComplete)

	Context() model.Context

	Queue(role model.Role)

	Start()

	Success(role model.Role)

	Fail(role model.Role)
}

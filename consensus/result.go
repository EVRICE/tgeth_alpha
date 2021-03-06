package consensus

import (
	"context"

	"github.com/EVRICE/tgeth_alpha/common/debug"
	"github.com/EVRICE/tgeth_alpha/core/types"
	"github.com/EVRICE/tgeth_alpha/log"
)

type ResultWithContext struct {
	Cancel
	*types.Block
}

type Cancel struct {
	context.Context
	cancel context.CancelFunc
}

func (c *Cancel) CancelFunc() {
	log.Debug("cancel mining task", "callers", debug.Callers(10))
	c.cancel()
}

func NewCancel(ctxs ...context.Context) Cancel {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		ctx = context.Background()
	}
	ctx, cancelFn := context.WithCancel(ctx)
	return Cancel{ctx, cancelFn}
}

func StabCancel() Cancel {
	return Cancel{
		context.Background(),
		func() {},
	}
}

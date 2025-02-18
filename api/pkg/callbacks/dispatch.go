package callbacks

import (
	"context"
	"github.com/pkg/errors"
	"strings"
)

type DispatchFunc func(ctx *Context, args []byte) (any, error)

var dispatchFunctions = map[string]DispatchFunc{}

func (s *service) Dispatch(ctx context.Context, actionName string, args []byte, metadata Metadata) (any, error) {
	myCtx := &Context{
		Context:  ctx,
		Metadata: metadata,

		config: s.config,
	}

	actionName = strings.ToLower(actionName)
	if dispatchFunctions[actionName] != nil {
		return dispatchFunctions[actionName](myCtx, args)
	}

	return nil, errors.Errorf("unknown action: %s", actionName)
}

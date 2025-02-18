package callbacks

import (
	"context"
	"github.com/pkg/errors"
	"strings"
)

type DispatchFunc func(s *service, ctx context.Context, args []byte, metadata Metadata) (any, error)

var dispatchFunctions = map[string]DispatchFunc{}

func (s *service) Dispatch(ctx context.Context, actionName string, args []byte, metadata Metadata) (any, error) {
	actionName = strings.ToLower(actionName)
	if dispatchFunctions[actionName] != nil {
		return dispatchFunctions[actionName](s, ctx, args, metadata)
	}

	return nil, errors.Errorf("unknown action: %s", actionName)
}

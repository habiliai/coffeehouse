package action

import (
	"context"
	"github.com/pkg/errors"
	"strings"
)

func (s *service) Dispatch(ctx context.Context, actionName string, args []byte) ([]byte, error) {
	switch strings.ToLower(actionName) {
	case "send_post_to_twitter":
		return s.SendPostToTwitter(ctx, args)
	default:
		return nil, errors.Errorf("unknown action: %s", actionName)
	}
}

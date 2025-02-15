package callbacks

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
)

func SendPostToTwitter(s *service, ctx context.Context, args []byte, metadata Metadata) (any, error) {
	request := struct {
		Content  string   `json:"content"`
		Hashtags []string `json:"hashtags"`
		Media    []string `json:"media"`
	}{}

	if err := json.Unmarshal(args, &request); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal args")
	}

	logger.Debug("sendPostToTwitter", "request", request)

	return nil, nil
}

func init() {
	dispatchFunctions["send_post_to_twitter"] = SendPostToTwitter
}

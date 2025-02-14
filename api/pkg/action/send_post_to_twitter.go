package action

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
)

func (s *service) SendPostToTwitter(ctx context.Context, args []byte) ([]byte, error) {
	request := struct {
		Content  string   `json:"content"`
		Hashtags []string `json:"hashtags"`
		Media    []string `json:"media"`
	}{}

	response := struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}

	if err := json.Unmarshal(args, &request); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal args")
	}

	logger.Debug("sendPostToTwitter", "request", request)

	resJson, err := json.Marshal(&response)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal response")
	}

	return resJson, nil
}

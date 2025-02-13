package habapi

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
)

func (s *server) getRun(ctx context.Context, threadId string, runId string) (*openai.Run, error) {
	run, err := s.openai.Beta.Threads.Runs.Get(ctx, threadId, runId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get run with id %s", runId)
	}

	return run, nil
}

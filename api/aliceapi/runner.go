package aliceapi

import (
	"context"
	_ "embed"
	"github.com/habiliai/agentruntime/runtime"
)

func (s *server) run(ctx context.Context, threadId uint, agentIds []uint32) error {
	thread, err := s.getThread(ctx, int32(threadId))
	if err != nil {
		return err
	}

	ctx = context.WithoutCancel(ctx)

	req := &runtime.RunRequest{
		ThreadId: thread.AgentRuntimeThreadId,
		AgentIds: agentIds,
	}

	_, err = s.runtime.Run(ctx, req)
	return err
}

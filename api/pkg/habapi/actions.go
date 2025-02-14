package habapi

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
)

func (s *server) requiresAction(ctx context.Context, run *openai.Run) (*openai.Run, error) {
	if len(run.RequiredAction.SubmitToolOutputs.ToolCalls) == 0 {
		return nil, errors.Errorf("no tool calls found")
	}

	toolCall := run.RequiredAction.SubmitToolOutputs.ToolCalls[0]

	output, err := s.actionService.Dispatch(ctx, toolCall.Function.Name, []byte(toolCall.Function.Arguments))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute action %s", toolCall.Function.Name)
	}

	toolCallRun, err := s.openai.Beta.Threads.Runs.SubmitToolOutputs(ctx, run.ThreadID, run.ID, openai.BetaThreadRunSubmitToolOutputsParams{
		ToolOutputs: openai.F([]openai.BetaThreadRunSubmitToolOutputsParamsToolOutput{
			{
				ToolCallID: openai.F(toolCall.ID),
				Output:     openai.F(string(output)),
			},
		}),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to submit tool outputs")
	}

	return toolCallRun, nil
}

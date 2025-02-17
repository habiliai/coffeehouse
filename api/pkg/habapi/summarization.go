package habapi

import (
	"context"
	_ "embed"
	"github.com/Masterminds/sprig/v3"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
	"text/template"
)

type SummarizeInstructionValues struct {
	Context []struct {
		Name string
		Text string
	}
	Agents []struct {
		Name string
		Role string
	}
}

var (
	//go:embed data/instructions/thread.summarize.md.tmpl
	summarizeInstruction string

	summarizeInstructionTemplate = template.Must(
		template.New("thread_summarize_instructions").Funcs(sprig.FuncMap()).Parse(summarizeInstruction),
	)
)

func (s *server) SummarizeThread(
	ctx context.Context,
	req *ThreadId,
) (*emptypb.Empty, error) {
	if err := s.summarizeThread(ctx, uint(req.Id)); err != nil {
		return nil, errors.Wrap(err, "failed to summarize thread")
	}

	return &emptypb.Empty{}, nil
}

func (s *server) summarizeThread(
	ctx context.Context,
	threadId uint,
) error {
	_, thread, err := s.getThread(ctx, int32(threadId))
	if err != nil {
		return errors.Wrap(err, "failed to get thread")
	}

	allMessages, err := s.getAllMessages(ctx, thread.OpenaiThreadId)
	if err != nil {
		return errors.Wrapf(err, "failed to get all messages")
	}

	instructionValues := SummarizeInstructionValues{}
	appendedAgentIds := map[int32]struct{}{}
	for _, msg := range allMessages {
		switch msg.Role {
		case Message_USER:
			var prefixBuilder strings.Builder
			for _, mention := range msg.Mentions {
				prefixBuilder.WriteString(mention)
				prefixBuilder.WriteString(" ")
			}
			instructionValues.Context = append(instructionValues.Context, struct {
				Name string
				Text string
			}{
				Name: "User",
				Text: prefixBuilder.String() + msg.Text,
			})
		case Message_ASSISTANT:
			instructionValues.Context = append(instructionValues.Context, struct {
				Name string
				Text string
			}{
				Name: msg.Agent.Name,
				Text: msg.Text,
			})
			if _, ok := appendedAgentIds[msg.Agent.Id]; !ok {
				instructionValues.Agents = append(instructionValues.Agents, struct {
					Name string
					Role string
				}{
					Name: msg.Agent.Name,
					Role: msg.Agent.Role,
				})
				appendedAgentIds[msg.Agent.Id] = struct{}{}
			}
		default:
			return errors.Wrapf(err, "unexpected message role: %s", msg.Role)
		}
	}

	var instructionBuilder strings.Builder
	if err := summarizeInstructionTemplate.Execute(&instructionBuilder, instructionValues); err != nil {
		return errors.Wrapf(err, "failed to execute instruction template")
	}

	completion, err := s.openai.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: openai.F("gpt-4o"),
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.ChatCompletionSystemMessageParam{
					Role: openai.F(openai.ChatCompletionSystemMessageParamRoleSystem),
					Content: openai.F([]openai.ChatCompletionContentPartTextParam{
						{
							Text: openai.F(instructionBuilder.String()),
							Type: openai.F(openai.ChatCompletionContentPartTextTypeText),
						},
					}),
				},
			}),
		},
	)
	if err != nil {
		return errors.Wrapf(err, "failed to get completion")
	}

	summary := completion.Choices[0].Message.Content

	if strings.HasPrefix(summary, "```markdown") {
		summary, _ = strings.CutPrefix(summary, "```markdown")
	} else if strings.HasPrefix(summary, "```") {
		summary, _ = strings.CutPrefix(summary, "```")
	}
	if strings.HasSuffix(summary, "```") {
		summary, _ = strings.CutSuffix(summary, "```")
	}
	summary = strings.TrimSpace(summary)

	if err := helpers.GetTx(ctx).Model(&domain.Thread{}).
		Where("id = ?", thread.ID).
		Update("summary", summary).
		Error; err != nil {
		return errors.Wrapf(err, "failed to update thread")
	}

	return nil
}

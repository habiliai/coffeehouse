package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	haberrors "github.com/habiliai/habiliai/api/pkg/errors"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/mokiat/gog"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
	"time"
)

func (s *server) AddMessage(ctx context.Context, req *AddMessageRequest) (*emptypb.Empty, error) {
	tx := helpers.GetTx(ctx)
	thread, threadData, err := s.getThread(ctx, req.ThreadId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get thread with id %s", req.ThreadId)
	}

	message := strings.TrimSpace(req.Message)
	if message == "" {
		return nil, errors.New("message is empty")
	}

	isMention := strings.HasPrefix(message, "@")
	if !isMention {
		return nil, errors.New("you must mention at least an agent")
	}

	messagePieces := strings.Split(message, " ")
	if len(messagePieces) == 0 {
		return nil, errors.New("you must mention at least an agent")
	}

	textPieces := make([]string, 0, len(messagePieces))
	mentionedAgents := make([]domain.Agent, 0, len(messagePieces))
	for _, messagePiece := range messagePieces {
		if !strings.HasPrefix(messagePiece, "@") {
			textPieces = append(textPieces, messagePiece)
			continue
		} else {
			agentName := strings.TrimPrefix(messagePiece, "@")
			var agent domain.Agent
			if err := tx.First(&agent, "name = ?", agentName).Error; err != nil {
				return nil, errors.Wrapf(err, "failed to find agent with name %s", agentName)
			}
			mentionedAgents = append(mentionedAgents, agent)
		}
	}

	messageData := NewEmptyMessageData(s.openai, thread.ID)
	messageData.SetAgentIds(gog.Map(mentionedAgents, func(agent domain.Agent) uint {
		return agent.ID
	}))

	text := strings.Join(textPieces, " ")
	params := openai.BetaThreadMessageNewParams{
		Content: openai.F([]openai.MessageContentPartParamUnion{
			openai.TextContentBlockParam{
				Text: openai.F(text),
				Type: openai.F(openai.TextContentBlockParamTypeText),
			},
		}),
		Role:     openai.F(openai.BetaThreadMessageNewParamsRoleUser),
		Metadata: openai.F(messageData.ToParam()),
	}

	if _, err := s.openai.Beta.Threads.Messages.New(ctx, req.ThreadId, params); err != nil {
		return nil, errors.Wrapf(err, "failed to add message")
	}

	for _, agent := range mentionedAgents {
		run, err := s.openai.Beta.Threads.Runs.New(ctx, req.ThreadId, openai.BetaThreadRunNewParams{
			AssistantID: openai.F(agent.AssistantId),
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to run thread")
		}
		for {
			run, err = s.openai.Beta.Threads.Runs.PollStatus(ctx, req.ThreadId, run.ID, 1000)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get run")
			}

			threadData.SetCurrentRunId(run.ID)
			if err := threadData.Save(ctx); err != nil {
				return nil, err
			}

			switch run.Status {
			case openai.RunStatusCompleted:
				logger.Debug("run completed", "run", run)
				return &emptypb.Empty{}, nil
			case openai.RunStatusFailed:
				return nil, errors.Wrapf(haberrors.ErrRuntime, "run failed: %s", run.LastError.Message)
			case openai.RunStatusIncomplete:
				return nil, errors.Wrapf(haberrors.ErrRuntime, "Run incomplete: %s", run.IncompleteDetails.Reason)
			case openai.RunStatusExpired:
				return nil, errors.Wrapf(haberrors.ErrRuntime, "Run expired. expires_at: %s", time.Unix(run.ExpiresAt, 0))
			case openai.RunStatusCancelled:
				return nil, errors.Wrapf(haberrors.ErrRuntime, "Run cancelled")
			case openai.RunStatusRequiresAction:
				return nil, errors.Wrapf(haberrors.ErrRuntime, "Run requires action")
			default:
				return nil, errors.Wrapf(haberrors.ErrBadRequest, "invalid thread run status: received %s", run.Status)
			}
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *server) getAllMessages(ctx context.Context, threadId string) ([]*Message, error) {
	tx := helpers.GetTx(ctx)

	var messages []*Message
	page, err := s.openai.Beta.Threads.Messages.List(ctx, threadId, openai.BetaThreadMessageListParams{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list messages")
	}

	for {
		for _, data := range page.Data {
			if len(data.Content) == 0 {
				continue
			}
			messageData := NewMessageData(s.openai, threadId, data.ID, data.Metadata)

			text := data.Content[0].Text.Value
			msg := &Message{
				Text: text,
				Id:   data.ID,
			}
			switch data.Role {
			case openai.MessageRoleAssistant:
				{
					msg.Role = Message_ASSISTANT
					run, err := s.getRun(ctx, threadId, data.RunID)
					if err != nil {
						return nil, err
					}

					var agent domain.Agent
					if err := tx.First(&agent, "assistant_id = ?", run.AssistantID).Error; err != nil {
						return nil, errors.Wrapf(err, "failed to find agent with assistant id %s", run.AssistantID)
					}

					msg.Agent = newAgentPbFromDb(&agent)
				}
			case openai.MessageRoleUser:
				{
					msg.Role = Message_USER

					var agents []domain.Agent
					if err := tx.Find(&agents, "agent_id IN (?)", messageData.GetAgentIds()).Error; err != nil {
						return nil, errors.Wrapf(err, "failed to find agents with ids %v", messageData.GetAgentIds())
					}

					mentions := gog.Map(agents, func(a domain.Agent) string {
						return a.Name
					})
					msg.Mentions = mentions
				}
			default:
				return nil, errors.Errorf("unknown message role: %s", data.Role)
			}
			messages = append(messages, msg)
		}

		if !page.HasMore {
			break
		}

		page, err = page.GetNextPage()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get next page")
		}
	}

	return messages, nil
}

package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/mokiat/gog"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (s *server) AddMessage(ctx context.Context, req *AddMessageRequest) (*emptypb.Empty, error) {
	tx := helpers.GetTx(ctx)
	_, thread, err := s.getThread(ctx, req.ThreadId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get thread with id %d", req.ThreadId)
	}

	if req.Message == nil || *req.Message == "" {
		return nil, errors.New("message is required")
	}

	message := strings.TrimSpace(*req.Message)

	isMention := strings.HasPrefix(message, "@")
	if !isMention {
		return nil, errors.New("you must mention at least an agent")
	}

	messagePieces := strings.Split(message, " ")
	if len(messagePieces) == 0 {
		return nil, errors.New("you must mention at least an agent")
	}

	mentionedAgents := make([]domain.Agent, 0, len(messagePieces))
	for _, messagePiece := range messagePieces {
		if !strings.HasPrefix(messagePiece, "@") {
			continue
		}

		agentName := strings.ToLower(strings.TrimPrefix(messagePiece, "@"))
		var agent domain.Agent
		if err := tx.First(&agent, "name = ?", agentName).Error; err != nil {
			return nil, errors.Wrapf(err, "failed to find agent with name %s", agentName)
		}
		mentionedAgents = append(mentionedAgents, agent)
	}

	messageData := NewEmptyMessageData(s.openai, thread.OpenaiThreadId)
	messageData.SetAgentIds(gog.Map(mentionedAgents, func(agent domain.Agent) uint {
		return agent.ID
	}))

	params := openai.BetaThreadMessageNewParams{
		Content: openai.F([]openai.MessageContentPartParamUnion{
			openai.TextContentBlockParam{
				Text: openai.F(message),
				Type: openai.F(openai.TextContentBlockParamTypeText),
			},
		}),
		Role:     openai.F(openai.BetaThreadMessageNewParamsRoleUser),
		Metadata: openai.F(messageData.ToParam()),
	}

	if _, err := s.openai.Beta.Threads.Messages.New(ctx, thread.OpenaiThreadId, params); err != nil {
		return nil, errors.Wrapf(err, "failed to add message")
	}

	if err := s.run(ctx, thread.ID); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) convMessage(ctx context.Context, threadId string, data *openai.Message) (*Message, error) {
	tx := helpers.GetTx(ctx)

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
			if err := tx.Find(&agents, "id IN ?", messageData.GetAgentIds()).Error; err != nil {
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

	return msg, nil
}

func (s *server) getAllMessages(ctx context.Context, threadId string) ([]*Message, error) {
	var messages []*Message
	page, err := s.openai.Beta.Threads.Messages.List(ctx, threadId, openai.BetaThreadMessageListParams{
		Order: openai.F(openai.BetaThreadMessageListParamsOrderAsc),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list messages")
	}

	for {
		for _, data := range page.Data {
			if len(data.Content) == 0 {
				continue
			}
			msg, err := s.convMessage(ctx, threadId, &data)
			if err != nil {
				return nil, err
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

func (s *server) getLastMessage(ctx context.Context, threadId string) (*Message, error) {
	page, err := s.openai.Beta.Threads.Messages.List(ctx, threadId, openai.BetaThreadMessageListParams{
		Order: openai.F(openai.BetaThreadMessageListParamsOrderDesc),
		Limit: openai.F(int64(1)),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list messages")
	}

	if len(page.Data) == 0 || len(page.Data[0].Content) == 0 {
		return nil, nil
	}

	msg, err := s.convMessage(ctx, threadId, &page.Data[0])
	if err != nil {
		return nil, err
	}

	return msg, nil
}

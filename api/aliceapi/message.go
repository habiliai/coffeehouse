package aliceapi

import (
	"context"
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/agentruntime/thread"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"strings"
)

func (s *server) AddMessage(ctx context.Context, req *AddMessageRequest) (*emptypb.Empty, error) {
	thd, err := s.getThread(ctx, req.ThreadId)
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

	mentionedAgentIds := make([]uint32, 0, len(messagePieces))
	for _, messagePiece := range messagePieces {
		if !strings.HasPrefix(messagePiece, "@") {
			continue
		}

		agentName := strings.ToLower(strings.TrimPrefix(messagePiece, "@"))
		ag, err := s.agentManager.GetAgentByName(ctx, &agent.GetAgentByNameRequest{
			Name: agentName,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find agent with name %s", agentName)
		}
		mentionedAgentIds = append(mentionedAgentIds, ag.Id)
	}

	if _, err := s.threadManager.AddMessage(ctx, &thread.AddMessageRequest{
		ThreadId: thd.AgentRuntimeThreadId,
		Message:  message,
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to add message")
	}

	if err := s.run(ctx, thd.ID, mentionedAgentIds); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) getAllMessages(ctx context.Context, threadId uint32) (res []*thread.Message, err error) {
	client, err := s.threadManager.GetMessages(ctx, &thread.GetMessagesRequest{
		ThreadId: threadId,
		Limit:    100,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get messages")
	}

	for {
		resp, err := client.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrapf(err, "failed to receive message")
		}
		res = append(res, resp.Messages...)
	}

	return
}

func (s *server) convertAllMessagesFromAgentRuntime(ctx context.Context, messages []*thread.Message) (res []*Message, err error) {
	memo := map[string]*agent.Agent{}
	for _, msg := range messages {
		row := Message{
			Id:   msg.Id,
			Text: msg.Content,
			User: msg.Sender,
		}
		if msg.Sender != "USER" {
			ag, ok := memo[msg.Sender]
			if !ok {
				ag, err = s.agentManager.GetAgentByName(ctx, &agent.GetAgentByNameRequest{
					Name: msg.Sender,
				})
				if err != nil {
					return nil, errors.Wrapf(err, "failed to find agent with name %s", msg.Sender)
				}
				memo[msg.Sender] = ag
			}
			row.Agent = newAgentPbFromAgentRuntime(ag)
		}
		res = append(res, &row)
	}

	return
}

func (s *server) getLastMessage(ctx context.Context, threadId uint32) (*Message, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := s.threadManager.GetMessages(ctx, &thread.GetMessagesRequest{
		ThreadId: threadId,
		Order:    thread.GetMessagesRequest_LATEST,
		Limit:    1,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list messages")
	}

	resp, err := stream.Recv()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to receive messages")
	}

	if len(resp.Messages) == 0 {
		return nil, nil
	}

	return &Message{
		Id:   resp.Messages[0].Id,
		Text: resp.Messages[0].Content,
		User: resp.Messages[0].Sender,
	}, nil
}

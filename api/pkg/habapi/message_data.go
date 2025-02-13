package habapi

import (
	"context"
	"github.com/mokiat/gog"
	"github.com/openai/openai-go"
)

type MessageData struct {
	data      openai.Metadata
	client    *openai.Client
	messageId string
	threadId  string
}

func NewEmptyMessageData(client *openai.Client, threadId string) *MessageData {
	return &MessageData{
		data:      openai.Metadata{},
		client:    client,
		messageId: "",
		threadId:  threadId,
	}
}

func NewMessageData(client *openai.Client, threadId string, messageId string, data openai.Metadata) *MessageData {
	return &MessageData{
		data:      data,
		client:    client,
		messageId: messageId,
		threadId:  threadId,
	}
}

func (m *MessageData) GetAgentIds() []uint {
	values := getIntSliceData(m.data, "agent_ids")

	return gog.Map(values, func(v int) uint {
		return uint(v)
	})
}

func (m *MessageData) SetAgentIds(agentIds []uint) {
	values := gog.Map(agentIds, func(v uint) int {
		return int(v)
	})

	setIntSliceData(m.data, "agent_ids", values)
}

func (m *MessageData) Save(ctx context.Context) error {
	_, err := m.client.Beta.Threads.Messages.Update(ctx, m.threadId, m.messageId, openai.BetaThreadMessageUpdateParams{
		Metadata: openai.F(openai.MetadataParam(m.data)),
	})
	return err
}

func (m *MessageData) ToParam() openai.MetadataParam {
	return openai.MetadataParam(m.data)
}

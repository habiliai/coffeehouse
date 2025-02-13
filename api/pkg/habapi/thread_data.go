package habapi

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
)

type ThreadData struct {
	data     openai.Metadata
	threadId string
	openai   *openai.Client
}

func NewThreadData(client *openai.Client, threadId string, data openai.Metadata) *ThreadData {
	return &ThreadData{
		data:     data,
		threadId: threadId,
		openai:   client,
	}
}

func (m *ThreadData) GetMissionId() uint {
	return uint(getIntData(m.data, "mission_id"))
}

func (m *ThreadData) SetMissionId(missionId uint) {
	setIntData(m.data, "mission_id", int(missionId))
}

func (m *ThreadData) SetCurrentRunId(id string) {
	setStringData(m.data, "current_run_id", id)
}

func (m *ThreadData) GetCurrentRunId() string {
	return getStringData(m.data, "current_run_id")
}

func (m *ThreadData) Save(ctx context.Context) error {
	_, err := m.openai.Beta.Threads.Update(ctx, m.threadId, openai.BetaThreadUpdateParams{
		Metadata: openai.F(openai.MetadataParam(m.data)),
	})
	return errors.Wrapf(err, "failed to update thread")
}

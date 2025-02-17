package habapi_test

import (
	domaintest "github.com/habiliai/habiliai/api/pkg/domain/testing"
	"github.com/habiliai/habiliai/api/pkg/habapi"
	"os"
)

func (s *HabApiTestSuite) TestMissionStepStatus() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		s.T().Skip("OPENAI_API_KEY not set")
	}

	s.Require().NoError(domaintest.SeedForTest(s.db))

	threadId, err := s.client.CreateThread(s, &habapi.CreateThreadRequest{
		MissionId: 1,
	})
	s.Require().NoError(err)

	missionStepStatus, err := s.client.GetMissionStepStatus(
		s,
		&habapi.GetMissionStepStatusRequest{
			StepSeqNo: 1,
			ThreadId:  threadId.Id,
		},
	)
	s.Require().NoError(err)

	s.T().Logf("MissionStepStatus: %v", missionStepStatus)
}

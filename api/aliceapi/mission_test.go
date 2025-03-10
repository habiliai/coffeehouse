package aliceapi_test

import (
	"github.com/habiliai/alice/api/aliceapi"
	"os"
)

func (s *AliceApiTestSuite) TestMissionStepStatus() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		s.T().Skip("OPENAI_API_KEY not set")
	}

	threadId, err := s.client.CreateThread(s, &aliceapi.CreateThreadRequest{
		MissionId: 1,
	})
	s.Require().NoError(err)

	missionStepStatus, err := s.client.GetMissionStepStatus(
		s,
		&aliceapi.GetMissionStepStatusRequest{
			StepSeqNo: 1,
			ThreadId:  threadId.Id,
		},
	)
	s.Require().NoError(err)

	s.T().Logf("MissionStepStatus: %v", missionStepStatus)
}

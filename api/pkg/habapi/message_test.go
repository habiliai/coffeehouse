package habapi_test

import (
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/habapi"
	"os"
)

func (s *HabApiTestSuite) TestAddMessage() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		s.T().Skip("OPENAI_API_KEY is not set")
	}

	seed, err := domain.SeedForTest(s.db)
	s.Require().NoError(err)

	threadId, err := s.client.CreateThread(
		s,
		&habapi.CreateThreadRequest{
			MissionId: int32(seed.Missions[0].ID),
		},
	)
	s.Require().NoError(err)
	thread, err := s.client.GetThread(s, threadId)
	s.Require().NoError(err)

	{
		_, err := s.client.AddMessage(
			s,
			&habapi.AddMessageRequest{
				ThreadId: thread.Id,
				Message:  "@engineer hello.",
			},
		)
		s.Require().NoError(err)

		thread, err = s.client.GetThread(s, threadId)
		s.Require().NoError(err)

		s.T().Logf("messages: %+v", thread.Messages)
	}
}

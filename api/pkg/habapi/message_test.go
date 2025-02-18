package habapi_test

import (
	domaintest "github.com/habiliai/habiliai/api/pkg/domain/testing"
	"github.com/habiliai/habiliai/api/pkg/habapi"
	"github.com/mokiat/gog"
	"os"
)

func (s *HabApiTestSuite) TestAddMessage() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		s.T().Skip("OPENAI_API_KEY is not set")
	}

	err := domaintest.SeedForTest(s.db)
	s.Require().NoError(err)

	threadId, err := s.client.CreateThread(
		s,
		&habapi.CreateThreadRequest{
			MissionId: int32(1),
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
				Message:  gog.PtrOf("@engineer hello."),
			},
		)
		s.Require().NoError(err)

		thread, err = s.client.GetThread(s, threadId)
		s.Require().NoError(err)

		s.T().Logf("messages: %+v", thread.Messages)
	}
}

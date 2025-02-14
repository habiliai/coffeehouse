package habapi_test

import (
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/habapi"
	"os"
)

func (s *HabApiTestSuite) TestGivenOneThreadWhenGetThreadShouldBeOK() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		s.T().Skip("OPENAI_API_KEY is not set")
	}

	seed, err := domain.SeedForTest(s.db)
	s.Require().NoError(err)

	// Create thread by first seed's mission
	threadId, err := s.client.CreateThread(s.Context, &habapi.CreateThreadRequest{
		MissionId: int32(seed.Missions[0].ID),
	})
	s.Require().NoError(err)

	// Get thread
	thread, err := s.client.GetThread(s.Context, threadId)
	s.Require().NoError(err)

	// Check response
	s.T().Logf("Thread: %+v", thread)
	s.Require().Equal(threadId.Id, thread.Id)
}

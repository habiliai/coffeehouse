package habapi_test

import (
	"github.com/habiliai/alice/api/pkg/domain"
	domaintest "github.com/habiliai/alice/api/pkg/domain/testing"
	"github.com/habiliai/alice/api/pkg/habapi"
	"os"
)

func (s *HabApiTestSuite) TestGivenOneThreadWhenGetThreadShouldBeOK() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		s.T().Skip("OPENAI_API_KEY is not set")
	}

	err := domaintest.SeedForTest(s.db)
	s.Require().NoError(err)

	mission := domain.Mission{}
	s.Require().NoError(s.db.First(&mission, "name = ?", "Mission 1").Error)

	// Create thread by first seed's mission
	threadId, err := s.client.CreateThread(s.Context, &habapi.CreateThreadRequest{
		MissionId: int32(1),
	})
	s.Require().NoError(err)

	// Get thread
	thread, err := s.client.GetThread(s.Context, threadId)
	s.Require().NoError(err)

	// Check response
	s.T().Logf("Thread: %+v", thread)
	s.Require().Equal(threadId.Id, thread.Id)
	s.Require().Equal(int32(mission.ID), thread.MissionId)
}

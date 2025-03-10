package aliceapi_test

import (
	"github.com/habiliai/alice/api/aliceapi"
	"github.com/habiliai/alice/api/domain"
	"github.com/habiliai/alice/api/domain/seed"
	"time"
)

func (s *AliceApiTestSuite) TestGivenOneThreadWhenGetThreadShouldBeOK() {
	s.Require().NoError(seed.Seed(s, false))

	mission := domain.Mission{}
	s.Require().NoError(s.db.First(&mission).Error)

	// Create thread by first seed's mission
	threadId, err := s.client.CreateThread(s.Context, &aliceapi.CreateThreadRequest{
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

	time.Sleep(1 * time.Second)
}

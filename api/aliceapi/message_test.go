package aliceapi_test

import (
	"github.com/habiliai/alice/api/aliceapi"
	"github.com/mokiat/gog"
)

func (s *AliceApiTestSuite) TestAddMessage() {
	threadId, err := s.client.CreateThread(
		s,
		&aliceapi.CreateThreadRequest{
			MissionId: int32(1),
		},
	)
	s.Require().NoError(err)
	thread, err := s.client.GetThread(s, threadId)
	s.Require().NoError(err)

	{
		_, err := s.client.AddMessage(
			s,
			&aliceapi.AddMessageRequest{
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

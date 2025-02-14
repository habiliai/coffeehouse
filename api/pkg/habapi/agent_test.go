package habapi_test

import (
	"github.com/habiliai/habiliai/api/pkg/domain"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *HabApiTestSuite) TestWhenGetAgentsShouldBeOKWithOneAgent() {
	seed, err := domain.SeedForTest(s.db)
	s.Require().NoError(err)

	// Get agents
	resp, err := s.client.GetAgents(s.Context, &emptypb.Empty{})
	s.Require().NoError(err)

	// Check response
	s.Require().Len(resp.Agents, 3)
	s.Require().Equal(seed.Agents[0].Name, resp.Agents[0].Name)
}

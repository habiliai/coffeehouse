package habapi_test

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/digo"
	"github.com/habiliai/habiliai/api/pkg/domain"
	habgrpc "github.com/habiliai/habiliai/api/pkg/grpc"
	"github.com/habiliai/habiliai/api/pkg/services"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
	"testing"
)

type HabApiTestSuite struct {
	suite.Suite
	context.Context

	db     *gorm.DB
	server *grpc.Server
	eg     errgroup.Group
}

func (s *HabApiTestSuite) SetupTest() {
	s.Context = context.TODO()

	container := digo.NewContainer(s, digo.EnvTest, nil)
	s.db = digo.MustGet[*gorm.DB](container, services.ServiceKeyDB)
	s.server = digo.MustGet[*grpc.Server](container, habgrpc.ServerKey)

	s.eg.Go(func() error {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			return errors.WithStack(err)
		}

		return errors.WithStack(s.server.Serve(lis))
	})

	s.Require().NoError(domain.SeedForTest(s.db))
	s.Require().NoError(habgrpc.WaitForServing(s, "localhost:50051"))
}

func (s *HabApiTestSuite) TearDownTest() {
	s.server.Stop()
	s.eg.Wait()
}

func TestHabApi(t *testing.T) {
	suite.Run(t, new(HabApiTestSuite))
}

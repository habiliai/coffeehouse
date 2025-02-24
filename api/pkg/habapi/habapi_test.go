package habapi_test

import (
	"context"
	"github.com/habiliai/alice/api/pkg/digo"
	habgrpc "github.com/habiliai/alice/api/pkg/grpc"
	"github.com/habiliai/alice/api/pkg/habapi"
	"github.com/habiliai/alice/api/pkg/services"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	conn   *grpc.ClientConn
	client habapi.HabiliApiClient
}

func (s *HabApiTestSuite) SetupTest() {
	var err error
	s.Context = context.TODO()

	container := digo.NewContainer(s, digo.EnvTest, nil)
	s.db = digo.MustGet[*gorm.DB](container, services.ServiceKeyDB)
	s.server = digo.MustGet[*grpc.Server](container, habgrpc.ServerKey)

	s.eg.Go(func() error {
		lis, err := net.Listen("tcp", "0.0.0.0:50051")
		if err != nil {
			return errors.WithStack(err)
		}

		return errors.WithStack(s.server.Serve(lis))
	})

	s.Require().NoError(habgrpc.WaitForServing(s, "127.0.0.1:50051"))

	s.conn, err = grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.client = habapi.NewHabiliApiClient(s.conn)
}

func (s *HabApiTestSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.server != nil {
		s.server.Stop()
	}
	s.eg.Wait()
}

func TestHabApi(t *testing.T) {
	suite.Run(t, new(HabApiTestSuite))
}

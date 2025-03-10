package aliceapi_test

import (
	"context"
	"github.com/habiliai/alice/api/aliceapi"
	"github.com/habiliai/alice/api/internal/db"
	"github.com/habiliai/alice/api/internal/di"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type AliceApiTestSuite struct {
	suite.Suite
	context.Context

	db *gorm.DB

	client aliceapi.AliceApiServer
}

func (s *AliceApiTestSuite) SetupTest() {
	s.Context = di.WithContainer(context.TODO(), di.EnvTest)

	s.db = di.MustGet[*gorm.DB](s, db.Key)
	s.client = di.MustGet[aliceapi.AliceApiServer](s, aliceapi.ServerKey)
}

func TestAliceApi(t *testing.T) {
	suite.Run(t, new(AliceApiTestSuite))
}

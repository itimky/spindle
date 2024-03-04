package test_test

import (
	"testing"

	"github.com/itimky/spindle/pkg/sys/log"
	"github.com/itimky/spindle/test"
	"github.com/stretchr/testify/suite"
)

type FactorySuite struct {
	suite.Suite
}

func (s *FactorySuite) Test_NewContext() {
	ctx := test.NewContext(s.T())

	s.NotNil(ctx)

	_, err := log.FromContext(ctx)
	s.NoError(err)
}

func TestFactorySuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(FactorySuite))
}

package logadapter_test

import (
	"testing"

	logadapter "github.com/itimky/spindle/pkg/sys/log/adapter"
	"github.com/stretchr/testify/suite"
)

type NopSuite struct {
	suite.Suite

	nop logadapter.Nop
}

func (s *NopSuite) SetupTest() {
	s.nop = logadapter.NewNop()
}

func (s *NopSuite) Test_Info() {
	s.nop.Info("smoke test")
}

func (s *NopSuite) Test_Infof() {
	s.nop.Infof("smoke test")
}

func (s *NopSuite) Test_Error() {
	s.nop.Error("smoke test")
}

func (s *NopSuite) Test_Errorf() {
	s.nop.Errorf("smoke test")
}

func (s *NopSuite) Test_WithField() {
	_ = s.nop.WithField("key", "value")
}

func TestNopSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(NopSuite))
}

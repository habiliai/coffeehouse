package hablog_test

import (
	aflog "github.com/habiliai/alice/api/pkg/log"
	"github.com/pkg/errors"
	"testing"
)

func TestLogger(t *testing.T) {
	var logger = aflog.GetLogger()
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error", aflog.Err(errors.New("test 123")))
}

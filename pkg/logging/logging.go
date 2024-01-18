package logging

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var Logger logr.Logger
var Level = zap.NewAtomicLevel()

func init() {
	opts := zap.NewDevelopmentConfig()
	opts.Level = Level
	opts.OutputPaths = []string{"stdout"}

	zapLogger := zap.Must(opts.Build())
	Logger = zapr.NewLogger(zapLogger)
}

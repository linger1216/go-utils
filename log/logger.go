package log

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	l      *zap.SugaredLogger
)

const (
	ProjectKey   = "project"
	ProjectValue = "schedule"
	_Mod         = "module"
)

func init() {
	logger, _ = zap.NewDevelopment()
	//logger, _ = zap.NewProduction()
	l = logger.Sugar().With(ProjectKey, ProjectValue)
}

// func (s *SugaredLogger) Debugf(template string, args ...interface{}) {
func _MOD(name string) *zap.SugaredLogger {
	return l.With(_Mod, name)
}

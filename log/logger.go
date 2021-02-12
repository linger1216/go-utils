package log

import (
	"go.uber.org/zap"
)

type Log struct {
	logger *zap.Logger
	l      *zap.SugaredLogger
}

func NewLog() *Log {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	l := logger.Sugar()
	return &Log{logger: logger, l: l}
}

func (l *Log) Debugf(template string, args ...interface{}) {
	l.l.Debugf(template, args...)
}

func (l *Log) Infof(template string, args ...interface{}) {
	l.l.Infof(template, args...)
}

func (l *Log) Warnf(template string, args ...interface{}) {
	l.l.Warnf(template, args...)
}

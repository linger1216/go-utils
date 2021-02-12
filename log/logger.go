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

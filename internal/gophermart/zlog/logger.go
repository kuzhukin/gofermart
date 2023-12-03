package zlog

import (
	"fmt"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("new zap development, err=%w", err))
	}

	Logger = logger.Sugar()
	_ = logger.Level().Enabled(zap.ErrorLevel)
}

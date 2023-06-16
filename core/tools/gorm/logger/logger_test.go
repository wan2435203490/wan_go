package logger

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm/logger"
	logCore "wan_go/core/logger"
)

func TestNew(t *testing.T) {
	l := New(logger.Config{
		SlowThreshold: time.Second,
		Colorful:      true,
		LogLevel: logger.LogLevel(
			logCore.DefaultLogger.Options().Level.LevelForGorm()),
	})
	l.Info(context.TODO(), "test")
}

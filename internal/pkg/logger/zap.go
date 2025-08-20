package logger

import (
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
)

var _ log.Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	log *zap.Logger
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{log: logger}
}

func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn("Keyvalues must appear in pairs", zap.Any("keyvals", keyvals))
		return nil
	}

	var fields []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		fields = append(fields, zap.Any(keyvals[i].(string), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug("", fields...)
	case log.LevelInfo:
		l.log.Info("", fields...)
	case log.LevelWarn:
		l.log.Warn("", fields...)
	case log.LevelError:
		l.log.Error("", fields...)
	case log.LevelFatal:
		l.log.Fatal("", fields...)
	}
	return nil
}

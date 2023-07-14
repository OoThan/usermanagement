package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

const (
	Console = "console"
	File    = "file"
)

var (
	Level  = zap.DebugLevel
	Target = Console
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func init() {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/usermanagement.log",
		MaxSize:    1024,
		MaxBackups: 10,
		MaxAge:     7,
	})

	var writerSyncer zapcore.WriteSyncer
	if Target == Console {
		writerSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}

	if Target == File {
		writerSyncer = zapcore.NewMultiWriteSyncer(w)
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewEncoderConfig()),
		writerSyncer,
		Level,
	)

	Logger = zap.New(core, zap.AddCaller())
	Sugar = Logger.Sugar()
}

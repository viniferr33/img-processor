package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	l    *zap.Logger
	once sync.Once
	cfg  Config
)

type Config struct {
	Development bool
	Level       string
	Color       bool
}

func Init(config Config) error {
	var err error
	once.Do(func() {
		l, err = newLogger(config)
		cfg = config
	})

	return err
}

func newLogger(config Config) (*zap.Logger, error) {
	var encoderConfig zapcore.EncoderConfig

	if config.Development {
		// Dev env config
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		if config.Color {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	} else {
		// Prod config
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.TimeKey = "timestamp"
	}

	level := getZapLevel(config.Level)
	if config.Development {
		level = zapcore.DebugLevel
	}

	var encoder zapcore.Encoder
	if config.Development {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	var options []zap.Option
	if config.Development {
		options = append(options, zap.AddCaller())
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	}

	return zap.New(core, options...), nil
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func Get() *zap.Logger {
	if l == nil {
		// Set defaults to dev mode
		Init(Config{
			Development: true,
			Level:       "debug",
			Color:       true,
		})
	}
	return l
}

func Sync() error {
	if l != nil {
		return l.Sync()
	}

	return nil
}

func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return Get().With(fields...)
}

func FromError(err error) zap.Field {
	return zap.Error(err)
}

package logster

import (
	"context"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapAdapter struct {
	prefix string
	*zap.SugaredLogger
}

func (z *ZapAdapter) Printf(format string, args ...interface{}) {
	z.Infof(format, args...)
}

func (z *ZapAdapter) Print(args ...interface{}) {
	z.Info(args...)
}

func (z *ZapAdapter) Println(args ...interface{}) {
	z.Info(args...)
}

func (z *ZapAdapter) WithPrefix(prefix string) Logger {
	return &ZapAdapter{SugaredLogger: z.SugaredLogger, prefix: prefix}
}

func (z *ZapAdapter) WithField(key string, value interface{}) Logger {
	return z.with(z.prefix+key, value)
}

func (z *ZapAdapter) WithContext(ctx context.Context) Logger {
	if f, ok := FieldsFromContext(ctx); ok {
		return z.WithFields(f)
	}
	return z
}

func (z *ZapAdapter) WithFields(fields Fields) Logger {
	s := make([]interface{}, 0, len(fields))
	for k, v := range fields {
		s = append(s, z.prefix+k, v)
	}
	return z.with(s...)
}

func (z *ZapAdapter) WithError(err error) Logger {
	if err != nil {
		return z.with(zap.String("error", err.Error()))
	}
	// check it! it means u have a strange code in ur project
	return z.with(zap.String("error", "<nil>"))
}

func (z *ZapAdapter) Write(p []byte) (n int, err error) {
	z.Warnw(string(p))
	return len(p), nil
}

func (z *ZapAdapter) with(args ...interface{}) Logger {
	return &ZapAdapter{SugaredLogger: z.With(args...), prefix: z.prefix}
}

var textToZapLevelMap = map[string]zapcore.Level{
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
	"error": zapcore.ErrorLevel,
	"warn":  zapcore.WarnLevel,
	"info":  zapcore.InfoLevel,
	"debug": zapcore.DebugLevel,
}

var zapToProdLevelMap = map[zapcore.Level]string{
	zapcore.PanicLevel: "FATAL",
	zapcore.FatalLevel: "FATAL",
	zapcore.ErrorLevel: "ERROR",
	zapcore.WarnLevel:  "WARN",
	zapcore.InfoLevel:  "INFO",
	zapcore.DebugLevel: "DEBUG",
}

func New(w zapcore.WriteSyncer, cfg Config) *ZapAdapter {
	if cfg.Project == "" {
		panic("logster: project field should be nonempty")
	}

	env := cfg.Env
	if env == "" {
		env = "local"
	}

	level, ok := textToZapLevelMap[cfg.Level]
	if !ok {
		level = zapcore.InfoLevel
	}

	fields := []zap.Field{
		zap.String("go_env", env),
		zap.String("go_project", cfg.Project),
	}

	var encoderCfg zapcore.EncoderConfig
	var enc zapcore.Encoder

	switch cfg.Format {
	case "text":
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		enc = zapcore.NewConsoleEncoder(encoderCfg)
	case "prod":
		if cfg.System == "" {
			panic("logster: system field should not be empty")
		}

		if cfg.Inst == "" {
			hostname, err := os.Hostname()
			if err != nil {
				log.Printf("logster: unable to get hostname: %v; using project_name as instance_name\n", err)
				hostname = cfg.Project
			}

			cfg.Inst = hostname
		}

		fields = append(fields,
			zap.String("inst", cfg.Inst),
			zap.String("env", env),
			zap.String("system", cfg.System),
		)

		fallthrough
	default:
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.MessageKey = "message"
		encoderCfg.CallerKey = "caller"
		encoderCfg.StacktraceKey = "stacktrace"

		encoderCfg.LevelKey = "level"
		encoderCfg.TimeKey = "@timestamp"
		encoderCfg.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(zapToProdLevelMap[l])
		}
		encoderCfg.EncodeTime = func(time time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(time.UTC().Format("2006-01-02T15:04:05.999Z07:00"))
		}
		enc = NewJSONEncoder(encoderCfg)
	}

	options := []zap.Option{
		zap.Fields(fields...),
		zap.AddCaller(),
	}

	if !cfg.DisableStackTrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	core := zapcore.NewCore(enc, w, level)
	sugar := zap.New(core).WithOptions(options...).Sugar()

	return &ZapAdapter{SugaredLogger: sugar, prefix: UserPrefix}
}

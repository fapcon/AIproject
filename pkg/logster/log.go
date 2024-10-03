package logster

import "context"

const (
	UserPrefix = "u_"
	LibPrefix  = "l_"
)

type Config struct {
	Project           string `yaml:"project"`
	Format            string `yaml:"format"`
	Level             string `yaml:"level"`
	Env               string `yaml:"env"`
	DisableStackTrace bool   `yaml:"disableStackTrace"`

	System string `yaml:"system"`
	Inst   string `yaml:"inst"`
}

type Fields map[string]interface{}

type Logger interface {
	WithPrefix(string) Logger
	WithField(key string, value interface{}) Logger
	WithFields(Fields) Logger
	WithContext(ctx context.Context) Logger
	WithError(error) Logger

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Write(p []byte) (n int, err error) // http/server logs interface
	Printer
}

type Printer interface {
	Printf(format string, args ...interface{}) // sarama logger interface
	Print(args ...interface{})                 // sarama logger interface
	Println(args ...interface{})               // prometheus logger interface
}

type discard struct{}

func (t discard) Sync() error {
	return nil
}

func (t discard) Write(p []byte) (n int, err error) {
	return 0, nil
}

// Discard is good logger for test purposes
var Discard = New(discard{}, Config{Project: "discard"})

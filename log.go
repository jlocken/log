package log

import (
	"github.com/jlocken/log/seq"
	"github.com/pkg/errors"
)

type Logger struct {
	hooks []hook
}

var log *Logger
var applicationName string

// Returns pointer to a Logger type that can be build on
// e.g configure to use a file, console or seq server hooks
// to write log messages
func BuildLogger(appName string) *Logger {
	applicationName = appName
	log = &Logger{}
	return log
}

type hook interface {
	Error(err error, s interface{})
	Fatal(err error, s interface{})
	Info(msg string, s interface{})
	Warning(msg string, s interface{})
}

// UseFile adds a file hook where message will be logged
// It returns the logger instance
func (logger *Logger) UseFile(filePath string) *Logger {
	// to be implemented
	return logger
}

// UseSeq add a hook to logging events on Seq Server running
// on the provided url e.g localhost:5431
func (logger *Logger) UseSeq(url, apikey string) *Logger {
	seqHook := &seq.SeqHook{
		BaseUrl: url,
		ApiKey:  apikey,
	}
	logger.hooks = append(logger.hooks, seqHook)
	return logger
}

func (logger *Logger) UseConsole() *Logger {
	// to be implemented
	return logger
}

func Error(err error, s interface{}) {
	err = errors.Wrap(err, applicationName)
	for _, h := range log.hooks {
		h.Error(err, s)
	}
}

func Fatal(err error, s interface{}) {
	err = errors.Wrap(err, applicationName)
	for _, h := range log.hooks {
		h.Fatal(err, s)
	}
}

func Info(msg string, s interface{}) {
	msg = applicationName + ": " + msg
	for _, h := range log.hooks {
		h.Info(msg, s)
	}
}

func Warning(msg string, s interface{}) {
	msg = applicationName + ": " + msg
	for _, h := range log.hooks {
		h.Warning(msg, s)
	}
}

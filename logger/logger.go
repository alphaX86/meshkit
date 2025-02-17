package logger

import (
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/layer5io/meshkit/errors"
	"github.com/sirupsen/logrus"
	gormlogger "gorm.io/gorm/logger"
)

type Handler interface {
	Info(description ...interface{})
	Debug(description ...interface{})
	Warn(err error)
	Error(err error)

	// Kubernetes Controller compliant logger
	ControllerLogger() logr.Logger
	DatabaseLogger() gormlogger.Interface
}

type Logger struct {
	handler *logrus.Entry
}

//TerminalFormatter is exported
type TerminalFormatter struct{}

// Format defined the format of output for Logrus logs
// Format is exported
func (f *TerminalFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return append([]byte(entry.Message), '\n'), nil
}

func New(appname string, opts Options) (Handler, error) {

	log := logrus.New()

	switch opts.Format {
	case JsonLogFormat:
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	case SyslogLogFormat:
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		})
	case TerminalLogFormat:
		log.SetFormatter(new(TerminalFormatter))
	}

	// log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	if opts.Output != nil {
		log.SetOutput(opts.Output)
	}

	log.SetLevel(logrus.InfoLevel)
	if opts.DebugLevel {
		log.SetLevel(logrus.DebugLevel)
	}

	entry := log.WithFields(logrus.Fields{"app": appname})
	return &Logger{handler: entry}, nil
}

func (l *Logger) Error(err error) {
	l.handler.WithFields(logrus.Fields{
		"code":                  errors.GetCode(err),
		"severity":              errors.GetSeverity(err),
		"short-description":     errors.GetSDescription(err),
		"probable-cause":        errors.GetCause(err),
		"suggested-remediation": errors.GetRemedy(err),
	}).Log(logrus.ErrorLevel, err.Error())
}

func (l *Logger) Info(description ...interface{}) {
	l.handler.Log(logrus.InfoLevel,
		description...,
	)
}

func (l *Logger) Debug(description ...interface{}) {
	l.handler.Log(logrus.DebugLevel,
		description...,
	)
}

func (l *Logger) Warn(err error) {
	l.handler.WithFields(logrus.Fields{
		"code":                  errors.GetCode(err),
		"severity":              errors.GetSeverity(err),
		"short-description":     errors.GetSDescription(err),
		"probable-cause":        errors.GetCause(err),
		"suggested-remediation": errors.GetRemedy(err),
	}).Log(logrus.WarnLevel, err.Error())
}

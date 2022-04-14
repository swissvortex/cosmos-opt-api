package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/Colm3na/cosmos-opt-api/constants"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	SetLoggingLevel(loggingLevel string)
	FileContext() string
	DebugWithContext(context string, arguments ...interface{})
	InfoMessageWithContext(context string, arguments ...interface{})
	WarningWithContext(context string, arguments ...interface{})
	EntryWithContext(context string, arguments ...interface{})
	ExitWithContext(context string, arguments ...interface{})
	InternalErrorWithContext(context string, arguments ...interface{})
	BadRequestErrorWithContext(context string, arguments ...interface{})
	ErrorWithContext(context string, arguments ...interface{})
}

type Event struct {
	id      int
	message string
}

type logger struct {
	*logrus.Logger
}

func NewLogger() Logger {
	var logrusLogger = logrus.New()
	var logger = &logger{logrusLogger}

	logger.Formatter = &logrus.TextFormatter{}
	logger.SetLoggingLevel(constants.DefaultLoggingLevel)

	return logger
}

var (
	genericErrorMessage    = Event{-1, "Error args: %s"}
	logOnEntryMessage      = Event{0, "> Entry args: %s"}
	logOnExitMessage       = Event{1, "< Exit args: %s"}
	invalidArgMessage      = Event{3, "Invalid arg: %s"}
	invalidArgValueMessage = Event{4, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{5, "Missing arg: %s"}
	debugMessage           = Event{6, "Message: %s"}
	infoMessage            = Event{7, "Message: %s"}
	warningMessage         = Event{8, "Message: %s"}
	internalErrorMessage   = Event{500, "Internal error: %s"}
	badRequestErrorMessage = Event{400, "Bad request error: %s"}
)

func (l *logger) SetLoggingLevel(loggingLevel string) {
	switch loggingLevel {
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "trace":
		l.SetLevel(logrus.TraceLevel)
	case "fatal":
		l.SetLevel(logrus.FatalLevel)
	case "panic":
		l.SetLevel(logrus.PanicLevel)
	default:
		l.SetLevel(logrus.ErrorLevel)
	}
}

func (l *logger) FileContext() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	splittedFuncName := strings.Split(frame.Function, ".")
	splittedFileNamePath := strings.Split(frame.File, "/")

	funcName := splittedFuncName[len(splittedFuncName)-1]
	filePath := splittedFileNamePath[len(splittedFileNamePath)-2] + "/" + splittedFileNamePath[len(splittedFileNamePath)-1]

	return fmt.Sprintf("%s() -> %s:%d", funcName, filePath, frame.Line)
}

func (l *logger) DebugWithContext(context string, arguments ...interface{}) {
	l.WithField("func", context).Debugf(debugMessage.message, arguments)
}

func (l *logger) InfoMessageWithContext(context string, arguments ...interface{}) {
	l.WithField("func", context).Infof(infoMessage.message, arguments)
}

func (l *logger) WarningWithContext(context string, arguments ...interface{}) {
	l.WithField("func", context).Warnf(warningMessage.message, arguments)
}

func (l *logger) EntryWithContext(context string, arguments ...interface{}) {
	l.WithField("func", context).Debugf(logOnEntryMessage.message, arguments)
}

func (l *logger) ExitWithContext(context string, arguments ...interface{}) {
	l.WithField("func", context).Debugf(logOnExitMessage.message, arguments)
}

func (l *logger) InternalErrorWithContext(context string, arguments ...interface{}) {
	l.WithField("func", context).Errorf(internalErrorMessage.message, arguments)
}

func (l *logger) BadRequestErrorWithContext(context string, arguments ...interface{}) {
	l.WithField("func", context).Errorf(badRequestErrorMessage.message, arguments)
}

func (l *logger) ErrorWithContext(context string, arguments ...interface{}) {
	l.WithField("func", context).Errorf(genericErrorMessage.message, arguments)
}

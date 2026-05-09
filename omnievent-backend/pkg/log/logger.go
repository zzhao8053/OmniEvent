package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"omnievent-backend/pkg/context"
)

const logFieldRequestId = "REQUEST_ID"
const logFieldExtra = "EXTRA"

var bootLogger = logrus.New()
var cliLogger = logrus.New()
var defaultLogger = logrus.New()
var requestLogger = logrus.New()
var sqlQueryLogger = logrus.New()

func init() {
	bootLogger.SetFormatter(&LogFormatter{})
	bootLogger.SetOutput(os.Stdout)
	bootLogger.SetLevel(logrus.InfoLevel)

	cliLogger.SetFormatter(&LogFormatter{})
	cliLogger.SetOutput(os.Stdout)
	cliLogger.SetLevel(logrus.InfoLevel)

	defaultLogger.SetFormatter(&LogFormatter{})
	defaultLogger.SetOutput(os.Stdout)
	defaultLogger.SetLevel(logrus.InfoLevel)

	requestLogger.SetFormatter(&LogFormatter{Prefix: "[REQUEST]", DisableLevel: true})
	requestLogger.SetOutput(os.Stdout)
	requestLogger.SetLevel(logrus.InfoLevel)

	sqlQueryLogger.SetFormatter(&LogFormatter{Prefix: "[SQLQUERY]", DisableLevel: true})
	sqlQueryLogger.SetOutput(os.Stdout)
	sqlQueryLogger.SetLevel(logrus.InfoLevel)
}

// SetLoggerConfiguration sets the logger according to the config
func SetLoggerConfiguration(enableConsoleLog, enableFileLog bool, logLevel string, enableRequestLog, enableQueryLog bool, fileLogPath, requestFileLogPath, queryFileLogPath string, logFileRotate bool, logFileMaxSize int64, logFileMaxDays uint32, isDisableBootLog bool) error {
	var bootWriters []io.Writer
	var defaultWriters []io.Writer
	var requestWriters []io.Writer
	var queryWriters []io.Writer

	if !isDisableBootLog {
		bootWriters = append(bootWriters, os.Stdout)
	}

	if enableConsoleLog {
		defaultWriters = append(defaultWriters, os.Stdout)
		requestWriters = append(requestWriters, os.Stdout)
		queryWriters = append(queryWriters, os.Stdout)
	}

	if enableFileLog {
		defaultWriter, err := NewRotateFileWriter(fileLogPath, logFileRotate, logFileMaxSize, logFileMaxDays)

		if err != nil {
			return err
		}

		if !isDisableBootLog {
			bootWriters = append(bootWriters, defaultWriter)
		}

		defaultWriters = append(defaultWriters, defaultWriter)

		if enableRequestLog {
			if requestFileLogPath != "" && requestFileLogPath != fileLogPath {
				requestWriter, err := NewRotateFileWriter(requestFileLogPath, logFileRotate, logFileMaxSize, logFileMaxDays)

				if err != nil {
					return err
				}

				requestWriters = append(requestWriters, requestWriter)
			} else {
				requestWriters = append(requestWriters, defaultWriter)
			}
		}

		if enableQueryLog {
			if queryFileLogPath != "" && queryFileLogPath != fileLogPath {
				queryWriter, err := NewRotateFileWriter(queryFileLogPath, logFileRotate, logFileMaxSize, logFileMaxDays)

				if err != nil {
					return err
				}

				queryWriters = append(queryWriters, queryWriter)
			} else {
				queryWriters = append(queryWriters, defaultWriter)
			}
		}
	}

	bootMultipleWriter := io.MultiWriter(bootWriters...)
	defaultMultipleWriter := io.MultiWriter(defaultWriters...)
	requestMultipleWriter := io.MultiWriter(requestWriters...)
	queryMultipleWriter := io.MultiWriter(queryWriters...)

	bootLogger.SetOutput(bootMultipleWriter)
	defaultLogger.SetOutput(defaultMultipleWriter)
	requestLogger.SetOutput(requestMultipleWriter)
	sqlQueryLogger.SetOutput(queryMultipleWriter)

	if logLevel == "DEBUG" {
		cliLogger.SetLevel(logrus.DebugLevel)
		bootLogger.SetLevel(logrus.DebugLevel)
		defaultLogger.SetLevel(logrus.DebugLevel)
	} else if logLevel == "INFO" {
		cliLogger.SetLevel(logrus.InfoLevel)
		bootLogger.SetLevel(logrus.InfoLevel)
		defaultLogger.SetLevel(logrus.InfoLevel)
	} else if logLevel == "WARN" {
		cliLogger.SetLevel(logrus.WarnLevel)
		bootLogger.SetLevel(logrus.WarnLevel)
		defaultLogger.SetLevel(logrus.WarnLevel)
	} else if logLevel == "ERROR" {
		cliLogger.SetLevel(logrus.ErrorLevel)
		bootLogger.SetLevel(logrus.ErrorLevel)
		defaultLogger.SetLevel(logrus.ErrorLevel)
	}

	if !enableRequestLog {
		requestLogger = nil
	}

	if !enableQueryLog {
		sqlQueryLogger = nil
	}

	return nil
}

// Debugf logs debug log with custom format
func Debugf(c *context.WebContext, format string, args ...any) {
	if c == nil {
		defaultLogger.Debug(getFinalLog(format, args...))
	} else {
		defaultLogger.WithField(logFieldRequestId, c.GetContextID()).Debug(getFinalLog(format, args...))
	}
}

// Infof logs info log with custom format
func Infof(c *context.WebContext, format string, args ...any) {
	if c == nil {
		defaultLogger.Info(getFinalLog(format, args...))
	} else {
		defaultLogger.WithField(logFieldRequestId, c.GetContextID()).Info(getFinalLog(format, args...))
	}
}

// Warnf logs warn log with custom format
func Warnf(c *context.WebContext, format string, args ...any) {
	if c == nil {
		defaultLogger.Warn(getFinalLog(format, args...))
	} else {
		defaultLogger.WithField(logFieldRequestId, c.GetContextID()).Warn(getFinalLog(format, args...))
	}
}

// Errorf logs error log with custom format
func Errorf(c *context.WebContext, format string, args ...any) {
	if c == nil {
		defaultLogger.Error(getFinalLog(format, args...))
	} else {
		defaultLogger.WithField(logFieldRequestId, c.GetContextID()).Error(getFinalLog(format, args...))
	}
}

// ErrorfWithExtra logs error log with custom format and extra info
func ErrorfWithExtra(c *context.WebContext, extraString string, format string, args ...any) {
	if c == nil {
		defaultLogger.WithField(logFieldExtra, extraString).Error(getFinalLog(format, args...))
	} else {
		defaultLogger.WithField(logFieldRequestId, c.GetContextID()).WithField(logFieldExtra, extraString).Error(getFinalLog(format, args...))
	}
}

// BootInfof logs boot info log
func BootInfof(c *context.WebContext, format string, args ...any) {
	if bootLogger != nil {
		if c == nil {
			bootLogger.Info(getFinalLog(format, args...))
		} else {
			bootLogger.WithField(logFieldRequestId, c.GetContextID()).Info(getFinalLog(format, args...))
		}
	}
}

// BootWarnf logs boot warn log
func BootWarnf(c *context.WebContext, format string, args ...any) {
	if bootLogger != nil {
		if c == nil {
			bootLogger.Warn(getFinalLog(format, args...))
		} else {
			bootLogger.WithField(logFieldRequestId, c.GetContextID()).Warn(getFinalLog(format, args...))
		}
	}
}

// BootErrorf logs boot error log
func BootErrorf(c *context.WebContext, format string, args ...any) {
	if bootLogger != nil {
		if c == nil {
			bootLogger.Error(getFinalLog(format, args...))
		} else {
			bootLogger.WithField(logFieldRequestId, c.GetContextID()).Error(getFinalLog(format, args...))
		}
	}
}

// CliInfof logs boot info log
func CliInfof(c *context.WebContext, format string, args ...any) {
	if cliLogger != nil {
		if c == nil {
			cliLogger.Info(getFinalLog(format, args...))
		} else {
			cliLogger.WithField(logFieldRequestId, c.GetContextID()).Info(getFinalLog(format, args...))
		}
	}
}

// CliWarnf logs boot warn log
func CliWarnf(c *context.WebContext, format string, args ...any) {
	if cliLogger != nil {
		if c == nil {
			cliLogger.Warn(getFinalLog(format, args...))
		} else {
			cliLogger.WithField(logFieldRequestId, c.GetContextID()).Warn(getFinalLog(format, args...))
		}
	}
}

// CliErrorf logs boot error log
func CliErrorf(c *context.WebContext, format string, args ...any) {
	if cliLogger != nil {
		if c == nil {
			cliLogger.Error(getFinalLog(format, args...))
		} else {
			cliLogger.WithField(logFieldRequestId, c.GetContextID()).Error(getFinalLog(format, args...))
		}
	}
}

// Requestf logs http request log with custom format
func Requestf(c *context.WebContext, format string, args ...any) {
	if requestLogger != nil {
		requestLogger.WithField(logFieldRequestId, c.GetContextID()).Info(getFinalLog(format, args...))
	}
}

// SqlQuery logs sql query log
func SqlQuery(args ...any) {
	if sqlQueryLogger != nil {
		sqlQueryLogger.Info(args...)
	}
}

// SqlQueryf logs sql query log with custom format
func SqlQueryf(format string, args ...any) {
	if sqlQueryLogger != nil {
		sqlQueryLogger.Info(getFinalLog(format, args...))
	}
}

func getFinalLog(format string, args ...any) string {
	result := fmt.Sprintf(format, args...)
	result = strings.Replace(result, "\n", " ", -1)

	return result
}
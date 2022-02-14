package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	defaultFormatter = &Formatter{
		LogFormat:       "%time% [%lvl%] %msg%",
		TimestampFormat: "2006-01-02 15:04:05",
	}

	logLevels = map[string]logrus.Level{
		"DEBUG": logrus.DebugLevel,
		"INFO":  logrus.InfoLevel,
		"WARN":  logrus.WarnLevel,
		"ERROR": logrus.ErrorLevel,
	}

	logger = getDefaultLogger()
)

func Setup(logLevel string, logFilePath string) {
	formatter := defaultFormatter
	level, ok := logLevels[strings.ToUpper(logLevel)]
	if !ok {
		level = logrus.InfoLevel
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetFormatter(formatter)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(level)

	// Output to file
	rotateFileHook, err := NewRotateFileHook(RotateFileConfig{
		Filename:   logFilePath,
		MaxSize:    50,
		MaxBackups: 7,
		MaxAge:     7,
		LocalTime:  true,
		Level:      level,
		Formatter:  formatter,
	})
	if err != nil {
		logger.Fatalf("Create log rotate hook error: %s\n", err)
		return
	}
	logger.AddHook(rotateFileHook)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func getDefaultLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(defaultFormatter)
	log.SetLevel(logrus.InfoLevel)
	return log
}

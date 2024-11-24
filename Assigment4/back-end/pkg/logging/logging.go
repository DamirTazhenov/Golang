package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"runtime"
)

type writeHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writeHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (hook *writeHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{e}
}

func (l *Logger) GetLoggerWithField(key string, value interface{}) *Logger {
	return &Logger{l.Entry.WithField(key, value)}
}

func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}
	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %s", err))
	}

	//l.SetOutput(io.Discard)

	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/app.log", // Path to log file
		MaxSize:    10,               // Max megabytes before rotation
		MaxBackups: 5,                // Max number of old log files to keep
		MaxAge:     30,               // Max number of days to retain old files
		Compress:   true,             // Compress rotated log files
	})

	logrus.SetFormatter(&logrus.JSONFormatter{})

	l.AddHook(&writeHook{
		Writer:    []io.Writer{allFile}, //, os.Stdout
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}

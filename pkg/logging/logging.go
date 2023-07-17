package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

type writerHook struct {
	Writers   []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		panic(err)
	}
	for _, writer := range hook.Writers {
		writer.Write([]byte(line))
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func (logger *Logger) GetLoggerWithFields(k string, v interface{}) Logger {
	return Logger{logger.WithField(k, v)}
}

func init() {
	logger := logrus.New()

	logger.SetReportCaller(true)

	logger.Formatter = &logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
	}

	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("logs/all.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		panic(err)
	}

	logger.SetOutput(io.Discard)

	logger.AddHook(&writerHook{
		Writers:   []io.Writer{file, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	logger.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(logger)
}

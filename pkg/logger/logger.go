package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type LevelHook struct {
	LevelsList []logrus.Level
	Writer     map[logrus.Level]io.Writer
}

func (hook *LevelHook) Fire(entry *logrus.Entry) error {
	writer, ok := hook.Writer[entry.Level]
	if ok {
		logString, err := entry.String()
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(logString + "\n"))
		return err
	}
	return nil
}

func (hook *LevelHook) Levels() []logrus.Level {
	return hook.LevelsList
}

type CustomFormatter struct {
	logrus.TextFormatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// timestamp := entry.Time.Format("2006-01-02 15:04:05")
	// levelText := fmt.Sprintf("%s:", timestamp)
	// caller := fmt.Sprintf("%s:%d %s", entry.Caller.File, entry.Caller.Line, entry.Caller.Function)
	// msg := fmt.Sprintf("%s %s [%s]\n", levelText, entry.Message, caller)
	// return []byte(msg), nil

	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")
	levelText := fmt.Sprintf("[%s]", timestamp)
	caller := fmt.Sprintf(" [func: %s %s:%d]", entry.Caller.Function, filepath.Base(entry.Caller.File), entry.Caller.Line)

	// Convert the entry.Data map to a string
	fieldsStr := ""
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for key, value := range entry.Data {
			fields = append(fields, fmt.Sprintf("%s=%v", key, value))
		}
		fieldsStr = strings.Join(fields, ", ")
		fieldsStr = fmt.Sprintf(" [Fields: {%s}]", fieldsStr)
	}

	msg := fmt.Sprintf("%s %s%s%s", levelText, entry.Message, fieldsStr, caller)
	return []byte(msg), nil
}

type customWriter struct {
	console io.Writer
	file    io.Writer
}

func (cw *customWriter) Write(b []byte) (int, error) {
	// cw.console.Write(append(b, '\n'))
	cw.console.Write(b)
	return cw.file.Write(b)
}

// with color
// func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
// 	f.DisableColors = true // disable colors
// 	timestamp := entry.Time.Format("2006-01-02 15:04:05")
// 	levelColor := f.getLevelColor(entry.Level)
// 	// levelText := fmt.Sprintf("%s[%s]", strings.ToUpper(entry.Level.String()), timestamp)
// 	levelText := fmt.Sprintf("[%s]", timestamp)
// 	caller := fmt.Sprintf("%s:%d %s", entry.Caller.File, entry.Caller.Line, entry.Caller.Function)
// 	msg := fmt.Sprintf("%s%s\x1b[0m %s [%s]\n", levelColor, levelText, entry.Message, caller)
// 	return []byte(msg), nil
// }

// func (f *CustomFormatter) getLevelColor(level logrus.Level) string {
// 	switch level {
// 	case logrus.DebugLevel, logrus.TraceLevel:
// 		return "\x1b[37m" // white
// 	case logrus.WarnLevel:
// 		return "\x1b[33m" // yellow
// 	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
// 		return "\x1b[31m" // red
// 	default:
// 		return "\x1b[36m" // cyan
// 	}
// }

func Setup(path string, logLevel logrus.Level, isHook bool) error {
	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)
	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	ForceColors:     true,
	// 	FullTimestamp:   true,
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// })
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetFormatter(&CustomFormatter{})

	logNameSuffix := "%Y%m%d"

	panicPath := fmt.Sprintf("%s%s%s.panic.log", path, string(filepath.Separator), logNameSuffix)
	panicLinkPath := fmt.Sprintf("%s%spanic.log", path, string(filepath.Separator))
	panicWriter, err := rotatelogs.New(
		panicPath,
		rotatelogs.WithLinkName(panicLinkPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}

	fatalPath := fmt.Sprintf("%s%s%s.fatal.log", path, string(filepath.Separator), logNameSuffix)
	fatalLinkPath := fmt.Sprintf("%s%sfatal.log", path, string(filepath.Separator))
	fatalWriter, err := rotatelogs.New(
		fatalPath,
		rotatelogs.WithLinkName(fatalLinkPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}

	errorPath := fmt.Sprintf("%s%s%s.error.log", path, string(filepath.Separator), logNameSuffix)
	errorLinkPath := fmt.Sprintf("%s%serror.log", path, string(filepath.Separator))
	errorWriter, err := rotatelogs.New(
		errorPath,
		rotatelogs.WithLinkName(errorLinkPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}

	warnPath := fmt.Sprintf("%s%s%s.warn.log", path, string(filepath.Separator), logNameSuffix)
	warnLinkPath := fmt.Sprintf("%s%swarn.log", path, string(filepath.Separator))
	warnWriter, err := rotatelogs.New(
		warnPath,
		rotatelogs.WithLinkName(warnLinkPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}

	infoPath := fmt.Sprintf("%s%s%s.info.log", path, string(filepath.Separator), logNameSuffix)
	infoLinkPath := fmt.Sprintf("%s%sinfo.log", path, string(filepath.Separator))
	infoWriter, err := rotatelogs.New(
		infoPath,
		rotatelogs.WithLinkName(infoLinkPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}

	debugPath := fmt.Sprintf("%s%s%s.debug.log", path, string(filepath.Separator), logNameSuffix)
	debugLinkPath := fmt.Sprintf("%s%sdebug.log", path, string(filepath.Separator))
	debugWriter, err := rotatelogs.New(
		debugPath,
		rotatelogs.WithLinkName(debugLinkPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}

	tracePath := fmt.Sprintf("%s%s%s.trace.log", path, string(filepath.Separator), logNameSuffix)
	traceLinkPath := fmt.Sprintf("%s%strace.log", path, string(filepath.Separator))
	traceWriter, err := rotatelogs.New(
		tracePath,
		rotatelogs.WithLinkName(traceLinkPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}

	logrus.AddHook(&LevelHook{
		LevelsList: logrus.AllLevels,
		Writer: map[logrus.Level]io.Writer{
			logrus.PanicLevel: &customWriter{console: os.Stdout, file: panicWriter},
			logrus.FatalLevel: &customWriter{console: os.Stdout, file: fatalWriter},
			logrus.ErrorLevel: &customWriter{console: os.Stdout, file: errorWriter},
			logrus.WarnLevel:  &customWriter{console: os.Stdout, file: warnWriter},
			logrus.InfoLevel:  &customWriter{console: os.Stdout, file: infoWriter},
			logrus.DebugLevel: &customWriter{console: os.Stdout, file: debugWriter},
			logrus.TraceLevel: &customWriter{console: os.Stdout, file: traceWriter},
		},
	})
	return nil
}

package logger

import (
	"fmt"
	"io"
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
	timestamp := entry.Time.Format("2006-01-02 15:04:05.999")
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

// with color
//
//	func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
//		f.DisableColors = true // disable colors
//		timestamp := entry.Time.Format("2006-01-02 15:04:05")
//		levelColor := f.getLevelColor(entry.Level)
//		// levelText := fmt.Sprintf("%s[%s]", strings.ToUpper(entry.Level.String()), timestamp)
//		levelText := fmt.Sprintf("[%s]", timestamp)
//		caller := fmt.Sprintf("%s:%d %s", entry.Caller.File, entry.Caller.Line, entry.Caller.Function)
//		msg := fmt.Sprintf("%s%s\x1b[0m %s [%s]\n", levelColor, levelText, entry.Message, caller)
//		return []byte(msg), nil
//	}
//
//	func (f *CustomFormatter) getLevelColor(level logrus.Level) string {
//		switch level {
//		case logrus.DebugLevel, logrus.TraceLevel:
//			return "\x1b[37m" // white
//		case logrus.WarnLevel:
//			return "\x1b[33m" // yellow
//		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
//			return "\x1b[31m" // red
//		default:
//			return "\x1b[36m" // cyan
//		}
//	}
func Setup(path string, logLevel logrus.Level, isHook bool) error {
	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)
	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	ForceColors:     true,
	// 	FullTimestamp:   true,
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// })
	logrus.SetFormatter(&CustomFormatter{})
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logNameSuffix := ".%Y%m%d"
	infoPath := fmt.Sprintf("%s%sinfo.log", path, string(filepath.Separator))
	infoWriter, err := rotatelogs.New(
		infoPath+logNameSuffix, // fmt.Sprintf("%s%s%s.info.log", path, "%Y%m%d", string(filepath.Separator)),
		rotatelogs.WithLinkName(infoPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}
	debugPath := fmt.Sprintf("%s%sdebug.log", path, string(filepath.Separator))
	debugWriter, err := rotatelogs.New(
		debugPath+logNameSuffix,
		rotatelogs.WithLinkName(debugPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}
	errorPath := fmt.Sprintf("%s%serror.log", path, string(filepath.Separator))
	errorWriter, err := rotatelogs.New(
		errorPath+logNameSuffix,
		rotatelogs.WithLinkName(errorPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}
	warnPath := fmt.Sprintf("%s%swarn.log", path, string(filepath.Separator))
	warnWriter, err := rotatelogs.New(
		warnPath+logNameSuffix,
		rotatelogs.WithLinkName(warnPath),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		return err
	}
	logrus.AddHook(&LevelHook{
		LevelsList: logrus.AllLevels,
		Writer: map[logrus.Level]io.Writer{
			logrus.InfoLevel:  infoWriter,
			logrus.DebugLevel: debugWriter,
			logrus.ErrorLevel: errorWriter,
			logrus.WarnLevel:  warnWriter,
		},
	})
	return nil
}

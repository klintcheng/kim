package logger

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	Filename    string
	Level       string
	RollingDays uint
	Format      string
}

func Init(settings Settings) error {

	if settings.Level == "" {
		settings.Level = "debug"
	}
	ll, err := logrus.ParseLevel(settings.Level)
	if err == nil {
		std.SetLevel(ll)
	} else {
		std.Error("Invalid log level")
	}

	if settings.Filename == "" {
		return nil
	}

	if settings.RollingDays == 0 {
		settings.RollingDays = 7
	}

	writer, err := rotatelogs.New(
		settings.Filename+".%Y%m%d",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(settings.Filename),

		// WithRotationTime设置日志分割的时间
		rotatelogs.WithRotationTime(time.Hour*24),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(settings.RollingDays),
	)
	if err != nil {
		return err
	}

	var logfr logrus.Formatter
	if settings.Format == "json" {
		logfr = &logrus.JSONFormatter{
			DisableTimestamp: false,
		}
	} else {
		logfr = &logrus.TextFormatter{
			DisableColors: true,
		}
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		// logrus.FatalLevel: writer,
		// logrus.PanicLevel: writer,
	}, logfr)

	std.AddHook(lfsHook)
	return nil
}

package log

import (
	"github.com/sirupsen/logrus"
)

const (
	logPath = "/var/log/go-kunpeng/go-kunpeng"
)

var Logger = logrus.New()

func init() {
	// 严禁调换hook次序, 否则会导致部分hook失效

	// 时间
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000"
	Logger.SetFormatter(customFormatter)

	// 行数
	Logger.AddHook(newLineHook())

	// 切割并保存30d的日志
	Logger.AddHook(newDividerHook(logrus.DebugLevel, 30))
}

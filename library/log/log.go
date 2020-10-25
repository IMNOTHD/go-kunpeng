package log

import (
	"io"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	_debugLogPath = "/var/log/go-kunpeng/debug.log"
	_infoLogPath  = "/var/log/go-kunpeng/info.log"
	_warnLogPath  = "/var/log/go-kunpeng/warn.log"
	_errorLogPath = "/var/log/go-kunpeng/error.log"
	_fatalLogPath = "/var/log/go-kunpeng/fatal.log"
	_panicLogPath = "/var/log/go-kunpeng/panic.log"
)

var Logger *zap.Logger

func init() {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		TimeKey:       "ts",
		StacktraceKey: "stacktrace",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现三个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
	// info level  -> debug, info
	// warn level  -> warn,
	// error level -> error, panic, fatal
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel
	})

	infoWriter := getWriter(_infoLogPath)
	warnWriter := getWriter(_warnLogPath)
	errorWriter := getWriter(_errorLogPath)

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)

	Logger = zap.New(core, zap.AddStacktrace(zapcore.WarnLevel), zap.AddCaller())
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 filename.YYmmdd
	// filename 是指向最新日志的链接
	// 保存30天内的日志，每24小时分割一次日志
	hook, err := rotatelogs.New(
		filename+".%Y%m%d", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

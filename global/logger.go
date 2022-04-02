package global

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

func SetLogger() {
	logger := &lumberjack.Logger{
		// 日志输出文件路径
		Filename: Config.Get("log.logPath").(string) + "/system.log",
		// 日志文件最大 size, 单位是 MB
		MaxSize: int(Config.Get("log.maxSize").(int64)), // megabytes
		// 最大过期日志保留的个数
		MaxBackups: int(Config.Get("log.maxBackups").(int64)),
		// 保留过期文件的最大时间间隔,单位是天
		MaxAge: int(Config.Get("log.maxAge").(int64)), //days
		// 是否需要压缩滚动日志, 使用的 gzip 压缩
		Compress: Config.Get("log.compress").(bool), // disabled by default
	}
	Logger = logrus.New()
	Logger.SetOutput(logger) //调用 logrus 的 SetOutput()函数
}

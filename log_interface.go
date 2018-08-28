package mylog

// 定义一个日志接口用于 让写文件的类和控制台的类进行实现
type LogInterface interface {
	Init()
	SetLevel(level int)
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Close()
}
/*
控制台日志的输出
*/
package mylog

import (
	"os"
	"fmt"
)


type ConsoleLogger struct {
	level int
}

func NewConsoleLogger(config map[string] string) (log LogInterface,err error) {
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not found log_level")
		return
	}
	level := getLogLevel(logLevel)
	log =  &ConsoleLogger{
		level:level,
	}
	return
}

func (c *ConsoleLogger) Init(){
	
}

// 设置日志级别
func (c *ConsoleLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		c.level = LogLevelDebug
	}
	c.level = level
}

func (c *ConsoleLogger) Debug(format string, args...interface{}) {
	if c.level > LogLevelDebug {
		return
	}
	LogData := WriteLog(LogLevelDebug, format, args...)
	fmt.Fprintf(os.Stdout,"%s [%s] %s:%d %s %s\n", LogData.TimeStr,LogData.LevelStr, LogData.FileName,LogData.LineNo,LogData.FuncName, LogData.Message)
}

func (c *ConsoleLogger) Trace(format string, args...interface{}) {
	if c.level > LogLevelTrace {
		return
	}
	LogData := WriteLog(LogLevelTrace, format, args...)
	fmt.Fprintf(os.Stdout,"%s [%s] %s:%d %s %s\n", LogData.TimeStr,LogData.LevelStr, LogData.FileName,LogData.LineNo,LogData.FuncName, LogData.Message)
}

func (c *ConsoleLogger) Warn(format string, args...interface{}) {
	if c.level > LogLevelWarn {
		return
	}
	LogData := WriteLog(LogLevelWarn, format, args...)
	fmt.Fprintf(os.Stdout,"%s [%s] %s:%d %s %s\n", LogData.TimeStr,LogData.LevelStr, LogData.FileName,LogData.LineNo,LogData.FuncName, LogData.Message)
}

func (c *ConsoleLogger) Error(format string, args...interface{}) {
	if c.level > LogLevelError {
		return
	}
	LogData := WriteLog(LogLevelError, format, args...)
	fmt.Fprintf(os.Stdout,"%s [%s] %s:%d %s %s\n", LogData.TimeStr,LogData.LevelStr, LogData.FileName,LogData.LineNo,LogData.FuncName, LogData.Message)
}

func (c *ConsoleLogger) Fatal(format string, args...interface{}) {
	if c.level > LogLevelFatal {
		return
	}
	LogData := WriteLog(LogLevelFatal, format, args...)
	fmt.Fprintf(os.Stdout,"%s [%s] %s:%d %s %s\n", LogData.TimeStr,LogData.LevelStr, LogData.FileName,LogData.LineNo,LogData.FuncName, LogData.Message)
}

func (c *ConsoleLogger) Info(format string, args...interface{}) {
	if c.level > LogLevelInfo {
		return
	}
	LogData := WriteLog(LogLevelInfo, format, args...)
	fmt.Fprintf(os.Stdout,"%s [%s] %s:%d %s %s\n", LogData.TimeStr,LogData.LevelStr, LogData.FileName,LogData.LineNo,LogData.FuncName, LogData.Message)
}

func (c *ConsoleLogger) Close(){

}
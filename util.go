package mylog

import (
	"path"
	"runtime"
	"time"
	"fmt"
)

// 定义日志的消息格式
type LogData struct {
	Message 	string		// 日志消息体
	TimeStr 	string		// 日志的打印时间
	LevelStr 	string		// 日志的级别
	FileName 	string		// 调用的日志输出的文件名
	FuncName 	string		// 调用日志输出的当前函数名
	LineNo 		int			// 调用日志输出的当前行号
	WarnAndFatal bool		// 是否是Warn Fatal Error日志
}


// 用于获取文件名和方法名和当前行号
func GetLineInfo()(fileName string, funcName string, lineNo int) {
	// caller 返回的其实是函数调用的栈信息
	pc, file,line,ok := runtime.Caller(4)
	if ok {
		fileName = path.Base(file)
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}
	return
}

func WriteLog(level int,format string,args...interface{}) *LogData{
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.999")
	levelStr := getLevelText(level)
	fileName, funcName, lineNo := GetLineInfo()
	msg := fmt.Sprintf(format, args...)
	// 时间 日志级别 文件名 行号 方法名 日志内容
	//fmt.Fprintf(file,"%s [%s] %s:%d %s %s\n", nowStr, levelStr, fileName,lineNo,funcName, msg)
	logData := &LogData{
		Message:msg,
		TimeStr:nowStr,
		LevelStr:levelStr,
		FileName:fileName,
		FuncName:funcName,
		LineNo:lineNo,
		WarnAndFatal:false,
	}
	if level == LogLevelError || level == LogLevelWarn || level == LogLevelFatal {
		logData.WarnAndFatal = true
	}
	return logData
}
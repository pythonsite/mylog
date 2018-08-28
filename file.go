/*
日志输出到文件中
*/

package mylog

import (
	"strconv"
	"fmt"
	"os"
	"time"
)

type FileLogger struct {
	level    int
	logPath  string
	logName  string
	file     *os.File		
	warnFile *os.File			// 用于写debug，fatal ,error日志内容
	LogDataChan chan *LogData
	logSplitType int
	logSplitSize int64
	lastSplitHour int
}

func NewFileLogger(config map[string] string) (log LogInterface, err error){
	logPath, ok := config["log_path"]
	if !ok {
		err = fmt.Errorf("not found log_path")
		return
	}
	logName, ok := config["log_name"]
	if !ok {
		err = fmt.Errorf("not found log_name")
		return
	}
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not found log_level")
		return
	}
	level := getLogLevel(logLevel)

	logChanSize, ok := config["log_chan_size"]
	if !ok {
		logChanSize = "50000"
	}
	chanSize, err := strconv.Atoi(logChanSize)
	if err != nil{
		chanSize = 50000
	}

	//默认以小时做切割
	var logSplitType int = LogSplitTypeHour
	var logSplitSize int64
	logSplitStr, ok := config["log_split_type"]
	if !ok{
		logSplitStr = "hour"
	} else {
		if logSplitStr == "size" {
			logSplitSizeStr, ok := config["log_split_size"]
			if !ok {
				logSplitSizeStr = "104857600"	// 默认是100M
			}
			logSplitSize, err = strconv.ParseInt(logSplitSizeStr,10,64)
			if err != nil {
				logSplitSize = 104857600		// 默认是100M
			}
			logSplitType = LogSplitTypeSize
		} else {
			logSplitType = LogSplitTypeHour
		}
	}

	log = &FileLogger{
		level:   level,
		logPath: logPath,
		logName: logName,
		LogDataChan:make(chan *LogData,chanSize),
		logSplitSize:logSplitSize,
		logSplitType:logSplitType,
		lastSplitHour:time.Now().Hour(),
	}
	log.Init()
	return
}

func (f *FileLogger) Init() {
	filename := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed ,err is :%v", filename, err))
	}
	f.file = file

	// 写错误日志和fatal,error日志的文件
	filename = fmt.Sprintf("%s/%s.wf", f.logPath, f.logName)
	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed ,err is :%v", filename, err))
	}
	f.warnFile = file
	go f.writeLogBackgroup()

}

// 以时间中的小时做切割
func (f *FileLogger) splitFileHour(warnFile bool) {
	now := time.Now()
	hour := now.Hour()
	
	if hour == f.lastSplitHour {
		return
	}
	f.lastSplitHour = hour
	var backupFileName string
	var fileName string
	if warnFile {
		backupFileName = fmt.Sprintf("%s/%s.wf_%04d%02d%2d%02d", f.logPath, f.logName,now.Year(),now.Month(),now.Day(), f.lastSplitHour)
		fileName = fmt.Sprintf("%s/%s.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d%02d%2d%02d", f.logPath, f.logName,now.Year(),now.Month(),now.Day(), f.lastSplitHour)
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}

	file := f.file
	if warnFile {
		file = f.warnFile
	}
	file.Close()
	os.Rename(fileName, backupFileName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

// 以日志文件大小做切割
func (f *FileLogger) splitFileSize(warnFile bool) {
	file := f.file
	if warnFile {
		file = f.warnFile
	}

	statInfo, err := file.Stat()
	if err != nil {
		return
	}
	fileSize := statInfo.Size()
	if fileSize <= f.logSplitSize {
		return
	}
	var backupFileName string
	var fileName string
	now := time.Now()
	if warnFile {
		backupFileName = fmt.Sprintf("%s/%s.wf_%04d%02d%2d%02d%02d%0d", 
		f.logPath, f.logName,now.Year(),now.Month(),now.Day(), now.Hour(),now.Minute(),now.Second)
		fileName = fmt.Sprintf("%s/%s.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d%02d%2d%02d%02d%0d", 
		f.logPath, f.logName,now.Year(),now.Month(),now.Day(), now.Hour(),now.Minute(),now.Second)
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}
	file.Close()
	// 将文件进行备份
	os.Rename(fileName, backupFileName)
	// 重新打开新的文件
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

// 确定切割文件的方式
func (f *FileLogger) checkSplitFile(warnFile bool) {
	if f.logSplitType == LogSplitTypeHour {
		f.splitFileHour(warnFile)
		return
	}
	f.splitFileSize(warnFile)
}

// 开线程进行写日志
func (f *FileLogger) writeLogBackgroup(){
	for logData := range f.LogDataChan {
		var file *os.File = f.file
		if logData.WarnAndFatal{
			file = f.warnFile
		}
		f.checkSplitFile(logData.WarnAndFatal)
		fmt.Fprintf(file,"%s [%s] %s:%d %s %s\n", logData.TimeStr,logData.LevelStr, logData.FileName,logData.LineNo,logData.FuncName, logData.Message)
	}
}

func (f *FileLogger) SetLevel(level int) {
	if f.level < LogLevelDebug || level > LogLevelFatal {
		f.level = LogLevelDebug
	}
	f.level = level
}


func (f *FileLogger) Debug(format string, args ...interface{}) {
	if f.level > LogLevelDebug {
		return
	}
	logData := WriteLog(LogLevelDebug, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}

}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	logData := WriteLog(LogLevelInfo, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}
	logData := WriteLog(LogLevelTrace, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}
	logData := WriteLog(LogLevelWarn, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}
	logData := WriteLog(LogLevelFatal, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}
	logData := WriteLog(LogLevelError, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Close() {
	f.file.Close()
	f.warnFile.Close()
}

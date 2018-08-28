package mylog


const (
	LogLevelDebug = iota
	LogLevelTrace
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)


const (
	LogSplitTypeHour = iota
	LogSplitTypeSize
)


func getLevelText(level int) string {
	switch level{
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	case LogLevelInfo:
		return "INFO"
	}
	return "UNKNOWN"
}

func getLogLevel(level string) int {
	switch level{
	case "debug":
		return LogLevelDebug
	case "trace":
		return LogLevelTrace
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelFatal
	case "info":
		return LogLevelInfo
	}
	return LogLevelDebug
}
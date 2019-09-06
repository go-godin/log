package log

var std = NewLogger("info")

func SetLevel(level string) {
	std.SetLevel(level)
}

func Info(message string, keyvals ...interface{}) {
	std.Info(message, keyvals...)
}

func Debug(message string, keyvals ...interface{}) {
	std.Debug(message, keyvals...)
}

func Warning(message string, keyvals ...interface{}) {
	std.Warning(message, keyvals...)
}

func Error(message string, keyvals ...interface{}) {
	std.Error(message, keyvals...)
}

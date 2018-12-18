package cphalo

import (
	"fmt"
	"log"
)

func logWithLevel(level string, v ...interface{}) {
	log.Print(fmt.Sprintf("[%s] ", level), fmt.Sprint(v...))
}

func logWithLevelf(level, format string, v ...interface{}) {
	log.Printf(fmt.Sprintf("[%s] %s", level, format), v...)
}

func logTrace(v ...interface{}) {
	logWithLevel("TRACE", v...)
}

func logTracef(format string, v ...interface{}) {
	logWithLevelf("TRACE", format, v...)
}

func logDebug(v ...interface{}) {
	logWithLevel("DEBUG", v...)
}

func logDebugf(format string, v ...interface{}) {
	logWithLevelf("DEBUG", format, v...)
}

func logInfo(v ...interface{}) {
	logWithLevel("INFO", v...)
}

func logInfof(format string, v ...interface{}) {
	logWithLevelf("INFO", format, v...)
}

func logWarn(v ...interface{}) {
	logWithLevel("WARN", v...)
}

func logWarnf(format string, v ...interface{}) {
	logWithLevelf("WARN", format, v...)
}

func logError(v ...interface{}) {
	logWithLevel("ERROR", v...)
}

func logErrorf(format string, v ...interface{}) {
	logWithLevelf("ERROR", format, v...)
}

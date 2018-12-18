package cphalo

import (
	"fmt"
	"log"
)

func logWithLevel(level string, v ...interface{}) {
	log.Print(fmt.Sprintf("[%s] ", level), fmt.Sprint(v...))
}

func logTrace(v ...interface{}) {
	logWithLevel("TRACE", v...)
}

func logDebug(v ...interface{}) {
	logWithLevel("DEBUG", v...)
}

func logInfo(v ...interface{}) {
	logWithLevel("INFO", v...)
}

func logWarn(v ...interface{}) {
	logWithLevel("WARN", v...)
}

func logError(v ...interface{}) {
	logWithLevel("ERROR", v...)
}

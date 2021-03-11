package network

import (
	"log"
	"os"
)

const INFO = 1
const WARNING = 2
const ERROR = 3
const NONE = 4

var logLevel = INFO
var infoLogger = log.New(os.Stdout, "Cobweb INFO: ", log.Ldate|log.Ltime)
var warningLogger = log.New(os.Stdout, "Cobweb WARNING: ", log.Ldate|log.Ltime)
var errorLogger = log.New(os.Stdout, "Cobweb ERROR: ", log.Ldate|log.Ltime)

func Info(v ...interface{})  {
	if logLevel <= INFO {
		infoLogger.Println(v...)
	}
}

func Warning(v ...interface{})  {
	if logLevel <= WARNING {
		warningLogger.Println(v...)
	}
}

func Error(v ...interface{})  {
	if logLevel <= ERROR {
		errorLogger.Println(v...)
	}
}

func SetLogLevel(level int)  {
	if level <= NONE && level >= INFO {
		logLevel = level
	}
}
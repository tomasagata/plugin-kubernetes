package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var infologger *log.Logger
var errorlogger *log.Logger
var debuglogger *log.Logger
var infoonce sync.Once
var erroronce sync.Once
var debugonce sync.Once
// var debugMode = false
var debugMode = true


const logFilePath = "/var/log/oakestra/cni_log2.txt"

type EventType string

const (
	DEPLOYREQUEST     EventType = "DEPLOY_REQUEST"
	UNDEPLOYREQUEST   EventType = "UNDEPLOY_REQUEST"
	DEPLOYED          EventType = "DEPLOYED"
	SERVICE_RESOURCES EventType = "RESOURCES"
	NODE_RESOURCES    EventType = "RESOURCES"
	DEAD              EventType = "DEAD"
)

func SetDebugMode() {
	debugMode = true
}

// getLogFile returns a writer to the log file, opening and creating it if needed
func getLogFile() io.Writer {
	dir := filepath.Dir(logFilePath)
	_ = os.MkdirAll(dir, 0755)

	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// fallback to stderr if file can't be opened
		return os.Stderr
	}
	return f
}

func InfoLogger() *log.Logger {
	infoonce.Do(func() {
		output := getLogFile()
		infologger = log.New(output, "INFO- ", log.Ldate|log.Ltime|log.Lshortfile)
	})
	return infologger
}

func ErrorLogger() *log.Logger {
	erroronce.Do(func() {
		output := getLogFile()
		errorlogger = log.New(output, "ERROR- ", log.Ldate|log.Ltime|log.Lshortfile)
	})
	return errorlogger
}

func DebugLogger() *log.Logger {
	debugonce.Do(func() {
		output := getLogFile()
		if !debugMode {
			output = io.Discard
		}
		debuglogger = log.New(output, "DEBUG- ", log.Ldate|log.Ltime|log.Lshortfile)
	})
	return debuglogger
}

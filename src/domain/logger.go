package domain

import (
	"fmt"
	"os"
	"time"
)

type ICustomLogger struct {
	file *os.File
}

var Logger *ICustomLogger

func CustomLogger(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	Logger = &ICustomLogger{file: file}
	return nil
}

func (l *ICustomLogger) Close() {
	err := l.file.Close()
	if err != nil {
		return
	}
}

func (l *ICustomLogger) WriteLog(logMessage string) {
	timestamp := time.Now().Format("02-01-2006 15:04:05")
	finalMessage := fmt.Sprintf("%s :: %s\n", timestamp, logMessage)
	_, err := l.file.WriteString(finalMessage)
	if err != nil {
		return
	}
}

func (l *ICustomLogger) Info(message string) {
	l.WriteLog(fmt.Sprintf("INFO: %s", message))
}

func (l *ICustomLogger) Debug(message string) {
	l.WriteLog(fmt.Sprintf("DEBUG: %s", message))
}

func (l *ICustomLogger) Error(message string) {
	l.WriteLog(fmt.Sprintf("ERROR: %s", message))
}

func (l *ICustomLogger) Warning(message string) {
	l.WriteLog(fmt.Sprintf("WARNING: %s", message))
}

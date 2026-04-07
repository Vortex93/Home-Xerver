package log

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// ========================================
// Bootstrap
// ========================================
func init() {
    // Create `logs` folder
    if err := os.MkdirAll("logs", os.ModePerm); err != nil {
        log.Fatalf("failed to create logs folder: %v", err)
    }
}

// ========================================
// Types
// ========================================
type Logger struct {
    name string
    logFile  *os.File
    errorFile *os.File

    consoleLogger *log.Logger
    logLogger *log.Logger
    errorLogger *log.Logger
}

// ========================================
// Functions
// ========================================

// Creates a new Logger with the given name. It will write info messages to `logs/{name}.log` and error messages to both `logs/{name}.log` and `logs/error.log`
func Create(name string) (*Logger, error) {
    // logs/{name}.log
    logPath := fmt.Sprintf("logs/%s.log", name)
    logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }

    // logs/error.Log
    errorFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }

    // Create a logger that writes to both logFile and errorFile
    logger := &Logger{
        name: name,
        logFile: logFile,
        errorFile: errorFile,

        consoleLogger: log.New(os.Stdout, fmt.Sprintf("[%s] ", name), log.LstdFlags),
        logLogger: log.New(logFile, fmt.Sprintf("[%s] ", name), log.LstdFlags),
        errorLogger: log.New(errorFile, fmt.Sprintf("[%s] ", name), log.LstdFlags),
    }

    return logger, nil
}

// ========================================
// Methods
// ========================================

// Writes an info message to the console and the log file
func (l *Logger) Info(msg string) {
    l.consoleLogger.Println("[INFO]", msg)
    l.logLogger.Println("[INFO]", msg)
}

// Writes a formatted info message to the console and the log file
func (l *Logger) Infof(format string, args ...any) {
    msg := fmt.Sprintf(format, args...)
    l.Info(msg)
}

// Writes a debug message to the console and the log file
func (l *Logger) Debug(msg string) {
    l.consoleLogger.Println("[DEBUG]", msg)
    l.logLogger.Println("[DEBUG]", msg)
}

// Writes a formatted info message to the console and the log file
func (l *Logger) Debugf(format string, args ...any) {
    msg := fmt.Sprintf(format, args...)
    l.Debug(msg)
}

// Writes an error message to the console, the log file, and the error file
func (l *Logger) Error(msg string) {
    l.consoleLogger.Println("[ERROR]", msg)
    l.logLogger.Println("[ERROR]", msg)
    l.errorLogger.Println("[ERROR]", msg)
}

func (l* Logger) Errorf(format string, args ...any) {
    msg := fmt.Sprintf(format, args...)
    l.Error(msg)
}

// Writes a fatal message to the console, the log file, and the error file, then exits the program with status code 1
func (l *Logger) Fatal(msg string) {
    l.consoleLogger.Println("[FATAL]", msg)
    l.logLogger.Println("[FATAL]", msg)
    l.errorLogger.Println("[FATAL]", msg)
    os.Exit(1)
}

// Writes a formatted fatal message to the console, the log file, and the error file, then exits the program with status code 1
func (l *Logger) Fatalf(format string, args ...any) {
    msg := fmt.Sprintf(format, args...)
    l.Fatal(msg)
}

// Closes the log files
func (l *Logger) Close() {
    if l.logFile != nil {
        l.logFile.Close()
    }
    if l.errorFile != nil {
        l.errorFile.Close()
    }
}

// ========================================
// For Debugging
// ========================================
func D(v any) {
    spew.Dump(v)
}

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
	"vms_go/internal/config"
)

var (
	info      *log.Logger
	warn      *log.Logger
	wrror     *log.Logger
	logFile   *os.File
	logFormat string // "plain" or "ecs"
)

// ECS log structure (minimal)
type ecsLog struct {
	Timestamp string `json:"@timestamp"`
	LogLevel  string `json:"log.level"`
	Message   string `json:"message"`
	File      string `json:"log.origin.file.name,omitempty"`
	Line      int    `json:"log.origin.file.line,omitempty"`
	Function  string `json:"log.origin.function,omitempty"`
}

// InitLogger configures the logger based on config (stdout/file, format)
func InitLogger(_ string) {
	cfg := config.FromEnv()
	var output io.Writer = os.Stdout
	logFormat = strings.ToLower(cfg.LOG_FORMAT)
	if logFormat != "ecs" {
		logFormat = "plain"
	}

	var infoEnabled, warnEnabled, errorEnabled bool
	switch strings.ToUpper(cfg.LOG_LEVEL) {
	case "INFO":
		infoEnabled, warnEnabled, errorEnabled = true, true, true
	case "WARN":
		infoEnabled, warnEnabled, errorEnabled = false, true, true
	case "ERROR":
		infoEnabled, warnEnabled, errorEnabled = false, false, true
	case "NONE":
		infoEnabled, warnEnabled, errorEnabled = false, false, false
	default:
		infoEnabled, warnEnabled, errorEnabled = true, true, true
	}

	if cfg.LOG_FILE != "" {
		file, err := os.OpenFile(cfg.LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open log file:", err)
		}
		logFile = file
		output = io.MultiWriter(os.Stdout, file)
	}

	flags := 0
	if logFormat == "plain" {
		flags = log.Ldate | log.Ltime | log.Lshortfile
	}
	info = createLogger(output, "INFO: ", flags, infoEnabled)
	warn = createLogger(output, "WARN: ", flags, warnEnabled)
	wrror = createLogger(output, "ERROR: ", flags, errorEnabled)
	LogInfo("Logger Initialized [format:", logFormat, "]")
}

func createLogger(output io.Writer, prefix string, flags int, enabled bool) *log.Logger {
	if !enabled {
		return log.New(io.Discard, prefix, flags)
	}
	return log.New(output, prefix, flags)
}

func isLoggerDiscarded(l *log.Logger) bool {
	// This is a bit hacky but works for standard library loggers
	return fmt.Sprintf("%v", l.Writer()) == "io.discard"
}

func logECS(lvl, msg string, calldepth int) string {
	pc, file, line, ok := runtime.Caller(calldepth)
	fn := ""
	if ok {
		fn = runtime.FuncForPC(pc).Name()
	}
	ecs := ecsLog{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		LogLevel:  strings.ToLower(lvl),
		Message:   msg,
		File:      file,
		Line:      line,
		Function:  fn,
	}
	b, _ := json.Marshal(ecs)
	return string(b)
}

func LogInfo(v ...interface{}) {
	if info != nil && !isLoggerDiscarded(info) {
		if logFormat == "ecs" {
			info.Output(2, logECS("INFO", fmt.Sprint(v...), 2))
		} else {
			info.Output(2, fmt.Sprint(v...))
		}
	}
}

func LogInfof(format string, v ...interface{}) {
	if info != nil && !isLoggerDiscarded(info) {
		if logFormat == "ecs" {
			info.Output(2, logECS("INFO", fmt.Sprintf(format, v...), 2))
		} else {
			info.Output(2, fmt.Sprintf(format, v...))
		}
	}
}

func LogWarn(v ...interface{}) {
	if warn != nil && !isLoggerDiscarded(warn) {
		if logFormat == "ecs" {
			warn.Output(2, logECS("WARN", fmt.Sprint(v...), 2))
		} else {
			warn.Output(2, fmt.Sprint(v...))
		}
	}
}

func LogWarnf(format string, v ...interface{}) {
	if warn != nil && !isLoggerDiscarded(warn) {
		if logFormat == "ecs" {
			warn.Output(2, logECS("WARN", fmt.Sprintf(format, v...), 2))
		} else {
			warn.Output(2, fmt.Sprintf(format, v...))
		}
	}
}

func LogError(v ...interface{}) {
	if wrror != nil && !isLoggerDiscarded(wrror) {
		if logFormat == "ecs" {
			wrror.Output(2, logECS("ERROR", fmt.Sprint(v...), 2))
		} else {
			wrror.Output(2, fmt.Sprint(v...))
		}
	}
}

func LogErrorf(format string, v ...interface{}) {
	if wrror != nil && !isLoggerDiscarded(wrror) {
		if logFormat == "ecs" {
			wrror.Output(2, logECS("ERROR", fmt.Sprintf(format, v...), 2))
		} else {
			wrror.Output(2, fmt.Sprintf(format, v...))
		}
	}
}

func CloseLogger() {
	if logFile != nil {
		info.Printf("Closing Logger")
		logFile.Close()
	}
}

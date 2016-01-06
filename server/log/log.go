package log

import (
	"log"
	"os"
	"fmt"
	"strings"
)

var (
	Logger *log.Logger
)

func init()  {
	logFile := getLogFile()
	fmt.Printf("Pillar log file: %s\n\n", logFile)
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file [%s], %s", logFile, err)
	}

	Logger = log.New(file, "Pillar: ", log.Ldate|log.Ltime)
}

func isDir(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}

	return false
}

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func getLogFile() string {
	pillarHome := os.Getenv("PILLAR_HOME")
	if pillarHome == "" {
		pillarHome = strings.TrimSuffix(os.TempDir(), "/") + "/pillar"
	}

	logPath := pillarHome + "/logs"
	if !exists(logPath) {
		os.MkdirAll(logPath, 0700)
	}

	return logPath + "/server.log"
}

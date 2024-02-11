package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

const (
	EVALUATION_INVOCATION     = "/evaluation/invocation"
	EVALUATION_HEADERS        = "/evaluation/headers"
	SCHEDULER_NAME            = "x-profaastinate-scheduler-name"
	EVALUATION_FUNCTION_START = "/evaluation/function-start"
	EVALUATION_FUNCTION_END   = "/evaluation/function-end"
	ASYNC_DEADLINE            = "x-profaastinate-process-deadline"
	FUNCTION_NAME             = "x-nuclio-function-name"
	FUNCTION_STATUS           = "x-nuclio-function-status" // start, end, invocation
	RELATIVE_PATH_LOG_PATH    = "profaastinate/evaluation/counter-backend/logs"
	ASYNC_EVALUATION_NAME     = "async.log"
	NORMAL_EVALUATION_NAME    = "normal.log"
	CPU_USAGE                 = "cpu-usage.log"
)

// Counter struct holds the count and a mutex to ensure safe access
type Counter struct {
	count   int
	headers []http.Header
	mu      sync.RWMutex
}

var (
	logger  *log.Logger
	logFile map[string]*os.File
)

func initLogger(filename string) {
	logFilePath := RELATIVE_PATH_LOG_PATH + filename
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(filepath.Dir(logFilePath), 0755)
		if errDir != nil {
			log.Fatalf("Failed to create directory for log file: %v", errDir)
		}
	}

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logFile[filename] = file
	logger = log.New(file, "", 0)
}

func (c *Counter) handleFunctionInvocations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.mu.Lock()
		defer c.mu.Unlock()

		// Increment count
		c.count++

		// Save headers
		c.headers = append(c.headers, r.Header)

		fmt.Fprintf(w, "Count increased by 1, Current Count: %d", c.count)
	case http.MethodGet:
		c.mu.RLock()
		defer c.mu.RUnlock()
		fmt.Fprintf(w, "Current Count: %d", c.count)
	case http.MethodDelete:
		c.mu.Lock()
		defer c.mu.Unlock()

		// Reset count and headers
		c.count = 0
		c.headers = nil

		fmt.Fprintf(w, "Count reset to: %d", c.count)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func logToFile(deadline string, functionName string, schedulerName string) {
	if deadline == "" {
		log.SetOutput(logFile[NORMAL_EVALUATION_NAME])
		log.Printf("%s", functionName)
	} else {
		log.SetOutput(logFile[ASYNC_EVALUATION_NAME])
		log.Printf("%s - %s", functionName, schedulerName)
	}
}

func (c *Counter) handleFunctionHeaders(w http.ResponseWriter, r *http.Request) {
	functionName := r.Header.Get(FUNCTION_NAME)
	deadline := r.Header.Get(ASYNC_DEADLINE)
	schedulerName := r.Header.Get(SCHEDULER_NAME)

	logToFile(deadline, functionName, schedulerName)
}

func main() {

	logFile = make(map[string]*os.File)
	logsList := []string{ASYNC_EVALUATION_NAME, NORMAL_EVALUATION_NAME, CPU_USAGE}
	for _, logName := range logsList {
		initLogger(logName)
		defer logFile[logName].Close()
	}

	counter := &Counter{count: 0, headers: make([]http.Header, 0)}

	server := &http.Server{
		Addr:         ":8888",
		Handler:      nil,
		ReadTimeout:  1000 * time.Millisecond,
		WriteTimeout: 1000 * time.Millisecond,
	}

	// Set up HTTP server with two endpoints
	http.HandleFunc(EVALUATION_INVOCATION, counter.handleFunctionInvocations)
	http.HandleFunc(EVALUATION_HEADERS, counter.handleFunctionHeaders)

	// Log CPU usage
	//go logCPUUsage()

	// Start the server
	fmt.Println("Server listening on :8888")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func logCPUUsage() {
	for {
		cpu, err := cpu.Percent(time.Second, false)
		if err != nil {
			log.Fatalf("Failed to get CPU usage: %v", err)
		}

		log.SetOutput(logFile[CPU_USAGE])

		// cpu usage with time stamp
		log.Printf("CPU Usage: %v%%", cpu)
		// Sleep for a while before checking again
		time.Sleep(time.Second)
	}
}

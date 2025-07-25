package app

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestActivityLogService(t *testing.T) {
	// Create a test configuration
	config := LogConfig{
		Level:         LogLevelInfo,
		BufferSize:    10,
		RetentionDays: 1,
		EnableConsole: false,
		EnableBuffer:  true,
		FlushInterval: 0, // Disable periodic flush for testing
	}

	// Create service without app context for testing
	service := NewActivityLogService(nil, config)
	defer service.Close()

	// Test basic logging
	service.LogEntry(LogLevelInfo, "Test", "TestOp", "Test message")
	service.LogEntry(LogLevelWarn, "Test", "TestOp", "Warning message")
	service.LogEntry(LogLevelError, "Test", "TestOp", "Error message")

	// Test log retrieval
	logs := service.GetLogs(LogFilter{})
	if len(logs) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(logs))
	}

	// Test level filtering
	warnLevel := LogLevelWarn
	filteredLogs := service.GetLogs(LogFilter{Level: &warnLevel})
	if len(filteredLogs) != 1 {
		t.Errorf("Expected 1 warning log, got %d", len(filteredLogs))
	}

	// Test component filtering
	component := "Test"
	componentLogs := service.GetLogs(LogFilter{Component: &component})
	if len(componentLogs) != 3 {
		t.Errorf("Expected 3 component logs, got %d", len(componentLogs))
	}

	// Test search
	searchReq := LogSearchRequest{
		Query:  "Warning",
		Filter: LogFilter{},
		Limit:  10,
		Offset: 0,
	}

	searchResult, err := service.SearchLogs(context.Background(), searchReq)
	if err != nil {
		t.Errorf("Search failed: %v", err)
	}
	if len(searchResult.Entries) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(searchResult.Entries))
	}

	// Test error reporting
	report := service.ReportError(ErrorTypeSystemErr, ErrorSeverityHigh, "Test", "TestError", "Test error", nil)
	if report == nil {
		t.Error("Expected error report, got nil")
	}
	if report.Type != ErrorTypeSystemErr {
		t.Errorf("Expected error type %v, got %v", ErrorTypeSystemErr, report.Type)
	}

	// Test application status
	status := service.GetApplicationStatus()
	if status == nil {
		t.Error("Expected application status, got nil")
	}
	if status.SystemHealth.MemoryUsage == 0 {
		t.Error("Expected non-zero memory usage")
	}
}

func TestActivityLogCircularBuffer(t *testing.T) {
	// Test circular buffer behavior with small buffer
	config := LogConfig{
		Level:         LogLevelInfo,
		BufferSize:    3,
		RetentionDays: 1,
		EnableConsole: false,
		EnableBuffer:  true,
		FlushInterval: 0,
	}

	service := NewActivityLogService(nil, config)
	defer service.Close()

	// Add more entries than buffer size
	for i := 0; i < 5; i++ {
		service.LogEntry(LogLevelInfo, "Test", "TestOp", fmt.Sprintf("Message %d", i))
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	}

	// Should only have the last 3 entries
	logs := service.GetLogs(LogFilter{})
	if len(logs) != 3 {
		t.Errorf("Expected 3 logs in circular buffer, got %d", len(logs))
	}

	// Verify we have the latest entries (2, 3, 4)
	foundMessages := make(map[string]bool)
	for _, log := range logs {
		foundMessages[log.Message] = true
	}

	expectedMessages := []string{"Message 2", "Message 3", "Message 4"}
	for _, expected := range expectedMessages {
		if !foundMessages[expected] {
			t.Errorf("Expected to find message '%s' in circular buffer", expected)
		}
	}
}

func TestActivityLogWithApp(t *testing.T) {
	// Test integration with App
	app := NewApp()

	// Test activity logging through app
	app.LogActivity(LogLevelInfo, "App", "TestOp", "App test message")

	// Test getting logs through app API
	logs, err := app.GetActivityLogs(context.Background(), LogFilter{})
	if err != nil {
		t.Errorf("Failed to get activity logs: %v", err)
	}

	if len(logs) == 0 {
		t.Error("Expected at least one log entry")
	}

	// Test search through app API
	searchReq := LogSearchRequest{
		Query:  "test",
		Filter: LogFilter{},
		Limit:  10,
		Offset: 0,
	}

	searchResult, err := app.SearchActivityLogs(context.Background(), searchReq)
	if err != nil {
		t.Errorf("Failed to search activity logs: %v", err)
	}

	if len(searchResult.Entries) == 0 {
		t.Error("Expected at least one search result")
	}

	// Test application status
	status, err := app.GetApplicationStatus(context.Background())
	if err != nil {
		t.Errorf("Failed to get application status: %v", err)
	}

	if status == nil {
		t.Error("Expected application status, got nil")
	}
}

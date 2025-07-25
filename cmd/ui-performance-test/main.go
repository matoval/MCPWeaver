package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"MCPWeaver/internal/database"
	"MCPWeaver/internal/parser"
	"MCPWeaver/internal/validator"
)

func main() {
	fmt.Println("UI Performance & Responsiveness Test")
	fmt.Println("====================================")

	// Test 1: Database Operations (UI data loading)
	fmt.Println("\n1. Testing Database Operations (Target: <100ms)")
	fmt.Println("------------------------------------------------")

	// Initialize database
	dbPath := "./test_ui_performance.db"
	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(dbPath)
	}()

	// Test project repository operations
	projectRepo := database.NewProjectRepository(db)

	// Test GetAll (common UI operation)
	start := time.Now()
	projects, err := projectRepo.GetAll()
	duration := time.Since(start)

	fmt.Printf("GetAll projects response time: %v\n", duration)
	if err != nil {
		fmt.Printf("⚠️  Error: %v\n", err)
	} else {
		fmt.Printf("Projects retrieved: %d\n", len(projects))
	}

	if duration < 100*time.Millisecond {
		fmt.Printf("✅ PASS: GetAll projects under 100ms\n")
	} else {
		fmt.Printf("❌ FAIL: GetAll projects exceeds 100ms\n")
	}

	// Test 2: File Operations (UI file handling)
	fmt.Println("\n2. Testing File Operations (Target: <100ms)")
	fmt.Println("---------------------------------------------")

	// Create test file
	testSpec := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      summary: Test endpoint
      responses:
        '200':
          description: Success
`

	testSpecPath := "./test_spec.yaml"
	err = os.WriteFile(testSpecPath, []byte(testSpec), 0644)
	if err != nil {
		log.Fatalf("Failed to write test spec: %v", err)
	}
	defer os.Remove(testSpecPath)

	// Test file reading (UI file import)
	start = time.Now()
	content, err := os.ReadFile(testSpecPath)
	duration = time.Since(start)

	fmt.Printf("File reading response time: %v\n", duration)
	if err != nil {
		fmt.Printf("⚠️  Error: %v\n", err)
	} else {
		fmt.Printf("File content length: %d\n", len(content))
	}

	if duration < 100*time.Millisecond {
		fmt.Printf("✅ PASS: File reading under 100ms\n")
	} else {
		fmt.Printf("❌ FAIL: File reading exceeds 100ms\n")
	}

	// Test file existence check (UI validation)
	start = time.Now()
	_, err = os.Stat(testSpecPath)
	exists := err == nil
	duration = time.Since(start)

	fmt.Printf("File existence check response time: %v\n", duration)
	fmt.Printf("File exists: %t\n", exists)

	if duration < 100*time.Millisecond {
		fmt.Printf("✅ PASS: File existence check under 100ms\n")
	} else {
		fmt.Printf("❌ FAIL: File existence check exceeds 100ms\n")
	}

	// Test 3: Parser Operations (UI spec processing)
	fmt.Println("\n3. Testing Parser Operations (Target: <100ms)")
	fmt.Println("----------------------------------------------")

	parserService := parser.NewService()

	// Test parsing (UI spec import)
	start = time.Now()
	_, err = parserService.ParseFromFile(testSpecPath)
	duration = time.Since(start)

	fmt.Printf("Spec parsing response time: %v\n", duration)
	if err != nil {
		fmt.Printf("⚠️  Error: %v\n", err)
	} else {
		fmt.Printf("Spec parsed successfully\n")
	}

	if duration < 100*time.Millisecond {
		fmt.Printf("✅ PASS: Spec parsing under 100ms\n")
	} else {
		fmt.Printf("❌ FAIL: Spec parsing exceeds 100ms\n")
	}

	// Test 4: Validation Operations (UI real-time validation)
	fmt.Println("\n4. Testing Validation Operations (Target: <100ms)")
	fmt.Println("--------------------------------------------------")

	validatorService := validator.New()

	// Test validation (UI spec validation)
	start = time.Now()
	ctx := context.Background()
	_, err = validatorService.ValidateFile(ctx, testSpecPath)
	duration = time.Since(start)

	fmt.Printf("Spec validation response time: %v\n", duration)
	if err != nil {
		fmt.Printf("⚠️  Error: %v\n", err)
	} else {
		fmt.Printf("Spec validated successfully\n")
	}

	if duration < 100*time.Millisecond {
		fmt.Printf("✅ PASS: Spec validation under 100ms\n")
	} else {
		fmt.Printf("❌ FAIL: Spec validation exceeds 100ms\n")
	}

	// Test 5: Memory Operations (UI monitoring)
	fmt.Println("\n5. Testing Memory Operations (Target: <10ms)")
	fmt.Println("---------------------------------------------")

	// Test memory stats (UI memory display)
	start = time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	duration = time.Since(start)

	fmt.Printf("Memory stats response time: %v\n", duration)
	fmt.Printf("Memory allocated: %.2f MB\n", float64(m.Alloc)/1024/1024)

	if duration < 10*time.Millisecond {
		fmt.Printf("✅ PASS: Memory stats under 10ms\n")
	} else {
		fmt.Printf("❌ FAIL: Memory stats exceeds 10ms\n")
	}

	// Test 6: Concurrent Operations (UI concurrent tasks)
	fmt.Println("\n6. Testing Concurrent Operations (Target: <200ms)")
	fmt.Println("--------------------------------------------------")

	start = time.Now()

	// Simulate multiple concurrent UI operations
	done := make(chan bool, 5)

	// Concurrent operations that UI might perform
	go func() {
		defer func() { done <- true }()
		projectRepo.GetAll()
	}()

	go func() {
		defer func() { done <- true }()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
	}()

	go func() {
		defer func() { done <- true }()
		os.Stat(testSpecPath)
	}()

	go func() {
		defer func() { done <- true }()
		parserService.ParseFromFile(testSpecPath)
	}()

	go func() {
		defer func() { done <- true }()
		validatorService.ValidateFile(context.Background(), testSpecPath)
	}()

	// Wait for all operations to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	duration = time.Since(start)
	fmt.Printf("5 concurrent operations time: %v\n", duration)

	if duration < 200*time.Millisecond {
		fmt.Printf("✅ PASS: Concurrent operations under 200ms\n")
	} else {
		fmt.Printf("❌ FAIL: Concurrent operations exceed 200ms\n")
	}

	// Test 7: Memory Stability During Operations
	fmt.Println("\n7. Testing Memory Stability (Target: <5MB growth)")
	fmt.Println("--------------------------------------------------")

	// Initial memory reading
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	initialAlloc := m1.Alloc

	// Simulate intensive UI operations
	for i := 0; i < 50; i++ {
		projectRepo.GetAll()

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		os.Stat(testSpecPath)
		parserService.ParseFromFile(testSpecPath)

		// Periodic GC
		if i%10 == 0 {
			runtime.GC()
			time.Sleep(10 * time.Millisecond)
		}
	}

	// Final memory reading
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	finalAlloc := m2.Alloc

	memoryGrowth := float64(finalAlloc-initialAlloc) / 1024 / 1024
	fmt.Printf("Memory growth after 50 operations: %.2f MB\n", memoryGrowth)

	if memoryGrowth < 5.0 {
		fmt.Printf("✅ PASS: Memory growth under 5MB\n")
	} else {
		fmt.Printf("⚠️  WARN: Memory growth exceeds 5MB\n")
	}

	// Summary
	fmt.Println("\n=== UI PERFORMANCE SUMMARY ===")
	fmt.Printf("Database Operations: Measured\n")
	fmt.Printf("File Operations: Measured\n")
	fmt.Printf("Parser Operations: Measured\n")
	fmt.Printf("Validation Operations: Measured\n")
	fmt.Printf("Memory Operations: Measured\n")
	fmt.Printf("Concurrent Operations: Measured\n")
	fmt.Printf("Memory Stability: %.2f MB growth\n", memoryGrowth)

	fmt.Println("\n🎉 UI PERFORMANCE TESTING COMPLETED!")
	fmt.Println("\nNote: This test simulates UI operations and measures backend response times.")
	fmt.Println("Frontend rendering performance depends on React components and browser optimization.")
}

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
	fmt.Println("MCPWeaver Performance Test")
	fmt.Println("==========================")

	// Test 1: Database Startup Time
	fmt.Println("\n1. Testing Database Startup Time")
	fmt.Println("---------------------------------")

	dbStart := time.Now()

	// Test database initialization
	dbPath := "./test_performance.db"
	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(dbPath) // Clean up test database
	}()

	dbDuration := time.Since(dbStart)
	fmt.Printf("Database initialization time: %v\n", dbDuration)

	if dbDuration < 500*time.Millisecond {
		fmt.Printf("✅ PASS: Database initialization under 500ms\n")
	} else {
		fmt.Printf("❌ FAIL: Database initialization exceeds 500ms\n")
	}

	// Test 2: Memory Usage
	fmt.Println("\n2. Testing Memory Usage (Target: <50MB)")
	fmt.Println("----------------------------------------")

	// Force garbage collection to get accurate reading
	runtime.GC()
	time.Sleep(100 * time.Millisecond) // Allow GC to complete

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	allocMB := float64(m.Alloc) / 1024 / 1024
	sysMB := float64(m.Sys) / 1024 / 1024

	fmt.Printf("Memory allocated: %.2f MB\n", allocMB)
	fmt.Printf("Memory system: %.2f MB\n", sysMB)

	if allocMB < 50 {
		fmt.Printf("✅ PASS: Memory usage is under 50MB\n")
	} else {
		fmt.Printf("❌ FAIL: Memory usage exceeds 50MB\n")
	}

	// Test 3: Parser Performance
	fmt.Println("\n3. Testing Parser Performance")
	fmt.Println("------------------------------")

	parserService := parser.NewService()

	// Create a small test OpenAPI spec
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

	// Write test spec to file
	testSpecPath := "./test_spec.yaml"
	err = os.WriteFile(testSpecPath, []byte(testSpec), 0644)
	if err != nil {
		log.Fatalf("Failed to write test spec: %v", err)
	}
	defer os.Remove(testSpecPath)

	// Test parsing performance
	parseStart := time.Now()
	_, err = parserService.ParseFromFile(testSpecPath)
	parseDuration := time.Since(parseStart)

	fmt.Printf("Small spec parsing time: %v\n", parseDuration)

	if err != nil {
		fmt.Printf("❌ FAIL: Parser error: %v\n", err)
	} else if parseDuration < 100*time.Millisecond {
		fmt.Printf("✅ PASS: Small spec parsing under 100ms\n")
	} else {
		fmt.Printf("❌ FAIL: Small spec parsing exceeds 100ms\n")
	}

	// Test 4: Validator Performance
	fmt.Println("\n4. Testing Validator Performance")
	fmt.Println("---------------------------------")

	validatorService := validator.New()

	ctx := context.Background()
	validateStart := time.Now()
	_, err = validatorService.ValidateFile(ctx, testSpecPath)
	validateDuration := time.Since(validateStart)

	fmt.Printf("Validation time: %v\n", validateDuration)

	if err != nil {
		fmt.Printf("❌ FAIL: Validator error: %v\n", err)
	} else if validateDuration < 50*time.Millisecond {
		fmt.Printf("✅ PASS: Validation under 50ms\n")
	} else {
		fmt.Printf("❌ FAIL: Validation exceeds 50ms\n")
	}

	// Test 5: Memory Leak Detection
	fmt.Println("\n5. Memory Leak Detection")
	fmt.Println("------------------------")

	// Take initial memory reading
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	initialAlloc := m1.Alloc

	// Simulate operations that might cause memory leaks
	for i := 0; i < 100; i++ {
		// Simulate parsing operations
		_, _ = parserService.ParseFromFile(testSpecPath)

		// Simulate validation operations
		_, _ = validatorService.ValidateFile(context.Background(), testSpecPath)

		// Force GC periodically
		if i%20 == 0 {
			runtime.GC()
			time.Sleep(10 * time.Millisecond)
		}
	}

	// Take final memory reading
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	finalAlloc := m2.Alloc

	memoryGrowth := float64(finalAlloc-initialAlloc) / 1024 / 1024
	fmt.Printf("Memory growth after 100 operations: %.2f MB\n", memoryGrowth)

	if memoryGrowth < 5.0 { // Less than 5MB growth is acceptable for 100 operations
		fmt.Printf("✅ PASS: No significant memory leaks detected\n")
	} else {
		fmt.Printf("⚠️  WARN: Potential memory leak detected (%.2f MB growth)\n", memoryGrowth)
	}

	// Test 6: Concurrent Performance
	fmt.Println("\n6. Testing Concurrent Performance")
	fmt.Println("----------------------------------")

	concurrentStart := time.Now()

	// Run multiple operations concurrently
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			// Parse spec
			_, _ = parserService.ParseFromFile(testSpecPath)

			// Validate spec
			_, _ = validatorService.ValidateFile(context.Background(), testSpecPath)
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	concurrentDuration := time.Since(concurrentStart)
	fmt.Printf("10 concurrent operations time: %v\n", concurrentDuration)

	if concurrentDuration < 1*time.Second {
		fmt.Printf("✅ PASS: Concurrent operations under 1 second\n")
	} else {
		fmt.Printf("❌ FAIL: Concurrent operations exceed 1 second\n")
	}

	// Summary
	fmt.Println("\n=== PERFORMANCE TEST SUMMARY ===")
	fmt.Printf("Database Init: %v (Target: <500ms)\n", dbDuration)
	fmt.Printf("Memory Usage: %.2f MB (Target: <50MB)\n", allocMB)
	fmt.Printf("Small Spec Parse: %v (Target: <100ms)\n", parseDuration)
	fmt.Printf("Validation: %v (Target: <50ms)\n", validateDuration)
	fmt.Printf("Memory Growth: %.2f MB (Target: <5MB)\n", memoryGrowth)
	fmt.Printf("Concurrent Ops: %v (Target: <1s)\n", concurrentDuration)

	// Overall result
	overallPass := dbDuration < 500*time.Millisecond &&
		allocMB < 50 &&
		parseDuration < 100*time.Millisecond &&
		validateDuration < 50*time.Millisecond &&
		memoryGrowth < 5.0 &&
		concurrentDuration < 1*time.Second

	if overallPass {
		fmt.Println("\n🎉 ALL PERFORMANCE TESTS PASSED!")
		os.Exit(0)
	} else {
		fmt.Println("\n❌ SOME PERFORMANCE TESTS FAILED")
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Frontend Bundle Size Test")
	fmt.Println("==========================")

	// Change to frontend directory
	frontendDir := "./frontend"

	// Check if we can measure the existing dist directory
	distDir := filepath.Join(frontendDir, "dist")

	if _, err := os.Stat(distDir); err == nil {
		fmt.Println("Measuring existing build...")
		measureBundleSize(distDir)
	} else {
		fmt.Println("No existing build found, attempting to build...")

		// Try to build the frontend
		cmd := exec.Command("npm", "run", "build")
		cmd.Dir = frontendDir

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Build failed: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))

			// Try to measure any partial build
			if _, err := os.Stat(distDir); err == nil {
				fmt.Println("Measuring partial build...")
				measureBundleSize(distDir)
			} else {
				fmt.Println("No build artifacts found")
			}
		} else {
			fmt.Println("Build successful!")
			measureBundleSize(distDir)
		}
	}

	// Performance targets
	fmt.Println("\n=== BUNDLE SIZE TARGETS ===")
	fmt.Println("Target: < 500KB total")
	fmt.Println("Target: < 300KB JS")
	fmt.Println("Target: < 100KB CSS")
	fmt.Println("Target: < 100KB Assets")
}

func measureBundleSize(distDir string) {
	fmt.Printf("\nMeasuring bundle size in: %s\n", distDir)
	fmt.Println("----------------------------------")

	var totalSize int64
	var jsSize int64
	var cssSize int64
	var assetSize int64

	err := filepath.Walk(distDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			totalSize += info.Size()

			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case ".js":
				jsSize += info.Size()
				fmt.Printf("JS: %s (%.2f KB)\n", info.Name(), float64(info.Size())/1024)
			case ".css":
				cssSize += info.Size()
				fmt.Printf("CSS: %s (%.2f KB)\n", info.Name(), float64(info.Size())/1024)
			case ".png", ".jpg", ".jpeg", ".svg", ".ico", ".woff", ".woff2":
				assetSize += info.Size()
				fmt.Printf("Asset: %s (%.2f KB)\n", info.Name(), float64(info.Size())/1024)
			default:
				fmt.Printf("Other: %s (%.2f KB)\n", info.Name(), float64(info.Size())/1024)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error measuring bundle size: %v\n", err)
		return
	}

	fmt.Printf("\n=== BUNDLE SIZE SUMMARY ===\n")
	fmt.Printf("Total size: %.2f KB\n", float64(totalSize)/1024)
	fmt.Printf("JavaScript: %.2f KB\n", float64(jsSize)/1024)
	fmt.Printf("CSS: %.2f KB\n", float64(cssSize)/1024)
	fmt.Printf("Assets: %.2f KB\n", float64(assetSize)/1024)

	// Check against targets
	totalKB := float64(totalSize) / 1024
	jsKB := float64(jsSize) / 1024
	cssKB := float64(cssSize) / 1024
	assetKB := float64(assetSize) / 1024

	fmt.Printf("\n=== PERFORMANCE ASSESSMENT ===\n")

	if totalKB < 500 {
		fmt.Printf("✅ PASS: Total size under 500KB\n")
	} else {
		fmt.Printf("❌ FAIL: Total size exceeds 500KB\n")
	}

	if jsKB < 300 {
		fmt.Printf("✅ PASS: JavaScript size under 300KB\n")
	} else {
		fmt.Printf("❌ FAIL: JavaScript size exceeds 300KB\n")
	}

	if cssKB < 100 {
		fmt.Printf("✅ PASS: CSS size under 100KB\n")
	} else {
		fmt.Printf("❌ FAIL: CSS size exceeds 100KB\n")
	}

	if assetKB < 100 {
		fmt.Printf("✅ PASS: Asset size under 100KB\n")
	} else {
		fmt.Printf("❌ FAIL: Asset size exceeds 100KB\n")
	}

	overallPass := totalKB < 500 && jsKB < 300 && cssKB < 100 && assetKB < 100

	if overallPass {
		fmt.Println("\n🎉 BUNDLE SIZE TARGETS MET!")
	} else {
		fmt.Println("\n⚠️  BUNDLE SIZE NEEDS OPTIMIZATION")
	}
}

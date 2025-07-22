package testing

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// CompilationValidator validates that the generated server compiles successfully
type CompilationValidator struct {
	config *TestConfig
}

// NewCompilationValidator creates a new compilation validator
func NewCompilationValidator(config *TestConfig) *CompilationValidator {
	return &CompilationValidator{config: config}
}

// Name returns the validator name
func (v *CompilationValidator) Name() string {
	return "compilation"
}

// SupportsAsync returns whether this validator supports async execution
func (v *CompilationValidator) SupportsAsync() bool {
	return true
}

// Validate checks if the server compiles successfully
func (v *CompilationValidator) Validate(ctx context.Context, serverPath string) (*ValidationResult, error) {
	startTime := time.Now()
	result := &ValidationResult{
		ValidatorName: v.Name(),
		Success:       true,
		Errors:        make([]string, 0),
		Warnings:      make([]string, 0),
	}

	// Check if main.go exists
	mainFile := filepath.Join(serverPath, "main.go")
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		result.Success = false
		result.Errors = append(result.Errors, "main.go file not found")
		result.Duration = time.Since(startTime)
		return result, nil
	}

	// Check if go.mod exists
	goModFile := filepath.Join(serverPath, "go.mod")
	if _, err := os.Stat(goModFile); os.IsNotExist(err) {
		result.Success = false
		result.Errors = append(result.Errors, "go.mod file not found")
		result.Duration = time.Since(startTime)
		return result, nil
	}

	result.FilesValidated = 2

	// Try to compile the server
	ctx, cancel := context.WithTimeout(ctx, v.config.Timeout)
	defer cancel()

	// Use secure command execution
	sch := NewSecureCommandHelper()
	cmd, err := sch.SecureCompileCommand(ctx, serverPath, "test-server", "main.go")
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to create secure compile command: %v", err))
		result.Duration = time.Since(startTime)
		return result, nil
	}
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Compilation failed: %v", err))
		result.Errors = append(result.Errors, string(output))
		result.Duration = time.Since(startTime)
		return result, nil
	}

	// Clean up compiled binary
	compiledBinary := filepath.Join(serverPath, "test-server")
	if err := os.Remove(compiledBinary); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Failed to remove test binary: %v", err))
	}

	// Check for compilation warnings in output
	if len(output) > 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Compilation warnings: %s", string(output)))
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// SyntaxValidator validates Go syntax and structure
type SyntaxValidator struct {
	config *TestConfig
}

// NewSyntaxValidator creates a new syntax validator
func NewSyntaxValidator(config *TestConfig) *SyntaxValidator {
	return &SyntaxValidator{config: config}
}

// Name returns the validator name
func (v *SyntaxValidator) Name() string {
	return "syntax"
}

// SupportsAsync returns whether this validator supports async execution
func (v *SyntaxValidator) SupportsAsync() bool {
	return false
}

// Validate checks Go syntax and structure
func (v *SyntaxValidator) Validate(ctx context.Context, serverPath string) (*ValidationResult, error) {
	startTime := time.Now()
	result := &ValidationResult{
		ValidatorName: v.Name(),
		Success:       true,
		Errors:        make([]string, 0),
		Warnings:      make([]string, 0),
	}

	// Parse main.go
	mainFile := filepath.Join(serverPath, "main.go")
	if err := v.validateGoFile(mainFile, result); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to validate main.go: %v", err))
	}

	// Validate go.mod structure
	goModFile := filepath.Join(serverPath, "go.mod")
	if err := v.validateGoMod(goModFile, result); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to validate go.mod: %v", err))
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// validateGoFile validates a Go source file
func (v *SyntaxValidator) validateGoFile(filePath string, result *ValidationResult) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	result.FilesValidated++

	// Parse the Go file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("syntax error: %w", err)
	}

	// Validate package declaration
	if node.Name.Name != "main" {
		result.Errors = append(result.Errors, "Expected 'package main' declaration")
	}

	// Check for required imports
	requiredImports := map[string]bool{
		"context":     false,
		"encoding/json": false,
		"fmt":         false,
		"log":         false,
		"os":          false,
		"time":        false,
		"github.com/sourcegraph/jsonrpc2": false,
	}

	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")
		if _, exists := requiredImports[importPath]; exists {
			requiredImports[importPath] = true
		}
	}

	// Check for missing imports
	for imp, found := range requiredImports {
		if !found {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Missing import: %s", imp))
		}
	}

	// Check for main function
	hasMainFunc := false
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name.Name == "main" && fn.Recv == nil {
				hasMainFunc = true
			}
		}
		return true
	})

	if !hasMainFunc {
		result.Errors = append(result.Errors, "Missing main() function")
	}

	// Check for required types and structures
	v.validateMCPStructures(node, result)

	return nil
}

// validateMCPStructures checks for required MCP-related structures
func (v *SyntaxValidator) validateMCPStructures(node *ast.File, result *ValidationResult) {
	requiredTypes := map[string]bool{
		"MCPServer":    false,
		"Tool":         false,
		"InputSchema":  false,
		"Property":     false,
		"ToolRequest":  false,
		"ToolResponse": false,
		"Content":      false,
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if _, exists := requiredTypes[typeSpec.Name.Name]; exists {
				requiredTypes[typeSpec.Name.Name] = true
			}
		}
		return true
	})

	// Check for missing types
	for typeName, found := range requiredTypes {
		if !found {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Missing type definition: %s", typeName))
		}
	}
}

// validateGoMod validates go.mod file structure
func (v *SyntaxValidator) validateGoMod(filePath string, result *ValidationResult) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	result.FilesValidated++
	contentStr := string(content)

	// Check for module declaration
	moduleRegex := regexp.MustCompile(`module\s+[\w\-./]+`)
	if !moduleRegex.MatchString(contentStr) {
		result.Errors = append(result.Errors, "Invalid or missing module declaration in go.mod")
	}

	// Check for Go version
	goVersionRegex := regexp.MustCompile(`go\s+1\.\d+`)
	if !goVersionRegex.MatchString(contentStr) {
		result.Warnings = append(result.Warnings, "Missing or invalid Go version in go.mod")
	}

	// Check for required dependency
	if !strings.Contains(contentStr, "github.com/sourcegraph/jsonrpc2") {
		result.Errors = append(result.Errors, "Missing required dependency: github.com/sourcegraph/jsonrpc2")
	}

	return nil
}

// LintValidator runs linting tools on the generated code
type LintValidator struct {
	config *TestConfig
}

// NewLintValidator creates a new lint validator
func NewLintValidator(config *TestConfig) *LintValidator {
	return &LintValidator{config: config}
}

// Name returns the validator name
func (v *LintValidator) Name() string {
	return "lint"
}

// SupportsAsync returns whether this validator supports async execution
func (v *LintValidator) SupportsAsync() bool {
	return true
}

// Validate runs linting tools
func (v *LintValidator) Validate(ctx context.Context, serverPath string) (*ValidationResult, error) {
	startTime := time.Now()
	result := &ValidationResult{
		ValidatorName: v.Name(),
		Success:       true,
		Errors:        make([]string, 0),
		Warnings:      make([]string, 0),
	}

	// Run go fmt
	if err := v.runGoFmt(ctx, serverPath, result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("go fmt issues: %v", err))
	}

	// Run go vet
	if err := v.runGoVet(ctx, serverPath, result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("go vet issues: %v", err))
	}

	// Run golangci-lint if available
	if err := v.runGolangciLint(ctx, serverPath, result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("golangci-lint issues: %v", err))
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// runGoFmt checks code formatting
func (v *LintValidator) runGoFmt(ctx context.Context, serverPath string, result *ValidationResult) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	sch := NewSecureCommandHelper()
	cmd, err := sch.SecureExecCommand(ctx, serverPath, "go", "fmt", "./...")
	if err != nil {
		return fmt.Errorf("failed to create secure go fmt command: %w", err)
	}
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("go fmt failed: %v, output: %s", err, output)
	}

	if len(output) > 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Formatting issues found: %s", output))
	}

	result.FilesValidated++
	return nil
}

// runGoVet checks for common Go mistakes
func (v *LintValidator) runGoVet(ctx context.Context, serverPath string, result *ValidationResult) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	sch := NewSecureCommandHelper()
	cmd, err := sch.SecureExecCommand(ctx, serverPath, "go", "vet", "./...")
	if err != nil {
		return fmt.Errorf("failed to create secure go vet command: %w", err)
	}
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("go vet failed: %v, output: %s", err, output)
	}

	if len(output) > 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Vet issues found: %s", output))
	}

	result.FilesValidated++
	return nil
}

// runGolangciLint runs golangci-lint if available
func (v *LintValidator) runGolangciLint(ctx context.Context, serverPath string, result *ValidationResult) error {
	// Check if golangci-lint is available
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		result.Warnings = append(result.Warnings, "golangci-lint not available, skipping advanced linting")
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	sch := NewSecureCommandHelper()
	cmd, err := sch.SecureExecCommand(ctx, serverPath, "golangci-lint", "run", "--fast")
	if err != nil {
		return fmt.Errorf("failed to create secure golangci-lint command: %w", err)
	}
	output, err := cmd.CombinedOutput()

	if err != nil {
		// golangci-lint returns non-zero exit code when issues are found
		if len(output) > 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Linting issues: %s", output))
		}
		return nil
	}

	result.FilesValidated++
	return nil
}

// SecurityValidator checks for security issues
type SecurityValidator struct {
	config *TestConfig
}

// NewSecurityValidator creates a new security validator
func NewSecurityValidator(config *TestConfig) *SecurityValidator {
	return &SecurityValidator{config: config}
}

// Name returns the validator name
func (v *SecurityValidator) Name() string {
	return "security"
}

// SupportsAsync returns whether this validator supports async execution
func (v *SecurityValidator) SupportsAsync() bool {
	return true
}

// Validate checks for security issues
func (v *SecurityValidator) Validate(ctx context.Context, serverPath string) (*ValidationResult, error) {
	startTime := time.Now()
	result := &ValidationResult{
		ValidatorName: v.Name(),
		Success:       true,
		Errors:        make([]string, 0),
		Warnings:      make([]string, 0),
	}

	// Check for common security issues
	if err := v.checkSecurityIssues(serverPath, result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Security check failed: %v", err))
	}

	// Run gosec if available
	if err := v.runGosec(ctx, serverPath, result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Gosec issues: %v", err))
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// checkSecurityIssues performs basic security checks
func (v *SecurityValidator) checkSecurityIssues(serverPath string, result *ValidationResult) error {
	mainFile := filepath.Join(serverPath, "main.go")
	content, err := os.ReadFile(mainFile)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	result.FilesValidated++
	contentStr := string(content)

	// Check for potential security issues
	securityChecks := []struct {
		pattern string
		message string
	}{
		{`os\.Getenv\s*\(\s*"[^"]*PASSWORD[^"]*"\s*\)`, "Potential password exposure via environment variable"},
		{`os\.Getenv\s*\(\s*"[^"]*SECRET[^"]*"\s*\)`, "Potential secret exposure via environment variable"},
		{`fmt\.Print[fl]?\s*\([^)]*password[^)]*\)`, "Potential password logging"},
		{`log\.[^(]*\([^)]*password[^)]*\)`, "Potential password logging"},
		{`http\.Client\s*\{[^}]*Timeout:\s*0`, "HTTP client without timeout"},
		{`exec\.Command[^(]*\([^)]*\$`, "Potential command injection vulnerability"},
	}

	for _, check := range securityChecks {
		matched, err := regexp.MatchString(check.pattern, contentStr)
		if err != nil {
			continue
		}
		if matched {
			result.Warnings = append(result.Warnings, check.message)
		}
	}

	return nil
}

// runGosec runs gosec security scanner if available
func (v *SecurityValidator) runGosec(ctx context.Context, serverPath string, result *ValidationResult) error {
	// Check if gosec is available
	if _, err := exec.LookPath("gosec"); err != nil {
		result.Warnings = append(result.Warnings, "gosec not available, skipping advanced security scanning")
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	sch := NewSecureCommandHelper()
	cmd, err := sch.SecureExecCommand(ctx, serverPath, "gosec", "-fmt", "text", "./...")
	if err != nil {
		return fmt.Errorf("failed to create secure gosec command: %w", err)
	}
	output, err := cmd.CombinedOutput()

	if err != nil {
		// gosec returns non-zero exit code when issues are found
		if len(output) > 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Security issues found: %s", output))
		}
		return nil
	}

	result.FilesValidated++
	return nil
}

// DependencyValidator checks dependencies and vulnerabilities
type DependencyValidator struct {
	config *TestConfig
}

// NewDependencyValidator creates a new dependency validator
func NewDependencyValidator(config *TestConfig) *DependencyValidator {
	return &DependencyValidator{config: config}
}

// Name returns the validator name
func (v *DependencyValidator) Name() string {
	return "dependencies"
}

// SupportsAsync returns whether this validator supports async execution
func (v *DependencyValidator) SupportsAsync() bool {
	return true
}

// Validate checks dependencies for vulnerabilities
func (v *DependencyValidator) Validate(ctx context.Context, serverPath string) (*ValidationResult, error) {
	startTime := time.Now()
	result := &ValidationResult{
		ValidatorName: v.Name(),
		Success:       true,
		Errors:        make([]string, 0),
		Warnings:      make([]string, 0),
	}

	// Run go mod verify
	if err := v.runGoModVerify(ctx, serverPath, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Module verification failed: %v", err))
		result.Success = false
	}

	// Run govulncheck if available
	if err := v.runGovulncheck(ctx, serverPath, result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Vulnerability check issues: %v", err))
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// runGoModVerify verifies module dependencies
func (v *DependencyValidator) runGoModVerify(ctx context.Context, serverPath string, result *ValidationResult) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	sch := NewSecureCommandHelper()
	cmd, err := sch.SecureExecCommand(ctx, serverPath, "go", "mod", "verify")
	if err != nil {
		return fmt.Errorf("failed to create secure go mod verify command: %w", err)
	}
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("go mod verify failed: %v, output: %s", err, output)
	}

	result.FilesValidated++
	return nil
}

// runGovulncheck checks for known vulnerabilities
func (v *DependencyValidator) runGovulncheck(ctx context.Context, serverPath string, result *ValidationResult) error {
	// Check if govulncheck is available
	if _, err := exec.LookPath("govulncheck"); err != nil {
		result.Warnings = append(result.Warnings, "govulncheck not available, skipping vulnerability scanning")
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	sch := NewSecureCommandHelper()
	cmd, err := sch.SecureExecCommand(ctx, serverPath, "govulncheck", "./...")
	if err != nil {
		return fmt.Errorf("failed to create secure govulncheck command: %w", err)
	}
	output, err := cmd.CombinedOutput()

	if err != nil {
		// govulncheck returns non-zero exit code when vulnerabilities are found
		if len(output) > 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Vulnerabilities found: %s", output))
		}
		return nil
	}

	result.FilesValidated++
	return nil
}
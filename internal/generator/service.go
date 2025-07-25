package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"MCPWeaver/internal/mapping"
	"MCPWeaver/internal/parser"
)

// TemplateData holds all data needed for template generation
type TemplateData struct {
	PackageName string
	APITitle    string
	APIVersion  string
	BaseURL     string
	Tools       []mapping.MCPTool
}

// SecurityValidator interface for template security validation
type SecurityValidator interface {
	ValidateTemplateName(name string) error
	SanitizeTemplateVariable(variable string) string
	ValidateFilePath(path string) error
}

// Service handles MCP server code generation
type Service struct {
	outputDir         string
	securityValidator SecurityValidator
}

// ValidationResult represents the result of code validation
type ValidationResult struct {
	IsValid        bool     `json:"isValid"`
	Errors         []string `json:"errors,omitempty"`
	Warnings       []string `json:"warnings,omitempty"`
	FilesValidated int      `json:"filesValidated"`
}

// NewService creates a new code generator service
func NewService(outputDir string) *Service {
	return &Service{
		outputDir: outputDir,
	}
}

// NewServiceWithValidator creates a new code generator service with security validator
func NewServiceWithValidator(outputDir string, validator SecurityValidator) *Service {
	return &Service{
		outputDir:         outputDir,
		securityValidator: validator,
	}
}

// Generate creates a complete MCP server from parsed API and tools
func (s *Service) Generate(api *parser.ParsedAPI, tools []mapping.MCPTool, serverName string) error {
	// Create output directory structure
	if err := s.createOutputStructure(); err != nil {
		return fmt.Errorf("failed to create output structure: %w", err)
	}

	// Prepare template data with security sanitization
	data := TemplateData{
		PackageName: s.sanitizePackageName(serverName),
		APITitle:    s.sanitizeTemplateData(api.Title),
		APIVersion:  s.sanitizeTemplateData(api.Version),
		BaseURL:     s.sanitizeTemplateData(api.BaseURL),
		Tools:       s.sanitizeTools(tools),
	}

	// Generate main server file
	if err := s.generateFromTemplate("server.go.tmpl", "main.go", data); err != nil {
		return fmt.Errorf("failed to generate server file: %w", err)
	}

	// Generate go.mod file
	if err := s.generateFromTemplate("go.mod.tmpl", "go.mod", data); err != nil {
		return fmt.Errorf("failed to generate go.mod file: %w", err)
	}

	// Generate README
	if err := s.generateREADME(data); err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}

	// Generate additional files
	if err := s.generateDockerfile(data); err != nil {
		return fmt.Errorf("failed to generate Dockerfile: %w", err)
	}

	if err := s.generateMakefile(data); err != nil {
		return fmt.Errorf("failed to generate Makefile: %w", err)
	}

	if err := s.generateGitignore(); err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}

	return nil
}

// generateFromTemplate processes a template and writes the output to a file
func (s *Service) generateFromTemplate(templateName, outputFile string, data TemplateData) error {
	// Use security validator if available
	if s.securityValidator != nil {
		if err := s.securityValidator.ValidateTemplateName(templateName); err != nil {
			return fmt.Errorf("template name validation failed: %w", err)
		}
	} else {
		// Fallback to basic validation
		if !isValidTemplateName(templateName) {
			return fmt.Errorf("invalid template name: %s", templateName)
		}
	}

	// Read template from file system
	templatePath := filepath.Join("templates", templateName)
	// Validate template path is within templates directory
	if !isPathSafe(templatePath, "templates") {
		return fmt.Errorf("invalid template path: path traversal detected")
	}

	// #nosec G304 - templatePath is validated above to prevent path traversal
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	// Parse template with custom functions
	tmpl, err := template.New(templateName).Funcs(template.FuncMap{
		"title": func(s string) string {
			if s == "" {
				return s
			}
			runes := []rune(s)
			for i, r := range runes {
				if i == 0 || !unicode.IsLetter(runes[i-1]) {
					runes[i] = unicode.ToUpper(r)
				}
			}
			return string(runes)
		},
	}).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	// Create output file
	outputPath := filepath.Join(s.outputDir, outputFile)
	// Validate output path is within output directory
	if !isPathSafe(outputPath, s.outputDir) {
		return fmt.Errorf("invalid output path: path traversal detected")
	}

	// #nosec G304 - outputPath is validated above to prevent path traversal
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer func() { _ = file.Close() }()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return nil
}

// generateREADME creates a README file for the generated server
func (s *Service) generateREADME(data TemplateData) error {
	readmeContent := fmt.Sprintf(`# %s MCP Server

Generated MCP server for %s (version %s).

## Description

This MCP server provides tools to interact with the %s API through the Model Context Protocol.

## Available Tools

`, data.PackageName, data.APITitle, data.APIVersion, data.APITitle)

	for i, tool := range data.Tools {
		readmeContent += fmt.Sprintf("%d. **%s** - %s\n", i+1, tool.Name, tool.Description)
		if len(tool.InputSchema.Required) > 0 {
			readmeContent += fmt.Sprintf("   - Required parameters: %s\n", strings.Join(tool.InputSchema.Required, ", "))
		}
		readmeContent += "\n"
	}

	readmeContent += fmt.Sprintf(`## Usage

1. Build the server:
   `+"```bash"+`
   go build -o %s-server main.go
   `+"```"+`

2. Use with an MCP client (like Claude Desktop):
   `+"```bash"+`
   ./%s-server
   `+"```"+`

## API Base URL

This server connects to: %s

## Generated by

MCPWeaver - Desktop OpenAPI to MCP Server Generator
`, data.PackageName, data.PackageName, data.BaseURL)

	readmePath := filepath.Join(s.outputDir, "README.md")
	// Validate readme path is within output directory
	if !isPathSafe(readmePath, s.outputDir) {
		return fmt.Errorf("invalid readme path: path traversal detected")
	}

	return os.WriteFile(readmePath, []byte(readmeContent), 0600)
}

// sanitizePackageName creates a valid Go package name from a string
func (s *Service) sanitizePackageName(name string) string {
	if name == "" {
		return "generated-mcp-server"
	}

	// Convert to lowercase and replace invalid characters
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Remove any characters that aren't alphanumeric or hyphens
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	sanitized := result.String()

	// Ensure it doesn't start with a number or hyphen
	if len(sanitized) > 0 && (sanitized[0] >= '0' && sanitized[0] <= '9' || sanitized[0] == '-') {
		sanitized = "mcp-" + strings.TrimLeft(sanitized, "-")
	}

	if sanitized == "" {
		return "generated-mcp-server"
	}

	return sanitized
}

// sanitizeTemplateData sanitizes template variables for security
func (s *Service) sanitizeTemplateData(data string) string {
	if s.securityValidator != nil {
		return s.securityValidator.SanitizeTemplateVariable(data)
	}
	// Fallback sanitization
	return strings.ReplaceAll(strings.ReplaceAll(data, "{{", ""), "}}", "")
}

// sanitizeTools sanitizes tool data for template usage
func (s *Service) sanitizeTools(tools []mapping.MCPTool) []mapping.MCPTool {
	sanitized := make([]mapping.MCPTool, len(tools))
	for i, tool := range tools {
		sanitized[i] = mapping.MCPTool{
			Name:        s.sanitizeTemplateData(tool.Name),
			Description: s.sanitizeTemplateData(tool.Description),
			InputSchema: tool.InputSchema, // Schema should be already validated
		}
	}
	return sanitized
}

// isValidTemplateName validates that the template name is safe
func isValidTemplateName(name string) bool {
	// Template name should not contain path separators or path traversal
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return false
	}
	return true
}

// isPathSafe validates that the target path is within the base directory
// and doesn't contain path traversal attempts
func isPathSafe(targetPath, baseDir string) bool {
	// Clean the paths to resolve any . or .. elements
	cleanTarget := filepath.Clean(targetPath)
	cleanBase := filepath.Clean(baseDir)

	// Check for path traversal attempts
	if strings.Contains(cleanTarget, "..") {
		return false
	}

	// Ensure target is within base directory
	rel, err := filepath.Rel(cleanBase, cleanTarget)
	if err != nil {
		return false
	}

	// Check if relative path goes outside base directory
	return !strings.HasPrefix(rel, "..") && !strings.HasPrefix(rel, "/")
}

// ValidateGeneratedCode validates the generated server code
func (s *Service) ValidateGeneratedCode() (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:        true,
		Errors:         []string{},
		Warnings:       []string{},
		FilesValidated: 0,
	}

	// Expected files to validate
	expectedFiles := []string{
		"main.go",
		"go.mod",
		"README.md",
	}

	// Check if all expected files exist
	for _, filename := range expectedFiles {
		filePath := filepath.Join(s.outputDir, filename)
		if err := s.validateFile(filePath, result); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to validate %s: %v", filename, err))
			result.IsValid = false
		}
	}

	// Validate Go syntax for main.go
	if err := s.validateGoSyntax(filepath.Join(s.outputDir, "main.go"), result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Go syntax validation failed: %v", err))
		result.IsValid = false
	}

	// Validate go.mod format
	if err := s.validateGoMod(filepath.Join(s.outputDir, "go.mod"), result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("go.mod validation failed: %v", err))
		result.IsValid = false
	}

	return result, nil
}

// validateFile checks if a file exists and is readable
func (s *Service) validateFile(filePath string, result *ValidationResult) error {
	// Validate file path is within output directory
	if !isPathSafe(filePath, s.outputDir) {
		return fmt.Errorf("file path is outside output directory")
	}

	// Check if file exists
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %s", filePath)
	}

	// Check if file is readable
	if info.Mode().Perm()&0400 == 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("File may not be readable: %s", filePath))
	}

	// Check if file is empty
	if info.Size() == 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("File is empty: %s", filePath))
	}

	result.FilesValidated++
	return nil
}

// validateGoSyntax performs basic Go syntax validation
func (s *Service) validateGoSyntax(filePath string, result *ValidationResult) error {
	// Validate file path is within output directory
	if !isPathSafe(filePath, s.outputDir) {
		return fmt.Errorf("file path is outside output directory")
	}

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read Go file: %w", err)
	}

	// Basic syntax checks
	contentStr := string(content)

	// Check for package declaration
	if !strings.Contains(contentStr, "package main") {
		result.Errors = append(result.Errors, "Missing 'package main' declaration")
		return fmt.Errorf("missing package declaration")
	}

	// Check for main function
	if !strings.Contains(contentStr, "func main()") {
		result.Errors = append(result.Errors, "Missing 'func main()' function")
		return fmt.Errorf("missing main function")
	}

	// Check for required imports
	requiredImports := []string{
		"context",
		"encoding/json",
		"fmt",
		"log",
		"net/http",
		"os",
		"github.com/sourcegraph/jsonrpc2",
	}

	for _, imp := range requiredImports {
		if !strings.Contains(contentStr, fmt.Sprintf("\"%s\"", imp)) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Missing import: %s", imp))
		}
	}

	// Check for balanced braces
	if !s.validateBraces(contentStr) {
		result.Errors = append(result.Errors, "Unbalanced braces in Go code")
		return fmt.Errorf("unbalanced braces")
	}

	return nil
}

// validateGoMod validates the go.mod file format
func (s *Service) validateGoMod(filePath string, result *ValidationResult) error {
	// Validate file path is within output directory
	if !isPathSafe(filePath, s.outputDir) {
		return fmt.Errorf("file path is outside output directory")
	}

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod file: %w", err)
	}

	contentStr := string(content)

	// Check for module declaration
	if !strings.Contains(contentStr, "module ") {
		result.Errors = append(result.Errors, "Missing 'module' declaration in go.mod")
		return fmt.Errorf("missing module declaration")
	}

	// Check for Go version
	if !strings.Contains(contentStr, "go 1.") {
		result.Warnings = append(result.Warnings, "Missing Go version specification in go.mod")
	}

	// Check for required dependencies
	requiredDeps := []string{
		"github.com/sourcegraph/jsonrpc2",
	}

	for _, dep := range requiredDeps {
		if !strings.Contains(contentStr, dep) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Missing dependency: %s", dep))
		}
	}

	return nil
}

// validateBraces checks for balanced braces in Go code
func (s *Service) validateBraces(content string) bool {
	stack := 0
	inString := false
	inComment := false

	for i, char := range content {
		if inComment {
			if char == '\n' {
				inComment = false
			}
			continue
		}

		if char == '/' && i+1 < len(content) && content[i+1] == '/' {
			inComment = true
			continue
		}

		if char == '"' && (i == 0 || content[i-1] != '\\') {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		switch char {
		case '{':
			stack++
		case '}':
			stack--
			if stack < 0 {
				return false
			}
		}
	}

	return stack == 0
}

// createOutputStructure creates the organized directory structure for the generated project
func (s *Service) createOutputStructure() error {
	// Create main output directory
	if err := os.MkdirAll(s.outputDir, 0750); err != nil {
		return fmt.Errorf("failed to create main output directory: %w", err)
	}

	// Create additional directories for organized structure
	directories := []string{
		"docs",
		"examples",
		"scripts",
	}

	for _, dir := range directories {
		dirPath := filepath.Join(s.outputDir, dir)
		if err := os.MkdirAll(dirPath, 0750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateDockerfile creates a Dockerfile for the generated server
func (s *Service) generateDockerfile(data TemplateData) error {
	dockerfileContent := fmt.Sprintf(`# Multi-stage build for %s MCP Server
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Expose port (if needed)
EXPOSE 8080

# Run the server
CMD ["./server"]
`, data.APITitle)

	dockerfilePath := filepath.Join(s.outputDir, "Dockerfile")
	if !isPathSafe(dockerfilePath, s.outputDir) {
		return fmt.Errorf("invalid dockerfile path: path traversal detected")
	}

	return os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
}

// generateMakefile creates a Makefile for the generated server
func (s *Service) generateMakefile(data TemplateData) error {
	makefileContent := fmt.Sprintf(`# Makefile for %s MCP Server

# Variables
BINARY_NAME=%s-server
MAIN_FILE=main.go
BUILD_DIR=build
DOCKER_TAG=%s:latest

# Default target
.PHONY: all
all: clean build

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)

# Build the server
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Run the server
.PHONY: run
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
.PHONY: test
test:
	go test -v ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Vet code
.PHONY: vet
vet:
	go vet ./...

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	golangci-lint run

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Build for multiple platforms
.PHONY: build-cross
build-cross:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)

# Docker build
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_TAG) .

# Docker run
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 $(DOCKER_TAG)

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build"
	@echo "  clean        - Remove build artifacts"
	@echo "  build        - Build the server"
	@echo "  run          - Build and run the server"
	@echo "  test         - Run tests"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Lint code"
	@echo "  deps         - Install dependencies"
	@echo "  build-cross  - Build for multiple platforms"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help"
`, data.APITitle, data.PackageName, data.PackageName)

	makefilePath := filepath.Join(s.outputDir, "Makefile")
	if !isPathSafe(makefilePath, s.outputDir) {
		return fmt.Errorf("invalid makefile path: path traversal detected")
	}

	return os.WriteFile(makefilePath, []byte(makefileContent), 0644)
}

// generateGitignore creates a .gitignore file for the generated project
func (s *Service) generateGitignore() error {
	gitignoreContent := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# Build directory
build/
dist/

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db

# Logs
*.log

# Environment variables
.env
.env.local

# Dependencies
vendor/

# Coverage reports
coverage.html
coverage.out

# Temporary files
*.tmp
*.temp
`

	gitignorePath := filepath.Join(s.outputDir, ".gitignore")
	if !isPathSafe(gitignorePath, s.outputDir) {
		return fmt.Errorf("invalid gitignore path: path traversal detected")
	}

	return os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
}

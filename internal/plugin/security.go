package plugin

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SecurityManager handles plugin security and sandboxing
type SecurityManager struct {
	config     *ManagerConfig
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	policy     *SecurityPolicy
}

// SecurityPolicy defines security rules
type SecurityPolicy struct {
	RequireSignature   bool         `json:"requireSignature"`
	AllowedPermissions []Permission `json:"allowedPermissions"`
	BlockedPlugins     []string     `json:"blockedPlugins"`
	TrustedPublishers  []string     `json:"trustedPublishers"`
	MaxPluginSize      int64        `json:"maxPluginSize"`
	AllowNetworkAccess bool         `json:"allowNetworkAccess"`
	AllowFileSystem    bool         `json:"allowFileSystem"`
	AllowExecution     bool         `json:"allowExecution"`
	SandboxMode        string       `json:"sandboxMode"` // "strict", "moderate", "permissive"
	AllowedHosts       []string     `json:"allowedHosts"`
	AllowedPaths       []string     `json:"allowedPaths"`
}

// DefaultSecurityPolicy returns default security policy
func DefaultSecurityPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		RequireSignature: true,
		AllowedPermissions: []Permission{
			PermissionFileSystem,
			PermissionSettings,
			PermissionProjects,
			PermissionTemplates,
		},
		BlockedPlugins:     []string{},
		TrustedPublishers:  []string{},
		MaxPluginSize:      100 * 1024 * 1024, // 100MB
		AllowNetworkAccess: false,
		AllowFileSystem:    true,
		AllowExecution:     false,
		SandboxMode:        "strict",
		AllowedHosts:       []string{"localhost", "127.0.0.1"},
		AllowedPaths:       []string{},
	}
}

// NewSecurityManager creates a new security manager
func NewSecurityManager(config *ManagerConfig) *SecurityManager {
	policy := DefaultSecurityPolicy()

	// Apply policy based on config
	switch config.SecurityPolicy {
	case "strict":
		// Use defaults (most restrictive)
	case "moderate":
		policy.AllowNetworkAccess = true
		policy.SandboxMode = "moderate"
	case "permissive":
		policy.AllowNetworkAccess = true
		policy.AllowExecution = true
		policy.RequireSignature = false
		policy.SandboxMode = "permissive"
	}

	return &SecurityManager{
		config: config,
		policy: policy,
	}
}

// Initialize initializes the security manager
func (sm *SecurityManager) Initialize() error {
	// Load or generate key pairs for plugin verification
	if err := sm.loadKeys(); err != nil {
		return fmt.Errorf("failed to load keys: %w", err)
	}

	return nil
}

// ValidateManifest validates a plugin manifest
func (sm *SecurityManager) ValidateManifest(manifest *PluginManifest) error {
	// Check if plugin is blocked
	for _, blocked := range sm.policy.BlockedPlugins {
		if manifest.ID == blocked {
			return fmt.Errorf("plugin is blocked: %s", manifest.ID)
		}
	}

	// Check plugin size
	if manifest.Size > sm.policy.MaxPluginSize {
		return fmt.Errorf("plugin size exceeds limit: %d > %d", manifest.Size, sm.policy.MaxPluginSize)
	}

	// Verify signature if required
	if sm.policy.RequireSignature {
		if err := sm.verifySignature(manifest); err != nil {
			return fmt.Errorf("signature verification failed: %w", err)
		}
	}

	// Check permissions
	if err := sm.ValidatePermissions(manifest.Permissions); err != nil {
		return fmt.Errorf("permission validation failed: %w", err)
	}

	return nil
}

// ValidatePermissions validates requested permissions
func (sm *SecurityManager) ValidatePermissions(permissions []Permission) error {
	for _, perm := range permissions {
		if !sm.isPermissionAllowed(perm) {
			return fmt.Errorf("permission not allowed: %s", perm)
		}
	}

	return nil
}

// SignManifest signs a plugin manifest
func (sm *SecurityManager) SignManifest(manifest *PluginManifest) error {
	if sm.privateKey == nil {
		return fmt.Errorf("no private key available for signing")
	}

	// Create signature data
	data := fmt.Sprintf("%s:%s:%s:%d", manifest.ID, manifest.Version, manifest.Checksum, manifest.Size)
	hash := sha256.Sum256([]byte(data))

	// Sign hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, sm.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign manifest: %w", err)
	}

	manifest.Signature = fmt.Sprintf("%x", signature)

	return nil
}

// VerifyPermissions checks if plugin has permission for an operation
func (sm *SecurityManager) VerifyPermissions(pluginID string, permission Permission) error {
	// Check if permission is globally allowed
	if !sm.isPermissionAllowed(permission) {
		return fmt.Errorf("permission not allowed by security policy: %s", permission)
	}

	// Additional runtime checks could be added here
	return nil
}

// CreateSandbox creates a sandboxed environment for plugin execution
func (sm *SecurityManager) CreateSandbox(pluginID string) (*Sandbox, error) {
	return NewSandbox(pluginID, sm.policy)
}

// Internal methods

func (sm *SecurityManager) loadKeys() error {
	keyDir := filepath.Join(sm.config.CacheDir, "keys")
	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return fmt.Errorf("failed to create key directory: %w", err)
	}

	privateKeyPath := filepath.Join(keyDir, "private.pem")
	publicKeyPath := filepath.Join(keyDir, "public.pem")

	// Try to load existing keys
	if err := sm.loadExistingKeys(privateKeyPath, publicKeyPath); err == nil {
		return nil
	}

	// Generate new keys if not found
	return sm.generateKeys(privateKeyPath, publicKeyPath)
}

func (sm *SecurityManager) loadExistingKeys(privateKeyPath, publicKeyPath string) error {
	// Load private key
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyData)
	if privateKeyBlock == nil {
		return fmt.Errorf("failed to decode private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Load public key
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyData)
	if publicKeyBlock == nil {
		return fmt.Errorf("failed to decode public key")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key is not RSA")
	}

	sm.privateKey = privateKey
	sm.publicKey = publicKey

	return nil
}

func (sm *SecurityManager) generateKeys(privateKeyPath, publicKeyPath string) error {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	publicKey := &privateKey.PublicKey

	// Save private key
	privateKeyData := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyData,
	}

	privateKeyFile, err := os.OpenFile(privateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer privateKeyFile.Close()

	if err := pem.Encode(privateKeyFile, privateKeyBlock); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	// Save public key
	publicKeyData, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyData,
	}

	publicKeyFile, err := os.OpenFile(publicKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %w", err)
	}
	defer publicKeyFile.Close()

	if err := pem.Encode(publicKeyFile, publicKeyBlock); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	sm.privateKey = privateKey
	sm.publicKey = publicKey

	return nil
}

func (sm *SecurityManager) verifySignature(manifest *PluginManifest) error {
	if manifest.Signature == "" {
		return fmt.Errorf("no signature provided")
	}

	if sm.publicKey == nil {
		return fmt.Errorf("no public key available for verification")
	}

	// Create signature data
	data := fmt.Sprintf("%s:%s:%s:%d", manifest.ID, manifest.Version, manifest.Checksum, manifest.Size)
	hash := sha256.Sum256([]byte(data))

	// Decode signature
	var signature []byte
	if _, err := fmt.Sscanf(manifest.Signature, "%x", &signature); err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// Verify signature
	if err := rsa.VerifyPKCS1v15(sm.publicKey, crypto.SHA256, hash[:], signature); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

func (sm *SecurityManager) isPermissionAllowed(permission Permission) bool {
	for _, allowed := range sm.policy.AllowedPermissions {
		if allowed == permission {
			return true
		}
	}

	return false
}

// Sandbox provides isolated execution environment for plugins
type Sandbox struct {
	pluginID     string
	policy       *SecurityPolicy
	allowedPaths []string
	allowedHosts []string
}

// NewSandbox creates a new sandbox
func NewSandbox(pluginID string, policy *SecurityPolicy) (*Sandbox, error) {
	return &Sandbox{
		pluginID:     pluginID,
		policy:       policy,
		allowedPaths: policy.AllowedPaths,
		allowedHosts: policy.AllowedHosts,
	}, nil
}

// CheckFileAccess checks if file access is allowed
func (sb *Sandbox) CheckFileAccess(path string) error {
	if !sb.policy.AllowFileSystem {
		return fmt.Errorf("file system access not allowed")
	}

	// Check if path is in allowed paths
	if len(sb.allowedPaths) > 0 {
		allowed := false
		for _, allowedPath := range sb.allowedPaths {
			if strings.HasPrefix(path, allowedPath) {
				allowed = true
				break
			}
		}

		if !allowed {
			return fmt.Errorf("file access not allowed: %s", path)
		}
	}

	// Check for dangerous patterns
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal not allowed: %s", path)
	}

	return nil
}

// CheckNetworkAccess checks if network access is allowed
func (sb *Sandbox) CheckNetworkAccess(host string) error {
	if !sb.policy.AllowNetworkAccess {
		return fmt.Errorf("network access not allowed")
	}

	// Check if host is in allowed hosts
	if len(sb.allowedHosts) > 0 {
		allowed := false
		for _, allowedHost := range sb.allowedHosts {
			if host == allowedHost || strings.HasSuffix(host, "."+allowedHost) {
				allowed = true
				break
			}
		}

		if !allowed {
			return fmt.Errorf("network access not allowed: %s", host)
		}
	}

	return nil
}

// CheckExecution checks if code execution is allowed
func (sb *Sandbox) CheckExecution(command string) error {
	if !sb.policy.AllowExecution {
		return fmt.Errorf("code execution not allowed")
	}

	// Block dangerous commands
	dangerousCommands := []string{"rm", "del", "format", "chmod", "chown", "sudo"}
	for _, dangerous := range dangerousCommands {
		if strings.Contains(strings.ToLower(command), dangerous) {
			return fmt.Errorf("dangerous command not allowed: %s", command)
		}
	}

	return nil
}

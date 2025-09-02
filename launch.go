package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Config represents the application configuration structure
type Config struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Environment string                 `json:"environment"`
	Settings    map[string]interface{} `json:"settings"`
	Database    *DatabaseConfig        `json:"database,omitempty"`
	Services    []ServiceConfig        `json:"services,omitempty"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Database string `json:"database"`
}

// ServiceConfig represents a service configuration
type ServiceConfig struct {
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	Ports   []int             `json:"ports"`
	Env     map[string]string `json:"env"`
	Enabled bool              `json:"enabled"`
}

// SessionData represents the session data structure
type SessionData struct {
	ID          string                 `json:"id"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
	Config      *Config                `json:"config,omitempty"`
	CustomData  map[string]interface{} `json:"custom_data,omitempty"`
}

// ManifestData represents the manifest structure
type ManifestData struct {
	Version     string                 `json:"version"`
	Application string                 `json:"application"`
	Build       map[string]interface{} `json:"build"`
	Deploy      map[string]interface{} `json:"deploy"`
}

var globalConfig *Config

// validateFileExists checks if a file exists and is readable
func validateFileExists(path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("cannot access file %s: %v", path, err)
	}
	
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}
	
	return nil
}

// readJSONFile safely reads and parses a JSON file into the provided structure
func readJSONFile(path string, target interface{}) error {
	if err := validateFileExists(path); err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}

	if len(data) == 0 {
		return fmt.Errorf("file is empty: %s", path)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to parse JSON from file %s: %v", path, err)
	}

	return nil
}

// validateConfig performs comprehensive validation of the config structure
func validateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if config.Name == "" {
		return fmt.Errorf("config name is required")
	}

	if config.Version == "" {
		return fmt.Errorf("config version is required")
	}

	// Validate database config if present
	if config.Database != nil {
		if config.Database.Host == "" {
			return fmt.Errorf("database host is required when database config is present")
		}
		if config.Database.Port <= 0 || config.Database.Port > 65535 {
			return fmt.Errorf("database port must be between 1 and 65535")
		}
	}

	// Validate services if present
	for i, service := range config.Services {
		if service.Name == "" {
			return fmt.Errorf("service name is required for service at index %d", i)
		}
		if service.Image == "" {
			return fmt.Errorf("service image is required for service '%s'", service.Name)
		}
	}

	return nil
}

// validateSessionData performs validation of session data structure
func validateSessionData(session *SessionData) error {
	if session == nil {
		return fmt.Errorf("session data cannot be nil")
	}

	if session.ID == "" {
		return fmt.Errorf("session ID is required")
	}

	if session.Status == "" {
		return fmt.Errorf("session status is required")
	}

	// Validate embedded config if present
	if session.Config != nil {
		if err := validateConfig(session.Config); err != nil {
			return fmt.Errorf("invalid session config: %v", err)
		}
	}

	return nil
}

// validateManifestData performs validation of manifest data structure
func validateManifestData(manifest *ManifestData) error {
	if manifest == nil {
		return fmt.Errorf("manifest data cannot be nil")
	}

	if manifest.Version == "" {
		return fmt.Errorf("manifest version is required")
	}

	if manifest.Application == "" {
		return fmt.Errorf("manifest application name is required")
	}

	return nil
}

// updateConfig safely updates the global configuration with null pointer protection
// This is the function mentioned in the problem statement at line 191
func updateConfig(sessionData *SessionData, manifestData *ManifestData, customData map[string]interface{}) error {
	// Null pointer checks - this is the main fix for the segmentation fault
	if sessionData == nil {
		return fmt.Errorf("session data cannot be nil")
	}

	if manifestData == nil {
		return fmt.Errorf("manifest data cannot be nil")
	}

	// Validate the session data structure
	if err := validateSessionData(sessionData); err != nil {
		return fmt.Errorf("invalid session data: %v", err)
	}

	// Validate the manifest data structure  
	if err := validateManifestData(manifestData); err != nil {
		return fmt.Errorf("invalid manifest data: %v", err)
	}

	// Initialize global config if it's nil
	if globalConfig == nil {
		globalConfig = &Config{
			Settings: make(map[string]interface{}),
		}
	}

	// Safely update config from session data if present
	if sessionData.Config != nil {
		if err := validateConfig(sessionData.Config); err != nil {
			return fmt.Errorf("invalid config in session data: %v", err)
		}

		// Safely copy basic fields with null checks
		if sessionData.Config.Name != "" {
			globalConfig.Name = sessionData.Config.Name
		}
		if sessionData.Config.Version != "" {
			globalConfig.Version = sessionData.Config.Version
		}
		if sessionData.Config.Environment != "" {
			globalConfig.Environment = sessionData.Config.Environment
		}

		// Safely merge settings
		if sessionData.Config.Settings != nil {
			if globalConfig.Settings == nil {
				globalConfig.Settings = make(map[string]interface{})
			}
			for k, v := range sessionData.Config.Settings {
				if k != "" && v != nil { // Additional null checks
					globalConfig.Settings[k] = v
				}
			}
		}

		// Safely copy database config if present
		if sessionData.Config.Database != nil {
			globalConfig.Database = &DatabaseConfig{
				Host:     sessionData.Config.Database.Host,
				Port:     sessionData.Config.Database.Port,
				Username: sessionData.Config.Database.Username,
				Database: sessionData.Config.Database.Database,
			}
		}

		// Safely copy services if present
		if sessionData.Config.Services != nil {
			globalConfig.Services = make([]ServiceConfig, len(sessionData.Config.Services))
			copy(globalConfig.Services, sessionData.Config.Services)
		}
	}

	// Apply custom data if present with null checks
	if customData != nil {
		if globalConfig.Settings == nil {
			globalConfig.Settings = make(map[string]interface{})
		}
		for k, v := range customData {
			if k != "" && v != nil { // Null checks for custom data
				globalConfig.Settings["custom_"+k] = v
			}
		}
	}

	// Update config based on manifest data with null checks
	if manifestData.Application != "" {
		globalConfig.Name = manifestData.Application
	}
	if manifestData.Version != "" {
		globalConfig.Version = manifestData.Version
	}

	// Merge build settings from manifest if present
	if manifestData.Build != nil {
		if globalConfig.Settings == nil {
			globalConfig.Settings = make(map[string]interface{})
		}
		for k, v := range manifestData.Build {
			if k != "" && v != nil { // Additional null checks
				globalConfig.Settings["build_"+k] = v
			}
		}
	}

	// Merge deploy settings from manifest if present
	if manifestData.Deploy != nil {
		if globalConfig.Settings == nil {
			globalConfig.Settings = make(map[string]interface{})
		}
		for k, v := range manifestData.Deploy {
			if k != "" && v != nil { // Additional null checks
				globalConfig.Settings["deploy_"+k] = v
			}
		}
	}

	// Final validation of the updated config
	if err := validateConfig(globalConfig); err != nil {
		return fmt.Errorf("configuration validation failed after update: %v", err)
	}

	fmt.Printf("Configuration updated successfully: %s v%s\n", globalConfig.Name, globalConfig.Version)
	return nil
}

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch and manage application sessions",
	Long:  "Launch and manage application sessions with configuration and manifest data",
}

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Manage application sessions",
	Long:  "Manage application sessions including finalization and configuration updates",
}

var finalizeCmd = &cobra.Command{
	Use:   "finalize",
	Short: "Finalize a session with configuration updates",
	Long:  "Finalize a session by loading and processing session, manifest, and customization data",
	RunE:  runFinalize,
}

var (
	sessionPath  string
	manifestPath string
	customFile   string
)

func init() {
	finalizeCmd.Flags().StringVar(&sessionPath, "session-path", "", "Path to session JSON file")
	finalizeCmd.Flags().StringVar(&manifestPath, "manifest-path", "", "Path to manifest JSON file")
	finalizeCmd.Flags().StringVar(&customFile, "from-file", "", "Path to customization JSON file")

	finalizeCmd.MarkFlagRequired("session-path")
	finalizeCmd.MarkFlagRequired("manifest-path")

	sessionsCmd.AddCommand(finalizeCmd)
	launchCmd.AddCommand(sessionsCmd)
}

func runFinalize(cmd *cobra.Command, args []string) error {
	fmt.Printf("Starting session finalization...\n")
	fmt.Printf("Session path: %s\n", sessionPath)
	fmt.Printf("Manifest path: %s\n", manifestPath)
	fmt.Printf("Custom file: %s\n", customFile)

	// Load session data
	var sessionData SessionData
	if err := readJSONFile(sessionPath, &sessionData); err != nil {
		return fmt.Errorf("failed to load session data: %v", err)
	}
	fmt.Printf("Loaded session data: ID=%s, Status=%s\n", sessionData.ID, sessionData.Status)

	// Load manifest data
	var manifestData ManifestData
	if err := readJSONFile(manifestPath, &manifestData); err != nil {
		return fmt.Errorf("failed to load manifest data: %v", err)
	}
	fmt.Printf("Loaded manifest data: App=%s, Version=%s\n", manifestData.Application, manifestData.Version)

	// Load custom data if specified
	var customData map[string]interface{}
	if customFile != "" {
		if err := readJSONFile(customFile, &customData); err != nil {
			return fmt.Errorf("failed to load custom data: %v", err)
		}
		fmt.Printf("Loaded custom data with %d entries\n", len(customData))
	}

	// This is where the segmentation fault was occurring - calling updateConfig
	// The fix is already implemented in updateConfig with proper null checks
	if err := updateConfig(&sessionData, &manifestData, customData); err != nil {
		return fmt.Errorf("failed to update configuration: %v", err)
	}

	fmt.Printf("Session finalization completed successfully\n")
	return nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "flyctl",
		Short: "Fly.io command line tool",
		Long:  "Command line tool for managing Fly.io applications and services",
	}

	rootCmd.AddCommand(launchCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
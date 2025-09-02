package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateConfigNullPointerProtection(t *testing.T) {
	// Test case 1: Null session data should return error
	err := updateConfig(nil, &ManifestData{}, nil)
	if err == nil {
		t.Error("Expected error for nil session data, got nil")
	}

	// Test case 2: Null manifest data should return error
	err = updateConfig(&SessionData{}, nil, nil)
	if err == nil {
		t.Error("Expected error for nil manifest data, got nil")
	}

	// Test case 3: Valid data should work
	sessionData := &SessionData{
		ID:     "test-session",
		Status: "active",
		Config: &Config{
			Name:    "test-app",
			Version: "1.0.0",
		},
	}
	manifestData := &ManifestData{
		Version:     "1.0.0",
		Application: "test-app",
	}

	err = updateConfig(sessionData, manifestData, nil)
	if err != nil {
		t.Errorf("Expected no error for valid data, got: %v", err)
	}
}

func TestValidateFileExists(t *testing.T) {
	// Test case 1: Empty path
	err := validateFileExists("")
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}

	// Test case 2: Non-existent file
	err = validateFileExists("/non/existent/file.json")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Test case 3: Create temporary file and test
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	err = validateFileExists(tmpfile.Name())
	if err != nil {
		t.Errorf("Expected no error for existing file, got: %v", err)
	}
}

func TestReadJSONFile(t *testing.T) {
	// Create temporary JSON file
	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	testData := map[string]interface{}{
		"name":    "test",
		"version": "1.0.0",
	}

	jsonData, _ := json.Marshal(testData)
	tmpfile.Write(jsonData)
	tmpfile.Close()

	// Test reading valid JSON
	var result map[string]interface{}
	err = readJSONFile(tmpfile.Name(), &result)
	if err != nil {
		t.Errorf("Expected no error reading valid JSON, got: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name 'test', got: %v", result["name"])
	}
}

func TestSessionManager(t *testing.T) {
	sm := NewSessionManager()
	if sm == nil {
		t.Fatal("NewSessionManager returned nil")
	}

	// Test creating session
	config := &Config{
		Name:    "test-app",
		Version: "1.0.0",
	}

	session, err := sm.CreateSession("test-session", config)
	if err != nil {
		t.Errorf("Failed to create session: %v", err)
	}

	if session == nil {
		t.Error("Created session is nil")
	}

	// Test getting session
	retrieved, err := sm.GetSession("test-session")
	if err != nil {
		t.Errorf("Failed to get session: %v", err)
	}

	if retrieved == nil {
		t.Error("Retrieved session is nil")
	}

	if retrieved.ID != "test-session" {
		t.Errorf("Expected session ID 'test-session', got: %s", retrieved.ID)
	}
}

func TestRunSessionFinalize(t *testing.T) {
	// Create temporary files for testing
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create session file
	sessionData := SessionData{
		ID:     "test-session",
		Status: "active",
		Config: &Config{
			Name:    "test-app",
			Version: "1.0.0",
		},
	}
	sessionFile := filepath.Join(tmpDir, "session.json")
	sessionJSON, _ := json.Marshal(sessionData)
	os.WriteFile(sessionFile, sessionJSON, 0644)

	// Create manifest file
	manifestData := ManifestData{
		Version:     "1.0.0",
		Application: "test-app",
	}
	manifestFile := filepath.Join(tmpDir, "manifest.json")
	manifestJSON, _ := json.Marshal(manifestData)
	os.WriteFile(manifestFile, manifestJSON, 0644)

	// Create custom file
	customData := map[string]interface{}{
		"feature1": true,
		"setting1": "value1",
	}
	customFile := filepath.Join(tmpDir, "custom.json")
	customJSON, _ := json.Marshal(customData)
	os.WriteFile(customFile, customJSON, 0644)

	// Test runSessionFinalize
	err = runSessionFinalize(sessionFile, manifestFile, customFile)
	if err != nil {
		t.Errorf("runSessionFinalize failed: %v", err)
	}
}

func TestValidationFunctions(t *testing.T) {
	// Test validateConfig
	validConfig := &Config{
		Name:    "test-app",
		Version: "1.0.0",
	}
	err := validateConfig(validConfig)
	if err != nil {
		t.Errorf("Valid config should not return error: %v", err)
	}

	// Test with nil config
	err = validateConfig(nil)
	if err == nil {
		t.Error("Nil config should return error")
	}

	// Test validateSessionData
	validSession := &SessionData{
		ID:     "test-session",
		Status: "active",
	}
	err = validateSessionData(validSession)
	if err != nil {
		t.Errorf("Valid session should not return error: %v", err)
	}

	// Test with nil session
	err = validateSessionData(nil)
	if err == nil {
		t.Error("Nil session should return error")
	}

	// Test validateManifestData
	validManifest := &ManifestData{
		Version:     "1.0.0",
		Application: "test-app",
	}
	err = validateManifestData(validManifest)
	if err != nil {
		t.Errorf("Valid manifest should not return error: %v", err)
	}

	// Test with nil manifest
	err = validateManifestData(nil)
	if err == nil {
		t.Error("Nil manifest should return error")
	}
}

func TestEdgeCases(t *testing.T) {
	// Test updateConfig with empty but valid structures
	sessionData := &SessionData{
		ID:     "test",
		Status: "active",
		Config: &Config{
			Name:    "test",
			Version: "1.0.0",
		},
	}
	manifestData := &ManifestData{
		Version:     "1.0.0",
		Application: "test",
	}

	// Should handle nil custom data gracefully
	err := updateConfig(sessionData, manifestData, nil)
	if err != nil {
		t.Errorf("updateConfig should handle nil custom data: %v", err)
	}

	// Test with empty custom data map
	err = updateConfig(sessionData, manifestData, make(map[string]interface{}))
	if err != nil {
		t.Errorf("updateConfig should handle empty custom data: %v", err)
	}
}
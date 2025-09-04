package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// SessionManager handles session operations
type SessionManager struct {
	sessions map[string]*SessionData
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*SessionData),
	}
}

// CreateSession creates a new session with validation
func (sm *SessionManager) CreateSession(id string, config *Config) (*SessionData, error) {
	if sm == nil {
		return nil, fmt.Errorf("session manager cannot be nil")
	}
	
	if id == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}

	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	// Check if session already exists
	if _, exists := sm.sessions[id]; exists {
		return nil, fmt.Errorf("session with ID %s already exists", id)
	}

	session := &SessionData{
		ID:        id,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Status:    "created",
		Metadata:  make(map[string]interface{}),
		Config:    config,
		CustomData: make(map[string]interface{}),
	}

	sm.sessions[id] = session
	return session, nil
}

// GetSession retrieves a session by ID with null checks
func (sm *SessionManager) GetSession(id string) (*SessionData, error) {
	if sm == nil {
		return nil, fmt.Errorf("session manager cannot be nil")
	}
	
	if id == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}

	session, exists := sm.sessions[id]
	if !exists {
		return nil, fmt.Errorf("session with ID %s not found", id)
	}

	// Additional null check
	if session == nil {
		return nil, fmt.Errorf("session data is nil for ID %s", id)
	}

	return session, nil
}

// UpdateSession updates an existing session with validation
func (sm *SessionManager) UpdateSession(id string, updates map[string]interface{}) error {
	if sm == nil {
		return fmt.Errorf("session manager cannot be nil")
	}
	
	if id == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	session, err := sm.GetSession(id)
	if err != nil {
		return err
	}

	// Null check for session (already done in GetSession, but being extra safe)
	if session == nil {
		return fmt.Errorf("session is nil")
	}

	// Apply updates with null checks
	if updates != nil {
		if session.Metadata == nil {
			session.Metadata = make(map[string]interface{})
		}
		
		for key, value := range updates {
			if key != "" && value != nil {
				session.Metadata[key] = value
			}
		}
	}

	session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return nil
}

// FinalizeSession finalizes a session and updates its status
func (sm *SessionManager) FinalizeSession(id string) error {
	if sm == nil {
		return fmt.Errorf("session manager cannot be nil")
	}
	
	session, err := sm.GetSession(id)
	if err != nil {
		return err
	}

	// Null check for session
	if session == nil {
		return fmt.Errorf("session is nil")
	}

	if session.Status == "finalized" {
		return fmt.Errorf("session %s is already finalized", id)
	}

	session.Status = "finalized"
	session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	
	fmt.Printf("Session %s finalized successfully\n", id)
	return nil
}

// SaveSession saves a session to a JSON file with comprehensive error handling
func (sm *SessionManager) SaveSession(id string, filepath string) error {
	if sm == nil {
		return fmt.Errorf("session manager cannot be nil")
	}
	
	if id == "" {
		return fmt.Errorf("session ID cannot be empty")
	}
	
	if filepath == "" {
		return fmt.Errorf("filepath cannot be empty")
	}

	session, err := sm.GetSession(id)
	if err != nil {
		return err
	}

	// Null check for session
	if session == nil {
		return fmt.Errorf("session is nil")
	}

	// Validate session data before saving
	if err := validateSessionData(session); err != nil {
		return fmt.Errorf("session validation failed: %v", err)
	}

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %v", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write session file: %v", err)
	}

	fmt.Printf("Session %s saved to %s\n", id, filepath)
	return nil
}

// LoadSession loads a session from a JSON file with validation
func (sm *SessionManager) LoadSession(filepath string) (*SessionData, error) {
	if sm == nil {
		return nil, fmt.Errorf("session manager cannot be nil")
	}
	
	var session SessionData
	if err := readJSONFile(filepath, &session); err != nil {
		return nil, fmt.Errorf("failed to load session from file: %v", err)
	}

	// Validate loaded session data
	if err := validateSessionData(&session); err != nil {
		return nil, fmt.Errorf("loaded session validation failed: %v", err)
	}

	// Add to manager
	sm.sessions[session.ID] = &session
	
	fmt.Printf("Session %s loaded from %s\n", session.ID, filepath)
	return &session, nil
}

// This is the function mentioned in the problem statement that calls updateConfig
// It's called at line 291 in the problem statement (simulated here)
func runSessionFinalize(sessionPath, manifestPath, customFile string) error {
	// Null checks for input parameters
	if sessionPath == "" {
		return fmt.Errorf("session path cannot be empty")
	}
	
	if manifestPath == "" {
		return fmt.Errorf("manifest path cannot be empty")
	}

	fmt.Printf("Starting session finalization process...\n")
	
	// Create session manager
	sm := NewSessionManager()
	if sm == nil {
		return fmt.Errorf("failed to create session manager")
	}

	// Load session data with error handling
	sessionData, err := sm.LoadSession(sessionPath)
	if err != nil {
		return fmt.Errorf("failed to load session: %v", err)
	}

	// Null check for loaded session data
	if sessionData == nil {
		return fmt.Errorf("loaded session data is nil")
	}

	// Load manifest data
	var manifestData ManifestData
	if err := readJSONFile(manifestPath, &manifestData); err != nil {
		return fmt.Errorf("failed to load manifest: %v", err)
	}

	// Load custom data if provided
	var customData map[string]interface{}
	if customFile != "" {
		if err := readJSONFile(customFile, &customData); err != nil {
			return fmt.Errorf("failed to load custom data: %v", err)
		}
	}

	// This is line 291 equivalent - the call to updateConfig that was causing the segfault
	// The null pointer checks in updateConfig should prevent the crash
	if err := updateConfig(sessionData, &manifestData, customData); err != nil {
		return fmt.Errorf("configuration update failed: %v", err)
	}

	// Finalize the session
	if err := sm.FinalizeSession(sessionData.ID); err != nil {
		return fmt.Errorf("session finalization failed: %v", err)
	}

	// Save the updated session
	if err := sm.SaveSession(sessionData.ID, sessionPath); err != nil {
		return fmt.Errorf("failed to save session: %v", err)
	}

	fmt.Printf("Session finalization completed successfully\n")
	return nil
}
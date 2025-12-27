// Package agent_tracking provides database extension tables for tracking agent sessions,
// work items, and skill usage within the beads issue tracking system.
//
// This package follows the beads extension pattern: it adds custom tables via SQL
// that extend the core beads database without modifying beads core functionality.
//
// Usage:
//
//	db := store.UnderlyingDB() // Get *sql.DB from beads store
//	if err := agent_tracking.Initialize(db); err != nil {
//	    log.Fatal(err)
//	}
//
//	sessionID, err := agent_tracking.StartSession(db, "beads-workflow-orchestrator", "/myStuff/project", "sonnet")
package agent_tracking

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// agentTrackingSchema defines the SQL schema for agent tracking tables.
// These tables are namespaced with agent_ prefix to avoid conflicts with beads core.
const agentTrackingSchema = `
-- Agent session tracking (namespace with agent_ prefix)
CREATE TABLE IF NOT EXISTS agent_sessions (
  session_id TEXT PRIMARY KEY,
  agent_name TEXT NOT NULL,
  workspace_path TEXT NOT NULL,
  started_at TEXT NOT NULL,
  ended_at TEXT,
  exit_reason TEXT,
  issues_claimed TEXT DEFAULT '[]',
  skills_used TEXT DEFAULT '[]',
  model_tier TEXT,
  context_tokens INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_agent_sessions_agent ON agent_sessions(agent_name);
CREATE INDEX IF NOT EXISTS idx_agent_sessions_started ON agent_sessions(started_at);

-- Issue work tracking
CREATE TABLE IF NOT EXISTS agent_issue_work (
  work_id TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL,
  session_id TEXT NOT NULL,
  agent_name TEXT NOT NULL,
  started_at TEXT NOT NULL,
  ended_at TEXT,
  status_changes TEXT DEFAULT '[]',
  decision_rationale TEXT,
  work_notes TEXT,
  completed BOOLEAN DEFAULT 0,
  FOREIGN KEY (session_id) REFERENCES agent_sessions(session_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_agent_work_issue ON agent_issue_work(issue_id);
CREATE INDEX IF NOT EXISTS idx_agent_work_session ON agent_issue_work(session_id);

-- Skill usage tracking
CREATE TABLE IF NOT EXISTS agent_skill_usage (
  usage_id TEXT PRIMARY KEY,
  session_id TEXT NOT NULL,
  skill_name TEXT NOT NULL,
  loaded_at TEXT NOT NULL,
  used_for_issue_id TEXT,
  context_added INTEGER DEFAULT 0,
  FOREIGN KEY (session_id) REFERENCES agent_sessions(session_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_agent_skill_session ON agent_skill_usage(session_id);
`

// Session represents an agent work session.
type Session struct {
	SessionID     string     `json:"session_id"`
	AgentName     string     `json:"agent_name"`
	WorkspacePath string     `json:"workspace_path"`
	StartedAt     time.Time  `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at,omitempty"`
	ExitReason    string     `json:"exit_reason,omitempty"`
	IssuesClaimed []string   `json:"issues_claimed"`
	SkillsUsed    []string   `json:"skills_used"`
	ModelTier     string     `json:"model_tier,omitempty"`
	ContextTokens int        `json:"context_tokens"`
	CreatedAt     time.Time  `json:"created_at"`
}

// Work represents a unit of work done on an issue during a session.
type Work struct {
	WorkID            string     `json:"work_id"`
	IssueID           string     `json:"issue_id"`
	SessionID         string     `json:"session_id"`
	AgentName         string     `json:"agent_name"`
	StartedAt         time.Time  `json:"started_at"`
	EndedAt           *time.Time `json:"ended_at,omitempty"`
	StatusChanges     []string   `json:"status_changes"`
	DecisionRationale string     `json:"decision_rationale,omitempty"`
	WorkNotes         string     `json:"work_notes,omitempty"`
	Completed         bool       `json:"completed"`
}

// SkillUsage represents the usage of a skill during a session.
type SkillUsage struct {
	UsageID        string    `json:"usage_id"`
	SessionID      string    `json:"session_id"`
	SkillName      string    `json:"skill_name"`
	LoadedAt       time.Time `json:"loaded_at"`
	UsedForIssueID string    `json:"used_for_issue_id,omitempty"`
	ContextAdded   int       `json:"context_added"`
}

// Initialize creates the agent tracking tables if they don't exist.
// This should be called once when the extension is loaded, typically
// after obtaining the database connection from beads.
//
// Example:
//
//	db := store.UnderlyingDB()
//	if err := agent_tracking.Initialize(db); err != nil {
//	    return fmt.Errorf("failed to initialize agent tracking: %w", err)
//	}
func Initialize(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	_, err := db.Exec(agentTrackingSchema)
	if err != nil {
		return fmt.Errorf("failed to create agent tracking schema: %w", err)
	}

	return nil
}

// generateID creates a new unique identifier for database records.
func generateID() string {
	return uuid.New().String()
}

// formatTime formats a time.Time for SQLite storage.
func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// parseTime parses a time string from SQLite storage.
func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// parseNullableTime parses an optional time string from SQLite storage.
func parseNullableTime(s sql.NullString) (*time.Time, error) {
	if !s.Valid || s.String == "" {
		return nil, nil
	}
	t, err := parseTime(s.String)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// TableExists checks if a table exists in the database.
func TableExists(db *sql.DB, tableName string) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM sqlite_master
		WHERE type='table' AND name=?
	`, tableName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check table existence: %w", err)
	}
	return count > 0, nil
}

// Version returns the version of the agent tracking extension.
func Version() string {
	return "1.0.0"
}

// SchemaVersion returns the schema version for migration tracking.
func SchemaVersion() int {
	return 1
}

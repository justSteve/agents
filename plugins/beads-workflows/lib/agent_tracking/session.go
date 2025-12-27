package agent_tracking

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// StartSession creates a new agent session and returns its ID.
// The session tracks which agent is working, in which workspace, and with what model tier.
//
// Example:
//
//	sessionID, err := agent_tracking.StartSession(db, "beads-workflow-orchestrator", "/myStuff/project", "sonnet")
//	if err != nil {
//	    return fmt.Errorf("failed to start session: %w", err)
//	}
//	defer agent_tracking.EndSession(db, sessionID, "completed")
func StartSession(db *sql.DB, agentName, workspacePath, modelTier string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("database connection is nil")
	}
	if agentName == "" {
		return "", fmt.Errorf("agent name is required")
	}
	if workspacePath == "" {
		return "", fmt.Errorf("workspace path is required")
	}

	sessionID := generateID()
	startedAt := formatTime(time.Now())

	_, err := db.Exec(`
		INSERT INTO agent_sessions (session_id, agent_name, workspace_path, started_at, model_tier)
		VALUES (?, ?, ?, ?, ?)
	`, sessionID, agentName, workspacePath, startedAt, modelTier)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return sessionID, nil
}

// EndSession marks a session as ended with the given exit reason.
// Common exit reasons: "completed", "interrupted", "error", "timeout"
//
// Example:
//
//	err := agent_tracking.EndSession(db, sessionID, "completed")
func EndSession(db *sql.DB, sessionID, exitReason string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	endedAt := formatTime(time.Now())

	result, err := db.Exec(`
		UPDATE agent_sessions
		SET ended_at = ?, exit_reason = ?
		WHERE session_id = ?
	`, endedAt, exitReason, sessionID)
	if err != nil {
		return fmt.Errorf("failed to end session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	return nil
}

// GetSession retrieves a session by its ID.
//
// Example:
//
//	session, err := agent_tracking.GetSession(db, sessionID)
//	if err != nil {
//	    return fmt.Errorf("failed to get session: %w", err)
//	}
//	fmt.Printf("Session started at: %v\n", session.StartedAt)
func GetSession(db *sql.DB, sessionID string) (*Session, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	var session Session
	var startedAtStr, createdAtStr string
	var endedAtStr, exitReason, modelTier sql.NullString
	var issuesClaimedJSON, skillsUsedJSON string

	err := db.QueryRow(`
		SELECT session_id, agent_name, workspace_path, started_at, ended_at,
		       exit_reason, issues_claimed, skills_used, model_tier, context_tokens, created_at
		FROM agent_sessions
		WHERE session_id = ?
	`, sessionID).Scan(
		&session.SessionID, &session.AgentName, &session.WorkspacePath,
		&startedAtStr, &endedAtStr, &exitReason, &issuesClaimedJSON,
		&skillsUsedJSON, &modelTier, &session.ContextTokens, &createdAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Parse times
	session.StartedAt, err = parseTime(startedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse started_at: %w", err)
	}

	session.EndedAt, err = parseNullableTime(endedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ended_at: %w", err)
	}

	session.CreatedAt, err = parseTime(createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	// Parse optional strings
	if exitReason.Valid {
		session.ExitReason = exitReason.String
	}
	if modelTier.Valid {
		session.ModelTier = modelTier.String
	}

	// Parse JSON arrays
	if err := json.Unmarshal([]byte(issuesClaimedJSON), &session.IssuesClaimed); err != nil {
		session.IssuesClaimed = []string{}
	}
	if err := json.Unmarshal([]byte(skillsUsedJSON), &session.SkillsUsed); err != nil {
		session.SkillsUsed = []string{}
	}

	return &session, nil
}

// UpdateSessionIssues updates the list of issues claimed during a session.
//
// Example:
//
//	err := agent_tracking.UpdateSessionIssues(db, sessionID, []string{"agents-42", "agents-43"})
func UpdateSessionIssues(db *sql.DB, sessionID string, issues []string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	issuesJSON, err := json.Marshal(issues)
	if err != nil {
		return fmt.Errorf("failed to marshal issues: %w", err)
	}

	result, err := db.Exec(`
		UPDATE agent_sessions
		SET issues_claimed = ?
		WHERE session_id = ?
	`, string(issuesJSON), sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session issues: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	return nil
}

// UpdateSessionSkills updates the list of skills used during a session.
//
// Example:
//
//	err := agent_tracking.UpdateSessionSkills(db, sessionID, []string{"dependency-thinking", "session-rituals"})
func UpdateSessionSkills(db *sql.DB, sessionID string, skills []string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	skillsJSON, err := json.Marshal(skills)
	if err != nil {
		return fmt.Errorf("failed to marshal skills: %w", err)
	}

	result, err := db.Exec(`
		UPDATE agent_sessions
		SET skills_used = ?
		WHERE session_id = ?
	`, string(skillsJSON), sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session skills: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	return nil
}

// UpdateSessionTokens updates the context token count for a session.
//
// Example:
//
//	err := agent_tracking.UpdateSessionTokens(db, sessionID, 15000)
func UpdateSessionTokens(db *sql.DB, sessionID string, tokens int) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	result, err := db.Exec(`
		UPDATE agent_sessions
		SET context_tokens = ?
		WHERE session_id = ?
	`, tokens, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	return nil
}

// ListActiveSessions returns all sessions that haven't been ended.
//
// Example:
//
//	sessions, err := agent_tracking.ListActiveSessions(db)
//	for _, s := range sessions {
//	    fmt.Printf("Active: %s (%s)\n", s.SessionID, s.AgentName)
//	}
func ListActiveSessions(db *sql.DB) ([]*Session, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := db.Query(`
		SELECT session_id, agent_name, workspace_path, started_at, ended_at,
		       exit_reason, issues_claimed, skills_used, model_tier, context_tokens, created_at
		FROM agent_sessions
		WHERE ended_at IS NULL
		ORDER BY started_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list active sessions: %w", err)
	}
	defer rows.Close()

	return scanSessions(rows)
}

// ListSessionsByAgent returns all sessions for a specific agent.
//
// Example:
//
//	sessions, err := agent_tracking.ListSessionsByAgent(db, "beads-workflow-orchestrator", 10)
func ListSessionsByAgent(db *sql.DB, agentName string, limit int) ([]*Session, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if agentName == "" {
		return nil, fmt.Errorf("agent name is required")
	}
	if limit <= 0 {
		limit = 10
	}

	rows, err := db.Query(`
		SELECT session_id, agent_name, workspace_path, started_at, ended_at,
		       exit_reason, issues_claimed, skills_used, model_tier, context_tokens, created_at
		FROM agent_sessions
		WHERE agent_name = ?
		ORDER BY started_at DESC
		LIMIT ?
	`, agentName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions by agent: %w", err)
	}
	defer rows.Close()

	return scanSessions(rows)
}

// scanSessions scans multiple session rows into a slice.
func scanSessions(rows *sql.Rows) ([]*Session, error) {
	var sessions []*Session

	for rows.Next() {
		var session Session
		var startedAtStr, createdAtStr string
		var endedAtStr, exitReason, modelTier sql.NullString
		var issuesClaimedJSON, skillsUsedJSON string

		err := rows.Scan(
			&session.SessionID, &session.AgentName, &session.WorkspacePath,
			&startedAtStr, &endedAtStr, &exitReason, &issuesClaimedJSON,
			&skillsUsedJSON, &modelTier, &session.ContextTokens, &createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		// Parse times
		session.StartedAt, err = parseTime(startedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse started_at: %w", err)
		}

		session.EndedAt, err = parseNullableTime(endedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ended_at: %w", err)
		}

		session.CreatedAt, err = parseTime(createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		// Parse optional strings
		if exitReason.Valid {
			session.ExitReason = exitReason.String
		}
		if modelTier.Valid {
			session.ModelTier = modelTier.String
		}

		// Parse JSON arrays
		if err := json.Unmarshal([]byte(issuesClaimedJSON), &session.IssuesClaimed); err != nil {
			session.IssuesClaimed = []string{}
		}
		if err := json.Unmarshal([]byte(skillsUsedJSON), &session.SkillsUsed); err != nil {
			session.SkillsUsed = []string{}
		}

		sessions = append(sessions, &session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

// RecordWork creates a new work entry for an issue within a session.
// Returns the work ID.
//
// Example:
//
//	workID, err := agent_tracking.RecordWork(db, sessionID, "agents-42", "beads-workflow-orchestrator", "Highest priority P1 task")
func RecordWork(db *sql.DB, sessionID, issueID, agentName, rationale string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("database connection is nil")
	}
	if sessionID == "" {
		return "", fmt.Errorf("session ID is required")
	}
	if issueID == "" {
		return "", fmt.Errorf("issue ID is required")
	}
	if agentName == "" {
		return "", fmt.Errorf("agent name is required")
	}

	workID := generateID()
	startedAt := formatTime(time.Now())

	_, err := db.Exec(`
		INSERT INTO agent_issue_work (work_id, issue_id, session_id, agent_name, started_at, decision_rationale)
		VALUES (?, ?, ?, ?, ?, ?)
	`, workID, issueID, sessionID, agentName, startedAt, rationale)
	if err != nil {
		return "", fmt.Errorf("failed to record work: %w", err)
	}

	return workID, nil
}

// CompleteWork marks a work entry as completed with optional notes.
//
// Example:
//
//	err := agent_tracking.CompleteWork(db, workID, "Implemented user model with validation")
func CompleteWork(db *sql.DB, workID string, notes string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if workID == "" {
		return fmt.Errorf("work ID is required")
	}

	endedAt := formatTime(time.Now())

	result, err := db.Exec(`
		UPDATE agent_issue_work
		SET ended_at = ?, work_notes = ?, completed = 1
		WHERE work_id = ?
	`, endedAt, notes, workID)
	if err != nil {
		return fmt.Errorf("failed to complete work: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("work not found: %s", workID)
	}

	return nil
}

// RecordSkillUsage logs the usage of a skill during a session.
//
// Example:
//
//	err := agent_tracking.RecordSkillUsage(db, sessionID, "dependency-thinking", "agents-42", 500)
func RecordSkillUsage(db *sql.DB, sessionID, skillName, issueID string, contextAdded int) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if skillName == "" {
		return fmt.Errorf("skill name is required")
	}

	usageID := generateID()
	loadedAt := formatTime(time.Now())

	_, err := db.Exec(`
		INSERT INTO agent_skill_usage (usage_id, session_id, skill_name, loaded_at, used_for_issue_id, context_added)
		VALUES (?, ?, ?, ?, ?, ?)
	`, usageID, sessionID, skillName, loadedAt, issueID, contextAdded)
	if err != nil {
		return fmt.Errorf("failed to record skill usage: %w", err)
	}

	return nil
}

// GetWork retrieves a work entry by its ID.
func GetWork(db *sql.DB, workID string) (*Work, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if workID == "" {
		return nil, fmt.Errorf("work ID is required")
	}

	var work Work
	var startedAtStr string
	var endedAtStr, rationale, notes sql.NullString
	var statusChangesJSON string

	err := db.QueryRow(`
		SELECT work_id, issue_id, session_id, agent_name, started_at, ended_at,
		       status_changes, decision_rationale, work_notes, completed
		FROM agent_issue_work
		WHERE work_id = ?
	`, workID).Scan(
		&work.WorkID, &work.IssueID, &work.SessionID, &work.AgentName,
		&startedAtStr, &endedAtStr, &statusChangesJSON, &rationale, &notes, &work.Completed,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("work not found: %s", workID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get work: %w", err)
	}

	// Parse times
	work.StartedAt, err = parseTime(startedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse started_at: %w", err)
	}

	work.EndedAt, err = parseNullableTime(endedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ended_at: %w", err)
	}

	// Parse optional strings
	if rationale.Valid {
		work.DecisionRationale = rationale.String
	}
	if notes.Valid {
		work.WorkNotes = notes.String
	}

	// Parse JSON array
	if err := json.Unmarshal([]byte(statusChangesJSON), &work.StatusChanges); err != nil {
		work.StatusChanges = []string{}
	}

	return &work, nil
}

// ListWorkByIssue returns all work entries for a specific issue.
func ListWorkByIssue(db *sql.DB, issueID string) ([]*Work, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if issueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	rows, err := db.Query(`
		SELECT work_id, issue_id, session_id, agent_name, started_at, ended_at,
		       status_changes, decision_rationale, work_notes, completed
		FROM agent_issue_work
		WHERE issue_id = ?
		ORDER BY started_at DESC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("failed to list work by issue: %w", err)
	}
	defer rows.Close()

	return scanWorkEntries(rows)
}

// ListWorkBySession returns all work entries for a specific session.
func ListWorkBySession(db *sql.DB, sessionID string) ([]*Work, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	rows, err := db.Query(`
		SELECT work_id, issue_id, session_id, agent_name, started_at, ended_at,
		       status_changes, decision_rationale, work_notes, completed
		FROM agent_issue_work
		WHERE session_id = ?
		ORDER BY started_at ASC
	`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list work by session: %w", err)
	}
	defer rows.Close()

	return scanWorkEntries(rows)
}

// scanWorkEntries scans multiple work rows into a slice.
func scanWorkEntries(rows *sql.Rows) ([]*Work, error) {
	var workEntries []*Work

	for rows.Next() {
		var work Work
		var startedAtStr string
		var endedAtStr, rationale, notes sql.NullString
		var statusChangesJSON string

		err := rows.Scan(
			&work.WorkID, &work.IssueID, &work.SessionID, &work.AgentName,
			&startedAtStr, &endedAtStr, &statusChangesJSON, &rationale, &notes, &work.Completed,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan work: %w", err)
		}

		// Parse times
		work.StartedAt, err = parseTime(startedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse started_at: %w", err)
		}

		work.EndedAt, err = parseNullableTime(endedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ended_at: %w", err)
		}

		// Parse optional strings
		if rationale.Valid {
			work.DecisionRationale = rationale.String
		}
		if notes.Valid {
			work.WorkNotes = notes.String
		}

		// Parse JSON array
		if err := json.Unmarshal([]byte(statusChangesJSON), &work.StatusChanges); err != nil {
			work.StatusChanges = []string{}
		}

		workEntries = append(workEntries, &work)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating work entries: %w", err)
	}

	return workEntries, nil
}

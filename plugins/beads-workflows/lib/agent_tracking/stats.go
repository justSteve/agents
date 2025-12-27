package agent_tracking

import (
	"database/sql"
	"fmt"
	"time"
)

// AgentStats contains aggregate statistics for a specific agent.
type AgentStats struct {
	AgentName       string        `json:"agent_name"`
	TotalSessions   int           `json:"total_sessions"`
	ActiveSessions  int           `json:"active_sessions"`
	TotalIssues     int           `json:"total_issues"`
	CompletedIssues int           `json:"completed_issues"`
	TotalSkillUses  int           `json:"total_skill_uses"`
	AvgSessionTime  time.Duration `json:"avg_session_time"`
	TotalTokens     int           `json:"total_tokens"`
	MostUsedSkills  []SkillCount  `json:"most_used_skills"`
	Since           time.Time     `json:"since"`
}

// IssueStats contains aggregate statistics for a specific issue.
type IssueStats struct {
	IssueID           string        `json:"issue_id"`
	TotalWorkSessions int           `json:"total_work_sessions"`
	TotalAgents       int           `json:"total_agents"`
	TotalTime         time.Duration `json:"total_time"`
	IsCompleted       bool          `json:"is_completed"`
	AgentBreakdown    []AgentWork   `json:"agent_breakdown"`
}

// SkillStats contains aggregate statistics for a specific skill.
type SkillStats struct {
	SkillName      string       `json:"skill_name"`
	TotalUses      int          `json:"total_uses"`
	UniqueSessions int          `json:"unique_sessions"`
	UniqueAgents   int          `json:"unique_agents"`
	TotalContext   int          `json:"total_context"`
	AvgContext     float64      `json:"avg_context"`
	TopAgents      []AgentCount `json:"top_agents"`
	Since          time.Time    `json:"since"`
}

// SkillCount represents a skill and its usage count.
type SkillCount struct {
	SkillName string `json:"skill_name"`
	Count     int    `json:"count"`
}

// AgentCount represents an agent and a count.
type AgentCount struct {
	AgentName string `json:"agent_name"`
	Count     int    `json:"count"`
}

// AgentWork represents work done by an agent on an issue.
type AgentWork struct {
	AgentName    string        `json:"agent_name"`
	WorkSessions int           `json:"work_sessions"`
	TotalTime    time.Duration `json:"total_time"`
	Completed    int           `json:"completed"`
}

// OverallStats contains aggregate statistics across all agents.
type OverallStats struct {
	TotalSessions   int          `json:"total_sessions"`
	ActiveSessions  int          `json:"active_sessions"`
	UniqueAgents    int          `json:"unique_agents"`
	TotalIssues     int          `json:"total_issues"`
	CompletedIssues int          `json:"completed_issues"`
	UniqueSkills    int          `json:"unique_skills"`
	TotalTokens     int          `json:"total_tokens"`
	TopAgents       []AgentCount `json:"top_agents"`
	Since           time.Time    `json:"since"`
}

// SessionDuration represents a session with its duration for visualization.
type SessionDuration struct {
	SessionID   string        `json:"session_id"`
	AgentName   string        `json:"agent_name"`
	StartedAt   time.Time     `json:"started_at"`
	Duration    time.Duration `json:"duration"`
	IsCompleted bool          `json:"is_completed"`
}

// GetAgentStats returns aggregate statistics for a specific agent since a given time.
//
// Example:
//
//	stats, err := agent_tracking.GetAgentStats(db, "beads-workflow-orchestrator", time.Now().AddDate(0, -1, 0))
//	fmt.Printf("Sessions: %d, Issues: %d\n", stats.TotalSessions, stats.TotalIssues)
func GetAgentStats(db *sql.DB, agentName string, since time.Time) (*AgentStats, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if agentName == "" {
		return nil, fmt.Errorf("agent name is required")
	}

	stats := &AgentStats{
		AgentName: agentName,
		Since:     since,
	}

	sinceStr := formatTime(since)

	// Get session counts
	err := db.QueryRow(`
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN ended_at IS NULL THEN 1 ELSE 0 END) as active,
			COALESCE(SUM(context_tokens), 0) as tokens
		FROM agent_sessions
		WHERE agent_name = ? AND started_at >= ?
	`, agentName, sinceStr).Scan(&stats.TotalSessions, &stats.ActiveSessions, &stats.TotalTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to get session counts: %w", err)
	}

	// Get average session time (only for completed sessions)
	var avgSeconds sql.NullFloat64
	err = db.QueryRow(`
		SELECT AVG(
			JULIANDAY(ended_at) - JULIANDAY(started_at)
		) * 86400 as avg_seconds
		FROM agent_sessions
		WHERE agent_name = ? AND started_at >= ? AND ended_at IS NOT NULL
	`, agentName, sinceStr).Scan(&avgSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to get average session time: %w", err)
	}
	if avgSeconds.Valid {
		stats.AvgSessionTime = time.Duration(avgSeconds.Float64 * float64(time.Second))
	}

	// Get issue counts
	err = db.QueryRow(`
		SELECT
			COUNT(DISTINCT issue_id) as total,
			COUNT(DISTINCT CASE WHEN completed = 1 THEN issue_id END) as completed
		FROM agent_issue_work w
		JOIN agent_sessions s ON w.session_id = s.session_id
		WHERE w.agent_name = ? AND s.started_at >= ?
	`, agentName, sinceStr).Scan(&stats.TotalIssues, &stats.CompletedIssues)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue counts: %w", err)
	}

	// Get skill usage counts
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM agent_skill_usage u
		JOIN agent_sessions s ON u.session_id = s.session_id
		WHERE s.agent_name = ? AND s.started_at >= ?
	`, agentName, sinceStr).Scan(&stats.TotalSkillUses)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill usage count: %w", err)
	}

	// Get most used skills
	rows, err := db.Query(`
		SELECT u.skill_name, COUNT(*) as cnt
		FROM agent_skill_usage u
		JOIN agent_sessions s ON u.session_id = s.session_id
		WHERE s.agent_name = ? AND s.started_at >= ?
		GROUP BY u.skill_name
		ORDER BY cnt DESC
		LIMIT 5
	`, agentName, sinceStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get most used skills: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sc SkillCount
		if err := rows.Scan(&sc.SkillName, &sc.Count); err != nil {
			return nil, fmt.Errorf("failed to scan skill count: %w", err)
		}
		stats.MostUsedSkills = append(stats.MostUsedSkills, sc)
	}

	return stats, nil
}

// GetIssueStats returns aggregate statistics for a specific issue.
//
// Example:
//
//	stats, err := agent_tracking.GetIssueStats(db, "agents-42")
//	fmt.Printf("Work sessions: %d, Agents: %d\n", stats.TotalWorkSessions, stats.TotalAgents)
func GetIssueStats(db *sql.DB, issueID string) (*IssueStats, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if issueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	stats := &IssueStats{
		IssueID: issueID,
	}

	// Get work session and agent counts
	var isCompleted int
	err := db.QueryRow(`
		SELECT
			COUNT(*) as total_work,
			COUNT(DISTINCT agent_name) as agents,
			MAX(completed) as is_completed
		FROM agent_issue_work
		WHERE issue_id = ?
	`, issueID).Scan(&stats.TotalWorkSessions, &stats.TotalAgents, &isCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue counts: %w", err)
	}
	stats.IsCompleted = isCompleted == 1

	// Get total time spent
	var totalSeconds sql.NullFloat64
	err = db.QueryRow(`
		SELECT SUM(
			JULIANDAY(COALESCE(ended_at, datetime('now'))) - JULIANDAY(started_at)
		) * 86400 as total_seconds
		FROM agent_issue_work
		WHERE issue_id = ?
	`, issueID).Scan(&totalSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to get total time: %w", err)
	}
	if totalSeconds.Valid {
		stats.TotalTime = time.Duration(totalSeconds.Float64 * float64(time.Second))
	}

	// Get agent breakdown
	rows, err := db.Query(`
		SELECT
			agent_name,
			COUNT(*) as work_sessions,
			SUM(JULIANDAY(COALESCE(ended_at, datetime('now'))) - JULIANDAY(started_at)) * 86400 as total_seconds,
			SUM(completed) as completed_count
		FROM agent_issue_work
		WHERE issue_id = ?
		GROUP BY agent_name
		ORDER BY work_sessions DESC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var aw AgentWork
		var totalSecs float64
		if err := rows.Scan(&aw.AgentName, &aw.WorkSessions, &totalSecs, &aw.Completed); err != nil {
			return nil, fmt.Errorf("failed to scan agent work: %w", err)
		}
		aw.TotalTime = time.Duration(totalSecs * float64(time.Second))
		stats.AgentBreakdown = append(stats.AgentBreakdown, aw)
	}

	return stats, nil
}

// GetSkillStats returns aggregate statistics for a specific skill since a given time.
//
// Example:
//
//	stats, err := agent_tracking.GetSkillStats(db, "dependency-thinking", time.Now().AddDate(0, -1, 0))
//	fmt.Printf("Uses: %d, Avg context: %.1f\n", stats.TotalUses, stats.AvgContext)
func GetSkillStats(db *sql.DB, skillName string, since time.Time) (*SkillStats, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if skillName == "" {
		return nil, fmt.Errorf("skill name is required")
	}

	stats := &SkillStats{
		SkillName: skillName,
		Since:     since,
	}

	sinceStr := formatTime(since)

	// Get usage counts
	var avgContext sql.NullFloat64
	err := db.QueryRow(`
		SELECT
			COUNT(*) as total,
			COUNT(DISTINCT u.session_id) as sessions,
			COUNT(DISTINCT s.agent_name) as agents,
			COALESCE(SUM(u.context_added), 0) as total_context,
			AVG(u.context_added) as avg_context
		FROM agent_skill_usage u
		JOIN agent_sessions s ON u.session_id = s.session_id
		WHERE u.skill_name = ? AND u.loaded_at >= ?
	`, skillName, sinceStr).Scan(
		&stats.TotalUses, &stats.UniqueSessions, &stats.UniqueAgents,
		&stats.TotalContext, &avgContext,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill usage counts: %w", err)
	}
	if avgContext.Valid {
		stats.AvgContext = avgContext.Float64
	}

	// Get top agents using this skill
	rows, err := db.Query(`
		SELECT s.agent_name, COUNT(*) as cnt
		FROM agent_skill_usage u
		JOIN agent_sessions s ON u.session_id = s.session_id
		WHERE u.skill_name = ? AND u.loaded_at >= ?
		GROUP BY s.agent_name
		ORDER BY cnt DESC
		LIMIT 5
	`, skillName, sinceStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get top agents: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ac AgentCount
		if err := rows.Scan(&ac.AgentName, &ac.Count); err != nil {
			return nil, fmt.Errorf("failed to scan agent count: %w", err)
		}
		stats.TopAgents = append(stats.TopAgents, ac)
	}

	return stats, nil
}

// GetOverallStats returns aggregate statistics across all agents.
//
// Example:
//
//	stats, err := agent_tracking.GetOverallStats(db, time.Now().AddDate(0, -1, 0))
func GetOverallStats(db *sql.DB, since time.Time) (*OverallStats, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	stats := &OverallStats{
		Since: since,
	}

	sinceStr := formatTime(since)

	// Get session counts
	err := db.QueryRow(`
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN ended_at IS NULL THEN 1 ELSE 0 END) as active,
			COUNT(DISTINCT agent_name) as agents,
			COALESCE(SUM(context_tokens), 0) as tokens
		FROM agent_sessions
		WHERE started_at >= ?
	`, sinceStr).Scan(&stats.TotalSessions, &stats.ActiveSessions, &stats.UniqueAgents, &stats.TotalTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to get session counts: %w", err)
	}

	// Get issue counts
	err = db.QueryRow(`
		SELECT
			COUNT(DISTINCT issue_id) as total,
			COUNT(DISTINCT CASE WHEN completed = 1 THEN issue_id END) as completed
		FROM agent_issue_work w
		JOIN agent_sessions s ON w.session_id = s.session_id
		WHERE s.started_at >= ?
	`, sinceStr).Scan(&stats.TotalIssues, &stats.CompletedIssues)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue counts: %w", err)
	}

	// Get skill usage count
	err = db.QueryRow(`
		SELECT COUNT(DISTINCT skill_name)
		FROM agent_skill_usage
		WHERE loaded_at >= ?
	`, sinceStr).Scan(&stats.UniqueSkills)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill count: %w", err)
	}

	// Get top agents
	rows, err := db.Query(`
		SELECT agent_name, COUNT(*) as cnt
		FROM agent_sessions
		WHERE started_at >= ?
		GROUP BY agent_name
		ORDER BY cnt DESC
		LIMIT 5
	`, sinceStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get top agents: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ac AgentCount
		if err := rows.Scan(&ac.AgentName, &ac.Count); err != nil {
			return nil, fmt.Errorf("failed to scan agent count: %w", err)
		}
		stats.TopAgents = append(stats.TopAgents, ac)
	}

	return stats, nil
}

// GetSessionDurations returns session duration statistics for visualization.
func GetSessionDurations(db *sql.DB, agentName string, since time.Time, limit int) ([]SessionDuration, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if limit <= 0 {
		limit = 50
	}

	sinceStr := formatTime(since)

	var query string
	var args []interface{}

	if agentName != "" {
		query = `
			SELECT
				session_id,
				agent_name,
				started_at,
				(JULIANDAY(COALESCE(ended_at, datetime('now'))) - JULIANDAY(started_at)) * 86400 as duration_seconds,
				ended_at IS NOT NULL as is_completed
			FROM agent_sessions
			WHERE agent_name = ? AND started_at >= ?
			ORDER BY started_at DESC
			LIMIT ?
		`
		args = []interface{}{agentName, sinceStr, limit}
	} else {
		query = `
			SELECT
				session_id,
				agent_name,
				started_at,
				(JULIANDAY(COALESCE(ended_at, datetime('now'))) - JULIANDAY(started_at)) * 86400 as duration_seconds,
				ended_at IS NOT NULL as is_completed
			FROM agent_sessions
			WHERE started_at >= ?
			ORDER BY started_at DESC
			LIMIT ?
		`
		args = []interface{}{sinceStr, limit}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get session durations: %w", err)
	}
	defer rows.Close()

	var durations []SessionDuration
	for rows.Next() {
		var sd SessionDuration
		var startedAtStr string
		var durationSecs float64

		if err := rows.Scan(&sd.SessionID, &sd.AgentName, &startedAtStr, &durationSecs, &sd.IsCompleted); err != nil {
			return nil, fmt.Errorf("failed to scan session duration: %w", err)
		}

		sd.StartedAt, err = parseTime(startedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse started_at: %w", err)
		}
		sd.Duration = time.Duration(durationSecs * float64(time.Second))

		durations = append(durations, sd)
	}

	return durations, nil
}

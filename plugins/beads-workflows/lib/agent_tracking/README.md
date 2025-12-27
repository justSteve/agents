# Agent Tracking Extension

Go library that extends the beads database with agent tracking tables. This follows the beads extension pattern - adding custom tables via SQL without modifying beads core.

## Overview

This extension adds three tables to track agent activity:

- **agent_sessions** - Tracks agent work sessions (start/end times, model tier, issues claimed, skills used)
- **agent_issue_work** - Tracks work done on specific issues (start/end, rationale, completion status)
- **agent_skill_usage** - Tracks skill loading and usage during sessions

## Usage

### 1. Initialize the Extension

Call `Initialize()` once when your application starts, after obtaining the database connection from beads:

```go
import (
    "github.com/steveyegge/beads/store"
    "github.com/justSteve/agents/plugins/beads-workflows/lib/agent_tracking"
)

func main() {
    // Get the beads store
    beadsStore, err := store.Open(".beads/beads.db")
    if err != nil {
        log.Fatal(err)
    }
    defer beadsStore.Close()

    // Get the underlying *sql.DB
    db := beadsStore.UnderlyingDB()

    // Initialize agent tracking tables
    if err := agent_tracking.Initialize(db); err != nil {
        log.Fatal(err)
    }
}
```

### 2. Session Management

```go
// Start a new session
sessionID, err := agent_tracking.StartSession(db, "beads-workflow-orchestrator", "/myStuff/project", "sonnet")
if err != nil {
    return err
}

// End the session when done
defer agent_tracking.EndSession(db, sessionID, "completed")

// Update session metadata
agent_tracking.UpdateSessionIssues(db, sessionID, []string{"agents-42"})
agent_tracking.UpdateSessionSkills(db, sessionID, []string{"dependency-thinking"})
agent_tracking.UpdateSessionTokens(db, sessionID, 15000)

// Get session details
session, err := agent_tracking.GetSession(db, sessionID)

// List active sessions
activeSessions, err := agent_tracking.ListActiveSessions(db)

// List sessions by agent
sessions, err := agent_tracking.ListSessionsByAgent(db, "beads-workflow-orchestrator", 10)
```

### 3. Work Tracking

```go
// Record starting work on an issue
workID, err := agent_tracking.RecordWork(db, sessionID, "agents-42", "beads-workflow-orchestrator", "Highest priority P1 task")
if err != nil {
    return err
}

// Complete the work
err = agent_tracking.CompleteWork(db, workID, "Implemented user model with validation")

// Get work details
work, err := agent_tracking.GetWork(db, workID)

// List all work on an issue
workEntries, err := agent_tracking.ListWorkByIssue(db, "agents-42")

// List all work in a session
workEntries, err := agent_tracking.ListWorkBySession(db, sessionID)
```

### 4. Skill Usage Tracking

```go
// Record loading a skill
err := agent_tracking.RecordSkillUsage(db, sessionID, "dependency-thinking", "agents-42", 500)
```

### 5. Statistics

```go
import "time"

// Get agent statistics for the last month
since := time.Now().AddDate(0, -1, 0)
agentStats, err := agent_tracking.GetAgentStats(db, "beads-workflow-orchestrator", since)
fmt.Printf("Sessions: %d, Issues: %d, Avg Time: %v\n",
    agentStats.TotalSessions,
    agentStats.TotalIssues,
    agentStats.AvgSessionTime)

// Get issue statistics
issueStats, err := agent_tracking.GetIssueStats(db, "agents-42")
fmt.Printf("Work sessions: %d, Agents: %d, Time: %v\n",
    issueStats.TotalWorkSessions,
    issueStats.TotalAgents,
    issueStats.TotalTime)

// Get skill statistics
skillStats, err := agent_tracking.GetSkillStats(db, "dependency-thinking", since)
fmt.Printf("Uses: %d, Avg context: %.1f\n",
    skillStats.TotalUses,
    skillStats.AvgContext)

// Get overall statistics
overallStats, err := agent_tracking.GetOverallStats(db, since)
fmt.Printf("Total sessions: %d, Active: %d, Unique agents: %d\n",
    overallStats.TotalSessions,
    overallStats.ActiveSessions,
    overallStats.UniqueAgents)

// Get session durations for visualization
durations, err := agent_tracking.GetSessionDurations(db, "", since, 50)
for _, d := range durations {
    fmt.Printf("%s: %v (completed: %v)\n", d.AgentName, d.Duration, d.IsCompleted)
}
```

## Schema

### agent_sessions

| Column | Type | Description |
|--------|------|-------------|
| session_id | TEXT PK | Unique session identifier |
| agent_name | TEXT | Name of the agent (e.g., "beads-workflow-orchestrator") |
| workspace_path | TEXT | Path to the workspace directory |
| started_at | TEXT | ISO 8601 timestamp when session started |
| ended_at | TEXT | ISO 8601 timestamp when session ended (NULL if active) |
| exit_reason | TEXT | Reason for session end ("completed", "interrupted", "error", "timeout") |
| issues_claimed | TEXT | JSON array of issue IDs claimed during session |
| skills_used | TEXT | JSON array of skill names used during session |
| model_tier | TEXT | Model tier used ("sonnet", "opus", etc.) |
| context_tokens | INTEGER | Number of context tokens used |
| created_at | TEXT | ISO 8601 timestamp when record was created |

### agent_issue_work

| Column | Type | Description |
|--------|------|-------------|
| work_id | TEXT PK | Unique work identifier |
| issue_id | TEXT FK | Reference to beads issue |
| session_id | TEXT FK | Reference to agent_sessions |
| agent_name | TEXT | Name of the agent doing the work |
| started_at | TEXT | ISO 8601 timestamp when work started |
| ended_at | TEXT | ISO 8601 timestamp when work ended |
| status_changes | TEXT | JSON array of status transitions |
| decision_rationale | TEXT | Why this issue was selected |
| work_notes | TEXT | Notes about the work done |
| completed | BOOLEAN | Whether work was completed |

### agent_skill_usage

| Column | Type | Description |
|--------|------|-------------|
| usage_id | TEXT PK | Unique usage identifier |
| session_id | TEXT FK | Reference to agent_sessions |
| skill_name | TEXT | Name of the skill loaded |
| loaded_at | TEXT | ISO 8601 timestamp when skill was loaded |
| used_for_issue_id | TEXT | Issue the skill was used for (optional) |
| context_added | INTEGER | Number of tokens added by the skill |

## Extension Pattern

This library follows the beads extension pattern:

1. **No core modifications** - Tables are added via SQL CREATE TABLE IF NOT EXISTS
2. **Namespaced tables** - All tables prefixed with `agent_` to avoid conflicts
3. **Foreign keys** - References to agent_sessions use ON DELETE CASCADE
4. **Initialization** - Call `Initialize(db)` once at startup
5. **Access via UnderlyingDB()** - Get `*sql.DB` from beads store

## Exit Reasons

Standard exit reasons for sessions:

- `completed` - Session ended normally after finishing work
- `interrupted` - Session was interrupted by user or system
- `error` - Session ended due to an error
- `timeout` - Session exceeded time limit
- `context_limit` - Session ended due to context window limit

## Dependencies

- `github.com/google/uuid` - UUID generation for record IDs
- Standard library `database/sql` - Database operations

---
name: beads-cli-reference
description: Quick reference for all beads CLI commands with syntax and examples. Use when you need to recall beads command syntax, options, or usage patterns.
---

# Beads CLI Reference

Complete reference for the `bd` (beads) command-line interface. Use this when you need to recall command syntax, available options, or usage examples.

## When to Use This Skill

- Need to recall exact syntax for a beads command
- Looking up available options for a command
- Finding examples of common operations
- Troubleshooting command errors
- Learning new beads features

## Command Overview

| Command | Purpose |
|---------|---------|
| `bd ready` | List issues ready to work on |
| `bd list` | List all issues with filters |
| `bd show` | Display issue details |
| `bd create` | Create a new issue |
| `bd update` | Modify an existing issue |
| `bd close` | Close a completed issue |
| `bd dep add` | Add dependency between issues |
| `bd dep remove` | Remove a dependency |
| `bd blocked` | Show blocked issues and their blockers |
| `bd sync` | Synchronize beads state with git |
| `bd stats` | Display tracking statistics |

---

## bd ready

**Purpose**: List issues that are ready to be worked on (open and unblocked).

**Syntax**:
```bash
bd ready [options]
```

**Options**:
| Option | Description |
|--------|-------------|
| `--type`, `-t` | Filter by issue type (task, bug, epic, feature) |
| `--priority`, `-p` | Filter by priority (0-4, 0 is highest) |
| `--limit`, `-l` | Maximum number of results |

**Examples**:
```bash
# List all ready issues
bd ready

# List only ready bugs
bd ready --type bug

# List high-priority ready issues
bd ready --priority 0
bd ready --priority 1

# List top 5 ready issues
bd ready --limit 5
```

**Output**:
```
Ready issues (unblocked, open):
  auth-login-fix     [task] [P1] Fix authentication timeout issue
  api-docs-update    [task] [P2] Update API documentation
  db-migration-v2    [task] [P1] Add user preferences table
```

---

## bd list

**Purpose**: List all issues with optional filters.

**Syntax**:
```bash
bd list [options]
```

**Options**:
| Option | Description |
|--------|-------------|
| `--status`, `-s` | Filter by status (open, in_progress, closed) |
| `--type`, `-t` | Filter by issue type |
| `--priority`, `-p` | Filter by priority |
| `--assignee`, `-a` | Filter by assignee |
| `--label` | Filter by label |
| `--limit`, `-l` | Maximum number of results |
| `--all` | Include closed issues |

**Examples**:
```bash
# List all open issues
bd list

# List issues by status
bd list --status open
bd list --status in_progress
bd list --status closed

# List by type
bd list --type bug
bd list --type epic

# List by priority
bd list --priority 0  # Critical
bd list --priority 1  # High

# Combine filters
bd list --status open --type bug --priority 1

# Include closed issues
bd list --all
```

**Output**:
```
Issues:
  auth-login-fix     [task] [P1] [open]        Fix authentication timeout
  api-docs-update    [task] [P2] [in_progress] Update API documentation
  user-profile       [epic] [P2] [open]        User profile system
```

---

## bd show

**Purpose**: Display detailed information about an issue.

**Syntax**:
```bash
bd show <issue-id>
```

**Arguments**:
| Argument | Description |
|----------|-------------|
| `issue-id` | The unique identifier of the issue |

**Examples**:
```bash
# Show issue details
bd show auth-login-fix

# Show epic with dependencies
bd show user-profile
```

**Output**:
```
Issue: auth-login-fix
Type: task
Priority: P1 (High)
Status: open
Created: 2024-01-15 10:30:00
Updated: 2024-01-16 14:22:00

Description:
  Authentication tokens expire after 1 hour with no refresh mechanism.
  Users lose unsaved work during long form submissions.

Acceptance Criteria:
  [ ] Tokens refresh automatically when less than 15 minutes remain
  [ ] Failed refresh attempts redirect to login
  [ ] Token refresh completes in < 100ms

Dependencies:
  Blocked by: (none)
  Blocks: api-session-management

Related Commits:
  abc1234 - Initial investigation
  def5678 - Added token refresh endpoint
```

---

## bd create

**Purpose**: Create a new issue.

**Syntax**:
```bash
bd create "<title>" [options]
```

**Options**:
| Option | Description |
|--------|-------------|
| `--description`, `-d` | Issue description (multiline supported) |
| `--type`, `-t` | Issue type: task, bug, epic, feature |
| `--priority`, `-p` | Priority: 0 (critical) to 4 (low) |
| `--assignee`, `-a` | Assign to user/agent |
| `--labels`, `-l` | Comma-separated labels |
| `--deps` | Dependencies in format `type:issue-id` |

**Dependency Types**:
| Type | Meaning |
|------|---------|
| `blocks` | This issue cannot start until dependency is done |
| `parent-child` | This issue is a sub-task of the parent |
| `discovered-from` | This issue was found while working on another |
| `related` | Issues are related but not blocking |

**Examples**:
```bash
# Simple issue creation
bd create "Fix login button alignment"

# Full issue with all options
bd create "Implement user preferences API" \
  --description="Add REST endpoints for user preferences management. Include GET, PUT, DELETE operations." \
  --type feature \
  --priority 2 \
  --assignee claude \
  --labels "api,user-system"

# Issue with dependencies
bd create "Add preference validation" \
  --description="Validate preference values before saving" \
  --type task \
  --deps blocks:user-preferences-api

# Discovered work (filed while working on another issue)
bd create "Fix null pointer in user lookup" \
  --description="Found while implementing preferences. get_user() can return null but callers don't check." \
  --type bug \
  --priority 1 \
  --deps discovered-from:user-preferences-api

# Child of an epic
bd create "Design preferences schema" \
  --type task \
  --deps parent-child:user-preferences-epic

# Multiple dependencies
bd create "Integration tests for preferences" \
  --type task \
  --deps blocks:preferences-api,blocks:preferences-validation
```

**Output**:
```
Created issue: user-preferences-api
Type: feature
Priority: P2
Status: open
```

---

## bd update

**Purpose**: Modify an existing issue.

**Syntax**:
```bash
bd update <issue-id> [options]
```

**Options**:
| Option | Description |
|--------|-------------|
| `--status`, `-s` | New status (open, in_progress, closed) |
| `--priority`, `-p` | New priority (0-4) |
| `--title` | New title |
| `--description`, `-d` | New description |
| `--assignee`, `-a` | New assignee |
| `--labels`, `-l` | Replace labels |
| `--notes` | Add notes (appended to description) |

**Examples**:
```bash
# Claim an issue (start working on it)
bd update auth-login-fix --status=in_progress

# Change priority
bd update auth-login-fix --priority 0

# Update title
bd update auth-login-fix --title="Fix authentication token refresh"

# Add notes
bd update auth-login-fix --notes="WIP: Implemented refresh logic, still need error handling"

# Reassign
bd update auth-login-fix --assignee=other-agent

# Mark as blocked (back to open from in_progress)
bd update auth-login-fix --status=open --notes="Blocked waiting for API review"
```

**Output**:
```
Updated issue: auth-login-fix
Status: open -> in_progress
```

---

## bd close

**Purpose**: Close a completed issue.

**Syntax**:
```bash
bd close <issue-id> [options]
```

**Options**:
| Option | Description |
|--------|-------------|
| `--resolution` | Resolution type (fixed, wontfix, duplicate, invalid) |
| `--notes` | Closing notes |

**Examples**:
```bash
# Close as fixed (default)
bd close auth-login-fix

# Close with notes
bd close auth-login-fix --notes="Implemented token refresh with 15-minute threshold"

# Close as won't fix
bd close outdated-feature --resolution=wontfix --notes="Feature no longer needed after pivot"

# Close as duplicate
bd close duplicate-bug --resolution=duplicate --notes="Duplicate of auth-login-fix"
```

**Output**:
```
Closed issue: auth-login-fix
Resolution: fixed
```

---

## bd dep add

**Purpose**: Add a dependency between issues.

**Syntax**:
```bash
bd dep add <dependent> <required> [--type TYPE]
```

**Arguments**:
| Argument | Description |
|----------|-------------|
| `dependent` | Issue that depends on (needs) the other |
| `required` | Issue that must be completed first |

**Dependency Types**:
| Type | Description | Default |
|------|-------------|---------|
| `blocks` | Hard blocker - dependent cannot start until required is done | Yes |
| `parent-child` | Hierarchical - dependent is sub-task of required | No |
| `discovered-from` | Discovery - dependent was found while working on required | No |
| `related` | Soft link - informational only, no blocking | No |

**Examples**:
```bash
# Add blocking dependency (api-endpoints needs db-schema first)
bd dep add api-endpoints db-schema
bd dep add api-endpoints db-schema --type blocks

# Add parent-child relationship
bd dep add login-endpoint user-auth-epic --type parent-child

# Add discovered-from relationship
bd dep add fix-null-check user-preferences --type discovered-from

# Add related link (informational)
bd dep add password-reset login-fix --type related
```

**Mental model**:
> "Y needs X" = `bd dep add Y X`

**Common mistake**:
```bash
# WRONG: "Phase 1 before Phase 2" thinking
bd dep add phase1 phase2  # This says phase1 NEEDS phase2!

# RIGHT: "Phase 2 needs Phase 1"
bd dep add phase2 phase1  # Phase 2 is blocked by Phase 1
```

**Output**:
```
Added dependency: api-endpoints -> db-schema (blocks)
```

---

## bd dep remove

**Purpose**: Remove a dependency between issues.

**Syntax**:
```bash
bd dep remove <dependent> <required>
```

**Examples**:
```bash
# Remove a blocking dependency
bd dep remove api-endpoints db-schema

# Remove after realizing dependency was added incorrectly
bd dep remove phase1 phase2  # Fix inverted dependency
```

**Output**:
```
Removed dependency: api-endpoints -> db-schema
```

---

## bd blocked

**Purpose**: Show all blocked issues and what's blocking them.

**Syntax**:
```bash
bd blocked [options]
```

**Options**:
| Option | Description |
|--------|-------------|
| `--issue`, `-i` | Show blockers for specific issue only |
| `--tree` | Show full dependency tree |

**Examples**:
```bash
# Show all blocked issues
bd blocked

# Show blockers for specific issue
bd blocked --issue api-endpoints

# Show full dependency tree
bd blocked --tree
```

**Output**:
```
Blocked issues:

api-endpoints is blocked by:
  - db-schema (blocks) [open]

integration-tests is blocked by:
  - api-endpoints (blocks) [open]
  - frontend-components (blocks) [in_progress]

Total: 2 blocked issues
```

**Use cases**:
- Verify dependencies after adding them
- Find circular dependencies
- Understand why an issue isn't in `bd ready`

---

## bd sync

**Purpose**: Synchronize beads state with git repository.

**Syntax**:
```bash
bd sync [options]
```

**Options**:
| Option | Description |
|--------|-------------|
| `--force` | Force sync even if no changes detected |
| `--dry-run` | Show what would be synced without doing it |

**Examples**:
```bash
# Standard sync
bd sync

# Sync after commit (links commit to issues)
git commit -m "fix: resolve timeout"
bd sync

# Dry run to preview
bd sync --dry-run
```

**When to sync**:
1. After claiming an issue (`bd update --status=in_progress`)
2. Before committing (captures current state)
3. After committing (links commit hash to issue)
4. After closing an issue
5. When starting a new session

**Output**:
```
Syncing beads state...
Updated: auth-login-fix (linked commit abc1234)
Updated: api-docs-update (status change)
Sync complete.
```

---

## bd stats

**Purpose**: Display tracking statistics.

**Syntax**:
```bash
bd stats [options]
```

**Options**:
| Option | Description |
|--------|-------------|
| `--period` | Time period (today, week, month, all) |
| `--assignee`, `-a` | Filter by assignee |

**Examples**:
```bash
# Overall statistics
bd stats

# This week's stats
bd stats --period week

# Stats for specific agent
bd stats --assignee claude
```

**Output**:
```
Beads Statistics (all time)

Issues:
  Total: 47
  Open: 12
  In Progress: 3
  Closed: 32

By Type:
  Tasks: 28
  Bugs: 11
  Features: 6
  Epics: 2

By Priority:
  P0 (Critical): 2
  P1 (High): 8
  P2 (Medium): 25
  P3 (Low): 12

Velocity (last 7 days):
  Closed: 8 issues
  Created: 5 issues
```

---

## Common Workflows

### Starting a Work Session

```bash
# Find available work
bd ready

# Review issue details
bd show <issue-id>

# Claim the issue
bd update <issue-id> --status=in_progress

# Sync state
bd sync
```

### Completing Work

```bash
# Check git state
git status

# Stage and sync
git add <files>
bd sync

# Commit
git commit -m "feat(scope): description

Resolves: <issue-id>"

# Sync and push
bd sync
git push

# Close issue
bd close <issue-id>
bd sync
```

### Filing Discovered Work

```bash
# While working on issue X, discover Y
bd create "Description of discovered issue" \
  --description="Context. Found while working on X." \
  --deps discovered-from:X \
  --type task

# Continue working on X
```

### Creating an Epic with Children

```bash
# Create epic
bd create "User Authentication System" \
  --type epic \
  --priority 1

# Create child tasks
bd create "Create user model" \
  --type task \
  --deps parent-child:user-authentication-system

bd create "Add login endpoint" \
  --type task \
  --deps parent-child:user-authentication-system,blocks:create-user-model

bd create "Add logout endpoint" \
  --type task \
  --deps parent-child:user-authentication-system,blocks:create-user-model
```

### Checking Why Issue Is Blocked

```bash
# See all blockers
bd blocked --issue <issue-id>

# Or show full details
bd show <issue-id>
```

---

## Quick Reference Card

**Find Work**:
```bash
bd ready                    # Ready issues
bd list --status open       # All open issues
bd show <id>               # Issue details
```

**Manage Issues**:
```bash
bd create "<title>"         # Create issue
bd update <id> --status=X   # Update status
bd close <id>              # Close issue
```

**Dependencies**:
```bash
bd dep add <needs> <first>  # Add dependency
bd dep remove <a> <b>       # Remove dependency
bd blocked                  # Show blocked issues
```

**Sync & Stats**:
```bash
bd sync                     # Sync with git
bd stats                    # Show statistics
```

**Status Values**: `open`, `in_progress`, `closed`

**Issue Types**: `task`, `bug`, `epic`, `feature`

**Priorities**: `0` (critical), `1` (high), `2` (medium), `3` (low), `4` (lowest)

**Dependency Types**: `blocks`, `parent-child`, `discovered-from`, `related`

---

## Troubleshooting

### Issue not appearing in bd ready

**Cause**: Issue is blocked or not open.

**Fix**:
```bash
bd show <issue-id>     # Check status and blockers
bd blocked --issue <issue-id>  # See what's blocking
```

### Dependency seems inverted

**Cause**: Used temporal thinking ("A before B") instead of requirement thinking ("B needs A").

**Fix**:
```bash
bd dep remove <wrong-dependent> <wrong-required>
bd dep add <correct-dependent> <correct-required>
bd blocked  # Verify
```

### Sync not capturing commits

**Cause**: Sync run before commit, not after.

**Fix**:
```bash
# Always sync AFTER committing
git commit -m "..."
bd sync  # Now it captures the commit
```

### Circular dependency detected

**Cause**: A blocks B, B blocks A (or longer cycle).

**Fix**:
```bash
bd blocked  # Identify the cycle
bd dep remove <one-side> <other-side>  # Break the cycle
```

## Summary

The `bd` CLI provides complete issue tracking integrated with git. Key commands:

- **bd ready**: Find work to do
- **bd create**: File new issues
- **bd update**: Claim and update issues
- **bd close**: Complete issues
- **bd dep add/remove**: Manage dependencies
- **bd sync**: Keep git and beads in sync

Master these commands and you'll have full control over your issue tracking workflow!

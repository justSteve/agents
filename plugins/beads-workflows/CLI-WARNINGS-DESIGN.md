# Beads CLI Soft Validation Warnings Design

> **Status**: Design Proposal (Layer 4 Enforcement)
> **Implementation**: Requires beads core changes
> **Purpose**: Non-blocking warnings to guide best practices

## Overview

This document proposes soft validation warnings for the beads CLI. These warnings:
- Display guidance when patterns are violated
- Do NOT block the operation
- Educate users about best practices
- Can be suppressed with `--quiet` or `--no-warnings`

## Design Principles

1. **Non-blocking**: Warnings never prevent the command from executing
2. **Educational**: Explain WHY the pattern matters
3. **Actionable**: Tell user how to fix or improve
4. **Suppressible**: Can be disabled for automation/scripting
5. **Consistent**: Same format across all warnings

## Warning Message Format

```
⚠️  [WARNING-CODE]: [Short description]

[Explanation of the issue]

Recommendation:
[How to fix or improve]

Tip: Use --quiet to suppress warnings
```

## Proposed Warnings

### 1. Description Quality Warnings

#### WARN-DESC-001: Missing Description

**Trigger**: `bd create "title"` without `--description`

**Message**:
```
⚠️  WARN-DESC-001: Issue created without description

Issues without descriptions are harder to understand and work on.

Recommendation:
Add a description with Why/What/How structure:
  bd update <id> --description="## Context
[Why this matters]

## Problem
[What needs to change]

## Acceptance Criteria
- [ ] [Verifiable criterion]"

Issue created: <id>
```

#### WARN-DESC-002: Short Description

**Trigger**: `bd create` with description < 50 characters

**Message**:
```
⚠️  WARN-DESC-002: Description may be too brief (N chars)

Short descriptions often lack essential context.

Minimum recommended:
- Tasks: 50+ characters
- Bugs: 100+ characters
- Features: 100+ characters
- Epics: 200+ characters

Recommendation:
Include Why (context), What (problem), and How (verification criteria).

Issue created: <id>
```

#### WARN-DESC-003: Vague Title

**Trigger**: Title matches vague patterns: "fix bug", "update", "refactor", "changes", single words

**Message**:
```
⚠️  WARN-DESC-003: Title may be too vague

Title "<title>" doesn't clearly describe the work.

Examples of better titles:
- "Fix null pointer in user lookup when email missing"
- "Update auth token expiry from 1h to 24h"
- "Refactor UserService into auth and profile modules"

Recommendation:
Update with: bd update <id> --title="<specific title>"

Issue created: <id>
```

### 2. Single-Issue Discipline Warnings

#### WARN-SID-001: Multiple Issues In Progress

**Trigger**: `bd update <id> --status=in_progress` when another issue is already in_progress

**Message**:
```
⚠️  WARN-SID-001: Multiple issues now in progress

You now have 2 issues in progress:
- <existing-id>: <existing-title>
- <new-id>: <new-title>

Beads best practice: Work on ONE issue at a time.

Problems with multiple in_progress:
- Context switching overhead
- Unclear which issue changes belong to
- Risk of orphaned work

Recommendation:
Release one issue: bd update <id> --status=open

Status updated: <new-id> is now in_progress
```

#### WARN-SID-002: Many Issues In Progress

**Trigger**: `bd update --status=in_progress` when 3+ issues already in_progress

**Message**:
```
⚠️  WARN-SID-002: Too many issues in progress (N total)

You have N issues in progress:
- <id1>: <title1>
- <id2>: <title2>
- <id3>: <title3>
...

This strongly violates single-issue discipline.

Recommendation:
1. Review each in-progress issue
2. Close completed ones: bd close <id>
3. Release others: bd update <id> --status=open
4. Focus on ONE issue

Status updated: <new-id> is now in_progress
```

### 3. Dependency Warnings

#### WARN-DEP-001: Temporal Language Detected

**Trigger**: Cannot be detected at CLI level (requires NLP), but could check issue titles/descriptions for patterns

**Note**: This warning is better suited for agent-level enforcement (beads-disciplinarian) rather than CLI.

#### WARN-DEP-002: Self-Dependency Attempt

**Trigger**: `bd dep add <id> <id>` (same ID twice)

**Message**:
```
⚠️  WARN-DEP-002: Cannot create self-dependency

Issue <id> cannot depend on itself.

Command rejected.
```

**Note**: This should be an ERROR, not a warning (blocking).

#### WARN-DEP-003: Circular Dependency Created

**Trigger**: `bd dep add` would create a cycle

**Message**:
```
⚠️  WARN-DEP-003: Circular dependency detected

Adding this dependency would create a cycle:
<id-a> → <id-b> → <id-c> → <id-a>

Circular dependencies cause all issues in the cycle to be permanently blocked.

Command rejected.
```

**Note**: This should be an ERROR, not a warning (blocking).

### 4. Session Ritual Warnings

#### WARN-SESS-001: Claiming Blocked Issue

**Trigger**: `bd update <id> --status=in_progress` on a blocked issue

**Message**:
```
⚠️  WARN-SESS-001: Issue is blocked by dependencies

Issue <id> is blocked by:
- <blocker-1>: <title1> (status: <status>)
- <blocker-2>: <title2> (status: <status>)

Working on blocked issues may lead to waiting or rework.

Recommendation:
1. Work on blockers first: bd update <blocker-1> --status=in_progress
2. Or use bd ready to find unblocked work

Status updated: <id> is now in_progress (blocked)
```

#### WARN-SESS-002: Closing Issue with Open Dependencies

**Trigger**: `bd close <id>` when issue has unclosed blocking dependencies

**Message**:
```
⚠️  WARN-SESS-002: Closing issue with unresolved blockers

Issue <id> has dependencies that are still open:
- <blocker-1>: <title1> (status: open)

This is unusual - typically blockers are completed first.

Possible scenarios:
1. Blockers were not actually needed (remove them)
2. Work was done out of order (acceptable)
3. Closing by mistake (re-open with bd update --status=open)

Issue closed: <id>
```

#### WARN-SESS-003: Orphaned In-Progress on Ready

**Trigger**: `bd ready` when issues are in_progress

**Message**:
```
⚠️  WARN-SESS-003: You have work in progress

Currently in progress:
- <id>: <title> (started: <timestamp>)

Consider completing or releasing this work before starting new issues.

Ready issues shown below...
```

### 5. Sync Warnings

#### WARN-SYNC-001: Uncommitted Changes

**Trigger**: `bd sync` when git has uncommitted changes in .beads/

**Message**:
```
⚠️  WARN-SYNC-001: Uncommitted beads changes detected

.beads/issues.jsonl has uncommitted changes.

Recommendation:
Commit beads changes to preserve history:
  git add .beads/
  git commit -m "chore(beads): sync issue state"

Sync completed.
```

#### WARN-SYNC-002: Remote Ahead

**Trigger**: `bd sync` when remote has changes not pulled

**Message**:
```
⚠️  WARN-SYNC-002: Remote has newer changes

Remote branch has N commits not in local.

Recommendation:
Pull before continuing:
  git pull

Sync completed with local state.
```

## Warning Severity Levels

| Level | Prefix | Meaning | Blocks? |
|-------|--------|---------|---------|
| INFO | ℹ️ | Helpful suggestion | No |
| WARN | ⚠️ | Best practice violation | No |
| ERROR | ❌ | Invalid operation | Yes |

## Suppression Options

```bash
# Suppress all warnings
bd create "title" --quiet

# Suppress specific warning
bd create "title" --suppress WARN-DESC-001

# Suppress warnings in config
bd config set warnings.suppress "WARN-DESC-001,WARN-DESC-002"

# Disable all warnings globally
bd config set warnings.enabled false
```

## Implementation Notes

### For beads Core Team

1. **Warning Registry**: Create central registry of warning codes and messages
2. **Hook Points**: Add warning emission points in relevant commands
3. **Suppression Logic**: Check config and flags before emitting
4. **Output Format**: Warnings go to stderr, command output to stdout
5. **Exit Codes**: Warnings don't affect exit codes (always 0 if command succeeds)

### Integration with beads-workflows

These warnings complement the agent-level enforcement:
- **CLI Warnings (Layer 4)**: Immediate feedback during command execution
- **Agent Enforcement (Layer 2-3)**: Deeper validation with context awareness

### Testing Recommendations

For each warning:
1. Test that warning triggers correctly
2. Test that command still completes
3. Test suppression works
4. Test warning message is accurate

## Future Enhancements

1. **Interactive Mode**: Prompt user to fix issues before proceeding
2. **Auto-fix Suggestions**: Offer to run fix commands automatically
3. **Warning Statistics**: Track which warnings fire most often
4. **Custom Warnings**: Allow plugins to register custom warnings
5. **Warning Levels**: Allow setting minimum severity to display

## Summary

| Warning Code | Trigger | Blocking? |
|--------------|---------|-----------|
| WARN-DESC-001 | No description | No |
| WARN-DESC-002 | Short description | No |
| WARN-DESC-003 | Vague title | No |
| WARN-SID-001 | 2 issues in_progress | No |
| WARN-SID-002 | 3+ issues in_progress | No |
| WARN-DEP-002 | Self-dependency | Yes (Error) |
| WARN-DEP-003 | Circular dependency | Yes (Error) |
| WARN-SESS-001 | Claiming blocked issue | No |
| WARN-SESS-002 | Closing with open blockers | No |
| WARN-SESS-003 | In-progress during ready | No |
| WARN-SYNC-001 | Uncommitted changes | No |
| WARN-SYNC-002 | Remote ahead | No |

This design provides a foundation for improving beads CLI user experience through gentle, educational guidance while maintaining command functionality.

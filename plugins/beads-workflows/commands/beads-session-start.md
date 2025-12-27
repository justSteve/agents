# Beads Session Start

Initialize a beads work session with proper rituals, work discovery, and context loading.

## Overview

This command orchestrates the complete session initialization workflow:
1. Environment verification
2. Database synchronization
3. Work discovery
4. Issue selection (using beads-workflow-orchestrator)
5. Issue claim
6. Context loading

## Arguments

**Optional**:
- `--issue <id>` - Specify issue ID to work on (skips work selection)
- `--type <type>` - Filter ready work by type (bug, feature, task, epic, chore)
- `--priority <0-4>` - Filter ready work by priority

## Phase 1: Environment Verification

Verify working environment and check for potential issues.

**Commands to run**:
```bash
pwd
git status
```

**Validation**:
- Verify current directory is expected workspace
- Check for uncommitted changes:
  - Clean: ✅ Proceed
  - Uncommitted changes: ⚠️ Ask user if intentional
  - Merge conflicts: ❌ Must resolve first

**Output**:
```
✓ Environment verified
  Working directory: <path>
  Git status: <clean|uncommitted changes|conflicts>
```

## Phase 2: Database Synchronization

Import remote updates to ensure working with latest issue state.

**Commands to run**:
```bash
bd sync
```

**Validation**:
- If sync succeeds: ✅ Proceed
- If merge conflicts in `.beads/issues.jsonl`:
  - Show conflict details
  - Provide resolution guidance
  - Wait for user to resolve

**Output**:
```
✓ Database synchronized
  Imported: <count> new/updated issues
  Exported: <count> local changes
```

## Phase 3: Check for Orphaned Work

Verify no in-progress issues from previous sessions.

**Commands to run**:
```bash
bd list --status in_progress --json
```

**Validation**:
- If 0 in_progress: ✅ Proceed to work discovery
- If 1 in_progress:
  - Show issue details
  - Ask: "Resume this work or start fresh?"
  - If resume: Skip to Phase 6 (context loading)
  - If fresh: Ask to update status first
- If 2+ in_progress: ⚠️ Single-issue discipline violation
  - Show all in_progress issues
  - Require user to clean up before proceeding

**Output** (if orphan found):
```
⚠️  Found in-progress work from previous session

Issue: <id> - <title>
Started: <timestamp>
Last updated: <timestamp>

Options:
[1] Resume this work (recommended)
[2] Close this issue and start fresh
[3] Mark as blocked and start new work

Your choice: _
```

## Phase 4: Work Discovery

Find unblocked issues ready to work on.

**Commands to run**:
```bash
bd ready --json
```

**Parse results**:
- If `--type` argument provided: Filter by issue type
- If `--priority` argument provided: Filter by priority
- If both provided: Apply both filters

**Validation**:
- If 0 ready issues:
  - Run `bd blocked` to show blocking chains
  - Suggest creating new work or unblocking existing
  - Exit session start
- If 1+ ready issues: ✅ Proceed to selection

**Output**:
```
✓ Ready work discovered
  Total: <count> issues
  P0: <count> | P1: <count> | P2: <count> | P3: <count> | P4: <count>
  Bugs: <count> | Features: <count> | Tasks: <count>
```

**Output** (if no ready work):
```
⚠️  No ready work available

Current status:
• Total open: <count>
• Ready: 0
• Blocked: <count>
• In progress: <count>

Top blockers:
• <id> - <title> (blocks <count> issues)
• <id> - <title> (blocks <count> issues)

Suggestions:
1. Run: bd blocked (view full blocking chain)
2. Continue in-progress work
3. Create new work: bd create "..."

Session start aborted.
```

## Phase 5: Work Selection

Use beads-workflow-orchestrator agent to select optimal issue.

**If `--issue <id>` provided**:
- Verify issue exists and is ready
- If ready: Select it
- If blocked: Show blockers and exit
- If not found: Error and exit

**If no issue specified**:
Use Task tool to invoke beads-workflow-orchestrator agent:

**Agent prompt**:
```
You are helping select which ONE issue to work on.

Ready issues (JSON):
<bd ready --json output>

User preferences:
<--type and --priority filters if provided>

Use your work selection algorithm to recommend ONE issue.

Provide:
1. Recommended issue ID and title
2. Rationale (priority, type, complexity, leverage)
3. Top 2 alternative options (if available)

Output format:
Recommended: <id> - <title>

Rationale:
• Priority: P<N> (<reason>)
• Type: <type> (<reason>)
• Impact: <count> issues blocked by this
• Complexity: <low|medium|high>

Alternatives:
1. <id> - <title> (P<N>, <type>)
2. <id> - <title> (P<N>, <type>)
```

**Parse agent response**:
- Extract recommended issue ID
- Present to user with alternatives
- Ask for confirmation or alternative selection

**User interaction**:
```
Recommended: <id> - <title>
<rationale from agent>

Proceed with recommended issue? [Y/n/1/2]
  Y - Work on recommended issue
  n - Abort session start
  1 - Work on alternative 1
  2 - Work on alternative 2
```

## Phase 6: Issue Claim

Mark selected issue as in_progress.

**Commands to run**:
```bash
bd update <selected-id> --status in_progress
```

**Validation**:
- If update succeeds: ✅ Proceed
- If fails (issue doesn't exist, already closed): ❌ Error and exit

**Double-check single-issue discipline**:
```bash
bd list --status in_progress --json
```

**Validation**:
- If exactly 1 in_progress (the one we just claimed): ✅ Proceed
- If 2+ in_progress: ⚠️ Violation - show warning but proceed

**Output**:
```
✓ Issue claimed
  ID: <id>
  Title: <title>
  Status: in_progress
```

## Phase 7: Context Loading

Load comprehensive context for the selected issue.

**Commands to run**:
```bash
bd show <selected-id> --json
bd dep tree <selected-id>
```

**Parse and display**:

### 7a. Issue Details
```
Issue: <id> - <title>

Description:
<description text>

Metadata:
• Type: <type>
• Priority: P<priority>
• Created: <timestamp>
• Updated: <timestamp>
```

### 7b. Dependencies
```
Dependencies:
• Blocking this issue: <count>
  - <id>: <title>
  - <id>: <title>

• Blocked by this issue: <count>
  - <id>: <title> (will unblock when complete)
  - <id>: <title> (will unblock when complete)
```

### 7c. Related Issues

**If `discovered-from` dependency exists**:
```bash
bd show <parent-id> --json
```

Display:
```
Related Work:
• Discovered from: <parent-id> - <parent-title>
  Context: <parent-description snippet>
```

### 7d. Skill Recommendations

Based on issue characteristics, recommend skills to load:

**Logic**:
- Description mentions "depends", "blocks", "requires" → `dependency-thinking`
- Issue type is bug with <100 char description → `description-quality`
- Multiple related issues or parent-child → `session-rituals`
- Issue is part of epic → Load epic skill if available

**Present recommendations**:
```
Recommended Skills:
• dependency-thinking: Issue involves dependency management
• description-quality: Ensure clear issue descriptions

Load recommended skills? [Y/n/select]
```

**If user confirms**:
For each skill, output:
```
Loading skill: dependency-thinking
<skill content loaded into context>
```

## Phase 8: Session Summary

Provide complete session start summary.

**Output**:
```
✅ Session started successfully

Working on: <id> - <title>
Priority: P<priority> | Type: <type>

Context loaded:
✓ Issue details and description
✓ Dependencies (<count> blocking, <count> blocked by)
✓ Related issues (<count>)
✓ Skills: <skill-list>

Next steps:
1. Review issue description above
2. Check dependencies (what's blocking, what this unblocks)
3. Begin implementation

When done:
• File discovered work: bd create --deps discovered-from:<id>
• Close issue: bd close <id> --reason "..."
• Sync: bd sync && git push

Happy coding! 🚀
```

## Error Handling

### Error: bd not installed

**Detection**: `bd` command not found

**Output**:
```
❌ Beads CLI not installed

This command requires beads (https://github.com/steveyegge/beads)

Installation:
1. Clone: git clone https://github.com/steveyegge/beads
2. Build: cd beads && go build -o bd ./cmd/bd
3. Add to PATH

Or see: https://github.com/steveyegge/beads#installation

Session start aborted.
```

### Error: Not in beads workspace

**Detection**: `.beads/beads.db` doesn't exist

**Output**:
```
❌ Not in a beads workspace

This directory has not been initialized with beads.

Initialize now:
bd init --prefix <project-prefix>
bd hooks install

Or cd to a directory with existing beads workspace.

Session start aborted.
```

### Error: Git repository issues

**Detection**: Merge conflicts, detached HEAD, etc.

**Output**:
```
❌ Git repository issues detected

<specific issue description>

Please resolve git issues before starting beads session:
• Merge conflicts: Resolve in .beads/issues.jsonl
• Detached HEAD: git checkout <branch>
• Uncommitted critical files: git add && git commit

Session start aborted.
```

## Success Criteria

Session start is successful when:
- ✅ Environment verified (git status checked)
- ✅ Database synchronized (bd sync completed)
- ✅ ONE issue claimed (status=in_progress)
- ✅ Context loaded (issue details, dependencies, skills)
- ✅ No orphaned in_progress issues (single-issue discipline)
- ✅ Agent ready to begin implementation

## Notes

**Performance**: This command prioritizes thoroughness over speed. Expect 30-60 seconds for complete session initialization due to:
- Database sync with remote
- Work discovery and parsing
- Agent invocation for work selection
- Context loading from multiple sources

**Integration**: This command is designed to work with:
- `beads-workflow-orchestrator` agent for intelligent work selection
- `beads-disciplinarian` agent for validation (future enhancement)
- Other beads-workflows commands (session-end, issue-create)

**Future enhancements**:
- `--quick` mode: Skip orchestrator, use simple priority sort
- `--resume` mode: Automatically resume last in_progress issue
- `--status` flag: Show session status without starting new work

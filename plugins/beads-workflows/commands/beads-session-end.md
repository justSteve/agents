# Beads Session End

Properly close a beads work session with structured handoff, ensuring all work is committed, synced, and tracked.

## Overview

This command orchestrates the complete session close workflow:
1. Check current work status
2. Handle issue completion or handoff
3. Run session close checklist
4. Verify clean state

## Arguments

**Optional**:
- `--close <id>` - Close a specific issue (mark as completed)
- `--reason "..."` - Reason for closing (used with --close)
- `--leave-open` - Leave in_progress issue open for future continuation
- `--handoff` - Prepare explicit handoff notes for next agent

## Phase 1: Check Current Work Status

Identify what work is currently in progress.

**Commands to run**:
```bash
bd list --status in_progress --json
git status
```

**Validation**:
- If 0 in_progress:
  - Check git status for uncommitted changes
  - If clean: Session is already ended
  - If uncommitted changes: ⚠️ Orphaned work - ask user what to do
- If 1 in_progress: ✅ Proceed with normal close
- If 2+ in_progress: ⚠️ Single-issue discipline violation
  - Show all in_progress issues
  - Must address each one before proceeding

**Output**:
```
Current session status:
• In-progress issues: <count>
  - <id>: <title> (started: <timestamp>)
• Uncommitted changes: <count> files

Proceeding with session close...
```

**Output** (if no work):
```
✅ No active work session detected

Git status: <clean|uncommitted changes>

Options:
[1] Run git cleanup only (commit/push orphaned changes)
[2] Exit - session already ended
```

## Phase 2: Handle Issue Completion

Determine whether to close or leave issues open.

### Option A: Close Issue (--close flag or user confirms)

**If closing with `--close <id>` or user chooses to close**:

**Commands to run**:
```bash
bd close <id> --reason "<reason>"
```

**Validation**:
- Verify issue was in_progress (warn if not)
- Verify close command succeeds

**Output**:
```
Closing issue: <id> - <title>
Reason: <reason or "completed">

✓ Issue closed
```

### Option B: Leave Open (--leave-open flag or user chooses)

**If leaving open for continuation**:

**Commands to run**:
```bash
bd update <id> --status=open --notes="Session ended. <handoff notes>"
```

**Validation**:
- Issue status updated to open (no longer in_progress)
- Handoff notes recorded

**Output**:
```
Pausing issue: <id> - <title>
Status: in_progress -> open

Notes recorded:
<handoff notes>

Issue will appear in bd ready for next session.
```

### Option C: Interactive Decision

**If neither --close nor --leave-open provided**:

**User interaction**:
```
Issue: <id> - <title>
Status: in_progress
Last updated: <timestamp>

What would you like to do?
[1] Close - Work is complete
[2] Leave open - Work to continue later
[3] Add handoff notes and leave open

Your choice: _
```

**If choice is [1]**:
```
Enter closing reason (or press Enter for "completed"):
> _
```

**If choice is [3]**:
```
Enter handoff notes for next agent:
> _
```

## Phase 3: Session Close Checklist

Execute the mandatory session close sequence.

**Critical**: This checklist must be run in order. Do not skip steps.

### Step 3a: Check Git Status

**Command**:
```bash
git status
```

**Validation**:
- Identify uncommitted changes
- Identify untracked files
- Check for staged vs unstaged changes

**Output**:
```
Git Status Check:
• Modified files: <count>
• Untracked files: <count>
• Staged files: <count>

Files to commit:
  - <file1>
  - <file2>
```

**If no changes**:
```
✓ No changes to commit
  Skipping to sync...
```

### Step 3b: Stage Changes

**Command** (if changes exist):
```bash
git add .
```

**Alternative** (if user wants selective staging):
```bash
git add <specific-files>
```

**Validation**:
- Verify files staged successfully
- Warn about any ignored files

**Output**:
```
✓ Staged <count> files for commit
```

### Step 3c: Pre-Commit Sync

**Command**:
```bash
bd sync
```

**Validation**:
- Sync should capture current issue state
- Handle sync conflicts if they occur

**Output**:
```
✓ Pre-commit sync complete
  Updated: <issue-id> (status captured)
```

### Step 3d: Create Commit

**Generate commit message** based on work done:

**Format**:
```
<type>(<scope>): <description>

<body - what was done>

Resolves: <issue-id>

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

**Command**:
```bash
git commit -m "<generated message>"
```

**Validation**:
- Commit succeeds
- Capture commit hash

**Output**:
```
✓ Committed: <short-hash>
  Message: <first line of commit message>
```

### Step 3e: Post-Commit Sync

**Command**:
```bash
bd sync
```

**Validation**:
- Sync links commit hash to issue
- Verify linkage succeeded

**Output**:
```
✓ Post-commit sync complete
  Linked: <issue-id> -> commit <short-hash>
```

### Step 3f: Push to Remote

**Command**:
```bash
git push
```

**Validation**:
- Push succeeds
- Handle push failures (behind remote, auth issues)

**Output**:
```
✓ Pushed to remote
  Branch: <branch-name>
  Commits: <count>
```

## Phase 4: Verify Clean State

Confirm session ended cleanly.

**Commands to run**:
```bash
git status
bd list --status in_progress
```

**Validation**:
- Git working tree is clean
- No issues in in_progress status

**Output** (success):
```
✅ Session ended cleanly

Final State:
• Git: clean working tree
• In-progress issues: 0
• Commits pushed: ✓

Issue Summary:
• <id>: <closed|paused> - <title>

Ready for next session!
```

**Output** (warnings):
```
⚠️  Session ended with warnings

Final State:
• Git: <clean|uncommitted changes>
• In-progress issues: <count>

Warnings:
• <warning description>

Please review before starting next session.
```

## Error Handling

### Error: Uncommitted Changes Won't Stage

**Detection**: `git add` fails or warnings appear

**Output**:
```
❌ Failed to stage changes

Error: <git error message>

Possible causes:
• File permissions issue
• .gitignore conflict
• Large file exceeding size limit

Resolution:
1. Check git error message above
2. Manually resolve the issue
3. Run: git add <files>
4. Re-run session end

Session end aborted at step 3b.
```

### Error: Sync Fails

**Detection**: `bd sync` returns error

**Output**:
```
❌ Beads sync failed

Error: <sync error message>

Possible causes:
• Merge conflict in .beads/issues.jsonl
• Database corruption
• File permission issue

Resolution:
1. Check .beads/ directory for conflicts
2. Resolve any merge markers
3. Run: bd sync
4. Re-run session end

Session end aborted at step 3c/3e.
```

### Error: Commit Fails

**Detection**: `git commit` returns error

**Output**:
```
❌ Commit failed

Error: <git error message>

Possible causes:
• Pre-commit hook failure
• Empty commit (no changes staged)
• Git configuration issue

Resolution:
1. Review error message above
2. Fix any hook issues
3. Ensure changes are staged: git status
4. Re-run: git commit -m "..."

Session end aborted at step 3d.
```

### Error: Push Fails

**Detection**: `git push` returns error

**Output**:
```
❌ Push failed

Error: <git error message>

Common causes and fixes:

If "rejected - non-fast-forward":
  Remote has new commits. Pull and merge:
  git pull --rebase
  git push

If "permission denied":
  Check your authentication:
  - SSH key configured?
  - Token valid?

If "remote not found":
  Verify remote is set:
  git remote -v

Session end paused at step 3f.
Retry: git push
```

### Error: Multiple In-Progress Issues

**Detection**: More than one issue with status in_progress

**Output**:
```
❌ Single-issue discipline violation

Found <count> issues in progress:
  - <id1>: <title1>
  - <id2>: <title2>
  ...

This violates single-issue discipline. You must resolve before ending session.

For each issue, choose one:
[1] Close (work complete)
[2] Pause (leave open for later)
[3] Mark as blocked

Resolve issue <id1> first:
Your choice: _
```

## Success Criteria

Session end is successful when:
- ✅ All in_progress issues addressed (closed or paused)
- ✅ All changes committed and pushed
- ✅ Beads state synced (before and after commit)
- ✅ Git working tree is clean
- ✅ No orphaned in_progress issues
- ✅ Handoff notes recorded (if applicable)

## Notes

**The "Never Skip" Rule**: Every step in the session close checklist exists to prevent common problems:
- Skipping git status → forgotten changes
- Skipping pre-commit sync → lost state
- Skipping post-commit sync → unlinked commits
- Skipping push → lost work

**Handoff Quality**: When leaving work for continuation, include:
- What was completed
- What remains to do
- Any blockers or concerns
- Relevant file locations or code references

**Performance**: Expect 15-30 seconds for complete session close:
- Git operations: 5-10s
- Sync operations: 5-10s
- User interaction: variable

**Integration**: This command works with:
- `beads-session-start` for session initialization
- `beads-issue-create` for filing discovered work before close
- `session-rituals` skill for detailed protocol guidance

**Future enhancements**:
- `--quick` mode: Skip user confirmations with defaults
- `--dry-run` mode: Show what would happen without executing
- Automatic handoff note generation from git diff

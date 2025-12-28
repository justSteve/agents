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

## Session Cleanup Verification

This section provides implementation-level verification to ensure complete session cleanup before claiming completion.

### Verification 1: In-Progress Issues Resolution

All in_progress issues must be resolved to either closed or open status.

**Verification Command**:
```bash
bd list --status in_progress
```

**Expected Output (PASS)**:
```
No issues found with status: in_progress

✓ All in_progress issues resolved
```

**Expected Output (FAIL)**:
```
Found 2 issues with status: in_progress:
  agents-abc: Some task title
  agents-xyz: Another task

❌ VERIFICATION FAILED: Orphaned in_progress issues

Each issue must be resolved:
- Close if complete: bd close <id>
- Return to open: bd update <id> --status=open

Session cleanup incomplete.
```

**Verification Logic**:
```python
def verify_in_progress_resolved():
    """Verify no issues remain in in_progress status."""

    result = run_command("bd list --status in_progress --json")
    issues = parse_json(result)

    if len(issues) == 0:
        return ('PASS', "All in_progress issues resolved")

    elif len(issues) == 1:
        # One issue is acceptable if it was just closed
        issue = issues[0]
        if recently_closed(issue):
            return ('PASS', f"Issue {issue['id']} recently closed")
        else:
            return ('FAIL', f"Issue {issue['id']} still in_progress",
                    f"Close with: bd close {issue['id']}")

    else:
        # Multiple in_progress is always a failure
        ids = [i['id'] for i in issues]
        return ('FAIL', f"Multiple issues in_progress: {ids}",
                "Resolve each before ending session")
```

**Resolution Prompt**:
```
Orphaned In-Progress Issue Found
=================================

Issue: <id> - <title>
Status: in_progress
Started: <timestamp>

This issue must be resolved before session can end.

Options:
[1] Close - Work is complete
    → bd close <id> --reason="..."

[2] Return to open - Work continues later
    → bd update <id> --status=open --notes="..."

[3] Mark blocked - Waiting on external factor
    → bd update <id> --status=blocked --reason="..."

Your choice: _
```

### Verification 2: bd sync Success

Both pre-commit and post-commit syncs must complete successfully.

**Verification Commands**:
```bash
# Pre-commit sync
bd sync
echo "Exit code: $?"

# Post-commit sync
bd sync
echo "Exit code: $?"
```

**Expected Output (PASS)**:
```
→ Exporting pending changes to JSONL...
→ No changes to commit
✓ Sync complete

Exit code: 0

✓ bd sync executed successfully
```

**Expected Output (FAIL)**:
```
→ Exporting pending changes to JSONL...
Error: Failed to write to .beads/issues.jsonl
Exit code: 1

❌ VERIFICATION FAILED: bd sync unsuccessful

Common causes:
- File permission issue
- Disk full
- .beads directory missing
- Merge conflict in issues.jsonl

Resolution:
1. Check .beads/ directory exists
2. Verify write permissions
3. Resolve any merge conflicts
4. Retry: bd sync
```

**Verification Logic**:
```python
def verify_sync_success():
    """Verify bd sync completes without errors."""

    # Pre-commit sync
    pre_result = run_command("bd sync")
    pre_exit = get_exit_code()

    if pre_exit != 0:
        return ('FAIL', "Pre-commit sync failed",
                f"Error: {pre_result}")

    # Post-commit sync (after git commit)
    post_result = run_command("bd sync")
    post_exit = get_exit_code()

    if post_exit != 0:
        return ('FAIL', "Post-commit sync failed",
                f"Error: {post_result}")

    return ('PASS', "Both syncs completed successfully")
```

**Sync Verification Display**:
```
Sync Status Verification
========================

Pre-commit sync:
  Command: bd sync
  Result: <success|failed>
  Output: <first 100 chars>

Post-commit sync:
  Command: bd sync
  Result: <success|failed>
  Output: <first 100 chars>

Verification: <PASS|FAIL>
```

### Verification 3: Git Clean State

Git working tree must be clean with all changes committed and pushed.

**Verification Commands**:
```bash
# Check for uncommitted changes
git status --porcelain

# Check for unpushed commits
git log origin/main..HEAD --oneline
```

**Expected Output (PASS)**:
```
$ git status --porcelain
(no output = clean)

$ git log origin/main..HEAD --oneline
(no output = all pushed)

✓ Git working tree clean
✓ All commits pushed
```

**Expected Output (FAIL - Uncommitted)**:
```
$ git status --porcelain
 M src/file.ts
?? new-file.js

❌ VERIFICATION FAILED: Uncommitted changes

Modified files:
  - src/file.ts

Untracked files:
  - new-file.js

Resolution:
1. Stage: git add .
2. Commit: git commit -m "..."
3. Push: git push
```

**Expected Output (FAIL - Unpushed)**:
```
$ git log origin/main..HEAD --oneline
abc1234 feat: some change
def5678 fix: another change

❌ VERIFICATION FAILED: Unpushed commits

2 commits not pushed to remote.

Resolution:
git push
```

**Verification Logic**:
```python
def verify_git_clean():
    """Verify git working tree is clean and all commits pushed."""

    # Check uncommitted changes
    status = run_command("git status --porcelain")
    if status.strip():
        files = status.strip().split('\n')
        return ('FAIL', f"Uncommitted changes: {len(files)} files",
                "Stage, commit, and push before claiming complete")

    # Check unpushed commits
    unpushed = run_command("git log origin/main..HEAD --oneline")
    if unpushed.strip():
        commits = unpushed.strip().split('\n')
        return ('FAIL', f"Unpushed commits: {len(commits)}",
                "Run: git push")

    return ('PASS', "Git clean: no uncommitted changes, all commits pushed")
```

**Git State Display**:
```
Git State Verification
======================

Working Tree:
  Modified:   <count> files
  Untracked:  <count> files
  Staged:     <count> files
  Status:     <clean|dirty>

Remote Sync:
  Branch:     <branch-name>
  Ahead:      <count> commits
  Behind:     <count> commits
  Status:     <synced|needs-push|needs-pull>

Verification: <PASS|FAIL>
```

### Verification 4: Discovered Work Filed

Any work discovered during the session should be filed as new issues with proper dependencies.

**Verification Approach**:

This verification is conversational - the agent must confirm no discovered work was left unfiled.

**Verification Prompt**:
```
Discovered Work Verification
============================

During this session, did you encounter any of the following?

• Bugs found while implementing the main issue
• Refactoring opportunities noticed but not addressed
• Technical debt identified
• Follow-up tasks that came to mind
• Edge cases that need separate handling

[Y] Yes, I found additional work
    → File each as a new issue before proceeding

[N] No, all discovered work has been filed
    → Verification complete

Your answer: _
```

**If Yes - Filing Guidance**:
```
Filing Discovered Work
======================

For each discovered item, file using:

bd create --title="<descriptive title>" \
  --description="## Context
Discovered while working on <current-issue-id>: <current-issue-title>

## Problem
<what you found>

## Notes
- Found in: <file:line or component>
- Priority: <your assessment>" \
  --deps discovered-from:<current-issue-id>

The --deps flag links this issue back to where it was discovered.

File discovered work now, then continue with session end.
```

**Verification Logic**:
```python
def verify_discovered_work_filed():
    """Verify all discovered work has been filed as issues."""

    # Check for recently created issues with discovered-from dependency
    recent_issues = run_command(
        "bd list --created-after='session-start' --json"
    )
    discovered = [i for i in parse_json(recent_issues)
                  if has_discovered_from_dep(i)]

    if discovered:
        return ('INFO', f"Filed {len(discovered)} discovered issues",
                [i['id'] for i in discovered])

    # Ask agent to confirm no unfiled work
    response = prompt_user("""
        Did you discover any work during this session that hasn't been filed?
        [Y/N]
    """)

    if response == 'Y':
        return ('INCOMPLETE', "Discovered work needs to be filed",
                "Use bd create with --deps discovered-from:<id>")

    return ('PASS', "No unfiled discovered work")
```

**Discovered Work Display**:
```
Discovered Work Summary
=======================

Issues filed during this session with discovered-from dependencies:

  <new-id-1>: <title>
    └─ discovered from: <parent-id>

  <new-id-2>: <title>
    └─ discovered from: <parent-id>

Total discovered work filed: <count>

Verification: <PASS|INCOMPLETE>
```

### Combined Cleanup Verification

Run all verifications before claiming session complete:

```
┌──────────────────────────────────────────────────────────────────────────┐
│                    SESSION CLEANUP VERIFICATION                          │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  1. In-Progress Issues                                                   │
│     Command: bd list --status in_progress                                │
│     Expected: 0 issues                                                   │
│     Status: [ ] PASS  [ ] FAIL                                           │
│                                                                          │
│  2. bd sync Success                                                      │
│     Commands: bd sync (pre-commit), bd sync (post-commit)                │
│     Expected: Exit code 0 for both                                       │
│     Status: [ ] PASS  [ ] FAIL                                           │
│                                                                          │
│  3. Git Clean State                                                      │
│     Commands: git status --porcelain, git log origin..HEAD               │
│     Expected: No output (clean, all pushed)                              │
│     Status: [ ] PASS  [ ] FAIL                                           │
│                                                                          │
│  4. Discovered Work Filed                                                │
│     Check: All found work filed with discovered-from dep                 │
│     Expected: User confirms none unfiled                                 │
│     Status: [ ] PASS  [ ] INCOMPLETE                                     │
│                                                                          │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ALL PASS: ✅ Session cleanup complete - safe to claim "done"            │
│  ANY FAIL: ❌ Resolve failures before claiming complete                  │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

**Final Verification Output**:
```
Session Cleanup Verification Complete
=====================================

✓ In-progress issues: 0 remaining (all resolved)
✓ bd sync: Pre-commit and post-commit successful
✓ Git state: Clean working tree, all commits pushed
✓ Discovered work: None unfiled (or N issues filed)

═══════════════════════════════════════════════════════════════════════════
                         SESSION END VERIFIED
═══════════════════════════════════════════════════════════════════════════

You may now claim "done" or "complete" for this session.

Summary:
- Issue closed: <id> - <title>
- Commits pushed: <count>
- Session duration: <time>
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

## Validation Checkpoints

This command enforces beads discipline through explicit validation checkpoints. Each checkpoint invokes the `beads-disciplinarian` agent for compliance validation.

### Checkpoint 1: Work Status Audit (Phase 1)

**Trigger**: At session end initiation

**Validation**:
```
Invoke beads-disciplinarian with context:
- All in_progress issues from bd list
- Uncommitted git changes

Expected response:
- PASS: 0-1 in_progress issues, git status known
- WARNING: Multiple in_progress (must address each)
- FAIL: Unknown state or critical issue
```

**On WARNING**: Force resolution of each in_progress issue before proceeding.

### Checkpoint 2: Session Close Checklist (Phase 3)

**Trigger**: Before executing close checklist

**Validation**:
```
Invoke beads-disciplinarian with context:
- Session close checklist template
- Request: Confirm all steps will be executed

Mandatory steps (NEVER skip):
- [ ] git status
- [ ] git add
- [ ] bd sync (pre-commit)
- [ ] git commit
- [ ] bd sync (post-commit)
- [ ] git push
- [ ] bd close (if applicable)

Expected response:
- PASS: Checklist understood, ready to execute
```

**On incomplete**: Block session end until all steps planned.

### Checkpoint 3: Clean State Verification (Phase 4)

**Trigger**: After all close steps complete

**Validation**:
```
Invoke beads-disciplinarian with context:
- Final git status
- Final bd list --status in_progress
- Push confirmation

Full compliance check:
- [ ] Git working tree clean
- [ ] No orphaned in_progress issues
- [ ] All changes pushed to remote
- [ ] Session end ritual complete

Expected response:
- PASS: Session ended cleanly
- WARNING: Ended with noted concerns
- FAIL: Critical step incomplete
```

**On FAIL**: Block completion claim until resolved.

### Agent Integration

When invoking beads-disciplinarian for validation:

```markdown
Validate session end for compliance:

Session state:
- In-progress issues: <list>
- Git status: <status>
- Uncommitted changes: <count>

Checklist execution:
- git status: <done|pending>
- git add: <done|pending>
- bd sync (pre): <done|pending>
- git commit: <done|pending>
- bd sync (post): <done|pending>
- git push: <done|pending>

Check:
1. Session end ritual completeness
2. No orphaned in_progress issues
3. All changes pushed

Return: PASS, WARNING, or FAIL with explanation

CRITICAL: Work is NOT done until git push completes!
```

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

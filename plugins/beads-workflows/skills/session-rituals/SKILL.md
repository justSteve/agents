---
name: session-rituals
description: Enforce proper session start and end procedures for beads work sessions. Use when starting work, claiming issues, or completing work to ensure consistent tracking and clean handoffs.
---

# Session Rituals

Enforce disciplined session start and end procedures to maintain accurate issue tracking and clean work handoffs in the beads system.

## When to Use This Skill

- Starting a new work session
- Claiming an issue to work on
- Completing work on an issue
- Taking a break or ending a session
- Handing off work to another agent

## Session Start Protocol

### Step 1: Find Available Work

```bash
bd ready
```

This shows all issues that are:
- Status: `open` (not in_progress or closed)
- Not blocked by incomplete dependencies
- Ready to be claimed

**Example output**:
```
Ready issues (unblocked, open):
  auth-login-fix     [task] [P1] Fix authentication timeout issue
  api-docs-update    [task] [P2] Update API documentation
  db-migration-v2    [task] [P1] Add user preferences table
```

### Step 2: Review Before Claiming

Before claiming an issue, review its details:

```bash
bd show <issue-id>
```

Check for:
- Clear description and acceptance criteria
- Reasonable scope for a single session
- Any soft dependencies or related issues
- Required context or prerequisites

### Step 3: Claim the Issue

```bash
bd update <issue-id> --status=in_progress
```

**Important**: Only claim ONE issue at a time. This is critical for:
- Accurate time tracking
- Clear accountability
- Preventing context-switching overhead
- Maintaining focus

### Step 4: Sync Tracking State

```bash
bd sync
```

This ensures the tracking system reflects your claimed work.

## Session End Protocol

**Critical Rule**: Never say "done" or "complete" without running this checklist.

### The Session Close Checklist

Execute these steps in order every time you finish work:

```bash
# 1. Check repository state
git status

# 2. Stage changes
git add <files>  # or: git add .

# 3. Sync beads state before commit
bd sync

# 4. Create commit with meaningful message
git commit -m "feat(scope): description

Resolves: <issue-id>"

# 5. Sync beads state after commit (captures commit hash)
bd sync

# 6. Push to remote
git push
```

### Why This Order Matters

| Step | Purpose |
|------|---------|
| `git status` | Verify all work is captured, no forgotten changes |
| `git add` | Stage intentional changes only |
| `bd sync` (pre-commit) | Capture in-progress state, validate issue status |
| `git commit` | Create atomic, traceable change |
| `bd sync` (post-commit) | Link commit hash to issue, update timestamps |
| `git push` | Share work with team, enable handoff |

### Completing an Issue

When work is fully done and verified:

```bash
# Update issue status to closed
bd close <issue-id>

# Final sync to capture completion
bd sync
```

**Before closing, verify**:
- All acceptance criteria met
- Tests passing (if applicable)
- No unintended side effects
- Documentation updated (if applicable)

## The "Never Say Done" Rule

**Problem**: Agents often claim completion without verification.

**Anti-pattern**:
```
Agent: "I've completed the authentication fix."
Reality: Changes uncommitted, not pushed, issue still in_progress
```

**Correct pattern**:
```
Agent: "Running session close checklist..."
Agent: "git status shows 3 modified files"
Agent: "Staged and committed: 'fix(auth): resolve timeout issue'"
Agent: "bd sync complete, issue linked to commit abc123"
Agent: "git push complete"
Agent: "bd close auth-login-fix"
Agent: "Session complete. Issue auth-login-fix is now closed."
```

## Session Interruption Protocol

If you must stop work mid-session:

### Planned Break

```bash
# Commit work-in-progress
git add .
git commit -m "wip(scope): partial progress on issue-id"

# Sync state
bd sync

# Push WIP branch
git push

# Document status
bd update <issue-id> --notes="WIP: Completed X, Y still needed. See commit abc123."
```

### Unexpected Interruption

If session ends unexpectedly, the next session should:

```bash
# Check current state
bd list --status=in_progress

# Review what was in progress
bd show <issue-id>

# Check git state
git status
git log -3
```

## Common Patterns

### Pattern 1: Clean Single-Issue Session

```bash
# Start
bd ready
bd show api-docs-update
bd update api-docs-update --status=in_progress
bd sync

# Work...

# End
git status
git add docs/api-reference.md
bd sync
git commit -m "docs(api): update endpoint documentation

Resolves: api-docs-update"
bd sync
git push
bd close api-docs-update
bd sync
```

### Pattern 2: Session with Discovered Work

```bash
# While working on issue X, discover issue Y
# DO NOT fix Y. File it:
bd create "Fix discovered bug in validation" \
  --description="Found null pointer in input validation while working on api-docs-update" \
  --deps discovered-from:api-docs-update

# Continue working on X only
# Complete X using normal close checklist
```

### Pattern 3: Handoff to Another Agent

```bash
# Current agent completing partial work
git add .
git commit -m "wip(auth): implement token refresh logic

Partial implementation. Remaining:
- Add token expiry check
- Handle refresh failures"
bd sync
git push

bd update auth-token-refresh --status=open \
  --notes="Partial implementation by Agent-A. See commit abc123. Remaining work documented in commit message."
bd sync

# Next agent picks up
bd ready  # Issue appears as ready
bd show auth-token-refresh  # Review notes and commits
bd update auth-token-refresh --status=in_progress
```

## Anti-Patterns to Avoid

### Anti-Pattern 1: Skipping git status

**Problem**: Forgetting uncommitted changes.

```bash
# Bad
bd close issue-id
# "Done!" - but changes not committed!
```

**Solution**: Always run `git status` first.

### Anti-Pattern 2: Claiming Multiple Issues

**Problem**: Working on several issues simultaneously.

```bash
# Bad
bd update issue-1 --status=in_progress
bd update issue-2 --status=in_progress
bd update issue-3 --status=in_progress
```

**Solution**: One issue at a time. File discovered work as new issues.

### Anti-Pattern 3: Skipping bd sync

**Problem**: Beads state gets out of sync with git.

```bash
# Bad
git commit -m "fix: something"
git push
bd close issue-id
# bd sync never called - no commit linkage!
```

**Solution**: Sync before commit, after commit, and after close.

### Anti-Pattern 4: Closing Without Verification

**Problem**: Marking done without checking work.

```bash
# Bad
bd close issue-id
# No verification that acceptance criteria met
```

**Solution**: Verify before closing. Run tests. Check acceptance criteria.

### Anti-Pattern 5: Verbal Completion Claims

**Problem**: Saying "done" without evidence.

```bash
# Bad
Agent: "I've fixed the bug and it's all done."
# No git commits, no bd commands shown
```

**Solution**: Show your work. Include command outputs in completion claims.

## Quick Reference Card

**Session Start**:
```bash
bd ready                           # Find available work
bd show <id>                       # Review issue details
bd update <id> --status=in_progress  # Claim issue
bd sync                            # Sync state
```

**Session End Checklist** (memorize this!):
```bash
git status                         # 1. Check state
git add <files>                    # 2. Stage changes
bd sync                            # 3. Pre-commit sync
git commit -m "..."                # 4. Commit
bd sync                            # 5. Post-commit sync
git push                           # 6. Push
bd close <id>                      # 7. Close issue
bd sync                            # 8. Final sync
```

**Interruption**:
```bash
git add . && git commit -m "wip: ..."
bd sync && git push
bd update <id> --notes="WIP status..."
```

**Golden Rule**:
> Never claim completion without running the checklist and showing the outputs.

## Pattern Enforcement Examples

Use these examples to detect and correct session ritual violations.

### Enforcement Example 1: Premature Completion Claim

**Agent says**: "I've fixed the bug. Done!"

**Detection**: Completion claimed without showing checklist execution.

**Enforcement dialogue**:
```
⚠️ PREMATURE COMPLETION CLAIM

You said "Done!" but haven't shown the session close checklist.

The "Never Say Done" rule: Show your work!

Required before claiming completion:
┌──────────────────────────────────────────────────────┐
│  SESSION CLOSE CHECKLIST                             │
├──────────────────────────────────────────────────────┤
│  [ ] git status         → Show what changed          │
│  [ ] git add            → Stage the changes          │
│  [ ] bd sync            → Pre-commit sync            │
│  [ ] git commit         → Create commit              │
│  [ ] bd sync            → Post-commit sync           │
│  [ ] git push           → Push to remote             │
│  [ ] bd close <id>      → Close the issue            │
└──────────────────────────────────────────────────────┘

Run the checklist now and show outputs:

git status
# Expected: Shows modified files from your work

git add <files>
# Expected: Files staged

bd sync
# Expected: Sync complete

git commit -m "fix(auth): resolve timeout issue

Resolves: auth-login-fix"
# Expected: Commit created with hash

bd sync
# Expected: Commit linked to issue

git push
# Expected: Pushed to origin

bd close auth-login-fix
# Expected: Issue closed

ONLY THEN can you say "Done!"
```

### Enforcement Example 2: Multiple Issues In Progress

**Agent runs**:
```
bd update issue-1 --status=in_progress
bd update issue-2 --status=in_progress
```

**Detection**: Single-issue discipline violation.

**Enforcement dialogue**:
```
⚠️ SINGLE-ISSUE DISCIPLINE VIOLATION

You have claimed multiple issues:
- issue-1: in_progress
- issue-2: in_progress

Beads discipline: ONE issue at a time!

Problems with multiple in_progress:
- Context switching overhead
- Unclear which issue changes belong to
- Inaccurate time tracking
- Risk of orphaned work

Resolution options:

Option 1: Keep issue-1, release issue-2
bd update issue-2 --status=open

Option 2: Keep issue-2, release issue-1
bd update issue-1 --status=open

Option 3: If issue-2 was discovered while working on issue-1
bd update issue-2 --status=open
bd dep add issue-2 issue-1 --type discovered-from
# Continue working on issue-1 only

Choose one issue and focus on it completely.
```

### Enforcement Example 3: Skipping bd sync

**Agent runs**:
```
git add .
git commit -m "fix something"
git push
bd close issue-id
```

**Detection**: Missing bd sync before and after commit.

**Enforcement dialogue**:
```
⚠️ BD SYNC SKIPPED

Your session close sequence is missing bd sync calls:

What you did:
1. git add .          ✓
2. git commit         ✓
3. git push           ✓
4. bd close           ✓

What you missed:
- bd sync BEFORE commit (captures in-progress state)
- bd sync AFTER commit (links commit hash to issue)

Why bd sync matters:
- Pre-commit: Validates issue state, exports to JSONL
- Post-commit: Links git commit hash to beads issue
- Without sync: Issue history is incomplete

Correct sequence:
git add .
bd sync              ← BEFORE commit
git commit -m "..."
bd sync              ← AFTER commit
git push
bd close issue-id
bd sync              ← AFTER close

Always sync: before commit, after commit, after close.
```

### Enforcement Example 4: Closing Without Verification

**Agent runs**: `bd close issue-id` immediately after making changes.

**Detection**: No verification that work is complete.

**Enforcement dialogue**:
```
⚠️ CLOSING WITHOUT VERIFICATION

You closed the issue without verifying completion.

Before bd close, you must verify:

1. ACCEPTANCE CRITERIA MET
   bd show issue-id
   # Review each acceptance criterion
   # All boxes should be checkable

2. TESTS PASSING (if applicable)
   npm test  # or your test command
   # All tests should pass

3. NO UNINTENDED SIDE EFFECTS
   git diff HEAD~1
   # Review changes are intentional
   # No debug code left behind

4. CHANGES COMMITTED AND PUSHED
   git status
   # Should show: "nothing to commit, working tree clean"
   git log origin/main..HEAD
   # Should show: 0 commits (already pushed)

Only after all verifications pass:
bd close issue-id --reason="All criteria met, tests pass"
```

### Enforcement Example 5: Orphaned In-Progress Issue

**Agent starts new session and finds**:
```
bd list --status in_progress
# old-issue-from-yesterday: in_progress
```

**Detection**: Previous session left issue in_progress.

**Enforcement dialogue**:
```
⚠️ ORPHANED IN-PROGRESS ISSUE FOUND

Previous session left work incomplete:
- Issue: old-issue-from-yesterday
- Status: in_progress
- Last updated: yesterday

This violates clean handoff principles.

Resolution required before starting new work:

Option 1: Resume the work
bd show old-issue-from-yesterday
# Review what was done
# Continue and complete it

Option 2: Work was actually completed
git log --oneline -5
# Check if work was committed
bd close old-issue-from-yesterday --reason="Completed in previous session"

Option 3: Work was abandoned
bd update old-issue-from-yesterday --status=open \
  --notes="Abandoned in previous session. Partial work in commit abc123."
# Now it's available for future sessions

Option 4: Work is blocked
bd update old-issue-from-yesterday --status=blocked \
  --reason="Waiting on external dependency X"

Choose an option to clean up before claiming new work.
```

### Enforcement Example 6: Verbal Completion Without Evidence

**Agent says**: "I've completed all the tasks and everything is working."

**Detection**: Completion claimed without command output evidence.

**Enforcement dialogue**:
```
⚠️ VERBAL COMPLETION WITHOUT EVIDENCE

You claimed completion but showed no evidence.

The "Show Your Work" rule requires command outputs:

❌ What you said:
"I've completed all the tasks and everything is working."

✅ What you should show:

$ git status
On branch main
nothing to commit, working tree clean

$ git log -1 --oneline
abc123 fix(auth): resolve timeout issue

$ bd show auth-login-fix
Status: closed
Closed: 2024-01-15T14:30:00Z
Commits: abc123

$ npm test
All 47 tests passed

Completion claim accepted only with evidence.
Run the commands and show the outputs.
```

### Session Ritual Verification Commands

Use these commands to verify session state:

**Start of Session**:
```bash
# Check for orphaned work
bd list --status in_progress
# Expected: 0 issues (or 1 you're resuming)

# Find available work
bd ready
# Expected: List of unblocked issues

# Claim one issue
bd update <id> --status in_progress
bd sync
```

**End of Session**:
```bash
# Verify git state
git status
# Expected: All work staged or committed

# Verify push state
git log origin/main..HEAD
# Expected: 0 commits (all pushed)

# Verify beads state
bd list --status in_progress
# Expected: 0 issues (all closed or released)
```

### Enforcement Summary Checklist

```
┌────────────────────────────────────────────────────────────┐
│              SESSION RITUAL ENFORCEMENT                    │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  SESSION START                                             │
│  [ ] bd ready              → Find available work           │
│  [ ] bd show <id>          → Review before claiming        │
│  [ ] Only ONE in_progress  → Single-issue discipline       │
│  [ ] bd sync               → Sync state                    │
│                                                            │
│  SESSION END (MANDATORY - NEVER SKIP)                      │
│  [ ] git status            → Verify all work captured      │
│  [ ] git add               → Stage changes                 │
│  [ ] bd sync               → Pre-commit sync               │
│  [ ] git commit            → Create commit                 │
│  [ ] bd sync               → Post-commit sync              │
│  [ ] git push              → Push to remote                │
│  [ ] bd close              → Close issue (if complete)     │
│  [ ] bd sync               → Final sync                    │
│                                                            │
│  COMPLETION CLAIM                                          │
│  [ ] Show command outputs  → Evidence required             │
│  [ ] Verify clean state    → git status shows clean        │
│  [ ] No orphaned work      → bd list shows 0 in_progress   │
│                                                            │
└────────────────────────────────────────────────────────────┘

"Done" = All boxes checked with evidence shown
```

## Summary

**Core principle**: Session discipline enables accurate tracking and clean handoffs.

**Key rituals**:
1. Start: `bd ready` → `bd show` → `bd update --status=in_progress` → `bd sync`
2. End: `git status` → `git add` → `bd sync` → `git commit` → `bd sync` → `git push` → `bd close` → `bd sync`
3. Never claim "done" without running the full checklist
4. One issue at a time - file discovered work as new issues

**Remember**:
- The checklist exists to prevent silent failures
- Syncing twice (pre and post commit) ensures linkage
- "Done" means checklist completed, not just "I think I'm done"
- Document handoffs clearly for the next agent

Master these rituals and you'll maintain a clean, traceable work history!

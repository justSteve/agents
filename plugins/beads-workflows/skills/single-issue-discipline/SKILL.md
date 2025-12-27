---
name: single-issue-discipline
description: Prevent scope creep by maintaining strict focus on the claimed issue. Use while working on an issue to resist the temptation to fix tangential problems, instead filing them as new issues with proper dependency tracking.
---

# Single Issue Discipline

Maintain laser focus on the claimed issue. When you discover related work, file it as a new issue instead of fixing it inline. This discipline is critical for tracking, accountability, and preventing scope creep.

## When to Use This Skill

- While actively working on a claimed issue
- When you notice a bug or improvement opportunity during implementation
- When tempted to "just quickly fix" something unrelated
- When refactoring reveals additional technical debt
- When tests expose issues outside the current scope

## The Core Principle

**One issue, one focus, one completion.**

When you claim issue X:
- Work on issue X only
- Complete issue X fully
- Close issue X

If you discover issue Y while working on X:
- **DO NOT** fix Y
- **DO** file Y as a new issue with `discovered-from` dependency
- **DO** continue working on X

## Why This Matters

### For Tracking

Each issue represents a unit of trackable work:
- Time spent is attributable
- Changes are traceable
- Progress is measurable

When you mix issues, you lose:
- Accurate time tracking
- Clear commit history
- Ability to revert specific changes

### For Accountability

When an issue is claimed:
- You are accountable for that scope
- Stakeholders know what's being worked on
- Other agents know what's taken

When you silently expand scope:
- Accountability becomes unclear
- Estimates become meaningless
- Coordination breaks down

### For Quality

Focused work produces better results:
- Full attention on one problem
- Complete solution, not partial fixes
- Proper testing of the change

Scattered work produces technical debt:
- Half-fixed issues
- Untested drive-by changes
- Forgotten side effects

## The Discovery Protocol

When you discover work while implementing:

### Step 1: Recognize the Discovery

**Signals that you've found new work**:
- "While I'm here, I should also..."
- "This would be cleaner if I also..."
- "I noticed this other bug..."
- "This code is messy, let me refactor..."
- "The tests are failing for a different reason..."

**Stop and evaluate**: Is this part of your claimed issue?

### Step 2: Evaluate Relevance

Ask yourself:

**Is this directly required to complete my claimed issue?**
- If YES: It's part of the current scope
- If NO: It's discovered work - file it

**Examples**:

| Current Issue | Discovery | Part of Scope? |
|--------------|-----------|----------------|
| Fix login timeout | Auth token format is wrong | YES - directly related |
| Fix login timeout | Password reset is also broken | NO - file new issue |
| Add user profile page | User model missing field | MAYBE - depends on field |
| Add user profile page | Admin panel needs same fix | NO - file new issue |
| Refactor auth module | Found SQL injection | NO - security issue, file immediately |

### Step 3: File Discovered Work

Use `bd create` with the `discovered-from` dependency:

```bash
bd create "Fix password reset token expiry" \
  --description="Found while working on login-timeout. Reset tokens use same flawed expiry logic." \
  --deps discovered-from:login-timeout \
  --type task \
  --priority 2
```

**Key elements**:
- **Clear title**: What needs to be done
- **Description**: Include discovery context
- **Dependency**: `discovered-from:<current-issue>` links the discovery
- **Priority**: Assess independently of current work

### Step 4: Continue Original Work

After filing, return focus to your claimed issue:

```bash
# Verify you're still on track
bd show <claimed-issue>

# Continue implementation
# ...
```

**Do not**:
- Switch to the new issue
- "Just quickly" fix it first
- Delay current work to assess the discovery

## Common Discovery Scenarios

### Scenario 1: Bug in Adjacent Code

**Situation**: While implementing feature X, you notice function Y has a bug.

```python
# Working on: add-user-preferences

def save_preferences(user_id, prefs):
    # Your new code
    validate_preferences(prefs)

    # You notice this existing function has a bug!
    user = get_user(user_id)  # BUG: No null check!
    user.preferences = prefs
    save(user)
```

**Wrong approach**:
```python
# "I'll just fix it while I'm here"
def save_preferences(user_id, prefs):
    validate_preferences(prefs)
    user = get_user(user_id)
    if user is None:  # Drive-by fix
        raise UserNotFoundError(user_id)
    user.preferences = prefs
    save(user)
```

**Right approach**:
```bash
# File the discovery
bd create "Fix null user handling in get_user callers" \
  --description="get_user() can return None but callers don't check. Found in save_preferences while implementing add-user-preferences. Multiple call sites affected." \
  --deps discovered-from:add-user-preferences \
  --type bug \
  --priority 1
```

Then continue with your original implementation, working around the bug:
```python
def save_preferences(user_id, prefs):
    validate_preferences(prefs)
    user = get_user(user_id)
    # Workaround for known bug (tracked in fix-null-user-handling)
    if user is None:
        raise UserNotFoundError(user_id)
    user.preferences = prefs
    save(user)
```

### Scenario 2: Refactoring Temptation

**Situation**: The code you're modifying is messy and could use refactoring.

**Wrong approach**:
```bash
# "While I'm here, let me clean this up"
# 3 hours later...
# Original issue still incomplete
# Massive refactor with unclear scope
# Tests broken in unexpected ways
```

**Right approach**:
```bash
# File refactoring as separate work
bd create "Refactor UserService for clarity" \
  --description="UserService has grown to 800 lines with mixed responsibilities. Found while implementing add-user-preferences. Consider splitting into UserPreferencesService." \
  --deps discovered-from:add-user-preferences \
  --type task \
  --priority 3

# Complete original issue with minimal changes
# Refactoring happens later, with proper scope
```

### Scenario 3: Test Reveals Unrelated Failure

**Situation**: Running tests for your change, an unrelated test fails.

**Wrong approach**:
```bash
# "I need all tests to pass, so I'll fix this one too"
# Fix unrelated test
# Commit mixed changes
# No tracking of the test fix
```

**Right approach**:
```bash
# File the test failure
bd create "Fix flaky user_auth_test" \
  --description="test_token_refresh intermittently fails with timing issue. Discovered when running test suite for add-user-preferences." \
  --deps discovered-from:add-user-preferences \
  --type bug \
  --priority 2

# For current work, skip or mark known-flaky if needed
pytest --ignore=tests/test_auth.py  # Run other tests
# Or document in PR that unrelated test is flaky
```

### Scenario 4: Security Issue Found

**Situation**: You discover a security vulnerability.

**Special handling**: Security issues may require immediate attention, but still file them properly:

```bash
# File immediately with high priority
bd create "SECURITY: SQL injection in search endpoint" \
  --description="User input not sanitized in /api/search. Discovered while implementing add-user-preferences. CRITICAL - allows arbitrary SQL execution." \
  --deps discovered-from:add-user-preferences \
  --type bug \
  --priority 0

# Notify appropriate channels (security team, etc.)
# Continue current work OR pause if security issue is critical
```

**Decision tree for security**:
- P0 actively exploited: Pause current work, fix immediately
- P0 not actively exploited: File, notify, continue current work
- P1-P2: File, continue current work

### Scenario 5: Scope Question

**Situation**: You're unsure if something is part of your issue.

**Decision framework**:

```
Is this explicitly mentioned in acceptance criteria?
├─ YES → It's in scope, do it
└─ NO → Is it technically required to complete the criteria?
         ├─ YES → It's in scope, do it
         └─ NO → Is it an improvement/nice-to-have?
                  ├─ YES → File as new issue, don't do it
                  └─ NO → Ask for clarification if truly ambiguous
```

## Anti-Patterns to Avoid

### Anti-Pattern 1: The "While I'm Here" Fix

**Problem**: Making small fixes to unrelated code while in the area.

```bash
# Working on issue: add-validation
# "While I'm in this file, I'll also fix this typo, update that comment,
# rename that variable, and add logging..."
```

**Why it's harmful**:
- Changes are untracked
- Code review scope becomes unclear
- Harder to revert specific changes
- Commit history becomes noisy

**Solution**: File each fix as its own issue or let it go.

### Anti-Pattern 2: The Scope Snowball

**Problem**: Discovered work triggers more discovered work.

```
Working on: A
Discover B while working on A
Start fixing B
Discover C while fixing B
Start fixing C
...
Original issue A never completed
```

**Solution**: Always file and return. Never chase the rabbit.

### Anti-Pattern 3: The "Quick" Fix

**Problem**: Underestimating discovered work.

```
"This is a tiny fix, I'll just do it real quick"
# 2 hours later...
# Turns out it wasn't quick
# Original issue blocked
# No tracking of time spent
```

**Solution**: If it's quick enough to not track, file it anyway. If it's worth doing, it's worth tracking.

### Anti-Pattern 4: The Silent Fix

**Problem**: Fixing things without any record.

```
# Just fix it, no one needs to know
# No commit message mentioning it
# No issue tracking it
# Future maintainers confused
```

**Solution**: Everything gets tracked. No silent fixes.

### Anti-Pattern 5: The Perfectionist Trap

**Problem**: Refusing to complete an issue until related problems are also fixed.

```
"I can't close this issue because the adjacent code is still messy"
"The feature works but the tests could be better"
"This is done but I noticed three other things..."
```

**Solution**: Complete what you claimed. File the rest. Perfect is the enemy of done.

## The Discipline Mindset

### Mantras

Repeat when tempted to expand scope:

- "File it, don't fix it"
- "One issue, one focus"
- "Tracked work is accountable work"
- "Future me will thank present me for filing this"

### Self-Check Questions

Before making any change, ask:

1. Is this explicitly required by my claimed issue?
2. Would this appear in my issue's acceptance criteria?
3. Am I changing code outside my issue's described scope?

If answers are NO, NO, YES - you're about to break discipline.

### The 30-Second Rule

If you can describe a discovered problem in 30 seconds, you can file it in 60 seconds:

```bash
bd create "Brief description" \
  --description="One sentence context. Found while working on X." \
  --deps discovered-from:current-issue \
  --type task
```

Filing is faster than fixing. Always.

## Quick Reference Card

**When you discover new work**:
```bash
# 1. Recognize: "This isn't part of my issue"
# 2. File immediately:
bd create "<title>" \
  --description="<context>. Found during <current-issue>." \
  --deps discovered-from:<current-issue>

# 3. Return focus:
bd show <current-issue>
# Continue original work
```

**Decision tree**:
```
Is it in my acceptance criteria?
├─ YES → Do it
└─ NO → File it (bd create --deps discovered-from:...)
```

**Discipline check**:
- Am I working on exactly one issue? ✓
- Are all my commits related to that issue? ✓
- Did I file discoveries instead of fixing them? ✓

**Golden rule**:
> Complete what you claimed. File what you found. Never mix the two.

## Summary

**Core principle**: One issue, complete focus, file everything else.

**Key practices**:
1. Work only on your claimed issue
2. Recognize when you've found new work
3. File discoveries with `bd create --deps discovered-from:<current>`
4. Return focus immediately after filing
5. Complete and close your original issue

**Why this matters**:
- Accurate tracking and accountability
- Clean commit history
- Predictable work completion
- Reduced cognitive overhead

**Remember**:
- The temptation to fix "just one more thing" is the enemy
- Filing takes 60 seconds; scope creep takes hours
- Discipline enables trust and velocity
- Future agents will thank you for proper tracking

Master single-issue discipline and you'll ship faster with cleaner history!

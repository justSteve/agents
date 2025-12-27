---
name: dependency-thinking
description: Master causal dependency reasoning to avoid temporal thinking traps when adding beads dependencies. Use when planning work sequences, adding blocking relationships, or structuring task hierarchies.
---

# Dependency Thinking

Master the art of causal dependency reasoning to avoid the temporal thinking trap that inverts dependency relationships in beads issue tracking.

## When to Use This Skill

- Planning work sequences and task ordering
- Adding blocking dependencies between issues
- Structuring task hierarchies (epics, features, tasks)
- Reviewing dependency graphs for correctness
- Debugging why tasks are unexpectedly blocked or unblocked

## The Cognitive Trap

**The Problem**: Temporal language triggers incorrect mental models.

When you think "Phase 1 comes before Phase 2", your brain naturally wants to express this as:
```bash
bd dep add phase1 phase2  # ❌ WRONG!
```

This says "phase1 NEEDS phase2" - exactly backwards!

**Why this happens**: The words "before", "first", "Phase 1", "Step 1" activate temporal reasoning:
- "X happens before Y"
- "X is a prerequisite for Y"
- Your brain shortcuts this to: "X blocks Y"

But beads dependency syntax is: `bd dep add DEPENDENT REQUIRED`
- Which reads: "DEPENDENT needs REQUIRED"
- The first argument is what gets blocked, not what blocks!

## The Correct Mental Model

**Think in requirements, not sequences:**

Ask yourself: **"Which issue NEEDS the other to be completed first?"**

### Example 1: Buffer and Rendering

**Temporal thinking (wrong)**:
- "First create buffer, then add rendering"
- "Buffer comes before rendering"
- Brain says: `bd dep add buffer rendering` ❌

**Requirement thinking (correct)**:
- "What does rendering NEED?" → "Rendering needs buffer"
- Brain says: `bd dep add rendering buffer` ✅

### Example 2: Authentication Flow

**Temporal thinking (wrong)**:
- "Step 1: Create user model"
- "Step 2: Add login endpoint"
- Brain says: `bd dep add user-model login-endpoint` ❌

**Requirement thinking (correct)**:
- "What does login NEED?" → "Login needs user model"
- Brain says: `bd dep add login-endpoint user-model` ✅

### Example 3: Database Migration

**Temporal thinking (wrong)**:
- "Phase 1: Write migration"
- "Phase 2: Update ORM models"
- "Phase 3: Update API handlers"
- Brain says:
  ```bash
  bd dep add migration orm-models      # ❌
  bd dep add orm-models api-handlers   # ❌
  ```

**Requirement thinking (correct)**:
- "ORM models NEED migration (schema must exist)"
- "API handlers NEED ORM models (to query data)"
- Brain says:
  ```bash
  bd dep add orm-models migration      # ✅
  bd dep add api-handlers orm-models   # ✅
  ```

## The Mental Shortcut

**Replace temporal words with requirement words:**

| Temporal (Wrong) | Requirement (Right) | Command |
|------------------|---------------------|---------|
| "X before Y" | "Y needs X" | `bd dep add Y X` |
| "X then Y" | "Y requires X" | `bd dep add Y X` |
| "First X, later Y" | "Y depends on X" | `bd dep add Y X` |
| "X → Y" | "Y ← X" | `bd dep add Y X` |

## Verification Strategy

After adding dependencies, **always verify** using `bd blocked`:

```bash
# Add dependency
bd dep add rendering buffer

# Verify
bd blocked

# Expected output:
# rendering is blocked by:
#   - buffer (blocks)
```

**Verification checklist**:
- ✅ Blocked tasks should be waiting for their **prerequisites**
- ✅ Blocking tasks should be the ones that must be **done first**
- ❌ If dependencies look inverted, you used temporal thinking

## Dependency Types in Beads

Beads supports four dependency types, each with different semantics:

### 1. `blocks` - Hard Blocker

**Syntax**: `bd dep add <dependent> <required> --type blocks`

**Meaning**: `<dependent>` cannot start until `<required>` is complete.

**Use when**: Strict sequential dependency (technical blocker).

**Example**:
```bash
# API handlers can't be written until database schema exists
bd dep add api-handlers db-schema --type blocks
```

### 2. `parent-child` - Hierarchical

**Syntax**: `bd dep add <child> <parent> --type parent-child`

**Meaning**: `<child>` is a sub-task of `<parent>` (epic/feature decomposition).

**Use when**: Breaking down large issues into smaller tasks.

**Example**:
```bash
# "Add login endpoint" is part of "User authentication" epic
bd dep add add-login-endpoint user-authentication --type parent-child
```

### 3. `discovered-from` - Work Discovery

**Syntax**: `bd dep add <new-issue> <original-issue> --type discovered-from`

**Meaning**: `<new-issue>` was discovered while working on `<original-issue>`.

**Use when**: Filing tangential work found during implementation.

**Example**:
```bash
# While implementing login, discovered auth token expiry bug
bd dep add fix-token-expiry add-login-endpoint --type discovered-from
```

### 4. `related` - Soft Reference

**Syntax**: `bd dep add <issue1> <issue2> --type related`

**Meaning**: Issues are related but not blocking (informational link).

**Use when**: Cross-referencing related work without enforcing order.

**Example**:
```bash
# Login and password reset both touch auth system
bd dep add password-reset login-endpoint --type related
```

## Common Patterns

### Pattern 1: Sequential Pipeline

**Scenario**: Tasks must be done in strict order.

```bash
# Task flow: Design → Implement → Test → Deploy
bd create "Design API schema" ...
bd create "Implement endpoints" ...
bd create "Add integration tests" ...
bd create "Deploy to staging" ...

# Dependencies (using requirement thinking):
bd dep add implement design --type blocks
bd dep add test implement --type blocks
bd dep add deploy test --type blocks
```

### Pattern 2: Epic Breakdown

**Scenario**: Large feature decomposed into sub-tasks.

```bash
# Epic with child tasks
bd create "User authentication system" -t epic ...
bd create "Add user model" -t task ...
bd create "Add login endpoint" -t task ...
bd create "Add registration endpoint" -t task ...

# Dependencies (parent-child relationships):
bd dep add user-model user-authentication --type parent-child
bd dep add login-endpoint user-authentication --type parent-child
bd dep add registration-endpoint user-authentication --type parent-child

# Plus technical blockers:
bd dep add login-endpoint user-model --type blocks
bd dep add registration-endpoint user-model --type blocks
```

### Pattern 3: Parallel Work with Convergence

**Scenario**: Independent tasks that converge on integration point.

```bash
# Parallel development
bd create "Backend API" -t task ...
bd create "Frontend UI" -t task ...
bd create "Integration tests" -t task ...

# Dependencies (integration needs both):
bd dep add integration backend --type blocks
bd dep add integration frontend --type blocks
```

### Pattern 4: Discovered Work Filing

**Scenario**: Filing issues found during implementation.

```bash
# While working on issue X, discover issue Y
bd update issue-x --status in_progress

# File discovered work immediately (don't fix it now!)
bd create "Fix discovered bug in auth" \
  --description="Found while implementing issue-x. Token validation fails for expired tokens." \
  --deps discovered-from:issue-x

# Continue working on issue-x only
```

## Anti-Patterns to Avoid

### Anti-Pattern 1: Temporal Naming

**Problem**: Using "Phase", "Step", "Part" in issue titles triggers temporal thinking.

**Example (bad)**:
```bash
bd create "Phase 1: Database setup" ...
bd create "Phase 2: API implementation" ...
bd dep add phase1 phase2  # ❌ Inverted!
```

**Solution**: Name by function, not sequence:
```bash
bd create "Database schema" ...
bd create "API endpoints" ...
bd dep add api-endpoints database-schema  # ✅ Requirement thinking
```

### Anti-Pattern 2: Bidirectional Dependencies

**Problem**: Creating circular dependencies.

**Example (bad)**:
```bash
bd dep add task-a task-b --type blocks
bd dep add task-b task-a --type blocks  # ❌ Circular!
```

**Detection**: `bd blocked` will show both tasks blocked indefinitely.

**Solution**: Identify the true dependency direction or split into more granular tasks.

### Anti-Pattern 3: Over-Blocking

**Problem**: Adding `blocks` dependencies when `related` would suffice.

**Example (bad)**:
```bash
# These are related but don't technically block each other
bd dep add frontend-refactor backend-refactor --type blocks  # ❌ Too strict
```

**Solution**: Use `related` for informational links:
```bash
bd dep add frontend-refactor backend-refactor --type related  # ✅ Informational
```

### Anti-Pattern 4: Implicit Parent-Child

**Problem**: Using `blocks` when you mean hierarchical decomposition.

**Example (bad)**:
```bash
# Epic and its sub-tasks
bd dep add subtask1 epic --type blocks  # ❌ Semantically wrong
bd dep add subtask2 epic --type blocks  # ❌ Semantically wrong
```

**Solution**: Use `parent-child` for hierarchies:
```bash
bd dep add subtask1 epic --type parent-child  # ✅ Semantic clarity
bd dep add subtask2 epic --type parent-child  # ✅ Semantic clarity
```

## Debugging Dependency Issues

### Issue: Task unexpectedly blocked

**Symptom**: `bd ready` doesn't show a task you expect to see.

**Diagnosis**:
```bash
# Check what's blocking it
bd show <task-id>

# Or see all blocked tasks
bd blocked
```

**Common causes**:
- Inverted dependency (temporal thinking trap)
- Transitive blocker (task A blocks B, B blocks C → C is blocked)
- Circular dependency

### Issue: Task not blocked when it should be

**Symptom**: `bd ready` shows a task that shouldn't be ready yet.

**Diagnosis**:
```bash
# Check task dependencies
bd show <task-id>

# View dependency tree
bd dep tree <task-id>
```

**Common causes**:
- Missing `blocks` dependency
- Used `related` instead of `blocks`
- Dependency direction inverted

### Issue: Circular dependency

**Symptom**: Both tasks show as blocked by each other.

**Diagnosis**:
```bash
bd blocked
# Look for tasks blocking each other mutually
```

**Solution**:
1. Identify the cycle: A blocks B blocks C blocks A
2. Determine true dependency direction
3. Remove incorrect dependency: `bd dep remove <task> <blocker>`
4. Verify: `bd blocked` should show resolution

## Quick Reference Card

**Command syntax**:
```bash
bd dep add <DEPENDENT> <REQUIRED> [--type TYPE]
```

**Mental model**:
- `<DEPENDENT>` = What gets blocked
- `<REQUIRED>` = What must be done first

**Verification**:
```bash
bd blocked           # Show all blocked tasks
bd show <id>         # Show task dependencies
bd dep tree <id>     # Show dependency tree
bd dep remove <id> <blocker>  # Remove dependency
```

**Golden rule**:
> "Y needs X" → `bd dep add Y X`

## Practice Exercises

Test your understanding:

### Exercise 1
You're building a REST API. You need to:
1. Design database schema
2. Implement ORM models
3. Create API endpoints
4. Write integration tests

**Question**: What dependencies should you add?

<details>
<summary>Answer</summary>

```bash
# Requirement thinking:
# - ORM models NEED schema
# - API endpoints NEED ORM models
# - Tests NEED API endpoints

bd dep add orm-models db-schema --type blocks
bd dep add api-endpoints orm-models --type blocks
bd dep add integration-tests api-endpoints --type blocks
```
</details>

### Exercise 2
While implementing "Add user registration", you discover the password hashing library has a security vulnerability.

**Question**: How do you file this without breaking single-issue discipline?

<details>
<summary>Answer</summary>

```bash
# File discovered work, don't fix it now
bd create "Upgrade password hashing library" \
  --description="Security vulnerability in bcrypt v1.2. Found while implementing registration. Need to upgrade to v2.0." \
  --deps discovered-from:user-registration \
  -t task -p 1

# Continue working on user-registration only
```
</details>

### Exercise 3
You have an epic "User authentication" with sub-tasks:
- Add user model
- Add login endpoint
- Add logout endpoint

Login and logout both need the user model.

**Question**: What dependencies do you add?

<details>
<summary>Answer</summary>

```bash
# Parent-child relationships (hierarchical):
bd dep add user-model user-authentication --type parent-child
bd dep add login-endpoint user-authentication --type parent-child
bd dep add logout-endpoint user-authentication --type parent-child

# Technical blockers (login/logout NEED user model):
bd dep add login-endpoint user-model --type blocks
bd dep add logout-endpoint user-model --type blocks
```
</details>

## Summary

**Core principle**: Think in requirements, not sequences.

**Key insights**:
1. Temporal language ("before", "first", "Phase 1") inverts dependencies
2. Ask "What does Y NEED?" → `bd dep add Y X`
3. Always verify with `bd blocked`
4. Use appropriate dependency types (`blocks`, `parent-child`, `discovered-from`, `related`)

**Remember**:
- The cognitive trap is real - even experienced developers fall for it
- Verification is mandatory - check your work with `bd blocked`
- When in doubt, draw a diagram showing "needs" arrows
- If dependencies look wrong, you probably used temporal thinking

Master this skill and you'll avoid the most common beads mistake!

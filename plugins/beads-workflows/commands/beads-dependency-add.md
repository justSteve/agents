# Beads Dependency Add

Add dependencies between issues with guided dependency-thinking to avoid the temporal thinking trap.

## Overview

This command orchestrates proper dependency addition:
1. Parse dependency request
2. Apply dependency-thinking guidance
3. Validate issues exist and no cycles
4. Add dependency with `bd dep add`
5. Verify with `bd blocked` or `bd show`

## Arguments

**Required**:
- `<dependent>` - Issue that NEEDS the other (gets blocked)
- `<required>` - Issue that must be done FIRST (does the blocking)

**Optional**:
- `--type <type>` - Dependency type (default: blocks)
  - `blocks` - Hard blocker, dependent cannot start until required is complete
  - `parent-child` - Hierarchical, dependent is sub-task of required (epic)
  - `discovered-from` - Discovery, dependent was found while working on required
  - `related` - Soft link, informational only, no blocking

## Phase 1: Parse Dependency Request

Extract and validate the dependency parameters.

**Parse arguments**:
```
<dependent> = First positional argument (issue that needs the other)
<required> = Second positional argument (issue that must be done first)
<type> = --type value or "blocks" if not specified
```

**Validation**:
- Both issue IDs provided
- Type is valid (blocks, parent-child, discovered-from, related)

**Output**:
```
Dependency Request:
  Dependent: <dependent> (will be blocked)
  Required: <required> (must complete first)
  Type: <type>

Applying dependency-thinking check...
```

**If arguments missing**:
```
❌ Missing required arguments

Usage: beads-dependency-add <dependent> <required> [--type TYPE]

Arguments:
  <dependent>  Issue that NEEDS the other (gets blocked)
  <required>   Issue that must be done FIRST (does the blocking)

Example:
  beads-dependency-add api-endpoints db-schema --type blocks

This means: api-endpoints NEEDS db-schema (api-endpoints is blocked until db-schema is done)
```

## Phase 2: Dependency-Thinking Guidance

Apply the dependency-thinking skill to prevent inverted dependencies.

### The Temporal Thinking Check

**Present the key question**:
```
================================================================================
DEPENDENCY-THINKING CHECK
================================================================================

You are adding: <dependent> depends on <required>

This means:
• <dependent> will be BLOCKED until <required> is complete
• <required> must be done FIRST before <dependent> can start

Verification question:
"Does <dependent> NEED <required> to be finished before it can start?"

Common trap: If you thought "I want to do <required> before <dependent>" - that's
temporal thinking. The command syntax is about REQUIREMENTS, not sequence.

================================================================================

Is this correct? [Y/n/flip]
  Y    - Yes, <dependent> needs <required>
  n    - No, abort
  flip - Swap them: <required> depends on <dependent>
```

**If user selects "flip"**:
```
Flipping dependency direction:
  Old: <dependent> -> <required>
  New: <required> -> <dependent>

This means: <required> NEEDS <dependent> to be finished first.

Proceeding with flipped dependency...
```

### Type-Specific Guidance

**For `blocks` type**:
```
Type: blocks (hard blocker)

This creates a HARD dependency:
• <dependent> will NOT appear in 'bd ready' until <required> is closed
• Attempting to start <dependent> while <required> is open is a violation

Use 'blocks' when:
• Technical dependency (can't write tests without API)
• Sequential phases (design must precede implementation)
• External dependency (waiting on another team/system)

Continue with 'blocks' type? [Y/n]
```

**For `parent-child` type**:
```
Type: parent-child (hierarchical)

This creates a HIERARCHY:
• <dependent> is a sub-task of <required> (the parent/epic)
• Does NOT block <dependent> from being worked on
• Links work for organization and rollup

Use 'parent-child' when:
• Breaking epic into tasks
• Creating sub-issues
• Organizing related work under a theme

Continue with 'parent-child' type? [Y/n]
```

**For `discovered-from` type**:
```
Type: discovered-from (work discovery)

This creates a DISCOVERY link:
• <dependent> was found while working on <required>
• Does NOT block either issue
• Tracks work lineage and context

Use 'discovered-from' when:
• Filing bugs found during feature work
• Creating follow-up tasks discovered during implementation
• Documenting scope expansions that were set aside

Continue with 'discovered-from' type? [Y/n]
```

**For `related` type**:
```
Type: related (informational)

This creates an INFORMATIONAL link:
• <dependent> and <required> are related
• Does NOT block either issue
• Useful for cross-referencing

Use 'related' when:
• Issues touch same subsystem
• Informational cross-reference
• Grouping without hierarchy

Continue with 'related' type? [Y/n]
```

### Example-Based Verification

**Show concrete example for clarity**:
```
Concrete example of your dependency:

Scenario: "<required-title>" must be done before "<dependent-title>"

Think of it as:
"I cannot start working on '<dependent-title>' until '<required-title>' is complete."

If this sounds RIGHT -> proceed
If this sounds WRONG -> you may have the direction inverted

[Press Enter to continue or 'flip' to swap]
```

## Phase 3: Validate Issues and Check for Cycles

Verify both issues exist and dependency won't create a cycle.

### Issue Existence Check

**Commands to run**:
```bash
bd show <dependent>
bd show <required>
```

**Validation**:
- Both issues exist in beads database
- Neither issue is closed (warn if adding dependency to closed issue)

**Output** (success):
```
✓ Issue validation passed

Dependent: <dependent>
  Title: <dependent-title>
  Status: <status>

Required: <required>
  Title: <required-title>
  Status: <status>
```

**Output** (issue not found):
```
❌ Issue not found: <id>

The issue "<id>" does not exist in the beads database.

Available similar issues:
  - <similar1>: <title1>
  - <similar2>: <title2>

Options:
[1] Enter correct issue ID
[2] Create the missing issue
[3] Cancel dependency addition

Your choice: _
```

**Output** (closed issue warning):
```
⚠️  Adding dependency to closed issue

Issue "<id>" is already closed.
Status: closed
Closed: <timestamp>

Adding a dependency to a closed issue is unusual.

Options:
[1] Continue anyway (rare case)
[2] Cancel (recommended)

Your choice: _
```

### Circular Dependency Check

**Logic**:
1. Get current dependencies of `<dependent>`
2. Check if adding this dependency creates a cycle
3. Trace through transitive dependencies

**Commands to run**:
```bash
bd dep tree <required>
```

**Check if `<dependent>` appears anywhere in the tree**:

**Output** (no cycle):
```
✓ No circular dependency detected

Dependency chain after addition:
<required> -> <dependent>
```

**Output** (cycle detected):
```
❌ Circular dependency detected!

Adding this dependency would create a cycle:

Current chain:
<required> is blocked by:
  <issue-a> is blocked by:
    <dependent>  ← This issue!

Adding <dependent> -> <required> would create:
<dependent> -> <required> -> ... -> <dependent> (CYCLE!)

This is not allowed. Dependencies must form a DAG (directed acyclic graph).

Solutions:
1. Re-evaluate the dependency direction
2. Break the cycle by restructuring work
3. Use 'related' type instead of 'blocks' if no true blocking exists

Dependency addition aborted.
```

## Phase 4: Add Dependency

Execute the dependency addition.

**Command**:
```bash
bd dep add <dependent> <required> --type <type>
```

**Validation**:
- Command succeeds
- No error messages

**Output** (success):
```
✓ Dependency added

<dependent> -> <required> (type: <type>)

This means:
• <dependent> now depends on <required>
• Effect: <effect based on type>
```

**Effect descriptions by type**:
- `blocks`: "<dependent> is blocked until <required> is closed"
- `parent-child`: "<dependent> is tracked as sub-task of <required>"
- `discovered-from`: "<dependent> is linked as discovered from <required>"
- `related`: "<dependent> and <required> are cross-referenced"

**Output** (failure):
```
❌ Failed to add dependency

Error: <bd error message>

This may indicate:
• Database issue
• Invalid issue ID
• Permission problem

Please check the error and retry.
```

## Phase 5: Verify Dependency

Confirm the dependency was added correctly.

**Commands to run**:
```bash
bd show <dependent>
bd blocked
```

**Verification checks**:
- Dependency appears in `bd show` output
- If type is `blocks`, issue appears in `bd blocked`
- No unexpected side effects

**Output** (success for `blocks` type):
```
✅ Dependency verified

<dependent>:
  Now blocked by: <required> (<type>)
  Status: <ready if required is done, blocked otherwise>

Blocked issues (updated):
  <dependent>: blocked by <required>

Dependency tree:
  <required> (must complete first)
    └── <dependent> (depends on above)
```

**Output** (success for non-blocking types):
```
✅ Dependency verified

<dependent>:
  <type> relationship with: <required>
  Status: <unchanged - no blocking for this type>

The dependency is recorded but does not affect work availability.
```

**Output** (verification mismatch):
```
⚠️  Verification issue

Dependency was added but verification shows unexpected state:

Expected: <dependent> blocked by <required>
Actual: <actual state>

This may indicate:
• Sync needed: run 'bd sync'
• Database delay
• Edge case in dependency tracking

Recommended: verify manually with 'bd show <dependent>'
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

Dependency addition aborted.
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

Dependency addition aborted.
```

### Error: Same issue specified twice

**Detection**: `<dependent>` == `<required>`

**Output**:
```
❌ Cannot add self-dependency

You specified the same issue for both dependent and required: <id>

An issue cannot depend on itself.

Check your command:
beads-dependency-add <dependent> <required>

Dependency addition aborted.
```

### Error: Dependency already exists

**Detection**: bd dep add returns "already exists" error

**Output**:
```
⚠️  Dependency already exists

<dependent> -> <required> (<type>) is already recorded.

Current state:
<bd show output for dependent>

Options:
[1] Change dependency type
[2] Remove and re-add
[3] Leave as-is (exit)

Your choice: _
```

## Success Criteria

Dependency addition is successful when:
- ✅ Dependency-thinking check passed (user confirmed direction)
- ✅ Both issues exist and are valid
- ✅ No circular dependency created
- ✅ Dependency added to beads database
- ✅ Verification shows expected blocking behavior
- ✅ User understands the effect of the dependency

## Notes

**The Golden Rule**:
> "Y needs X" = `bd dep add Y X`

This is the most important thing to remember. If you think "X comes before Y" or "do X first, then Y", you need to translate that to "Y needs X".

**Mental Translation Table**:

| Temporal Thinking (Trap) | Requirement Thinking (Correct) |
|--------------------------|-------------------------------|
| "X before Y" | "Y needs X" → `bd dep add Y X` |
| "X then Y" | "Y requires X" → `bd dep add Y X` |
| "Phase 1, Phase 2" | "Phase 2 needs Phase 1" → `bd dep add phase2 phase1` |
| "First X, later Y" | "Y depends on X" → `bd dep add Y X` |

**Verification Strategy**: Always run `bd blocked` after adding dependencies to confirm the blocking is as expected. If something looks wrong, you probably used temporal thinking.

**Choosing the Right Type**:
- **blocks**: Technical dependency, must be sequential
- **parent-child**: Epic/task hierarchy, organizational
- **discovered-from**: Work found during other work, tracking
- **related**: Cross-reference, informational only

**Integration**: This command works with:
- `dependency-thinking` skill for detailed guidance
- `beads-session-start` which checks ready work (affected by blocks dependencies)
- `bd blocked` for verification and debugging

**Common Patterns**:

Sequential pipeline:
```bash
beads-dependency-add step2 step1 --type blocks
beads-dependency-add step3 step2 --type blocks
```

Epic breakdown:
```bash
beads-dependency-add task1 epic --type parent-child
beads-dependency-add task2 epic --type parent-child
beads-dependency-add task2 task1 --type blocks  # if sequential
```

Discovered work:
```bash
beads-dependency-add new-bug current-work --type discovered-from
```

**Future enhancements**:
- Batch dependency addition
- Visual dependency graph output
- Automatic dependency suggestion based on issue content
- Undo/history for dependency changes

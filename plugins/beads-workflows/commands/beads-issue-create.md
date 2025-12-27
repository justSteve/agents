# Beads Issue Create

Create issues with structured quality checks, ensuring descriptions are actionable and dependencies are properly tracked.

## Overview

This command orchestrates high-quality issue creation:
1. Gather required information
2. Apply quality checks (description-quality skill principles)
3. Create issue with `bd create`
4. Add dependencies if needed
5. Display created issue

## Arguments

**Required** (one of):
- `--title "<title>"` - Issue title (can also be positional argument)
- Interactive mode (if no arguments provided)

**Optional**:
- `--description "<text>"` or `-d` - Full issue description
- `--type <type>` or `-t` - Issue type: task, bug, epic, feature (default: task)
- `--priority <0-4>` or `-p` - Priority: 0=critical, 1=high, 2=medium, 3=low, 4=lowest (default: 2)
- `--deps "<deps>"` - Dependencies in format `type:issue-id` (comma-separated)
- `--quick` - Skip quality review (use with caution)
- `--review` - Request quality review from beads-issue-reviewer agent

**Dependency shorthand for discovered work**:
- `--deps discovered-from:<id>` - This issue was discovered while working on `<id>`

## Phase 1: Gather Required Information

Collect issue details either from arguments or interactively.

### From Arguments

**If `--title` and `--description` provided**:
- Parse all provided arguments
- Proceed to Phase 2 (quality check)

**Example**:
```bash
beads-issue-create --title "Fix null pointer in user lookup" \
  --description="get_user() can return null but callers don't check" \
  --type bug --priority 1 \
  --deps discovered-from:user-preferences
```

### Interactive Mode

**If insufficient arguments provided**:

**Step 1a: Title**
```
Create New Issue
================

Enter issue title (brief, descriptive):
> _
```

**Validation**:
- Title must be non-empty
- Title should be < 100 characters
- Warn if title is vague (e.g., "fix bug", "update code")

**Step 1b: Type**
```
Issue type:
[1] task    - General work item
[2] bug     - Something is broken
[3] feature - New functionality
[4] epic    - Large multi-part work

Your choice (default: task): _
```

**Step 1c: Priority**
```
Priority:
[0] Critical - System down, data loss risk
[1] High     - Major impact, needs quick attention
[2] Medium   - Normal priority (default)
[3] Low      - Nice to have, when time permits
[4] Lowest   - Backlog, someday/maybe

Your choice (default: 2): _
```

**Step 1d: Description**
```
Issue description (end with empty line):

## Context
[Why this matters]

## Problem
[What needs to change]

## Acceptance Criteria
- [ ] [Verifiable criterion]

> _
```

**Provide template if user enters empty line immediately**:
```
Would you like to use the description template? [Y/n]

Template:
## Context
[Why this matters - 1-2 sentences]

## Problem
Current: [What's wrong or missing]
Desired: [What should happen]

## Acceptance Criteria
- [ ] [Verifiable criterion 1]
- [ ] [Verifiable criterion 2]

## Technical Notes
- Location: [file paths]
- [Any constraints]
```

**Step 1e: Dependencies (optional)**
```
Does this issue have dependencies? [y/N]

If yes:
Dependency type:
[1] blocks        - This issue needs another to complete first
[2] parent-child  - This is a sub-task of an epic/feature
[3] discovered-from - Found while working on another issue
[4] related       - Informational link only

Enter the required issue ID: _
```

**Output** (after gathering):
```
Issue Draft:
  Title: <title>
  Type: <type>
  Priority: P<priority>
  Description: <first 100 chars>...
  Dependencies: <list or "none">

Proceeding to quality check...
```

## Phase 2: Quality Check

Apply description-quality skill principles to ensure issue is actionable.

**Quality criteria** (from description-quality skill):

| Criterion | Check | Weight |
|-----------|-------|--------|
| Context | Does description explain WHY this matters? | Required |
| Problem Statement | Is current vs desired state clear? | Required |
| Acceptance Criteria | Are there verifiable completion conditions? | Required |
| Location | Are relevant files/components mentioned? | Recommended |
| Scope | Are boundaries clear (what's NOT included)? | Recommended |
| Standalone | Can someone work on this without asking questions? | Critical |

### Quality Score Calculation

**Scoring**:
- Required criteria missing: FAIL
- All required criteria present, some recommended missing: PASS with warnings
- All criteria present: PASS

**If FAIL** (missing required criteria):
```
⚠️  Quality check: NEEDS IMPROVEMENT

Missing required elements:
• Context: No explanation of why this matters
• Acceptance Criteria: No verifiable completion conditions

The description needs these elements for any agent to work on it effectively.

Options:
[1] Add missing elements now (recommended)
[2] Create anyway (issue may need clarification later)
[3] Cancel creation

Your choice: _
```

**If PASS with warnings**:
```
✓ Quality check: ACCEPTABLE

Optional improvements:
• Consider adding file locations
• Consider clarifying scope boundaries

Options:
[1] Proceed with creation
[2] Improve description first
[3] Cancel

Your choice: _
```

**If PASS** (all criteria):
```
✓ Quality check: GOOD

All quality criteria met:
✓ Context explained
✓ Problem statement clear
✓ Acceptance criteria defined
✓ Technical notes included

Proceeding to create issue...
```

### --review Flag (Optional Quality Review)

**If `--review` flag provided**:

Invoke beads-issue-reviewer agent (if available):

**Agent prompt**:
```
Review this issue for quality and completeness:

Title: <title>
Type: <type>
Priority: P<priority>

Description:
<description>

Dependencies: <deps>

Provide:
1. Quality score (1-10)
2. Strengths
3. Areas for improvement
4. Suggested edits (if any)
```

**Display review results**:
```
Issue Review (from beads-issue-reviewer):

Quality Score: <score>/10

Strengths:
• <strength 1>
• <strength 2>

Improvements suggested:
• <improvement 1>
• <improvement 2>

Apply suggested improvements? [Y/n/edit]
```

## Phase 3: Create Issue

Execute the bd create command.

**Build command**:
```bash
bd create "<title>" \
  --description="<description>" \
  --type <type> \
  --priority <priority>
```

**Execute and capture output**:

**Validation**:
- Command exits successfully
- Issue ID returned
- No duplicate title warning (if applicable)

**Output** (success):
```
Creating issue...

✓ Issue created: <issue-id>
  Title: <title>
  Type: <type>
  Priority: P<priority>
  Status: open
```

**Output** (failure):
```
❌ Failed to create issue

Error: <bd create error message>

Possible causes:
• Duplicate issue title
• Invalid type or priority value
• Database write error

Please resolve and retry.
```

## Phase 4: Add Dependencies

Add any specified dependencies to the created issue.

**If no dependencies**: Skip to Phase 5

**If dependencies specified**:

**For each dependency in `--deps`**:

Parse format: `<type>:<issue-id>`

**Command**:
```bash
bd dep add <new-issue-id> <required-issue-id> --type <dep-type>
```

### Special Case: discovered-from

**If `--deps discovered-from:<id>`**:

This is the "discovered work" pattern. Add context:

```bash
bd dep add <new-issue-id> <parent-id> --type discovered-from
```

**Output**:
```
Adding dependency: discovered-from <parent-id>

This issue will be linked as discovered work from <parent-id>.
```

### Dependency Validation

**Before adding, verify**:
1. Required issue exists: `bd show <required-id>`
2. No circular dependency would result
3. Dependency type is appropriate

**If issue doesn't exist**:
```
⚠️  Dependency target not found: <id>

Issue "<id>" does not exist.

Options:
[1] Skip this dependency
[2] Create the missing issue first
[3] Enter correct issue ID

Your choice: _
```

**If circular dependency detected**:
```
❌ Circular dependency detected

Adding this dependency would create a cycle:
<new-id> -> <req-id> -> ... -> <new-id>

Dependency NOT added.

Tip: Review your dependency structure. Consider using 'related' instead of 'blocks'.
```

**Output** (success):
```
✓ Dependencies added:
  - <new-id> -> <dep1-id> (blocks)
  - <new-id> -> <dep2-id> (discovered-from)
```

## Phase 5: Show Created Issue

Display the complete created issue for verification.

**Command**:
```bash
bd show <new-issue-id>
```

**Output**:
```
✅ Issue created successfully

================================================================================
Issue: <issue-id>
Type: <type>
Priority: P<priority>
Status: open
Created: <timestamp>

Description:
<full description>

Dependencies:
  Blocked by: <list or "none">
  Discovered from: <parent or "none">
  Related: <list or "none">
================================================================================

Next steps:
• Work on this issue: bd update <id> --status in_progress
• Add more dependencies: bd dep add <id> <required>
• View ready work: bd ready
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

Issue creation aborted.
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

Issue creation aborted.
```

### Error: Title too vague

**Detection**: Title matches known vague patterns

**Vague patterns**:
- "fix bug"
- "update code"
- "refactor"
- "changes"
- "misc"
- Single-word titles (without context)

**Output**:
```
⚠️  Title may be too vague: "<title>"

Vague titles make issues hard to find and understand.

Examples of better titles:
• "Fix null pointer in user lookup when email is missing"
• "Update authentication to support OAuth2"
• "Refactor UserService into separate auth and profile services"

Options:
[1] Enter new title
[2] Keep current title (not recommended)

Your choice: _
```

### Error: Description too short

**Detection**: Description < 50 characters and not using --quick flag

**Output**:
```
⚠️  Description is very brief: <length> characters

Short descriptions often lack essential context.

Minimum recommended sections:
- Context (why this matters)
- Problem (what needs to change)
- Acceptance criteria (how to verify completion)

Options:
[1] Expand description now
[2] Use description template
[3] Keep as-is (not recommended)

Your choice: _
```

## Success Criteria

Issue creation is successful when:
- ✅ Title is clear and specific
- ✅ Description passes quality check
- ✅ Issue created in beads database
- ✅ All dependencies added and validated
- ✅ No circular dependencies created
- ✅ Issue visible in `bd list` and `bd ready` (if unblocked)

## Notes

**The "Future Self" Rule**: Write descriptions for your future self who has forgotten all context. If you can't understand an issue 6 months later, it's not good enough.

**Discovered Work Pattern**: When you find issues during implementation:
1. Stop - don't fix the discovered issue now
2. File it with `--deps discovered-from:<current-issue>`
3. Continue working on your original issue

This maintains single-issue discipline while ensuring nothing is forgotten.

**Priority Guidelines**:
- P0 (Critical): Production down, data loss, security breach
- P1 (High): Major feature broken, significant user impact
- P2 (Medium): Normal bugs and features (default)
- P3 (Low): Minor issues, polish, nice-to-haves
- P4 (Lowest): Someday/maybe, backlog items

**Description Quality Checklist**:
- [ ] Context: Why does this matter?
- [ ] Problem: What's the current vs desired state?
- [ ] Scope: What's NOT included?
- [ ] Criteria: How do we verify completion?
- [ ] Location: What files/components are involved?

**Integration**: This command works with:
- `description-quality` skill for quality standards
- `beads-session-start` to pick up created issues
- `dependency-thinking` skill for proper dependency direction

**Future enhancements**:
- Template library for common issue types
- Auto-suggested dependencies from description text
- Integration with external issue trackers

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

### Quality Validation Logic

This section provides the implementation-level validation rules for automated quality checking.

#### Title Validation

**Length Check**:
```
Minimum: 10 characters
Maximum: 100 characters
Optimal: 30-70 characters

if length < 10:
    FAIL: "Title too short - provide more context"
if length > 100:
    WARNING: "Title may be too long - consider shortening"
```

**Vague Pattern Detection**:
```python
VAGUE_PATTERNS = [
    r'^fix\s*(bug)?$',           # "fix" or "fix bug"
    r'^update\s*(code)?$',       # "update" or "update code"
    r'^refactor$',               # just "refactor"
    r'^changes?$',               # "change" or "changes"
    r'^misc$',                   # "misc"
    r'^todo$',                   # "todo"
    r'^wip$',                    # "wip"
    r'^\w+$',                    # single word without context
]

for pattern in VAGUE_PATTERNS:
    if regex.match(pattern, title, IGNORECASE):
        WARNING: "Title may be too vague"
        break
```

**Good Title Indicators**:
```python
GOOD_INDICATORS = [
    r'\b(in|when|for|with|to|from)\b',  # prepositions add context
    r'\b\w+\.\w+',                       # file references (foo.ts)
    r'\b(add|remove|fix|update|implement|refactor)\s+\w+\s+\w+',  # verb + object + context
]

score = sum(1 for p in GOOD_INDICATORS if regex.search(p, title))
if score >= 2:
    PASS: "Title is descriptive"
```

**Corrective Prompt (Title)**:
```
⚠️  Title needs improvement

Current: "<title>"
Issues: <list of detected issues>

A good title should:
✓ Be specific about WHAT is changing
✓ Include context (where, when, why)
✓ Be findable via search

Examples:
• "Fix null pointer in UserService.getUser() when email is null"
• "Add OAuth2 support to authentication flow"
• "Refactor PaymentProcessor into separate validation and execution modules"

Enter improved title (or press Enter to keep current):
> _
```

#### Description Validation

**Length Check by Type**:
```
Type        Minimum   Warning   Good
----        -------   -------   ----
chore       30        50        100+
task        50        100       150+
bug         100       150       200+
feature     100       200       300+
epic        200       300       500+

if length < minimum:
    FAIL: "Description too short for <type> (need {minimum}+ chars)"
if length < warning:
    WARNING: "Description could use more detail"
```

**Why/What/How Structure Detection**:
```python
def detect_structure(description):
    indicators = {
        'why': {
            'headers': [r'##?\s*context', r'##?\s*background', r'##?\s*motivation'],
            'keywords': [r'\bbecause\b', r'\bsince\b', r'\bdue to\b', r'\bneeded for\b',
                        r'\brequired by\b', r'\bto support\b', r'\benables?\b'],
            'weight': 'required'
        },
        'what': {
            'headers': [r'##?\s*problem', r'##?\s*issue', r'##?\s*current', r'##?\s*desired'],
            'keywords': [r'\bcurrently?\b', r'\bshould\b', r'\bexpected\b', r'\bactual\b',
                        r'\bwant\b', r'\bneed\b', r'\bmust\b'],
            'weight': 'required'
        },
        'how': {
            'headers': [r'##?\s*acceptance', r'##?\s*criteria', r'##?\s*verification',
                       r'##?\s*done when', r'##?\s*success'],
            'keywords': [r'\[\s*\]', r'\bverify\b', r'\btest\b', r'\bconfirm\b',
                        r'\bcomplete when\b', r'\bdone when\b'],
            'weight': 'required'
        },
        'where': {
            'headers': [r'##?\s*location', r'##?\s*files?', r'##?\s*technical'],
            'keywords': [r'\b\w+\.(ts|js|py|go|rs|md)\b', r'\bsrc/', r'\blib/',
                        r'\bcomponents?/', r'\bservices?/'],
            'weight': 'recommended'
        }
    }

    results = {}
    for element, checks in indicators.items():
        found_header = any(regex.search(h, description, IGNORECASE) for h in checks['headers'])
        found_keyword = any(regex.search(k, description, IGNORECASE) for k in checks['keywords'])
        results[element] = {
            'present': found_header or found_keyword,
            'weight': checks['weight'],
            'method': 'header' if found_header else 'keyword' if found_keyword else None
        }

    return results
```

**Structure Validation**:
```
structure = detect_structure(description)

missing_required = [e for e, r in structure.items() if r['weight'] == 'required' and not r['present']]
missing_recommended = [e for e, r in structure.items() if r['weight'] == 'recommended' and not r['present']]

if missing_required:
    FAIL: "Missing required elements: " + ', '.join(missing_required)
elif missing_recommended:
    WARNING: "Consider adding: " + ', '.join(missing_recommended)
else:
    PASS: "Description has good structure"
```

**Corrective Prompt (Description - Missing Required)**:
```
⚠️  Description missing required elements

Missing:
• WHY: No context explaining why this matters
• WHAT: No clear problem statement
• HOW: No acceptance criteria for verification

Your description:
"<first 100 chars of description>..."

Add the missing elements using this template:

## Context (WHY)
[Explain why this work is needed - who is affected, what triggered it]

## Problem (WHAT)
Current: [What's happening now that shouldn't be]
Desired: [What should happen instead]

## Acceptance Criteria (HOW)
- [ ] [Specific, testable condition 1]
- [ ] [Specific, testable condition 2]

Options:
[1] Add missing elements now (recommended)
[2] Create anyway - I'll add details later
[3] Cancel

Your choice: _
```

**Corrective Prompt (Description - Short)**:
```
⚠️  Description is brief: <N> characters

For <type> issues, we recommend at least <minimum> characters.

Your description:
"<full description>"

Tips to expand:
• Add context: Why is this needed? Who reported it?
• Be specific: What exactly should change?
• Add criteria: How will you know it's done?
• Add location: What files/components are involved?

Options:
[1] Expand description now
[2] Use description template
[3] Keep as-is (not recommended for <type>)

Your choice: _
```

#### Combined Quality Score

```python
def calculate_quality_score(title, description, issue_type):
    score = 0
    issues = []

    # Title checks (max 30 points)
    title_len = len(title)
    if title_len >= 30:
        score += 15
    elif title_len >= 10:
        score += 10
    else:
        issues.append("Title too short")

    if not is_vague(title):
        score += 15
    else:
        issues.append("Title may be vague")

    # Description checks (max 70 points)
    desc_len = len(description)
    min_len = MIN_LENGTHS[issue_type]

    if desc_len >= min_len * 2:
        score += 20
    elif desc_len >= min_len:
        score += 15
    elif desc_len >= min_len / 2:
        score += 5
        issues.append("Description shorter than recommended")
    else:
        issues.append("Description too short")

    structure = detect_structure(description)
    for element in ['why', 'what', 'how']:
        if structure[element]['present']:
            score += 15
        else:
            issues.append(f"Missing {element.upper()}")

    if structure['where']['present']:
        score += 5

    # Determine result
    if score >= 80:
        return ('PASS', score, issues)
    elif score >= 50:
        return ('WARNING', score, issues)
    else:
        return ('FAIL', score, issues)
```

**Quality Score Display**:
```
Quality Analysis
================

Title:       "<title>"
Description: <N> characters
Type:        <type>

Score: <score>/100

Breakdown:
  Title length:     <points>/15  <status>
  Title clarity:    <points>/15  <status>
  Desc length:      <points>/20  <status>
  WHY (context):    <points>/15  <status>
  WHAT (problem):   <points>/15  <status>
  HOW (criteria):   <points>/15  <status>
  WHERE (location): <points>/5   <status>

Result: <PASS|WARNING|FAIL>

<if issues>
Issues to address:
• <issue 1>
• <issue 2>
</if>
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

## Validation Checkpoints

This command enforces beads discipline through explicit validation checkpoints. Each checkpoint invokes the `beads-issue-reviewer` agent for quality validation.

### Checkpoint 1: Title Quality (Phase 1)

**Trigger**: After title is provided

**Validation**:
```
Invoke beads-issue-reviewer with context:
- Proposed title
- Issue type

Title validation rules:
- Length: 10-100 characters (optimal: 30-70)
- Content: Describes WHAT, not HOW
- Specificity: Concrete, not vague
- No red flags: "Fix bug", "Update", "Refactor" alone

Expected response:
- PASS: Title is clear and specific
- WARNING: Title could be improved (suggest better)
- FAIL: Title is vague or inappropriate (block until fixed)
```

**On FAIL**: Block creation until title is improved. Show examples.

### Checkpoint 2: Description Quality (Phase 2)

**Trigger**: After description is provided

**Validation**:
```
Invoke beads-issue-reviewer with context:
- Full description
- Issue type
- Priority

Description validation rules:
- Why: Problem statement present (required)
- What: Scope/approach defined (required)
- How discovered: Context if applicable (recommended)
- Minimum length by type (task: 50, bug: 100, feature: 100, epic: 200)

Expected response:
- PASS: Description meets quality standards
- WARNING: Missing recommended elements (list them)
- FAIL: Missing required elements (block until added)
```

**On FAIL**: Block creation. Show template and missing elements.

### Checkpoint 3: Pre-Creation Validation (Phase 3)

**Trigger**: Before executing bd create

**Validation**:
```
Invoke beads-disciplinarian with context:
- Complete issue draft (title, description, type, priority)
- Dependencies to add

Full compliance check:
- [ ] Title passes quality check
- [ ] Description has Why/What/How
- [ ] Type is appropriate for content
- [ ] Priority is justified
- [ ] Dependencies use correct direction

Expected response:
- PASS: Ready to create
- WARNING: Create with noted concerns
- FAIL: Critical issues must be resolved
```

**On FAIL**: Block creation until all issues resolved.

### Checkpoint 4: Dependency Direction (Phase 4)

**Trigger**: When adding dependencies

**Validation**:
```
Invoke beads-disciplinarian with context:
- Proposed dependency (dependent, required, type)
- Request: Validate causal reasoning

Dependency validation:
- Question: "Does <dependent> NEED <required>?"
- Check: No temporal thinking ("first", "then", "before")
- Verify: Correct bd dep add syntax

Expected response:
- PASS: Dependency direction correct
- WARNING: Direction may be inverted (confirm with user)
- FAIL: Circular dependency or invalid issue
```

**On WARNING**: Force user confirmation of direction before proceeding.

### Agent Integration

When invoking beads-issue-reviewer for validation:

```markdown
Review issue quality before creation:

Issue draft:
- Title: "<title>"
- Type: <type>
- Priority: P<priority>
- Description: "<description>"
- Dependencies: <list>

Check:
1. Title quality (specific, actionable)
2. Description quality (Why/What/How structure)
3. Acceptance criteria (if present)
4. Metadata appropriateness

Return: Score (1-5), PASS/WARNING/FAIL, specific improvements
```

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

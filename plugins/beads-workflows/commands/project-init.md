# Project Init

Initialize a beads-native project with proper structure, templates, and initial issues.

Supports both **greenfield** (new projects) and **brownfield** (existing projects) scenarios.

## Overview

This command scaffolds a complete project setup:
1. Mode detection (greenfield vs brownfield)
2. Directory setup (create/verify, git init)
3. Beads initialization (bd init, hooks, config)
4. Template application (files, CLAUDE.md)
5. Initial issues (setup epic with tasks)
6. First session (optional beads-session-start)

## Arguments

**Required** (greenfield only):
- `<project-name>` - Name for the project (used for directory and module)

**Brownfield Flags**:
- `--existing` or `-e` - Initialize beads in current directory (brownfield mode)
- `--beads-only` - Only add beads, skip all template files
- `--no-overwrite` - Never overwrite existing files (safest for brownfield)

**Optional**:
- `--prefix <prefix>` - Issue ID prefix (default: derived from project/directory name)
- `--template <type>` - Project template: standard, python, typescript, go (default: auto-detect or standard)
- `--agents <list>` - Comma-separated agent plugins to enable (default: beads-workflows)
- `--path <dir>` - Parent directory for project (default: current directory)
- `--description <text>` - Project description for README and CLAUDE.md
- `--no-session` - Skip running beads-session-start after setup
- `--no-issues` - Skip creating initial issues (useful for brownfield with existing backlog)

## Mode Detection

The command automatically detects greenfield vs brownfield based on arguments and context.

### Greenfield Mode (New Project)
```bash
/project-init my-new-app --template python
```
- Creates new directory
- Full template application
- Creates all files from scratch

### Brownfield Mode (Existing Project)
```bash
/project-init --existing
/project-init -e
/project-init . --existing
```
- Works in current directory
- Auto-detects project type
- Preserves existing files
- Appends to .gitignore and CLAUDE.md

### Mode Selection Logic
```
if --existing flag OR project-name is "."
  → Brownfield mode
else
  → Greenfield mode
```

## Project Type Detection (Brownfield)

When running in brownfield mode without explicit `--template`, detect project type from existing files:

| File Found | Detected Type |
|------------|---------------|
| `package.json` | typescript |
| `pyproject.toml` | python |
| `setup.py` | python |
| `go.mod` | go |
| `Cargo.toml` | rust (future) |
| None of above | standard |

**Detection Command**:
```bash
# Check for project markers
if [ -f "package.json" ]; then
  PROJECT_TYPE="typescript"
elif [ -f "pyproject.toml" ] || [ -f "setup.py" ]; then
  PROJECT_TYPE="python"
elif [ -f "go.mod" ]; then
  PROJECT_TYPE="go"
else
  PROJECT_TYPE="standard"
fi
```

**Output**:
```
✓ Detected project type: python (from pyproject.toml)
```

## Phase 1: Directory Setup

Create or verify project directory with git initialization.

**Steps**:
```bash
# Create project directory
mkdir -p <path>/<project-name>
cd <path>/<project-name>

# Initialize git if not already a repo
git init
```

**Validation**:
- Directory created successfully
- Git initialized (or already a git repo)
- No conflicting .beads directory

**Output**:
```
✓ Directory created: <path>/<project-name>
✓ Git initialized
```

## Phase 2: Beads Initialization

Initialize beads issue tracking with hooks and agent configuration.

**Commands**:
```bash
# Initialize beads with prefix
bd init --prefix <prefix>

# Install git hooks
bd hooks install

# Configure agent marketplace (if available)
bd config set agents.marketplace.repo "<agents-repo-path>"
bd config set agents.marketplace.enabled "<agents-list>"
bd config set agents.default_model "sonnet"
```

**Validation**:
- .beads directory created
- beads.db exists
- Hooks installed
- Configuration set

**Output**:
```
✓ Beads initialized with prefix: <prefix>
✓ Git hooks installed
✓ Agent plugins enabled: <list>
```

## Phase 3: Template Application

Apply project template and generate CLAUDE.md.

**Template Selection**:
| Template | Files Created |
|----------|---------------|
| standard | .gitignore, README.md |
| python | + pyproject.toml |
| typescript | + package.json, tsconfig.json |
| go | + go.mod, Makefile |

**Steps**:
1. Copy template files from `templates/<type>/`
2. Substitute variables in templates:
   - `{{project_name}}` - Project name
   - `{{description}}` - Project description
   - `{{module_path}}` - Go module path (for go template)
   - `{{created_date}}` - Current date
3. Generate `.claude/CLAUDE.md` from template
4. Commit initial structure

**Variable Substitution**:
```
{{project_name}} → my-awesome-app
{{description}} → A new beads-native project
{{created_date}} → 2025-01-15
{{repo_url}} → https://github.com/<user>/<project>
{{license}} → MIT
```

**Commands**:
```bash
# Create .claude directory
mkdir -p .claude

# Copy and process template files
# (Agent performs file operations)

# Initial commit
git add .
git commit -m "chore: Initialize project with beads"
```

**Output**:
```
✓ Template applied: <type>
✓ CLAUDE.md generated
✓ Initial commit created
```

### Brownfield File Handling

In brownfield mode, files are handled differently to preserve existing work:

| File Type | Action | Details |
|-----------|--------|---------|
| `.gitignore` | **Append** | Add beads patterns if not present |
| `CLAUDE.md` | **Merge** | Append beads section to existing |
| `README.md` | **Skip** | Never overwrite existing README |
| `pyproject.toml` | **Skip** | Preserve existing config |
| `package.json` | **Skip** | Preserve existing config |
| `go.mod` | **Skip** | Preserve existing config |
| `.beads/` | **Create** | Always create (core functionality) |

**Append Logic for .gitignore**:
```bash
# Check if beads patterns already exist
if ! grep -q "\.beads/beads\.db" .gitignore 2>/dev/null; then
  echo "" >> .gitignore
  echo "# Beads issue tracking" >> .gitignore
  echo ".beads/beads.db" >> .gitignore
  echo ".beads/beads.db-*" >> .gitignore
fi
```

**Output (brownfield)**:
```
✓ Appended beads patterns to .gitignore
✓ Merged beads section into CLAUDE.md
⊘ Skipped README.md (exists)
⊘ Skipped pyproject.toml (exists)
```

### CLAUDE.md Merge Logic

When a `.claude/CLAUDE.md` or `CLAUDE.md` already exists, merge rather than replace:

**Detection**:
```bash
# Check for existing CLAUDE.md
if [ -f ".claude/CLAUDE.md" ]; then
  CLAUDE_PATH=".claude/CLAUDE.md"
elif [ -f "CLAUDE.md" ]; then
  CLAUDE_PATH="CLAUDE.md"
else
  CLAUDE_PATH=""  # Will create new
fi
```

**Merge Strategy**:
1. Read existing CLAUDE.md content
2. Check if beads section already exists (look for `## Beads Workflow` or `## Issue Tracking`)
3. If not present, append beads section at end
4. Preserve all existing content

**Beads Section Template** (appended to existing):
```markdown

## Beads Workflow Integration

This project uses **Beads** for AI-native issue tracking.

### Essential Commands

```bash
bd ready                    # Show unblocked issues
bd list --status=open       # All open issues
bd update <id> --status=in_progress  # Claim work
bd close <id>               # Complete work
bd sync                     # Sync with git
```

### Session Protocol

Before completing work:
```bash
git status && git add <files> && bd sync && git commit -m "..." && git push
```

See: [Beads Documentation](https://github.com/steveyegge/beads)
```

**Output**:
```
✓ Found existing CLAUDE.md at ./CLAUDE.md
✓ Appended beads workflow section
```

**With --beads-only flag**:
```
✓ Beads initialized (--beads-only mode)
⊘ Skipped all template files
```

## Phase 4: Initial Issues

Create setup epic and child tasks from template.

**Source**: `templates/<type>/initial-issues.json`

**Steps**:
1. Parse initial-issues.json
2. Create epic issue
3. Create child task issues
4. Add dependencies (tasks depend on epic)

**Commands**:
```bash
# Create epic
bd create --title="<epic.title>" --description="<epic.description>" \
  --type=epic --priority=<epic.priority>

# Create tasks (for each task in tasks[])
bd create --title="<task.title>" --description="<task.description>" \
  --type=<task.type> --priority=<task.priority>

# Add dependencies
bd dep add <task-id> <epic-id>
```

**Output**:
```
✓ Created epic: <prefix>-xxx - <epic.title>
✓ Created task: <prefix>-yyy - <task.title>
✓ Created task: <prefix>-zzz - <task.title>
...
✓ Dependencies configured
```

## Phase 5: First Session (Optional)

Start the first work session unless --no-session specified.

**Command**:
```bash
# Run session start
/beads-session-start
```

**Output**:
```
✓ Project ready!

📋 Ready work:
1. [P1] <prefix>-xxx: <first-task>
...

Run 'bd update <id> --status in_progress' to claim an issue.
```

## Complete Example

```bash
# Create a Python project
/project-init my-api --template python --prefix api \
  --description "REST API for my application" \
  --agents beads-workflows

# Output:
✓ Directory created: ./my-api
✓ Git initialized
✓ Beads initialized with prefix: api
✓ Git hooks installed
✓ Agent plugins enabled: beads-workflows
✓ Template applied: python
✓ CLAUDE.md generated
✓ Initial commit created
✓ Created epic: api-xxx - Python Project Setup
✓ Created task: api-yyy - Configure Python virtual environment
✓ Created task: api-zzz - Set up pytest testing
✓ Created task: api-aaa - Configure linting and formatting
✓ Created task: api-bbb - Add type checking
✓ Dependencies configured
✓ Project ready!

📋 Ready work (4 issues):
1. [P1] api-yyy: Configure Python virtual environment
2. [P1] api-zzz: Set up pytest testing
3. [P2] api-aaa: Configure linting and formatting
4. [P2] api-bbb: Add type checking
```

## Brownfield Example

```bash
# Add beads to an existing TypeScript project
cd ~/projects/my-existing-app
/project-init --existing --no-issues

# Output:
✓ Brownfield mode: existing project detected
✓ Detected project type: typescript (from package.json)
✓ Git repository found
✓ Beads initialized with prefix: mea
✓ Git hooks installed
✓ Agent plugins enabled: beads-workflows
✓ Appended beads patterns to .gitignore
✓ Found existing CLAUDE.md at ./CLAUDE.md
✓ Appended beads workflow section
⊘ Skipped README.md (exists)
⊘ Skipped package.json (exists)
⊘ Skipped tsconfig.json (exists)
⊘ Skipped initial issues (--no-issues)
✓ Project ready for beads!

Next steps:
  bd create --title="First beads issue" --type=task
  bd ready
```

```bash
# Minimal beads-only setup (preserves everything)
/project-init -e --beads-only

# Output:
✓ Brownfield mode: existing project detected
✓ Beads initialized with prefix: proj
✓ Git hooks installed
✓ Beads initialized (--beads-only mode)
⊘ Skipped all template files
✓ Project ready for beads!
```

## Error Handling

| Error | Resolution |
|-------|------------|
| Directory exists with .beads | Ask user to use different name or --force |
| Git init fails | Check permissions, directory state |
| Template not found | Default to 'standard' with warning |
| bd init fails | Check beads installation |
| Not a git repo (brownfield) | Run `git init` first or let command initialize |
| CLAUDE.md has beads section | Skip merge, report already configured |
| --existing without git repo | Initialize git automatically with confirmation |
| Conflicting --beads-only with --template | Warn: --beads-only ignores template files |

## Files Created

### Greenfield (New Project)

After running `/project-init my-app --template python`:

```
my-app/
├── .beads/
│   ├── beads.db
│   ├── config.yaml
│   └── issues.jsonl
├── .claude/
│   └── CLAUDE.md
├── .git/
│   └── hooks/
│       └── post-commit (beads hook)
├── .gitignore
├── pyproject.toml
└── README.md
```

### Brownfield (Existing Project)

After running `/project-init --existing` in an existing project:

```
existing-project/
├── .beads/                    # NEW - Created
│   ├── beads.db
│   ├── config.yaml
│   └── issues.jsonl
├── .git/
│   └── hooks/
│       └── post-commit        # NEW - Beads hook added
├── .gitignore                 # MODIFIED - Beads patterns appended
├── CLAUDE.md                  # MODIFIED - Beads section appended
├── package.json               # UNCHANGED
├── tsconfig.json              # UNCHANGED
├── src/                       # UNCHANGED
│   └── ...
└── README.md                  # UNCHANGED
```

## Integration with beads-workflows

This command is part of the beads-workflows plugin and follows these patterns:

- **Session Rituals**: Sets up hooks for automatic sync
- **Description Quality**: Initial issues have proper Why/What/How structure
- **Dependency Thinking**: Tasks correctly depend on epic
- **Single-Issue Discipline**: Ready issues shown for focused work
